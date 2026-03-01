## ADDED Requirements

### Requirement: Daemon exposes REST API
The system SHALL provide HTTP API for all operations.

#### Scenario: API responds on configured port
- **WHEN** daemon is started
- **THEN** HTTP server listens on configured host:port

#### Scenario: Health check endpoint
- **WHEN** user requests GET /health
- **THEN** returns 200 OK with {"status": "healthy"}

### Requirement: API uses OpenAPI 3 contract
The system SHALL document all endpoints in OpenAPI 3 format.

#### Scenario: OpenAPI spec available
- **WHEN** user requests GET /openapi.yaml
- **THEN** returns OpenAPI 3 specification document

### Requirement: Hosts API endpoints
The system SHALL provide CRUD operations for hosts.

#### Scenario: GET /api/v1/hosts
- **WHEN** user requests all hosts
- **THEN** returns 200 OK with array of hosts

#### Scenario: POST /api/v1/hosts
- **WHEN** user creates new host with valid data
- **THEN** returns 201 Created with host object

#### Scenario: GET /api/v1/hosts/{id}
- **WHEN** user requests specific host
- **THEN** returns 200 OK with host object or 404 if not found

#### Scenario: PUT /api/v1/hosts/{id}
- **WHEN** user updates existing host
- **THEN** returns 200 OK with updated host object

#### Scenario: DELETE /api/v1/hosts/{id}
- **WHEN** user deletes existing host
- **THEN** returns 204 No Content

### Requirement: Packages API endpoints
The system SHALL provide CRUD operations for packages.

#### Scenario: GET /api/v1/packages
- **WHEN** user requests all packages
- **THEN** returns 200 OK with array of packages (summary)

#### Scenario: POST /api/v1/packages
- **WHEN** user uploads new package with valid compose
- **THEN** returns 201 Created with package object

#### Scenario: GET /api/v1/packages/{id}
- **WHEN** user requests specific package
- **THEN** returns 200 OK with package including compose content

#### Scenario: DELETE /api/v1/packages/{id}
- **WHEN** user deletes existing package
- **THEN** returns 204 No Content

### Requirement: Deployments API endpoints
The system SHALL provide CRUD and execution operations for deployments.

#### Scenario: GET /api/v1/deployments
- **WHEN** user requests all deployments
- **THEN** returns 200 OK with array of deployments

#### Scenario: POST /api/v1/deployments
- **WHEN** user creates new deployment with valid host_id and package_id
- **THEN** returns 201 Created with deployment object

#### Scenario: GET /api/v1/deployments/{id}
- **WHEN** user requests specific deployment
- **THEN** returns 200 OK with deployment including host and package info

#### Scenario: DELETE /api/v1/deployments/{id}
- **WHEN** user deletes existing deployment
- **THEN** returns 204 No Content

#### Scenario: POST /api/v1/deployments/{id}/run
- **WHEN** user triggers deployment execution
- **THEN** returns 200 OK with updated deployment (status: running)

### Requirement: Logs API endpoints
The system SHALL provide access to deployment logs.

#### Scenario: GET /api/v1/deployments/{id}/logs
- **WHEN** user requests logs for deployment
- **THEN** returns 200 OK with array of log entries
