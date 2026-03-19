package teststand

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	envContent := `TEST_STAND_HOSTNAME=example.com
TEST_STAND_USERNAME=testuser
TEST_STAND_PASSWORD=testpass
`
	err := os.WriteFile(envPath, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test env file: %v", err)
	}

	cfg, err := parseEnvFile(envPath)
	if err != nil {
		t.Fatalf("failed to parse env file: %v", err)
	}

	if cfg.Hostname != "example.com" {
		t.Errorf("expected hostname 'example.com', got '%s'", cfg.Hostname)
	}
	if cfg.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", cfg.Username)
	}
	if cfg.Password != "testpass" {
		t.Errorf("expected password 'testpass', got '%s'", cfg.Password)
	}
}

func TestParseEnvFile_WithComments(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	envContent := `# This is a comment
TEST_STAND_HOSTNAME=example.com
# Another comment
TEST_STAND_USERNAME=testuser
TEST_STAND_PASSWORD=testpass
`
	err := os.WriteFile(envPath, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test env file: %v", err)
	}

	cfg, err := parseEnvFile(envPath)
	if err != nil {
		t.Fatalf("failed to parse env file: %v", err)
	}

	if cfg.Hostname != "example.com" {
		t.Errorf("expected hostname 'example.com', got '%s'", cfg.Hostname)
	}
}

func TestParseEnvFile_EmptyValues(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	envContent := `TEST_STAND_HOSTNAME=
TEST_STAND_USERNAME=
TEST_STAND_PASSWORD=
`
	err := os.WriteFile(envPath, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test env file: %v", err)
	}

	cfg, err := parseEnvFile(envPath)
	if err != nil {
		t.Fatalf("failed to parse env file: %v", err)
	}

	if cfg.Hostname != "" {
		t.Errorf("expected empty hostname, got '%s'", cfg.Hostname)
	}
}

func TestParseEnvFile_MissingFile(t *testing.T) {
	_, err := parseEnvFile("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestGetExecutionMode_Remote(t *testing.T) {
	os.Setenv("TEST_STAND_MODE", "remote")
	defer os.Unsetenv("TEST_STAND_MODE")

	mode := GetExecutionMode()
	if mode != ModeRemote {
		t.Errorf("expected ModeRemote, got %s", mode)
	}
}

func TestGetExecutionMode_Local(t *testing.T) {
	os.Setenv("TEST_STAND_MODE", "local")
	defer os.Unsetenv("TEST_STAND_MODE")

	mode := GetExecutionMode()
	if mode != ModeLocal {
		t.Errorf("expected ModeLocal, got %s", mode)
	}
}

func TestGetExecutionMode_Auto(t *testing.T) {
	os.Unsetenv("TEST_STAND_MODE")

	mode := GetExecutionMode()
	if mode != ModeAuto {
		t.Errorf("expected ModeAuto, got %s", mode)
	}
}

func TestGetExecutionMode_Invalid(t *testing.T) {
	os.Setenv("TEST_STAND_MODE", "invalid")
	defer os.Unsetenv("TEST_STAND_MODE")

	mode := GetExecutionMode()
	if mode != ModeAuto {
		t.Errorf("expected ModeAuto for invalid value, got %s", mode)
	}
}

func TestLoadEnv_NoEnvFile(t *testing.T) {
	originalCwd, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(originalCwd)

	for _, name := range []string{".", "..", "../.."} {
		envPath := filepath.Join(name, ".env")
		os.Remove(envPath)
	}

	cfg, err := LoadEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg != nil {
		t.Errorf("expected nil config when no env file, got %+v", cfg)
	}
}

func TestLoadEnv_WithEnvFile(t *testing.T) {
	originalCwd, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(originalCwd)

	envContent := `TEST_STAND_HOSTNAME=test.example.com
TEST_STAND_USERNAME=admin
TEST_STAND_PASSWORD=secret
`
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	cfg, err := LoadEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected config, got nil")
	}
	if cfg.Hostname != "test.example.com" {
		t.Errorf("expected 'test.example.com', got '%s'", cfg.Hostname)
	}
}

func TestTestStand_New(t *testing.T) {
	cfg := &Config{
		Hostname: "example.com",
		Username: "user",
		Password: "pass",
	}

	stand := NewTestStand(cfg)
	if stand == nil {
		t.Fatal("expected TestStand, got nil")
	}
	if stand.config != cfg {
		t.Error("config not set correctly")
	}
}

func TestTestStand_IsConnected(t *testing.T) {
	stand := &TestStand{}
	if stand.IsConnected() {
		t.Error("expected not connected initially")
	}
}

func TestRunLocal(t *testing.T) {
	output, err := RunLocal("echo hello")
	if err != nil {
		t.Fatalf("RunLocal failed: %v", err)
	}
	if output != "hello\n" {
		t.Errorf("expected 'hello\\n', got '%s'", output)
	}
}

func TestRunLocal_CommandFailure(t *testing.T) {
	_, err := RunLocal("exit 1")
	if err == nil {
		t.Error("expected error for failing command, got nil")
	}
}

func TestRunRemote_NoConnection(t *testing.T) {
	os.Setenv("TEST_STAND_MODE", "local")
	defer os.Unsetenv("TEST_STAND_MODE")

	stand, err := GetTestStand()
	if err != nil {
		t.Fatalf("GetTestStand failed: %v", err)
	}

	if stand.Target() != ModeLocal {
		t.Errorf("expected ModeLocal, got %s", stand.Target())
	}
}

func TestRunCommand(t *testing.T) {
	target := GetTarget()
	if target != ModeLocal && target != ModeRemote {
		t.Errorf("expected ModeLocal or ModeRemote, got %s", target)
	}

	output, err := Run("echo test")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if output != "test\n" {
		t.Errorf("expected 'test\\n', got '%s'", output)
	}
}

func TestDetectTarget_NoEnv(t *testing.T) {
	os.Unsetenv("TEST_STAND_MODE")

	cfg := &Config{}
	stand := &TestStand{config: cfg}
	if err := stand.Connect(); err == nil {
		t.Error("expected error when no config")
	}
}
