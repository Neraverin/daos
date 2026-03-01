package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Migrate(fsEmbed embed.FS) error {
	migrations, err := fs.ReadDir(fsEmbed, "pkg/db/migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return strings.Trim(migrations[i].Name(), ".sql") < strings.Trim(migrations[j].Name(), ".sql")
	})

	for _, migration := range migrations {
		if !strings.HasSuffix(migration.Name(), ".sql") {
			continue
		}

		content, err := fsEmbed.ReadFile(filepath.Join("pkg/db/migrations", migration.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", migration.Name(), err)
		}

		if _, err := m.db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Name(), err)
		}

		fmt.Printf("Applied migration: %s\n", migration.Name())
	}

	return nil
}

func (m *Migrator) EnsureTable() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TEXT NOT NULL DEFAULT (datetime('now'))
		)
	`)
	return err
}
