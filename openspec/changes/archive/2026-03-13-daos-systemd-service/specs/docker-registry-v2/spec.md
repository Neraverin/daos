## ADDED Requirements

### Requirement: Docker Registry v2 is available as local container image registry
The system SHALL provide a local Docker Registry v2 instance for storing and distributing role container images. The registry MUST be managed by systemd and automatically start when DAOS service starts.

#### Scenario: Registry container starts with DAOS service
- **WHEN** DAOS systemd service starts
- **THEN** Docker Registry v2 container SHALL start automatically
- **AND** Container SHALL be named `daos-registry`
- **AND** Container SHALL use `--restart=unless-stopped` policy

#### Scenario: Registry persists image data across restarts
- **WHEN** Registry container restarts
- **THEN** Previously pushed images SHALL remain available
- **AND** Registry data SHALL be stored at `/opt/daos/registry/` on the host

#### Scenario: Registry is accessible on configured port
- **WHEN** Client pushes or pulls image
- **THEN** Registry SHALL accept connections on port defined in `/opt/daos/registry/daos-registry.env`
- **AND** Default port SHALL be 5000

### Requirement: Docker and containerd are installed as package dependencies
DAOS packages MUST include Docker and containerd as dependencies to ensure the registry can run.

#### Scenario: Installing DAOS on Debian/Ubuntu
- **WHEN** Administrator installs DAOS .deb package
- **THEN** Package manager SHALL install docker.io and containerd packages
- **AND** Docker service SHALL be available

#### Scenario: Installing DAOS on RHEL/CentOS/Fedora
- **WHEN** Administrator installs DAOS .rpm package
- **AND** Package manager SHALL install docker-ce and containerd packages
- **AND** Docker service SHALL be available

### Requirement: Registry port is configurable
The registry port MUST be configurable via environment file without modifying the systemd service file.

#### Scenario: Changing registry port
- **WHEN** Administrator edits `/opt/daos/registry/daos-registry.env` and sets `DAOS_REGISTRY_PORT=5100`
- **AND** restarts DAOS service
- **THEN** Registry SHALL accept connections on port 5100

### Requirement: Registry image is pre-pulled during package installation
The registry:2 image MUST be pulled during package installation to ensure availability on first service start.

#### Scenario: Package installation
- **WHEN** DAOS package is installed
- **AND** Docker is available
- **THEN** Package post-install script SHALL pull `registry:2` image
- **AND** Image SHALL be available locally for container creation
