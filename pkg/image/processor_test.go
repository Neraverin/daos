package image

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateImageFile_EmptyImageFile(t *testing.T) {
	p := NewProcessor(Config{})

	err := p.ValidateImageFile("/role/path", "")
	if err != nil {
		t.Errorf("expected nil error for empty imageFile, got %v", err)
	}
}

func TestValidateImageFile_ValidTarFile(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.tar")
	if err := os.WriteFile(imagePath, []byte("fake tar content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	p := NewProcessor(Config{})

	err := p.ValidateImageFile(tmpDir, "test.tar")
	if err != nil {
		t.Errorf("expected nil error for valid tar file, got %v", err)
	}
}

func TestValidateImageFile_CaseInsensitiveExtension(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.TAR")
	if err := os.WriteFile(imagePath, []byte("fake tar content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	p := NewProcessor(Config{})

	err := p.ValidateImageFile(tmpDir, "test.TAR")
	if err != nil {
		t.Errorf("expected nil error for .TAR extension, got %v", err)
	}
}

func TestValidateImageFile_MissingFile(t *testing.T) {
	p := NewProcessor(Config{})

	err := p.ValidateImageFile("/nonexistent", "missing.tar")
	if err == nil {
		t.Errorf("expected error for missing file, got nil")
	}
}

func TestValidateImageFile_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	p := NewProcessor(Config{})

	err := p.ValidateImageFile(tmpDir, filepath.Base(tmpDir))
	if err == nil {
		t.Errorf("expected error for directory, got nil")
	}
}

func TestValidateImageFile_InvalidExtension(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(imagePath, []byte("text content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	p := NewProcessor(Config{})

	err := p.ValidateImageFile(tmpDir, "test.txt")
	if err == nil {
		t.Errorf("expected error for non-tar file, got nil")
	}
}

func TestValidateImageFile_InvalidExtensionUppercase(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.PNG")
	if err := os.WriteFile(imagePath, []byte("image content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	p := NewProcessor(Config{})

	err := p.ValidateImageFile(tmpDir, "test.PNG")
	if err == nil {
		t.Errorf("expected error for non-tar file, got nil")
	}
}

func TestExtractImageID_ImageID(t *testing.T) {
	output := "Loaded image ID: sha256:abc123def456"
	id := extractImageID(output)
	expected := "sha256:abc123def456"
	if id != expected {
		t.Errorf("expected %s, got %s", expected, id)
	}
}

func TestExtractImageID_ImageName(t *testing.T) {
	output := "Loaded image: myimage:latest"
	id := extractImageID(output)
	expected := "myimage:latest"
	if id != expected {
		t.Errorf("expected %s, got %s", expected, id)
	}
}

func TestExtractImageName_WithTag(t *testing.T) {
	name := extractImageName("myimage:v1.0")
	expected := "myimage"
	if name != expected {
		t.Errorf("expected %s, got %s", expected, name)
	}
}

func TestExtractImageName_WithoutTag(t *testing.T) {
	name := extractImageName("sha256:abc123")
	expected := "sha256"
	if name != expected {
		t.Errorf("expected %s, got %s", expected, name)
	}
}

func TestExtractImageName_NoColon(t *testing.T) {
	name := extractImageName("plainname")
	expected := "plainname"
	if name != expected {
		t.Errorf("expected %s, got %s", expected, name)
	}
}

func TestLoadImage_Timeout(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.tar")
	if err := os.WriteFile(imagePath, []byte("fake tar"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	slowExecutor := func(ctx context.Context, name string, args ...string) *exec.Cmd {
		cmd := exec.CommandContext(ctx, "sleep", "10")
		return cmd
	}

	p := &Processor{
		registry:         "",
		imageLoadTimeout: 1,
		executor:         slowExecutor,
	}

	_, err := p.loadImage(imagePath)
	if err == nil {
		t.Errorf("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("expected timeout error message, got: %v", err)
	}
}

func TestLoadImage_CommandFailure(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.tar")
	if err := os.WriteFile(imagePath, []byte("fake tar"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	failExecutor := func(ctx context.Context, name string, args ...string) *exec.Cmd {
		cmd := exec.CommandContext(ctx, "sh", "-c", "echo 'error loading image' >&2; exit 1")
		return cmd
	}

	p := &Processor{
		registry:         "",
		imageLoadTimeout: 300,
		executor:         failExecutor,
	}

	_, err := p.loadImage(imagePath)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to load image") {
		t.Errorf("expected 'failed to load image' error, got: %v", err)
	}
}

func TestLoadImage_InvalidOutput(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.tar")
	if err := os.WriteFile(imagePath, []byte("fake tar"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	noIdExecutor := func(ctx context.Context, name string, args ...string) *exec.Cmd {
		cmd := exec.CommandContext(ctx, "sh", "-c", "echo 'Loaded something but no ID'")
		return cmd
	}

	p := &Processor{
		registry:         "",
		imageLoadTimeout: 300,
		executor:         noIdExecutor,
	}

	_, err := p.loadImage(imagePath)
	if err == nil {
		t.Errorf("expected error for invalid output, got nil")
	}
	if !strings.Contains(err.Error(), "could not parse image ID") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestProcessImage_EmptyImageFile(t *testing.T) {
	p := NewProcessor(Config{
		Registry:         "registry.example.com",
		ImageLoadTimeout: 300,
	})

	result, err := p.ProcessImage("/some/path", "")
	if err != nil {
		t.Errorf("expected nil error for empty imageFile, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result for empty imageFile, got %v", result)
	}
}

func TestProcessImage_ValidationError(t *testing.T) {
	p := NewProcessor(Config{
		Registry:         "registry.example.com",
		ImageLoadTimeout: 300,
	})

	_, err := p.ProcessImage("/nonexistent", "missing.tar")
	if err == nil {
		t.Errorf("expected validation error, got nil")
	}
}

func TestProcessImage_SkipPushWhenNoRegistry(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.tar")
	if err := os.WriteFile(imagePath, []byte("fake tar"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	loadCalled := false
	loadExecutor := func(ctx context.Context, name string, args ...string) *exec.Cmd {
		loadCalled = true
		cmd := exec.CommandContext(ctx, "sh", "-c", "echo 'Loaded image: testimage:latest'")
		return cmd
	}

	p := &Processor{
		registry:         "",
		imageLoadTimeout: 300,
		executor:         loadExecutor,
	}

	result, err := p.ProcessImage(tmpDir, "test.tar")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if !loadCalled {
		t.Errorf("expected load to be called")
	}
	if result == nil {
		t.Errorf("expected result, got nil")
	}
	if result != nil && result.ImageName != "testimage" {
		t.Errorf("expected ImageName 'testimage', got %s", result.ImageName)
	}
}

func TestPushImage_TagFailure(t *testing.T) {
	failTagExecutor := func(ctx context.Context, name string, args ...string) *exec.Cmd {
		if name == "docker" && len(args) > 0 && args[0] == "tag" {
			cmd := exec.CommandContext(ctx, "sh", "-c", "echo 'tag failed' >&2; exit 1")
			return cmd
		}
		cmd := exec.CommandContext(ctx, name, args...)
		return cmd
	}

	p := &Processor{
		registry:         "registry.example.com",
		imageLoadTimeout: 300,
		executor:         failTagExecutor,
	}

	err := p.pushImage("sha256:abc123", "registry.example.com")
	if err == nil {
		t.Errorf("expected tag error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to tag image") {
		t.Errorf("expected tag error, got: %v", err)
	}
}

func TestPushImage_PushFailure(t *testing.T) {
	failPushExecutor := func(ctx context.Context, name string, args ...string) *exec.Cmd {
		if name == "docker" && len(args) > 0 && args[0] == "tag" {
			cmd := exec.CommandContext(ctx, "sh", "-c", "echo 'tagged'")
			return cmd
		}
		if name == "docker" && len(args) > 0 && args[0] == "push" {
			cmd := exec.CommandContext(ctx, "sh", "-c", "echo 'push failed' >&2; exit 1")
			return cmd
		}
		cmd := exec.CommandContext(ctx, name, args...)
		return cmd
	}

	p := &Processor{
		registry:         "registry.example.com",
		imageLoadTimeout: 300,
		executor:         failPushExecutor,
	}

	err := p.pushImage("testimage:latest", "registry.example.com")
	if err == nil {
		t.Errorf("expected push error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to push image") {
		t.Errorf("expected push error, got: %v", err)
	}
}

func TestNewProcessor_DefaultTimeout(t *testing.T) {
	p := NewProcessor(Config{})

	if p.imageLoadTimeout != 300 {
		t.Errorf("expected default timeout 300, got %d", p.imageLoadTimeout)
	}
}

func TestNewProcessor_CustomTimeout(t *testing.T) {
	p := NewProcessor(Config{
		ImageLoadTimeout: 600,
	})

	if p.imageLoadTimeout != 600 {
		t.Errorf("expected timeout 600, got %d", p.imageLoadTimeout)
	}
}

func TestErrorTypes(t *testing.T) {
	tests := []struct {
		err      error
		expected string
	}{
		{ErrFileNotFound, "image file not found"},
		{ErrInvalidExtension, "image file must have .tar extension"},
		{ErrLoadTimeout, "image load timed out"},
		{ErrLoadFailed, "failed to load image"},
		{ErrPushFailed, "failed to push image"},
		{ErrTagFailed, "failed to tag image"},
	}

	for _, tc := range tests {
		if tc.err.Error() != tc.expected {
			t.Errorf("expected '%s', got '%s'", tc.expected, tc.err.Error())
		}
	}
}

func TestIsImageError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		isImageErr bool
	}{
		{"ErrFileNotFound", ErrFileNotFound, true},
		{"ErrInvalidExtension", ErrInvalidExtension, true},
		{"ErrLoadTimeout", ErrLoadTimeout, true},
		{"ErrLoadFailed", ErrLoadFailed, true},
		{"Other error", errors.New("other error"), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			isErr := tc.err == ErrFileNotFound ||
				tc.err == ErrInvalidExtension ||
				tc.err == ErrLoadTimeout ||
				tc.err == ErrLoadFailed
			if isErr != tc.isImageErr {
				t.Errorf("expected %v, got %v", tc.isImageErr, isErr)
			}
		})
	}
}

func TestContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.tar")
	if err := os.WriteFile(imagePath, []byte("fake tar"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	slowExecutor := func(ctx context.Context, name string, args ...string) *exec.Cmd {
		cmd := exec.CommandContext(ctx, "sleep", "10")
		return cmd
	}

	p := &Processor{
		registry:         "",
		imageLoadTimeout: 300,
		executor:         slowExecutor,
	}

	_, err := p.loadImageWithContext(cancelledCtx, imagePath)
	if err == nil {
		t.Errorf("expected context cancelled error, got nil")
	}
}

func TestLoadImageWithContext(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.tar")
	if err := os.WriteFile(imagePath, []byte("fake tar"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	successExecutor := func(ctx context.Context, name string, args ...string) *exec.Cmd {
		cmd := exec.CommandContext(ctx, "sh", "-c", "echo 'Loaded image: testimage:v1'")
		return cmd
	}

	p := &Processor{
		registry:         "",
		imageLoadTimeout: 300,
		executor:         successExecutor,
	}

	ctx := context.Background()
	result, err := p.loadImageWithContext(ctx, imagePath)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if result == nil {
		t.Fatalf("expected result, got nil")
	}
	if result.ImageName != "testimage" {
		t.Errorf("expected ImageName 'testimage', got %s", result.ImageName)
	}
}
