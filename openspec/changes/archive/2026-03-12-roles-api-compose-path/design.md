## Context

The current Roles API accepts `compose_content` as inline content in POST/PUT requests. This requires clients to read the entire file and send it in the request body. The proposed change allows clients to specify a `role_path` - a folder path on the daemon's host filesystem - instead of inline content. The daemon will read all role files (including compose files) from the specified folder.

Current constraints:
- Daemon runs as a REST API server
- TUI client connects to daemon via HTTP
- SQLite database stores role metadata
- Role files are stored in a folder on the daemon's host

Stakeholders: DAOS users who manage Docker Compose deployments via the API and TUI.

## Goals / Non-Goals

**Goals:**
- Allow clients to specify role folder path instead of inline content
- Folder path must be absolute
- Read role files (including compose) from daemon's host filesystem
- Validate folder existence and readability
- Update TUI to use folder path selection
- Maintain the same API response format

**Non-Goals:**
- Support remote folder paths (URLs, remote mounts) - only local daemon filesystem
- Keep backward compatibility with `compose_content` - this is a breaking change

## Decisions

### 1. API Field Change: `compose_content` → `role_path`
**Decision:** Replace `compose_content` with `role_path` in role create/update requests.

**Rationale:** Simplifies the API contract and removes the need for clients to send large payloads. The daemon reads all role files from the folder directly on its filesystem.

**Alternative Considered:** Support both fields simultaneously. Rejected because it adds complexity and the proposal explicitly calls for a breaking change.

### 2. Folder Reading Location
**Decision:** Read role files in the API handler (daemon side).

**Rationale:** The daemon has access to the local filesystem. Reading in the handler keeps validation logic centralized.

**Alternative Considered:** Read in a separate service layer. Rejected for simplicity - the change is straightforward enough to handle in the handler.

### 3. TUI Implementation
**Decision:** Use native folder picker dialog for path selection.

**Rationale:** Provides better UX than manual path entry with validation.

**Alternative Considered:** Text input with manual path entry. Rejected - prone to typos and doesn't validate path exists.

### 4. Path Validation
**Decision:** Validate path is absolute and resolves to a directory.

**Rationale:** Security concern - prevent path traversal attacks.

**Alternative Considered:** Allow relative paths. Rejected - relative paths are ambiguous and can lead to security issues.

## Risks / Trade-offs

- **[Risk]** Folder not found → **Mitigation**: Return 404 with clear error message "role folder not found at path"
- **[Risk]** Permission denied → **Mitigation**: Return 403 with error message
- **[Risk]** Daemon and TUI on different hosts → **Mitigation**: Document requirement; users must ensure filesystem access

## Migration Plan

1. Deploy updated daemon with new API
2. Update TUI to use folder path selection
3. Update any API client libraries
4. Rollback: Revert daemon to previous version if issues arise

## Open Questions

- Should we add a configuration option to specify allowed base directories for role folders?
  - No.
