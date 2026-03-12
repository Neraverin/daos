# role-management Specification

## Purpose
Manages Docker Compose roles for deployment orchestration.

## MODIFIED Requirements

### Requirement: User can upload a role
The system SHALL store role content from a folder path on the daemon's host filesystem with a name.

#### Scenario: Upload valid role from folder
- **WHEN** user provides valid absolute folder path and name
- **AND** folder exists and is readable
- **THEN** role is created from files in the folder and persisted to database
- **AND** returns the created role with unique ID and timestamp

#### Scenario: Upload role without folder path
- **WHEN** user provides no folder path
- **THEN** returns validation error

#### Scenario: Upload role without name
- **WHEN** user provides folder path without name
- **THEN** returns validation error

#### Scenario: Upload role with invalid folder path
- **WHEN** user provides invalid or non-existent folder path
- **THEN** returns error as specified in compose-file-path spec
