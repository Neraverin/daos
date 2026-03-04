## 1. OpenAPI Specification

- [x] 1.1 Update all ID field types: `type: integer` → `type: string, format: uuid`
- [x] 1.2 Update all path parameter types: `type: integer` → `type: string`

## 2. API Generation

- [x] 2.1 Run `make generate-api` to regenerate pkg/api/openapi.gen.go

## 3. Database Schema

- [x] 3.1 Create new migration file: pkg/db/migrations/002_uuid_schema.sql
- [x] 3.2 Change all `id INTEGER PRIMARY KEY AUTOINCREMENT` to `id TEXT PRIMARY KEY`
- [x] 3.3 Change all foreign key references to `TEXT` type

## 4. SQL Queries

- [x] 4.1 Update hosts.sql: Add id parameter to CreateHost, remove RETURNING id
- [x] 4.2 Update packages.sql: Add id parameter to CreatePackage, remove RETURNING id
- [x] 4.3 Update deployments.sql: Add id parameter to CreateDeployment, remove RETURNING id
- [x] 4.4 Update logs.sql: Add id parameter to CreateLog, remove RETURNING id

## 5. Database Code Generation

- [x] 5.1 Run `make generate-sql` to regenerate db layer

## 6. Domain Models

- [x] 6.1 Update Host.ID: int64 → uuid.UUID
- [x] 6.2 Update Package.ID: int64 → uuid.UUID
- [x] 6.3 Update Deployment.ID, HostID, PackageID: int64 → uuid.UUID
- [x] 6.4 Update Log.ID, DeploymentID: int64 → uuid.UUID

## 7. HTTP Handlers

- [x] 7.1 Update hosts.go: Generate UUID on create, parse UUID from path
- [x] 7.2 Update packages.go: Generate UUID on create, parse UUID from path
- [x] 7.3 Update deployments.go: Generate UUID on create, parse UUID from path
- [x] 7.4 Add uuid import to all handler files

## 8. TUI Models

- [x] 8.1 Update hosts.go: Change ID field to uuid.UUID
- [x] 8.2 Update packages.go: Change ID field to uuid.UUID
- [x] 8.3 Update deployments.go: Change ID, HostID, PackageID to uuid.UUID
- [x] 8.4 Update message types that carry IDs (deleteHostMsg, etc.)

## 9. Verification

- [x] 9.1 Run `make test` to verify all tests pass
- [x] 9.2 Run `make build` to verify compilation
