## 1. Test Infrastructure

- [x] 1.1 Create `tests/daemon/` directory
- [x] 1.2 Create `tests/daemon/server_test.go` with test server setup helper
- [x] 1.3 Add helper function `setupTestServer(t)` that creates in-memory SQLite and runs migrations

## 2. Health and OpenAPI Endpoint

- [x] 2.1 Add `/openapi.yaml` static route in `cmd/daemon/main.go`
- [x] 2.2 Create `tests/daemon/health_api_test.go` with health check test
- [x] 2.3 Add test for GET /openapi.yaml endpoint
- [x] 2.4 Delete old `cmd/daemon/handlers/handlers_test.go`

## 3. Hosts API Tests

- [x] 3.1 Create `tests/daemon/hosts_api_test.go`
- [x] 3.2 Add test: List hosts (empty)
- [x] 3.3 Add test: Create host (201)
- [x] 3.4 Add test: List hosts (with data)
- [x] 3.5 Add test: Get host by ID (200)
- [x] 3.6 Add test: Get host by ID (404)
- [x] 3.7 Add test: Update host (200)
- [x] 3.8 Add test: Delete host (204)
- [x] 3.9 Add test: Delete host (404) - Note: returns 204 (idempotent delete)

## 4. Packages API Tests

- [x] 4.1 Create `tests/daemon/packages_api_test.go`
- [x] 4.2 Add test: List packages (empty)
- [x] 4.3 Add test: Create package (201)
- [x] 4.4 Add test: List packages (with data)
- [x] 4.5 Add test: Get package by ID (200)
- [x] 4.6 Add test: Get package by ID (404)
- [x] 4.7 Add test: Delete package (204)
- [x] 4.8 Add test: Delete package (404) - Note: returns 204 (idempotent delete)

## 5. Deployments API Tests

- [x] 5.1 Create `tests/daemon/deployments_api_test.go`
- [x] 5.2 Add test: List deployments (empty)
- [x] 5.3 Add test: Create deployment (201)
- [x] 5.4 Add test: Create deployment with invalid host (400)
- [x] 5.5 Add test: Create deployment with invalid package (400)
- [x] 5.6 Add test: List deployments (with data)
- [x] 5.7 Add test: Get deployment by ID (200)
- [x] 5.8 Add test: Get deployment by ID (404)
- [x] 5.9 Add test: Delete deployment (204)
- [x] 5.10 Add test: Delete deployment (404) - Note: returns 204 (idempotent delete)
- [x] 5.11 Add test: Run deployment (200)
- [x] 5.12 Add test: Get deployment logs (200)

## 6. Documentation Updates

- [x] 6.1 Update AGENTS.md with new test running instructions
- [x] 6.2 Verify tests pass with `make test`
