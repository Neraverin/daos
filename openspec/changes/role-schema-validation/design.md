## Context

DAOS currently lacks a standardized way to validate role YAML configurations. Roles are stored in filesystem directories with a `role.yaml` file, but there is no validation to ensure required fields are present. This leads to runtime errors when deploying roles with invalid configurations.

**Current State:**
- Roles are stored in filesystem directories specified by `role_path`
- Each role directory must contain a `role.yaml` file
- No validation is performed on role YAML structure before deployment
- Users need a JSON Schema file for IDE integration and external validation

**Constraints:**
- Must maintain backward compatibility with existing APIs
- Validation should be reusable both in daemon API and external tools
- Error messages should be clear and actionable

## Goals / Non-Goals

**Goals:**
- Create a JSON Schema file for role validation that can be used in IDEs and external tools
- Implement a Go validation library that validates role YAML against the schema
- Add a validation endpoint to the daemon API (`POST /api/v1/validate/role`)
- Provide clear, actionable error messages for validation failures

**Non-Goals:**
- Modifying existing role CRUD endpoints (create, update, delete)
- Adding validation to deployment execution (future enhancement)
- Supporting other configuration formats (JSON, TOML) - YAML only for now

## Decisions

### 1. JSON Schema library choice
**Decision:** Use `github.com/santhosh-tekuri/jsonschema/v6` as the JSON Schema validator
**Rationale:** Pure Go implementation, no external C dependencies, actively maintained, supports JSON Schema draft-07

**Alternatives considered:**
- `github.com/xeipuuv/gojsonschema`: Requires C bindings, more complex to build
- `github.com/everit-org/json schema`: Good but less maintained

### 2. Package location
**Decision:** Create new package `pkg/validator/` for validation logic
**Rationale:** Follows existing package structure (`pkg/ansible/`, `pkg/db/`, `pkg/config/`), keeps validation concerns separate

**Alternatives considered:**
- Add to existing `pkg/api/`: Too focused on HTTP handling
- Add to `pkg/models/`: Models should be data structures only

### 3. Schema storage
**Decision:** Embed JSON Schema in Go code as a constant string, serve via API endpoint
**Rationale:** Single binary deployment, no separate schema file needed, easy to version

**Alternatives considered:**
- Store schema file in filesystem: Requires file deployment, more complex setup
- Load schema at runtime: Adds startup complexity

### 4. API endpoint design
**Decision:** Add new endpoint `POST /api/v1/validate/role` that accepts role YAML as request body
**Rationale:** Simple, focused API; follows existing API patterns; returns validation result with errors

**Alternatives considered:**
- Add validation to existing `POST /roles` endpoint: Would change existing API behavior
- Query parameter for validation: Not suitable for YAML body

## Risks / Trade-offs

- **[Risk]** Schema changes may break existing valid roles
  - **Mitigation:** Version the schema, maintain backward compatibility, document breaking changes

- **[Risk]** Large role YAML files may cause performance issues
  - **Mitigation:** Set reasonable size limits (e.g., 1MB max), validate asynchronously for large files

- **[Risk]** JSON Schema validation errors may be confusing to end users
  - **Mitigation:** Map schema errors to user-friendly messages, provide field paths in errors

- **[Trade-off]** Embedding schema increases binary size by ~10KB
  - **Acceptable:** Worth the deployment simplicity

## Migration Plan

1. **Phase 1:** Create JSON Schema file and Go validation library
   - Create `pkg/validator/role.go` with validation logic
   - Create `pkg/validator/role.schema.json` (embedded)
   - Write unit tests

2. **Phase 2:** Add API endpoint
   - Add `POST /api/v1/validate/role` endpoint
   - Update OpenAPI spec
   - Generate API code

3. **Phase 3:** Deploy and monitor
   - Deploy to staging, test with existing roles
   - Monitor validation error rates
   - Gather user feedback

**Rollback:** Remove the validation endpoint and validator package. No database migrations or breaking API changes.

## Open Questions

- Should we validate roles automatically when they're created via the API, or only on explicit validation request?
  - **Answer:** Automatic validation when created via API

- Should we cache validation results for frequently validated roles?
  - **Answer:** No need for cache. Role creation and validation are rare operations

- Do we need to support schema version negotiation for forward compatibility?
  - **Answer:** Yes. The Version field in the role schema provides versioning capability
