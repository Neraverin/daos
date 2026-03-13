package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	version = "0.1.0"
)

func requireRoot(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Fatal("This test must be run as root")
	}
}

func detectOS(t *testing.T) string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		t.Skipf("Cannot detect OS: %v", err)
		return ""
	}

	content := string(data)
	if strings.Contains(content, `ID="debian"`) || strings.Contains(content, `ID="ubuntu"`) {
		return "deb"
	}
	if strings.Contains(content, `ID="rhel"`) || strings.Contains(content, `ID="centos"`) || strings.Contains(content, `ID="fedora"`) {
		return "rpm"
	}

	t.Skipf("Unsupported OS. Only debian/ubuntu and rhel/centos/fedora are supported")
	return ""
}

func getProjectRoot() string {
	wd, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return "."
		}
		wd = parent
	}
}

func getPackagePath(osType string) string {
	projectRoot := getProjectRoot()
	switch osType {
	case "deb":
		return filepath.Join(projectRoot, "bin", "daos_"+version+"_amd64.deb")
	case "rpm":
		return filepath.Join(projectRoot, "bin", "daos-"+version+"-1.x86_64.rpm")
	default:
		return ""
	}
}

func installDAOS(t *testing.T, osType string) {
	pkgPath := getPackagePath(osType)

	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		t.Fatalf("Package not found: %s", pkgPath)
	}

	var cmd *exec.Cmd
	switch osType {
	case "deb":
		cmd = exec.Command("dpkg", "-i", pkgPath)
	case "rpm":
		cmd = exec.Command("rpm", "-ivh", pkgPath)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to install package: %v", err)
	}

	// Reload systemd and start service
	exec.Command("systemctl", "daemon-reload").Run()
	exec.Command("systemctl", "start", "daos").Run()
}

func verifyInstallation(t *testing.T) {
	// Check /opt/daos exists
	if _, err := os.Stat("/opt/daos"); os.IsNotExist(err) {
		t.Fatal("/opt/daos directory does not exist")
	}

	// Check daemon executable
	if _, err := os.Stat("/opt/daos/daemon"); os.IsNotExist(err) {
		t.Fatal("/opt/daos/daemon does not exist")
	}

	// Check tui executable
	if _, err := os.Stat("/opt/daos/tui"); os.IsNotExist(err) {
		t.Fatal("/opt/daos/tui does not exist")
	}

	// Check registry directory
	if _, err := os.Stat("/opt/daos/registry"); os.IsNotExist(err) {
		t.Fatal("/opt/daos/registry does not exist")
	}

	// Check daemon is running
	cmd := exec.Command("systemctl", "is-active", "daos")
	output, err := cmd.Output()
	if err != nil || strings.TrimSpace(string(output)) != "active" {
		t.Fatalf("DAOS daemon is not running. Output: %s, Error: %v", string(output), err)
	}

	// Check docker registry container is running
	cmd = exec.Command("docker", "ps")
	output, err = cmd.Output()
	if err != nil {
		t.Fatalf("Failed to check docker containers: %v", err)
	}
	if !strings.Contains(string(output), "daos-registry") {
		t.Fatal("Docker registry container is not running")
	}
}

func uninstallDAOS(t *testing.T, osType string) {
	// Stop daemon first
	exec.Command("systemctl", "stop", "daos").Run()

	var cmd *exec.Cmd
	switch osType {
	case "deb":
		cmd = exec.Command("dpkg", "-r", "daos")
	case "rpm":
		cmd = exec.Command("rpm", "-e", "daos")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Logf("Warning: Failed to uninstall package: %v", err)
	}

	// Reload systemd
	exec.Command("systemctl", "daemon-reload").Run()
}

func verifyUninstallation(t *testing.T) {
	// Check /opt/daos is removed
	if _, err := os.Stat("/opt/daos"); !os.IsNotExist(err) {
		t.Fatal("/opt/daos directory still exists after uninstall")
	}

	// Check daemon is not running
	cmd := exec.Command("systemctl", "is-active", "daos")
	output, err := cmd.Output()
	if err == nil && strings.TrimSpace(string(output)) == "active" {
		t.Fatal("DAOS daemon is still running after uninstall")
	}
}

func TestPackageInstallUninstall(t *testing.T) {
	requireRoot(t)

	osType := detectOS(t)
	if osType == "" {
		t.Skip("OS detection failed or unsupported OS")
	}

	// Install
	installDAOS(t, osType)
	verifyInstallation(t)

	// Ensure cleanup on failure
	t.Cleanup(func() {
		uninstallDAOS(t, osType)
	})

	// Uninstall
	uninstallDAOS(t, osType)
	verifyUninstallation(t)
}
