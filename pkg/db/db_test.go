package db

import (
	"testing"
	"os"
)

func TestNew(t *testing.T) {
	tmpFile := "/tmp/test_daos.db"
	defer os.Remove(tmpFile)

	db, err := New(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestCreateAndQueryHost(t *testing.T) {
	tmpFile := "/tmp/test_daos_host.db"
	defer os.Remove(tmpFile)

	db, err := New(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS hosts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			hostname TEXT NOT NULL,
			port INTEGER NOT NULL DEFAULT 22,
			username TEXT NOT NULL,
			ssh_key_path TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now'))
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	result, err := db.Exec(
		"INSERT INTO hosts (name, hostname, port, username, ssh_key_path) VALUES (?, ?, ?, ?, ?)",
		"test-server", "192.168.1.100", 22, "root", "/home/user/.ssh/id_rsa",
	)
	if err != nil {
		t.Fatalf("Failed to insert host: %v", err)
	}

	id, _ := result.LastInsertId()
	if id != 1 {
		t.Fatalf("Expected id 1, got %d", id)
	}

	var name, hostname, username, sshKeyPath string
	var port int
	err = db.QueryRow("SELECT name, hostname, port, username, ssh_key_path FROM hosts WHERE id = ?", id).
		Scan(&name, &hostname, &port, &username, &sshKeyPath)
	if err != nil {
		t.Fatalf("Failed to query host: %v", err)
	}

	if name != "test-server" {
		t.Errorf("Expected name 'test-server', got '%s'", name)
	}
	if hostname != "192.168.1.100" {
		t.Errorf("Expected hostname '192.168.1.100', got '%s'", hostname)
	}
	if port != 22 {
		t.Errorf("Expected port 22, got %d", port)
	}
	if username != "root" {
		t.Errorf("Expected username 'root', got '%s'", username)
	}
}
