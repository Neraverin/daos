## Context

DAOS is a new installer service that provides:
- A daemon running as a systemd service with REST API
- A TUI client that communicates with the daemon via HTTP
- Deployment of Docker Compose applications via Ansible

### Current State
- Greenfield project - no existing codebase
- No authentication required (localhost-trusted)

### Constraints
- Go 1.25+ for both daemon and TUI
- SQLite for persistent storage
- bubbletea for TUI
- OpenAPI 3 for API contracts
- deb and rpm packaging for distribution
- systemd service for daemon

### Stakeholders
- System administrators managing remote servers
- DevOps engineers deploying applications

## Goals / Non-Goals

**Goals:**
- Implement daemon with REST API for host, package, and deployment management
- Implement TUI client with CRUD operations for all resources
- Define OpenAPI 3 contract and generate Go client/server
- Package daemon and TUI as deb and rpm packages
- Provide systemd service for daemon

**Non-Goals:**
- Authentication/authorization (localhost-trusted for v1)
- High availability / clustering
- Docker Swarm or Kubernetes integration
- Repository management for compose files
- Deployment scheduling/queuing

## Decisions

### 1. HTTP Framework: Gin over Fiber
- **Choice**: Gin
- **Rationale**: More mature, better documentation, widely used in Go ecosystem

### 2. TUI Framework: bubbletea
- **Choice**: bubbletea
- **Rationale**: Pure Go, Elm-like architecture, good for interactive CLIs

### 3. Database: SQLite with sqlc
- **Choice**: SQLite + sqlc for code generation
- **Rationale**: Single-file DB, no external dependencies, type-safe SQL

### 4. API Client Generation: oapi-codegen
- **Choice**: oapi-codegen
- **Rationale**: Generates both client and server from OpenAPI spec

### 5. Package Storage: Inline in SQLite
- **Choice**: Store compose content as TEXT in SQLite
- **Rationale**: Simplest for v1, no file sync needed

### 6. Ansible Execution: os/exec direct invocation
- **Choice**: Direct exec of ansible-playbook
- **Rationale**: Simple, no additional dependencies

### 7. Project Structure: Single module
- **Choice**: Single go.mod with cmd/daemon and cmd/tui
- **Rationale**: Simpler build/packaging for v1

## Risks / Trade-offs

### Risk: Ansible dependency on daemon host
- **Mitigation**: Document requirement; daemon host must have ansible installed

### Risk: SSH key security
- **Mitigation**: Only store path to key, not credentials; document proper permissions

### Risk: Parallel deployments to same host
- **Mitigation**: Allow for v1; let users manage coordination externally

### Risk: Large compose files
- **Mitigation**: Set reasonable size limit (10MB) for stored compose content

## Migration Plan

1. Build daemon and TUI binaries
2. Create packages (deb/rpm)
3. Install daemon package (installs binary, creates user, sets up systemd)
4. Install TUI package
5. Start daemon service
6. TUI connects to localhost:8080

### Rollback
- Uninstall packages removes binaries, service and all data in `/opt/daos`

## Open Questions

1. Should daemon listen on specific IP or all interfaces?
   - Recommendation: 127.0.0.1 (localhost) for security

2. Default port?
   - Recommendation: 8080 (configurable via config file)

3. Where to store SQLite database?
   - Recommendation: `/opt/daos/daos.db`

4. Config file location?
   - Recommendation: `/opt/daos/configs/daemon.yaml`

5. Log file location?
   - Recommendation: `/opt/daos/logs/daemon.log`
