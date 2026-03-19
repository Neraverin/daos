# api-testing Specification

## Purpose

This capability adds integration tests that verify all REST API endpoints match their OpenAPI specification defined in the `daemon-api` spec.

## MODIFIED Requirements

### Requirement: API integration tests exist
The system SHALL provide integration tests for all REST API endpoints defined in the daemon-api specification. **Changed**: Tests can now run against remote test stand via TestStand connection manager.

#### Scenario: Health check test exists
- **WHEN** tests are run for GET /health endpoint
- **THEN** test verifies 200 OK response with {"status": "healthy"}
- **AND** test runs against configured test stand (local or remote)

#### Scenario: OpenAPI spec endpoint test exists
- **WHEN** tests are run for GET /openapi.yaml endpoint
- **THEN** test verifies endpoint returns OpenAPI 3 specification

#### Scenario: Hosts API tests exist
- **WHEN** tests are run for hosts endpoints
- **THEN** tests verify List, Create, Get, Update, Delete operations

#### Scenario: Packages API tests exist
- **WHEN** tests are run for packages endpoints
- **THEN** tests verify List, Create, Get, Delete operations

#### Scenario: Deployments API tests exist
- **WHEN** tests are run for deployments endpoints
- **THEN** tests verify List, Create, Get, Delete, Run, and Logs operations
