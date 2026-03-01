## 1. Project Setup

- [x] 1.1 Initialize Go module with go.mod
- [x] 1.2 Create project directory structure (cmd/, pkg/, api/, packaging/)
- [x] 1.3 Set up Makefile for building daemon and TUI

## 2. Database Layer

- [x] 2.1 Define SQLite schema (hosts, packages, deployments, logs tables)
- [x] 2.2 Create database migration mechanism
- [x] 2.3 Implement db package with sqlc
- [x] 2.4 Write repository functions for each entity

## 3. OpenAPI Contract

- [x] 3.1 Define OpenAPI 3 spec for all endpoints (hosts, packages, deployments, logs)
- [x] 3.2 Generate Go server stubs using oapi-codegen
- [x] 3.3 Generate Go client stubs for TUI

## 4. Daemon Implementation

- [x] 4.1 Create daemon main.go entry point
- [x] 4.2 Implement HTTP handlers for hosts CRUD
- [x] 4.3 Implement HTTP handlers for packages CRUD
- [x] 4.4 Implement HTTP handlers for deployments CRUD
- [x] 4.5 Implement HTTP handlers for logs retrieval
- [x] 4.6 Implement health check endpoint
- [x] 4.7 Add configuration loading (config.yaml)

## 5. Ansible Integration

- [x] 5.1 Create ansible executor package
- [x] 5.2 Implement inventory generation per host
- [x] 5.3 Implement compose deployment execution
- [x] 5.4 Implement log capture from Ansible output

## 6. TUI Implementation

- [x] 6.1 Create TUI main.go entry point
- [x] 6.2 Implement hosts list screen
- [x] 6.3 Implement host add/edit form
- [x] 6.4 Implement packages list screen
- [x] 6.5 Implement package upload form
- [x] 6.6 Implement deployments list screen
- [x] 6.7 Implement deployment create form
- [x] 6.8 Implement deployment run action
- [x] 6.9 Implement logs view screen

## 7. Packaging

- [x] 7.1 Create systemd service unit file
- [x] 7.2 Create deb packaging (control, postinst, prerm)
- [x] 7.3 Create rpm packaging (spec file)
- [x] 7.4 Configure build for both architectures (amd64, arm64)

## 8. Testing

- [ ] 8.1 Write unit tests for db layer
- [ ] 8.2 Write unit tests for handlers
- [ ] 8.3 Write integration tests for API
- [ ] 8.4 Manual testing of TUI flows
