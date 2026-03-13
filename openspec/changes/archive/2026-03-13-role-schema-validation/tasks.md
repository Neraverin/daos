## 1. Create JSON Schema

- [x] 1.1 Create JSON Schema file for role validation (`pkg/validator/role.schema.json`)
- [x] 1.2 Define required fields: Version, Role.TypeId, Role.Name, Role.Version
- [x] 1.3 Define optional fields: Role.Description, Role.Definitions
- [x] 1.4 Define Definitions structure: ImageFile (string), Images (array), TemplateFiles (array)

## 2. Create Validator Package

- [x] 2.1 Create `pkg/validator/` package directory
- [x] 2.2 Add jsonschema/v6 dependency to go.mod
- [x] 2.3 Create validation function to validate role YAML against schema
- [x] 2.4 Create function to get schema as JSON string
- [x] 2.5 Implement error message formatting with field paths
- [x] 2.6 Write unit tests for validation functions

## 3. Add API Endpoint

- [x] 3.1 Update OpenAPI spec (`docs/openapi.yaml`) with `/validate/role` endpoint
- [x] 3.2 Run `make generate-api` to regenerate API types
- [x] 3.3 Create validation handler in `cmd/daemon/handlers/`
- [x] 3.4 Register endpoint in server setup
- [x] 3.5 Test endpoint with valid and invalid role YAML

## 4. Integration and Testing

- [x] 4.1 Test with sample valid role YAML
- [x] 4.2 Test with role missing required fields
- [x] 4.3 Test with role having invalid Definitions structure
- [x] 4.4 Verify error messages include field paths
- [x] 4.5 Run full test suite (`make test`)
