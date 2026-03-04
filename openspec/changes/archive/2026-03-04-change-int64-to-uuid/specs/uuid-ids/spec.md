## ADDED Requirements

### Requirement: All entity IDs use UUID type
All entities in the system SHALL use UUID type for their ID field instead of int64.

#### Scenario: Creating a new host
- **WHEN** a client creates a new host via POST /hosts
- **THEN** the system generates a UUID for the host ID and returns it in the response

#### Scenario: Retrieving an entity by ID
- **WHEN** a client requests an entity by ID via GET /{entity}/{id}
- **THEN** the ID in the path MUST be a valid UUID string

#### Scenario: Response contains UUID IDs
- **WHEN** a client receives a response containing entity data
- **THEN** all ID fields SHALL be UUID strings in format "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

### Requirement: UUID generation is deterministic
The system SHALL generate UUIDs using a cryptographically secure random UUID generator.

#### Scenario: Creating multiple entities
- **WHEN** multiple entities are created in sequence
- **THEN** each entity MUST have a unique UUID (no collisions)
