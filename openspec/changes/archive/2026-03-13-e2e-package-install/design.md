## Context

DAOS packages need automated end-to-end testing to validate the installation and uninstallation process. The test must work on both Debian-based and RHEL-based Linux distributions.

Current state:
- No E2E tests for package installation
- Manual verification required for each release

Constraints:
- Test must run as root (required for package installation and systemctl)
- Must detect OS type to use correct package manager
- Must build package before testing
- Test is destructive (installs/removes packages)

## Goals / Non-Goals

**Goals:**
- Create automated E2E test for package install/uninstall
- Support both .deb and .rpm packages
- Verify all installation artifacts are present
- Verify daemon and Docker registry are running after install
- Verify clean removal after uninstall

**Non-Goals:**
- Test upgrade scenarios (not covered)
- Test on multiple Linux distributions in CI
- Integration with existing unit tests

## Decisions

### 1. Test Location

**Decision:** `tests/e2e/package_install_test.go`

**Rationale:** Follows existing test structure under `tests/`, separate directory for E2E tests.

### 2. Root Requirement

**Decision:** Fail with `t.Fatal()` if not running as root

**Rationale:** Package installation requires root. Test cannot proceed without it.

```go
func requireRoot(t *testing.T) {
    if os.Geteuid() != 0 {
        t.Fatal("This test must be run as root")
    }
}
```

### 3. OS Detection

**Decision:** Check `/etc/os-release` for ID field

**Rationale:** Standard way to detect Linux distribution. Map to package type:
- `debian`, `ubuntu` → `.deb`
- `rhel`, `centos`, `fedora` → `.rpm`

### 4. Package Build

**Decision:** Build packages via Makefile before running test

**Rationale:** Test should not be responsible for building. Use `make e2e-test` which runs `make build deb rpm` first.

```makefile
e2e-test: build deb rpm
	sudo go test -v ./tests/e2e/...
```

### 5. Package Path

**Decision:** Use hardcoded paths matching Makefile output

**Rationale:** Simple and predictable:
- Deb: `bin/daos_0.1.0_amd64.deb`
- Rpm: `bin/daos-0.1.0-1.x86_64.rpm`

### 6. Verification Checks

**Decision:** Check specific items for install/uninstall

**Install verification:**
- `/opt/daos/` directory exists
- `/opt/daos/daemon` executable exists
- `/opt/daos/tui` executable exists
- `/opt/daos/registry/` directory exists
- `systemctl is-active daos` returns "active"
- `docker ps` contains "daos-registry"

**Uninstall verification:**
- `/opt/daos/` does not exist
- `systemctl is-active daos` returns error (not running)

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Test leaves system in dirty state on failure | Use `t.Cleanup()` for rollback |
| Package build fails | Fail test early with clear error |
| Systemd not available (container) | Skip test if systemd not available |

**Trade-offs:**
- Test is slow (builds package + install/uninstall)
- Test is destructive - should never run in CI on shared infrastructure
- No parallelization possible

## Test Structure

```go
package e2e

import (
    "testing"
)

func TestPackageInstallUninstall(t *testing.T) {
    requireRoot(t)
    
    osType := detectOS(t)
    
    // Install
    installDAOS(t, osType)
    verifyInstallation(t)
    
    t.Cleanup(func() {
        uninstallDAOS(t, osType)
    })
    
    // Uninstall
    uninstallDAOS(t, osType)
    verifyUninstallation(t)
}
```

Run with:
```bash
make e2e-test
```

## Open Questions

- Should we add a test for upgrading packages? (Not in scope for now)
- Should we verify Docker registry can actually push/pull images? (Out of scope - just verify container running)
