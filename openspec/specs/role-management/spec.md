# role-management Specification

## Purpose
Manages Docker Compose roles for deployment orchestration.

## Requirements

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

### Requirement: User can list all roles
The system SHALL return a list of all stored roles.

#### Scenario: List roles when empty
- **WHEN** no roles exist
- **THEN** returns empty array

#### Scenario: List roles with entries
- **WHEN** multiple roles exist
- **THEN** returns array of all roles

### Requirement: User can get a specific role
The system SHALL return a single role by ID.

#### Scenario: Get existing role
- **WHEN** user requests role by valid ID
- **THEN** returns role details

#### Scenario: Get non-existent role
- **WHEN** user requests role by invalid ID
- **THEN** returns 404 error

### Requirement: User can delete a role
The system SHALL allow removing a role.

#### Scenario: Delete existing role
- **WHEN** user requests deletion of existing role
- **THEN** role is removed from database

#### Scenario: Delete non-existent role
- **WHEN** user attempts to delete role with invalid ID
- **THEN** returns 404 error
