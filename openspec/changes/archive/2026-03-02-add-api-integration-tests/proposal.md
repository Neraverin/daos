## Why

The DAOS API lacks automated tests to verify that all endpoints work as defined in the `daemon-api` specification. Currently, only a single health check test exists. Without comprehensive API tests, we risk regressions going undetected and cannot confidently refactor or extend the system.

## What Changes

- Create new `tests/daemon/` directory with integration tests
- Add test server helper with in-memory SQLite setup
- Implement tests for all 15 API endpoints from `daemon-api` spec
- Add `/openapi.yaml` endpoint to serve the spec file
- Update AGENTS.md with new test running instructions
- Move existing health check test to new test location

## Capabilities

### New Capabilities
- `api-testing`: Integration tests that verify all REST API endpoints match their OpenAPI specification

### Modified Capabilities
- (none - testing existing functionality, no requirement changes)

## Impact

- **New Files**: `tests/daemon/` directory with 5 test files
- **Modified Files**: `cmd/daemon/main.go` (add openapi route), `AGENTS.md` (update test instructions)
- **Dependencies**: No new dependencies - uses existing httptest and sqlite3
- **Affected Teams**: Developers working on the daemon API
