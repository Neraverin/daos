## Context

The DAOS daemon exposes a REST API via Gin with endpoints for hosts, packages, deployments, and logs. Currently, only one health check test exists in `cmd/daemon/handlers/handlers_test.go`. The API is documented in `docs/openapi.yaml` and the specification is in `openspec/specs/daemon-api/spec.md`.

## Goals / Non-Goals

**Goals:**
- Create integration tests covering all 15 API endpoints defined in `daemon-api` spec
- Verify response status codes, JSON structures, and error handling
- Add `/openapi.yaml` endpoint to serve the spec file
- Update AGENTS.md with proper test running instructions

**Non-Goals:**
- Unit tests for internal handler logic (only API-level integration tests)
- Performance/load testing
- End-to-end tests with real SSH connections
- Test TUI (deferred to future change)

## Decisions

### 1. Test Organization: Hybrid Approach
**Decision**: Use lightweight httptest for simple endpoints, in-memory SQLite for CRUD endpoints.

**Rationale**: 
- Health check and OpenAPI endpoints don't need DB - keeps tests fast
- CRUD operations need real DB to test foreign keys, constraints, and data transformations
- Mirrors existing test pattern in `handlers_test.go`

### 2. Test File Structure
**Decision**: Separate test files per resource in `tests/daemon/` subdirectory.

**Rationale**:
- Clear organization by API domain
- Easy to run subsets: `go test ./tests/daemon/hosts...`
- Scalable as more tests are added

### 3. Database Migration in Tests
**Decision**: Execute `001_initial_schema.sql` directly in test setup.

**Rationale**:
- Simpler than adding golang-migrate dependency for tests
- Matches existing `db_test.go` pattern
- No migration version management needed for ephemeral test DB

### 4. Test Server Setup
**Decision**: Create reusable `setupTestServer(t)` helper that returns router + server + db.

**Rationale**:
- DRY - avoids repeating DB setup in each test file
- Ensures consistent test environment
- Easy to add cleanup in defer

### 5. Test Naming Convention
**Decision**: Use `*_api_test.go` suffix (e.g., `hosts_api_test.go`).

**Rationale**:
- Distinguishes from potential unit tests in same package
- Clear that these test HTTP API behavior

## Risks / Trade-offs

- **[Risk]** In-memory SQLite may behave differently from file-based SQLite
  - **Mitigation**: Use same driver, enable foreign keys, exec same schema
- **[Risk]** Tests may be flaky if handlers have hidden dependencies
  - **Mitigation**: Use fresh DB per test, verify test isolation
- **[Risk]** Test coverage may miss edge cases
  - **Mitigation**: Cover happy path + key error scenarios (404, 400)

## Migration Plan

1. Create `tests/daemon/` directory
2. Add `server_test.go` with test helper functions
3. Add test files: `health_api_test.go`, `hosts_api_test.go`, `packages_api_test.go`, `deployments_api_test.go`
4. Add openapi route in `cmd/daemon/main.go`
5. Update `AGENTS.md` with test instructions
6. Run `make test` to verify all tests pass
7. Remove old test file `cmd/daemon/handlers/handlers_test.go`

## Open Questions

- test concurrent requests (faster, more reliable)
- no rate limiting (deferred to future change)
- no API coverage (deferred to future change)
