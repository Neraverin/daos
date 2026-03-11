## 1. Database

- [x] 1.1 Rewrite migration file: rename table `packages` → `roles`, column `package_id` → `role_id`
- [x] 1.2 Rename SQL query file: `packages.sql` → `roles.sql`
- [x] 1.3 Update `deployments.sql` references to use `role_id` and `roles` table
- [x] 1.4 Run `make generate-sql` to regenerate db layer

## 2. API Specification

- [x] 2.1 Update OpenAPI spec: rename paths `/packages` → `/roles`
- [x] 2.2 Update OpenAPI spec: rename schemas `Package` → `Role`, `PackageSummary` → `RoleSummary`, `PackageInput` → `RoleInput`
- [x] 2.3 Update OpenAPI spec: rename Deployment fields `package_id` → `role_id`, `package_name` → `role_name`
- [x] 2.4 Run `make generate-api` to regenerate API types

## 3. Handlers

- [x] 3.1 Rename handler file: `packages.go` → `roles.go`
- [x] 3.2 Update handler functions: `ListPackages` → `ListRoles`, `CreatePackage` → `CreateRole`, etc.
- [x] 3.3 Update type references: `db.Package` → `db.Role`, `api.PackageSummary` → `api.RoleSummary`
- [x] 3.4 Update `deployments.go` handler to use `RoleID` instead of `PackageID`

## 4. TUI Models

- [x] 4.1 Rename model file: `packages.go` → `roles.go`
- [x] 4.2 Update types: `PackageSummary` → `RoleSummary`, `PackagesList` → `RolesList`
- [x] 4.3 Update API calls: `/packages` → `/roles`
- [x] 4.4 Update `deployments.go` model: `PackageID` → `RoleID`, `PackageName` → `RoleName`
- [x] 4.5 Update menu: "Packages" → "Roles"

## 5. Tests

- [x] 5.1 Rename test file: `packages_api_test.go` → `roles_api_test.go`
- [x] 5.2 Update test references: `/packages` → `/roles`, `package` → `role`
- [x] 5.3 Update inline schema in `server_test.go`
- [x] 5.4 Update `deployments_api_test.go` to use `createTestRole` helper

## 6. Verification

- [x] 6.1 Run `make test` to verify all tests pass
- [x] 6.2 Run `make build` to verify compilation
