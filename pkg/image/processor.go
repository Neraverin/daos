package image

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrFileNotFound     = fmt.Errorf("image file not found")
	ErrInvalidExtension = fmt.Errorf("image file must have .tar extension")
	ErrLoadTimeout      = fmt.Errorf("image load timed out")
	ErrLoadFailed       = fmt.Errorf("failed to load image")
	ErrPushFailed       = fmt.Errorf("failed to push image")
	ErrTagFailed        = fmt.Errorf("failed to tag image")
	ErrProcessCancelled = fmt.Errorf("process cancelled")
	ErrImageNotLoaded   = fmt.Errorf("image not loaded")
)

type CommandExecutor func(ctx context.Context, name string, args ...string) *exec.Cmd

func DefaultExecutor(ctx context.Context, name string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, args...)
}

type Processor struct {
	registry         string
	imageLoadTimeout int
	executor         CommandExecutor
}

type Config struct {
	Registry         string
	ImageLoadTimeout int
}

type Result struct {
	ImageID   string
	ImageName string
}

func NewProcessor(cfg Config) *Processor {
	timeout := cfg.ImageLoadTimeout
	if timeout == 0 {
		timeout = 300
	}
	return &Processor{
		registry:         cfg.Registry,
		imageLoadTimeout: timeout,
		executor:         DefaultExecutor,
	}
}

func (p *Processor) ValidateImageFile(rolePath, imageFile string) error {
	if imageFile == "" {
		return nil
	}

	fullPath := filepath.Join(rolePath, imageFile)

	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", ErrFileNotFound, fullPath)
	}
	if err != nil {
		return fmt.Errorf("failed to stat image file: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("image file path is a directory: %s", fullPath)
	}

	if strings.ToLower(filepath.Ext(fullPath)) != ".tar" {
		return fmt.Errorf("%w: %s", ErrInvalidExtension, filepath.Ext(fullPath))
	}

	return nil
}

func (p *Processor) ProcessImage(rolePath, imageFile string) (*Result, error) {
	if imageFile == "" {
		return nil, nil
	}

	if err := p.ValidateImageFile(rolePath, imageFile); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(rolePath, imageFile)

	result, err := p.loadImage(fullPath)
	if err != nil {
		return nil, err
	}

	if p.registry != "" {
		if err := p.pushImage(result.ImageID, p.registry); err != nil {
			return nil, err
		}
		result.ImageName = fmt.Sprintf("%s/%s", p.registry, result.ImageName)
	}

	return result, nil
}

func (p *Processor) loadImage(path string) (*Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.imageLoadTimeout)*time.Second)
	defer cancel()

	return p.loadImageWithContext(ctx, path)
}

func (p *Processor) loadImageWithContext(ctx context.Context, path string) (*Result, error) {
	cmd := p.executor(ctx, "docker", "load", "-i", path)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("%w: %d seconds", ErrLoadTimeout, p.imageLoadTimeout)
		}
		if ctx.Err() == context.Canceled {
			return nil, fmt.Errorf("%w", ErrProcessCancelled)
		}
		return nil, fmt.Errorf("%w: %s", ErrLoadFailed, string(output))
	}

	imageID := extractImageID(string(output))
	if imageID == "" {
		return nil, fmt.Errorf("%w: could not parse image ID from output", ErrLoadFailed)
	}

	imageName := extractImageName(imageID)

	return &Result{
		ImageID:   imageID,
		ImageName: imageName,
	}, nil
}

func (p *Processor) pushImage(imageID, registry string) error {
	imageName := extractImageName(imageID)
	taggedName := fmt.Sprintf("%s/%s", registry, imageName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tagCmd := p.executor(ctx, "docker", "tag", imageID, taggedName)
	tagOutput, err := tagCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrTagFailed, string(tagOutput))
	}

	pushCtx, pushCancel := context.WithTimeout(context.Background(), time.Duration(p.imageLoadTimeout)*time.Second)
	defer pushCancel()

	pushCmd := p.executor(pushCtx, "docker", "push", taggedName)
	pushOutput, err := pushCmd.CombinedOutput()
	if err != nil {
		if pushCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("%w: %d seconds", ErrLoadTimeout, p.imageLoadTimeout)
		}
		return fmt.Errorf("%w: %s", ErrPushFailed, string(pushOutput))
	}

	return nil
}

func extractImageID(output string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "Loaded image ID:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) >= 2 {
				id := strings.TrimSpace(parts[1])
				return id
			}
		}
		if strings.Contains(line, "Loaded image:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) >= 2 {
				name := strings.TrimSpace(parts[1])
				return name
			}
		}
	}
	return ""
}

func extractImageName(imageID string) string {
	parts := strings.Split(imageID, ":")
	if len(parts) >= 2 {
		return parts[0]
	}
	return imageID
}
