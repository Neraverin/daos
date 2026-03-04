## Why

The project currently uses `int64` for all ID fields (Host, Package, Deployment, Log). This needs to change to UUID because:
- UUID is the company standard for all ID fields
- Any future ID fields should also use UUID

## What Changes

- Change all ID fields in domain models from `int64` to `uuid.UUID`
- Update OpenAPI specification: `type: integer` → `type: string, format: uuid`
- Regenerate API types (generated from OpenAPI)
- Create new database migration with TEXT PRIMARY KEY for all tables
- Update SQL queries to accept ID as parameter (remove auto-increment)
- Update HTTP handlers to generate UUID on create and parse UUID from path parameters
- Update TUI models to use `uuid.UUID` for ID fields

## Capabilities

### New Capabilities
- `uuid-ids`: Replace all int64 ID fields with UUID type across the entire codebase

### Modified Capabilities
- (none - this is an implementation-only change with no behavioral changes)

## Impact

- **API**: All endpoint path parameters and response body IDs change from integer to string (no users yet, not a breaking change)
- **Database**: New migration required (can drop existing data, no users yet)
- **Generated code**: pkg/api/openapi.gen.go, pkg/db/*.sql.go will be regenerated
- **Models**: pkg/models/models.go, cmd/tui/models/*.go need ID type changes
- **Handlers**: cmd/daemon/handlers/*.go need UUID generation and parsing
