package teststand

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Mode string

const (
	ModeLocal  Mode = "local"
	ModeRemote Mode = "remote"
	ModeAuto   Mode = "auto"
)

type Config struct {
	Hostname string
	Username string
	Password string
}

type TestStand struct {
	config *Config
	client *ssh.Client
	mode   Mode
	target Mode
}

var defaultStand *TestStand

func LoadEnv() (*Config, error) {
	envPath, err := findEnvFile()
	if err != nil {
		return nil, err
	}
	if envPath == "" {
		return nil, nil
	}

	return parseEnvFile(envPath)
}

func findEnvFile() (string, error) {
	dirs := []string{".", "..", "../.."}

	for _, dir := range dirs {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}
	}

	return "", nil
}

func parseEnvFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file: %w", err)
	}

	cfg := &Config{}
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "TEST_STAND_HOSTNAME":
			cfg.Hostname = value
		case "TEST_STAND_USERNAME":
			cfg.Username = value
		case "TEST_STAND_PASSWORD":
			cfg.Password = value
		}
	}

	return cfg, nil
}

func GetExecutionMode() Mode {
	switch os.Getenv("TEST_STAND_MODE") {
	case "remote":
		return ModeRemote
	case "local":
		return ModeLocal
	default:
		return ModeAuto
	}
}

func DetectTarget() Mode {
	cfg, err := LoadEnv()
	if err != nil {
		fmt.Printf("Warning: failed to load .env: %v\n", err)
		return ModeLocal
	}

	if cfg == nil {
		return ModeLocal
	}

	stand := &TestStand{
		config: cfg,
		mode:   ModeAuto,
	}

	if err := stand.Connect(); err != nil {
		fmt.Printf("Warning: failed to connect to test stand: %v\n", err)
		return ModeLocal
	}
	stand.Close()

	return ModeRemote
}

func GetTarget() Mode {
	mode := GetExecutionMode()

	switch mode {
	case ModeRemote:
		return ModeRemote
	case ModeLocal:
		return ModeLocal
	case ModeAuto:
		return DetectTarget()
	default:
		return ModeLocal
	}
}

func GetTestStand() (*TestStand, error) {
	if defaultStand != nil {
		return defaultStand, nil
	}

	cfg, err := LoadEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	defaultStand = &TestStand{
		config: cfg,
		mode:   GetExecutionMode(),
		target: GetTarget(),
	}

	if defaultStand.target == ModeRemote {
		if err := defaultStand.Connect(); err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}
	}

	return defaultStand, nil
}

func NewTestStand(cfg *Config) *TestStand {
	return &TestStand{
		config: cfg,
		mode:   GetExecutionMode(),
		target: GetTarget(),
	}
}

func (ts *TestStand) Connect() error {
	if ts.config == nil {
		return fmt.Errorf("no config provided")
	}

	if ts.config.Hostname == "" || ts.config.Username == "" || ts.config.Password == "" {
		return fmt.Errorf("incomplete config: hostname, username, and password are required")
	}

	clientConfig := &ssh.ClientConfig{
		User: ts.config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(ts.config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ts.config.Hostname), clientConfig)
	if err != nil {
		return fmt.Errorf("failed to connect via SSH: %w", err)
	}

	ts.client = client
	return nil
}

func (ts *TestStand) Execute(cmd string) (string, error) {
	if ts.client == nil {
		return "", fmt.Errorf("not connected")
	}

	session, err := ts.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(cmd); err != nil {
		return "", fmt.Errorf("command failed: %w (stderr: %s)", err, stderr.String())
	}

	return stdout.String(), nil
}

func (ts *TestStand) Close() {
	if ts.client != nil {
		ts.client.Close()
		ts.client = nil
	}
}

func (ts *TestStand) IsConnected() bool {
	return ts.client != nil
}

func (ts *TestStand) Target() Mode {
	return ts.target
}

func runCommandLocal(cmd string) (string, error) {
	c := exec.Command("sh", "-c", cmd)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	if err := c.Run(); err != nil {
		return "", fmt.Errorf("command failed: %w (stderr: %s)", err, stderr.String())
	}

	return stdout.String(), nil
}

func RunLocal(cmd string) (string, error) {
	return runCommandLocal(cmd)
}

func RunRemote(cmd string) (string, error) {
	stand, err := GetTestStand()
	if err != nil {
		return "", fmt.Errorf("failed to get test stand: %w", err)
	}
	return stand.Execute(cmd)
}

func Run(cmd string) (string, error) {
	target := GetTarget()
	if target == ModeRemote {
		return RunRemote(cmd)
	}
	return RunLocal(cmd)
}
