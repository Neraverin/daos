## Why

DAOS daemon currently lacks systemd integration, requiring manual start/stop management. Additionally, DAOS requires Docker Registry v2 as a local container image registry for storing role images. Running DAOS as a systemd service provides automatic startup, restart on failure, and proper dependency management with Docker.

## What Changes

- Add systemd service file for DAOS daemon
- Add Docker and containerd as package dependencies (.deb and .rpm)
- Use systemd ExecStartPre to ensure Docker daemon is running
- Use systemd ExecStartPre to ensure Docker Registry v2 container is running (start if not running)
- Container uses `--restart=unless-stopped` policy
- Registry port configurable via `/opt/daos/registry/daos-registry.env` environment file

## Capabilities

### New Capabilities
- **docker-registry-v2**: Docker Registry v2 as package dependency with systemd-managed lifecycle

### Modified Capabilities
- None

## Impact

- **Package definitions**: Add docker.io/containerd to DEBIAN/control, docker-ce/containerd to daos.spec
- **New files**: packaging/systemd/daos.service
- **Configuration**: Registry port configurable via `/opt/daos/registry/daos-registry.env` (DAOS_REGISTRY_PORT=5000)
- **Dependencies**: docker.io (deb), docker-ce (rpm), containerd
- **Affected teams**: Operations (packaging), Backend

## Rollback Plan

- Remove systemd service file from packaging
- Remove Docker/containerd from package dependencies
- No database migrations required
