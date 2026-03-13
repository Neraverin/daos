## Context

DAOS daemon currently lacks systemd integration and must be started manually. The system also requires Docker Registry v2 as a local container image registry for storing and distributing role images. Both Docker daemon and the registry container need to be running before DAOS starts.

Current state:
- DAOS daemon is started manually or via init script
- No Docker dependency management at package level
- Registry must be started separately

Constraints:
- Must work on Debian/Ubuntu (.deb) and RHEL/CentOS/Fedora (.rpm)
- Registry port must be configurable via environment file
- Docker must be running before DAOS starts
- Registry container must restart automatically on failure

## Goals / Non-Goals

**Goals:**
- Create systemd service file for DAOS daemon
- Add Docker and containerd as package dependencies
- Ensure Docker daemon is running before DAOS starts
- Ensure Docker Registry v2 container is running before DAOS starts
- Make registry port configurable via environment file

**Non-Goals:**
- Modify DAOS Go code for Docker/registry management (handled by systemd)
- Implement registry authentication/SSL (use default HTTP for local)
- Support non-systemd init systems

## Decisions

### 1. Systemd Service Structure

**Decision:** Single service file with ExecStartPre directives

**Rationale:** Simpler than two separate services. Uses systemd's native dependency ordering.

**Alternative Considered:** Separate docker-registry.service
- Rejected: Adds complexity with service ordering and failure handling

### 2. Registry Port Configuration

**Decision:** Use environment file `/opt/daos/registry/daos-registry.env` with `DAOS_REGISTRY_PORT`

**Rationale:** systemd supports environment file loading via `EnvironmentFile`. Allows runtime configuration without modifying service file.

**Alternative Considered:** Hardcoded port 5000
- Rejected: Less flexible, doesn't meet requirement for configurability

### 3. Docker Dependency Handling

**Decision:** Use systemd `Requires=docker.service` and `After=docker.service`

**Rationale:** Ensures Docker socket is available. Does not automatically start Docker daemon (relies on Docker's own service configuration).

**Alternative Considered:** ExecStartPre to start Docker
- Rejected: Docker should manage its own daemon lifecycle

### 4. Registry Container Management

**Decision:** ExecStartPre with docker rm (ignore failure) then docker run

**Rationale:** Ensures clean state - removes any existing container with same name before creating new one. Uses `--restart=unless-stopped` for automatic restart.

```ini
ExecStartPre=-/usr/bin/docker rm -f daos-registry 2>/dev/null || true
ExecStartPre=/usr/bin/docker run -d --name daos-registry \
    --restart=unless-stopped \
    -v /opt/daos/registry:/var/lib/registry \
    -p 5000:5000 \
    registry:2
```

**Alternative Considered:** docker-compose
- Rejected: Adds dependency on docker-compose, more complex for single container

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Docker not installed | Add docker.io/docker-ce as package dependency |
| Registry port conflict | Make configurable via environment file |
| Registry image not present | Pull during package install or first start |
| Docker socket permission denied | Ensure daos user has docker group membership |

**Trade-offs:**
- Environment file approach requires manual service restart to apply port changes
- Registry data not persisted across container recreation (acceptable for local dev use)

## Migration Plan

1. Update package metadata (DEBIAN/control, daos.spec) with Docker dependencies
2. Pre-install script pulls registry:2 image (`docker pull registry:2`)
3. Create /opt/daos/registry/ directory for persistent storage
4. Create packaging/systemd/daos.service
5. Create packaging/systemd/daos-registry.env with default port
6. Package includes postinst/%post script to:
   - Install service file to /etc/systemd/system/
   - Run systemctl daemon-reload
   - Optionally: systemctl enable daos

## Persistence

Registry data is persisted to host volume at `/opt/daos/registry/` mounted to `/var/lib/registry` inside the container.

Registry:2 image is automatically pulled during package installation to ensure the image is available on first start.
