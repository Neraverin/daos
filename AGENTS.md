# DAOS Developer Guide

This document provides guidelines and instructions for agents working on the DAOS codebase.

## Project Overview

DAOS (Deployment and Orchestration System) is a Go-based system for managing host deployments and Docker Compose applications. It consists of:
- A daemon server (REST API via Gin)
- A TUI client (using Bubble Tea)
- SQLite database for persistence

## Build Commands

All build commands are in the Makefile:

```bash
# Build everything (daemon + TUI + packages)
make all

# Build only binaries
make build          # daemon + tui
make daemon         # daemon only
make tui            # TUI only

# Run tests
make test           # Run all tests
go test -v ./...    # Same as above

# Run a single test
go test -v ./pkg/db/...          # Test specific package
go test -v -run TestHealthCheck ./cmd/daemon/handlers/...  # Test specific function

# Code generation
make generate-api    # Generate Go types from OpenAPI spec (docs/openapi.yaml)
make generate-sql    # Generate sqlc code from SQL queries

# Install dependencies
make install-dependencies

# Clean build artifacts
make clean
```

## Code Style Guidelines

### General Principles

- Write clear, readable code with descriptive names
- Keep functions small and focused
- Handle errors explicitly and return meaningful error messages

### Naming Conventions

- **Variables/Functions**: `camelCase` (e.g., `createHost`, `parseTime`)
- **Constants**: `PascalCase` with descriptive names (e.g., `DeploymentStatusPending`)
- **Types/Structs**: `PascalCase` (e.g., `Server`, `Host`)
- **Database columns**: `snake_case` in schema, mapped to `PascalCase` in Go
- **Packages**: lowercase, single word or simple compound (e.g., `db`, `api`, `config`)
- **Files**: lowercase with underscores for multi-word names (e.g., `db_test.go`)

### Imports

Organize imports in three groups (standard library first, then third-party, then project):

```go
import (
    "database/sql"
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/Neraverin/daos/pkg/api"
    "github.com/Neraverin/daos/pkg/config"
    "github.com/Neraverin/daos/cmd/daemon/handlers"
    _ "github.com/mattn/go-sqlite3"
    "github.com/gin-gonic/gin"
)
```

### Error Handling

- Use `fmt.Errorf` with `%w` for wrapping errors
- Return errors with context (e.g., "failed to create host: %w")
- Handle errors at the appropriate level; don't ignore them
- Use helper functions like `api.ErrorJSON` for HTTP error responses

```go
// Good error handling
func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    // ...
}

// In HTTP handlers
func (s *Server) ListHosts(ctx *gin.Context) {
    hosts, err := s.db.GetAllHosts(ctx)
    if err != nil {
        api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    // ...
}
```

### Types and Structs

- Use struct tags for JSON/YAML serialization
- Use pointers for nullable fields (`*T`)
- Define constants for status values and enums

```go
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
}

const (
    DeploymentStatusPending  = "pending"
    DeploymentStatusRunning  = "running"
    DeploymentStatusSuccess  = "success"
    DeploymentStatusFailed   = "failed"
)
```

### Testing

- Test files are named `*_test.go` in the same package
- Use `testing.T` for assertions with `t.Fatalf` for setup failures and `t.Errorf` for assertions

### API Development

1. Define API in `docs/openapi.yaml` (OpenAPI 3.0)
2. Run `make generate-api` to generate Go types and handlers
3. Implement handler methods in `cmd/daemon/handlers/`
4. Use generated types; do not modify `pkg/api/openapi.gen.go`

### Database

- SQL schema in `pkg/db/migrations/`, queries in `pkg/db/queries/`
- Run `make generate-sql` to regenerate code
- Do not modify `pkg/db/*.sql.go` directly

### Configuration

- Example config in `config.example.yaml`
- Default values set programmatically (see `pkg/config/config.go`)
- Load config with `config.Load(path)` or use `config.Default()`

## Project Structure

```
cmd/
    daemon/           # REST API server
        handlers/     # HTTP handlers
        main.go
    tui/              # Terminal UI
        models/       # TUI models
        main.go
pkg/
    api/              # Generated API types (do not edit)
    ansible/          # Ansible execution logic
    config/           # Configuration loading
    db/               # Database layer (generated)
    models/           # Domain models
docs/
    openapi.yaml      # API specification
```
