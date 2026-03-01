# deployment-management Specification

## Purpose
TBD - created by archiving change daos-v1. Update Purpose after archive.
## Requirements
### Requirement: User can create a deployment
The system SHALL allow creating a deployment linking a host and package.

#### Scenario: Create deployment with valid host and package
- **WHEN** user provides valid host_id and package_id
- **THEN** deployment is created with status "pending"
- **AND** returns the created deployment with unique ID and timestamps

#### Scenario: Create deployment with non-existent host
- **WHEN** user provides invalid host_id
- **THEN** returns 400 error indicating host not found

#### Scenario: Create deployment with non-existent package
- **WHEN** user provides invalid package_id
- **THEN** returns 400 error indicating package not found

### Requirement: User can list all deployments
The system SHALL return a list of all deployments.

#### Scenario: List deployments when empty
- **WHEN** no deployments exist
- **THEN** returns empty array

#### Scenario: List deployments with entries
- **WHEN** multiple deployments exist
- **THEN** returns array of all deployments with host and package info

### Requirement: User can get a specific deployment
The system SHALL return a single deployment by ID.

#### Scenario: Get existing deployment
- **WHEN** user requests deployment by valid ID
- **THEN** returns deployment details with host and package info

#### Scenario: Get non-existent deployment
- **WHEN** user requests deployment by invalid ID
- **THEN** returns 404 error

### Requirement: User can delete a deployment
The system SHALL allow removing a deployment.

#### Scenario: Delete existing deployment
- **WHEN** user requests deletion of existing deployment
- **THEN** deployment is removed from database
- **AND** associated logs are also removed

#### Scenario: Delete non-existent deployment
- **WHEN** user attempts to delete deployment with invalid ID
- **THEN** returns 404 error

### Requirement: User can run a deployment
The system SHALL trigger Ansible execution for a deployment.

#### Scenario: Run pending deployment
- **WHEN** user triggers deployment run for pending deployment
- **AND** host has valid SSH configuration
- **THEN** deployment status changes to "running"
- **AND** Ansible executes compose deployment to remote host
- **AND** on success, status changes to "success"
- **AND** on failure, status changes to "failed"

#### Scenario: Run already running deployment
- **WHEN** user triggers deployment run for already running deployment
- **THEN** returns 400 error indicating deployment already in progress

#### Scenario: Run deployment with invalid SSH key path
- **WHEN** user triggers deployment run but SSH key file does not exist
- **THEN** deployment status changes to "failed"
- **AND** error is logged

### Requirement: User can view deployment status
The system SHALL provide current status of deployments.

#### Scenario: Check deployment status
- **WHEN** user requests deployment details
- **THEN** status field shows one of: pending, running, success, failed

