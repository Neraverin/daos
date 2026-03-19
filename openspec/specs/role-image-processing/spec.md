# role-image-processing Specification

## Purpose

Process Docker image archives bundled with roles by validating paths, loading via `docker load`, and pushing to the configured Docker registry.

## Requirements

### Requirement: System validates image file path
The system SHALL validate that the `Role.Definitions.ImageFile` path exists and has a `.tar` extension when specified.

#### Scenario: Valid image file path
- **WHEN** role contains `Definitions.ImageFile` with valid path
- **AND** the file exists and has `.tar` extension
- **THEN** validation passes

#### Scenario: Missing image file
- **WHEN** role contains `Definitions.ImageFile` pointing to non-existent file
- **THEN** returns error indicating file not found

#### Scenario: Invalid file extension
- **WHEN** role contains `Definitions.ImageFile` with non-`.tar` extension
- **THEN** returns error indicating only `.tar` files are allowed

### Requirement: System loads image from archive
The system SHALL load Docker images from `.tar` archives using the `docker load` command.

#### Scenario: Successful image load
- **WHEN** image archive exists and is valid
- **THEN** `docker load -i <path>` is executed successfully
- **AND** image is available in local Docker daemon

#### Scenario: Image load timeout
- **WHEN** `docker load` exceeds configured timeout
- **THEN** operation fails with timeout error
- **AND** partial state is cleaned up

#### Scenario: Invalid tar archive
- **WHEN** image archive exists but is not a valid Docker image tar
- **THEN** `docker load` fails
- **AND** error is returned indicating invalid archive format

### Requirement: System pushes image to registry
The system SHALL tag and push loaded images to the configured Docker registry.

#### Scenario: Successful image push
- **WHEN** image is loaded successfully
- **AND** registry is configured
- **THEN** image is tagged with `<registry>/<image>:<tag>`
- **AND** pushed to registry

#### Scenario: Registry not configured
- **WHEN** image processing is triggered
- **AND** Docker registry is not configured
- **THEN** image processing is skipped
- **AND** role is created without image processing

#### Scenario: Image push failure
- **WHEN** image push to registry fails
- **THEN** error is returned with failure details
- **AND** role creation fails with appropriate message

### Requirement: System processes images during role creation
The system SHALL process images when a role is created or updated via API.

#### Scenario: Image processing on role creation
- **WHEN** user creates a role with `Definitions.ImageFile`
- **THEN** image is validated, loaded, and pushed before role creation completes
- **AND** role is only created if image processing succeeds

#### Scenario: Role without image file
- **WHEN** user creates a role without `Definitions.ImageFile`
- **THEN** role is created without image processing
