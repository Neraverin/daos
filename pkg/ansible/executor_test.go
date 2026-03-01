package ansible

import (
	"testing"
)

func TestNewExecutor(t *testing.T) {
	exec := NewExecutor("192.168.1.100", 22, "root", "/home/user/.ssh/id_rsa", "version: '3'")

	if exec.hostname != "192.168.1.100" {
		t.Errorf("Expected hostname '192.168.1.100', got '%s'", exec.hostname)
	}
	if exec.port != 22 {
		t.Errorf("Expected port 22, got %d", exec.port)
	}
	if exec.username != "root" {
		t.Errorf("Expected username 'root', got '%s'", exec.username)
	}
	if exec.sshKeyPath != "/home/user/.ssh/id_rsa" {
		t.Errorf("Expected sshKeyPath '/home/user/.ssh/id_rsa', got '%s'", exec.sshKeyPath)
	}
	if exec.composeContent != "version: '3'" {
		t.Errorf("Expected composeContent \"version: '3'\", got '%s'", exec.composeContent)
	}
}
