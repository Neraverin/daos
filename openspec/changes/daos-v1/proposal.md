## Why

DAOS (Deployment and Orchestration Service) is needed to automate deployment of Docker Compose applications to remote servers. Currently, manual SSH and manual deployment process is error-prone and not scalable. A unified service with TUI interface will provide a simple way to manage hosts, store Docker Compose packages, and trigger deployments via Ansible.

## What Changes

- Create new Go-based daemon with REST API
- Create new Go-based TUI client using bubbletea
- Define OpenAPI 3 contract for daemon API
- Implement SQLite database for persistent storage
- Integrate Ansible execution for remote deployments
- Create deb and rpm packages for distribution
- Configure systemd service for daemon

## Capabilities

### New Capabilities
- `host-management`: CRUD operations for remote hosts (SSH connection details)
- `package-management`: Upload and store Docker Compose files
- `deployment-management`: Create deployments, trigger Ansible runs, track status
- `log-streaming`: View deployment logs in real-time
- `daemon-api`: RESTful HTTP API with OpenAPI 3 specification

### Modified Capabilities
- (none - this is a new project)

## Impact

- New repository: DAOS installer application
- Dependencies: Go 1.25+, SQLite, bubbletea, Ansible (on target hosts)
- No breaking changes to existing systems (new project)
