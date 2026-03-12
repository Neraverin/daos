## Why

Currently, the Roles API requires clients to send the full compose file content (`compose_content`) in the request body. This creates unnecessary bandwidth overhead and requires clients to read files before sending. By allowing clients to specify a folder path, containing role files, on the daemon's host, we can simplify client implementation and reduce payload sizes, assuming API clients know the filesystem structure of the host where the daemon runs.

## What Changes

- Modify Roles API to accept `role_path` (string) instead of `compose_content` in create/update requests
- The daemon will read the role files (including compose) from the specified path on its host filesystem
- **BREAKING**: Remove `compose_content` field from role create/update requests
- Update TUI to use folder path selection instead of paste/input of compose content
- Daemon and TUI must run on the same host (or at least share filesystem access)

## Capabilities

### New Capabilities
- `compose-file-path`: Allow specifying compose folder path on daemon host instead of inline content

### Modified Capabilities
- `role-management`: Change from accepting inline `compose_content` to accepting `compose_path` pointing to file on daemon's filesystem

## Impact

- **Affected APIs**: Roles API (POST/PUT /roles endpoints)
- **Affected Components**: Daemon handlers, TUI models, API client libraries
- **Deployment**: Requires daemon and client to share filesystem (same host or network mount)
- **Rollback Plan**: Revert API to accept `compose_content` field; preserve backward compatibility by supporting both fields temporarily during transition period
