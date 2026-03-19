## Why

Roles may include Docker images bundled as `.tar` archives. Currently, these images are not being processed during role deployment. We need to load these images into the configured Docker registry so containers can use them.

## What Changes

- Add validation that `Role.Definitions.ImageFile` path exists and has `.tar` extension when specified
- Add logic to run `docker load` to import the image archive
- Tag and push loaded images to configured Docker registry
- Handle errors gracefully with clear messaging for missing files or failed docker operations

## Capabilities

### New Capabilities
- `role-image-processing`: Process Docker image archives from roles by validating paths, loading via `docker load`, and pushing to registry

### Modified Capabilities
- `role-management`: Add image file validation and processing during role deployment

## Impact

- Affected code: Role processing/deployment logic (likely in `pkg/deployment/` or similar)
- New dependency: Docker CLI for `load` and `push` commands
- Configuration: Docker registry must be configured in daemon config
