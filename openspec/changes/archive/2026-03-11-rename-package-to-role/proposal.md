## Why

The term "package" is overloaded in software contexts and can cause confusion. In this system, a "package" refers to a Docker Compose file (role/configuration), not an installable software package. Renaming to "role" better reflects the semantic meaning in deployment orchestration contexts.

## What Changes

- Rename database table `packages` → `roles`
- Rename database column `package_id` → `role_id` in deployments table
- Rename API endpoint `/packages` → `/roles`
- Rename API schemas: `Package` → `Role`, `PackageSummary` → `RoleSummary`, `PackageInput` → `RoleInput`
- Rename API fields in Deployment: `package_id` → `role_id`, `package_name` → `role_name`
- Rename handler file `packages.go` → `roles.go`
- Rename TUI model file `packages.go` → `roles.go`
- Update all tests and inline schemas
- Regenerate API and SQL code

**BREAKING**: API endpoint `/packages` is renamed to `/roles`. Clients must update their integration.

## Capabilities

### New Capabilities
- None - this is a refactoring, not a new capability

### Modified Capabilities
- `package-management`: Renamed to `role-management` with updated terminology. Requirements remain functionally equivalent but terminology changes from "package" to "role".

## Impact

- **API**: All `/packages` endpoints now `/roles`. Existing clients must update.
- **Database**: Table and column renamed (no migration needed - clean schema rewrite)
- **Code**: Handlers, models, tests updated
- **Generated**: API types and SQL queries regenerated

### Affected Teams
- Frontend/TUI team: Update any references to packages
- API consumers: Update integrations to use `/roles` endpoint

### Rollback Plan
If rollback needed:
1. Rename database table `roles` → `packages` and column `role_id` → `package_id`
2. Rename API endpoints back to `/packages`
3. Regenerate code
4. Revert all source file changes
