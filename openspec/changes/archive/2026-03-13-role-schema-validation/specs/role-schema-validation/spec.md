## ADDED Requirements

### Requirement: Role YAML file structure validation
The system SHALL validate role YAML files to ensure they conform to the required structure with all mandatory fields present.

#### Scenario: Valid role with all required fields
- **WHEN** a role YAML file contains Version, Role.TypeId, Role.Name, and Role.Version
- **THEN** the validation SHALL pass

#### Scenario: Role missing Version field
- **WHEN** a role YAML file is missing the Version field
- **THEN** the validation SHALL fail with an error indicating Version is required

#### Scenario: Role missing Role.TypeId field
- **WHEN** a role YAML file is missing Role.TypeId
- **THEN** the validation SHALL fail with an error indicating Role.TypeId is required

#### Scenario: Role missing Role.Name field
- **WHEN** a role YAML file is missing Role.Name
- **THEN** the validation SHALL fail with an error indicating Role.Name is required

#### Scenario: Role missing Role.Version field
- **WHEN** a role YAML file is missing Role.Version
- **THEN** the validation SHALL fail with an error indicating Role.Version is required

### Requirement: Optional role fields validation
The system SHALL validate optional role fields (Description, Definitions) when they are present in the role YAML file.

#### Scenario: Role with optional Description field
- **WHEN** a role YAML file contains Role.Description
- **THEN** the validation SHALL pass and accept any string value

#### Scenario: Role with optional Definitions field
- **WHEN** a role YAML file contains Role.Definitions
- **THEN** the validation SHALL pass when Definitions contains valid structure (ImageFile, Images, TemplateFiles)

#### Scenario: Role with empty Definitions
- **WHEN** a role YAML file contains empty Role.Definitions
- **THEN** the validation SHALL pass

### Requirement: Role Definitions structure validation
The system SHALL validate the Definitions object when present, ensuring all nested fields conform to expected types.

#### Scenario: Valid Definitions with ImageFile
- **WHEN** Role.Definitions.ImageFile is a valid string path
- **THEN** the validation SHALL pass

#### Scenario: Valid Definitions with Images array
- **WHEN** Role.Definitions.Images is an array of objects with Id field
- **THEN** the validation SHALL pass

#### Scenario: Valid Definitions with TemplateFiles array
- **WHEN** Role.Definitions.TemplateFiles is an array of strings
- **THEN** the validation SHALL pass

#### Scenario: Invalid ImageFile type
- **WHEN** Role.Definitions.ImageFile is not a string
- **THEN** the validation SHALL fail with an error indicating ImageFile must be a string

#### Scenario: Invalid Images array format
- **WHEN** Role.Definitions.Images contains objects without Id field
- **THEN** the validation SHALL fail with an error indicating each image must have an Id

### Requirement: JSON Schema file availability
The system SHALL provide a JSON Schema file that can be used for external validation in tools like IDEs.

#### Scenario: JSON Schema file exists
- **WHEN** the role.schema.json file is requested
- **THEN** the system SHALL return a valid JSON Schema that validates role YAML files

#### Scenario: JSON Schema validates required fields
- **WHEN** a role YAML is validated against the JSON Schema
- **THEN** Version, Role.TypeId, Role.Name, and Role.Version SHALL be marked as required

### Requirement: API validation endpoint
The system SHALL provide an API endpoint to validate role YAML content.

#### Scenario: Valid role YAML submitted to validation endpoint
- **WHEN** a POST request with valid role YAML is sent to /api/v1/validate/role
- **THEN** the response SHALL indicate validation passed

#### Scenario: Invalid role YAML submitted to validation endpoint
- **WHEN** a POST request with invalid role YAML is sent to /api/v1/validate/role
- **THEN** the response SHALL indicate validation failed with specific error messages

### Requirement: Clear validation error messages
The system SHALL provide clear, actionable error messages when validation fails.

#### Scenario: Multiple validation errors
- **WHEN** a role YAML has multiple validation errors
- **THEN** all errors SHALL be returned in a list with field paths and error descriptions

#### Scenario: Error message includes field path
- **THEN** each error SHALL include the path to the field that failed validation (e.g., "Role.Version is required")
