## Why

DAOS packages need end-to-end testing to verify installation and uninstallation work correctly on different Linux distributions. Currently, there is no automated test that validates the full installation lifecycle, including Docker Registry v2 container management via systemd.

## What Changes

- Add E2E test in `tests/e2e/package_install_test.go`
- Detect OS type (Debian/Ubuntu vs RHEL/CentOS/Fedora)
- Build package before test (via `make deb` or `make rpm`)
- Test installation: verify all files present, daemon running, registry running
- Test uninstallation: verify files removed, daemon stopped, service removed
- Require root privileges to run (fail gracefully if not root)
- Run as separate test binary, not part of `make test`

## Capabilities

### New Capabilities
- **package-install-e2e**: End-to-end test for DAOS package installation and uninstallation

### Modified Capabilities
- None

## Impact

- **New file**: `tests/e2e/package_install_test.go`
- **Test execution**: `sudo go test -v ./tests/e2e/...`
- **Dependencies**: Requires root, systemctl, dpkg/rpm
- **Affected teams**: QA, Operations

## Rollback Plan

- Remove `tests/e2e/` directory
- No database migrations required
