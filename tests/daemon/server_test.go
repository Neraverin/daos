package daemon

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/Neraverin/daos/cmd/daemon/handlers"
	"github.com/Neraverin/daos/pkg/api"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type testServer struct {
	router *gin.Engine
	server *handlers.Server
	db     *sql.DB
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

func setupTestServer(t *testing.T) *testServer {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	db.Exec("PRAGMA foreign_keys = ON")

	projectRoot := getProjectRoot()

	// Create tables inline
	schema := `
	CREATE TABLE IF NOT EXISTS hosts (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		hostname TEXT NOT NULL,
		port INTEGER NOT NULL DEFAULT 22,
		username TEXT NOT NULL,
		ssh_key_path TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS packages (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		compose_content TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS deployments (
		id TEXT PRIMARY KEY,
		host_id TEXT NOT NULL,
		package_id TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		FOREIGN KEY (host_id) REFERENCES hosts(id) ON DELETE CASCADE,
		FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS logs (
		id TEXT PRIMARY KEY,
		deployment_id TEXT NOT NULL,
		timestamp TEXT NOT NULL DEFAULT (datetime('now')),
		message TEXT NOT NULL,
		FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE
	);
	`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	server := handlers.New(db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.StaticFile("/openapi.yaml", filepath.Join(projectRoot, "docs/openapi.yaml"))
	api.RegisterHandlers(router, server)

	return &testServer{
		router: router,
		server: server,
		db:     db,
	}
}

func (ts *testServer) close() {
	ts.db.Close()
}
