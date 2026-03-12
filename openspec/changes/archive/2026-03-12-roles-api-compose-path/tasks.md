## 1. API Specification Updates

- [x] 1.1 Update OpenAPI spec (docs/openapi.yaml) to replace `compose_content` with `role_path`
- [x] 1.2 Regenerate API types with `make generate-api`

## 2. Daemon Handler Implementation

- [x] 2.1 Update role create handler to read from folder path
- [x] 2.2 Add validation for absolute path
- [x] 2.3 Add validation for folder existence (return 404 if not found)
- [x] 2.4 Add validation for folder readability (return 403 if not accessible)
- [x] 2.5 Update role update handler to read from folder path (N/A - no PUT endpoint exists)
- [x] 2.6 Update error responses to match spec messages

## 3. Database Schema Updates

- [x] 3.1 Update database schema to store `role_path` instead of `compose_content`
- [x] 3.2 Create migration for existing roles (move content to files or clear)
- [x] 3.3 Regenerate SQL code with `make generate-sql`

## 4. TUI Updates

- [x] 4.1 Update role create view to use folder picker instead of text area
- [x] 4.2 Update role edit view to use folder picker
- [x] 4.3 Add validation feedback for invalid folder paths
- [x] 4.4 Update role display to show folder path instead of content

## 5. Testing

- [x] 5.1 Add unit tests for path validation
- [x] 5.2 Add integration tests for API endpoints
- [x] 5.3 Update existing tests to use folder path

## 6. Documentation

- [x] 6.1 Update API documentation (updated OpenAPI spec)
