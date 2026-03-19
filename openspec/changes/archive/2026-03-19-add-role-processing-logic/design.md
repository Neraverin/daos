## Context

Currently, DAOS can store and manage roles with Docker Compose configurations. Roles may include Docker images bundled as `.tar` archives referenced via `Role.Definitions.ImageFile`. These images are not being processed during role deployment.

The daemon runs on a host where Docker is available and can execute commands. Images need to be loaded into a Docker registry so deployed containers can pull them.

## Goals / Non-Goals

**Goals:**
- Validate `ImageFile` path exists and has `.tar` extension
- Load Docker images from role archives using `docker load`
- Tag and push loaded images to configured Docker registry
- Provide clear error messages for missing files or failed Docker operations

**Non-Goals:**
- Building images from Dockerfiles (role archives are pre-built)
- Managing image lifecycle (retention, cleanup)
- Remote Docker daemon support (operates on local Docker)

## Decisions

### 1. Add Docker registry configuration to daemon config

**Decision:** Add `Docker` section to `Config` with `Registry` field and `ImageLoadTimeout`.

**Rationale:** The registry URL must be configurable since different environments use different registries. Timeout prevents indefinite hangs on large image loads.

```go
type DockerConfig struct {
    Registry         string `yaml:"registry"`
    ImageLoadTimeout int    `yaml:"image_load_timeout"` // seconds
}
```

### 2. Process images during role upload, not during deployment

**Decision:** Load images when a role is created/updated via API, not during host deployment.

**Rationale:** Catches errors early at role creation time. Deployment failures are harder to diagnose. Images only need to be pushed once per role version.

**Alternative:** Load during deployment - rejected because it slows down deployments and duplicates work if the same role is deployed multiple times.

### 3. Use `docker load` command for image import

**Decision:** Use `docker load -i <path>` instead of importing via Docker API.

**Rationale:** The `docker load` CLI handles various tar formats automatically. Simpler than using the Docker SDK for this use case.

### 4. Tag images with registry prefix before push

**Decision:** After loading, tag the image with `docker tag <image> <registry>/<image>` before push.

**Rationale:** Standard Docker registry workflow. Allows local testing before registry push.

## Risks / Trade-offs

- **[Risk]** Large images may take significant time to load/push → **Mitigation:** `ImageLoadTimeout` config prevents indefinite hangs
- **[Risk]** Registry authentication not configured → **Mitigation:** Require registry to be accessible without auth initially; add auth support later
- **[Risk]** Image file path is relative and may not resolve → **Mitigation:** Validate path exists before loading; use role folder as base directory
