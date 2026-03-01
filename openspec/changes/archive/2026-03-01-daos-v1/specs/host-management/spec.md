## ADDED Requirements

### Requirement: User can add a new host
The system SHALL allow users to add remote hosts with SSH connection details.

#### Scenario: Add host with all fields
- **WHEN** user provides name, hostname, port, username, and SSH key path
- **THEN** host is created and persisted to database
- **AND** returns the created host with unique ID

#### Scenario: Add host with default port
- **WHEN** user provides name, hostname, username, and SSH key path without port
- **THEN** port defaults to 22
- **AND** host is created successfully

### Requirement: User can list all hosts
The system SHALL return a list of all registered hosts.

#### Scenario: List hosts when empty
- **WHEN** no hosts exist
- **THEN** returns empty array

#### Scenario: List hosts with entries
- **WHEN** multiple hosts exist
- **THEN** returns array of all hosts with their details (excluding sensitive data if any)

### Requirement: User can get a specific host
The system SHALL return a single host by ID.

#### Scenario: Get existing host
- **WHEN** user requests host by valid ID
- **THEN** returns host details

#### Scenario: Get non-existent host
- **WHEN** user requests host by invalid ID
- **THEN** returns 404 error

### Requirement: User can update a host
The system SHALL allow updating host connection details.

#### Scenario: Update host fields
- **WHEN** user provides updated name, hostname, port, username, or SSH key path
- **THEN** host is updated in database
- **AND** returns the updated host

#### Scenario: Update non-existent host
- **WHEN** user attempts to update host with invalid ID
- **THEN** returns 404 error

### Requirement: User can delete a host
The system SHALL allow removing a host.

#### Scenario: Delete existing host
- **WHEN** user requests deletion of existing host
- **THEN** host is removed from database

#### Scenario: Delete non-existent host
- **WHEN** user attempts to delete host with invalid ID
- **THEN** returns 404 error
