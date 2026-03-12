# compose-file-path Specification

## Purpose
Allows specifying a folder path on the daemon's host filesystem instead of inline compose content.

## ADDED Requirements

### Requirement: User can specify role folder path
The system SHALL accept a folder path on the daemon's host filesystem instead of inline content.

#### Scenario: Create role with valid folder path
- **WHEN** user provides valid absolute folder path and role name
- **THEN** daemon reads role files from the specified folder
- **AND** role is created and persisted to database
- **AND** returns the created role with unique ID and timestamp

#### Scenario: Create role with non-existent folder path
- **WHEN** user provides folder path that does not exist
- **THEN** returns 404 error with message "role folder not found at path"

#### Scenario: Create role with non-absolute folder path
- **WHEN** user provides a non-absolute (relative) folder path
- **THEN** returns validation error with message "folder path must be absolute"

#### Scenario: Create role with folder path that is not readable
- **WHEN** user provides folder path that exists but is not readable
- **THEN** returns 403 error with message "permission denied"

#### Scenario: Update role with valid folder path
- **WHEN** user provides valid absolute folder path for existing role
- **AND** folder contains valid role files
- **THEN** daemon reads updated role files from the folder
- **AND** role is updated in database
- **AND** returns the updated role
