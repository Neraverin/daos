## Why

Currently, there is no standardized way to define and validate role configurations in DAOS. Roles need a consistent schema to ensure all required fields (Version, TypeId, Name, Version) are present, while optional fields (Description, Definitions) can be added as needed. This will enable validation both in external tools (like IDE plugins) and in the daemon API.

## What Changes

- Add a JSON Schema file for role validation (`role.schema.json`)
- Create a role validation library in Go that can validate role YAML files
- Integrate role validation into the daemon API for hosts/deployments
- Provide clear error messages for validation failures

## Capabilities

### New Capabilities
- `role-schema-validation`: Validate role YAML files against JSON Schema, ensuring required fields (Version, Role.TypeId, Role.Name, Role.Version) are present and optional fields (Role.Description, Role.Definitions) are properly structured

### Modified Capabilities
- None

## Impact

- **New files**: `pkg/validator/role.go`, `pkg/validator/role.schema.json`
- **API changes**: Add validation endpoint or integrate into existing host/deployment endpoints
- **Dependencies**: Need to add JSON Schema validation library (e.g., gojsonschema)
- **Affected teams**: Backend team (daemon), TUI team (if displaying validation errors)

## Rollback Plan

- Rollback involves removing the validation library and reverting API changes
- No breaking changes to existing API contracts
- Role files that pass validation will continue to work as before

## Affected Teams

- Backend team (daemon API)
- TUI team (display validation errors)
- External tool developers (using the schema for IDE integration)
