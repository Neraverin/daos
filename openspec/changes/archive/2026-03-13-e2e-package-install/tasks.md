## 1. Test Setup

- [x] 1.1 Create tests/e2e/ directory
- [x] 1.2 Create package_install_test.go with basic test structure

## 2. Helper Functions

- [x] 2.1 Implement requireRoot() - fail if not root
- [x] 2.2 Implement detectOS() - read /etc/os-release for ID

## 3. Install & Verify Functions

- [x] 3.1 Implement installDAOS() - install deb or rpm based on OS
- [x] 3.2 Implement verifyInstallation() - check /opt/daos/, daemon, registry

## 4. Uninstall & Verify Functions

- [x] 4.1 Implement uninstallDAOS() - remove deb or rpm based on OS
- [x] 4.2 Implement verifyUninstallation() - check /opt/daos/ removed, daemon stopped

## 5. Main Test

- [x] 5.1 Implement TestPackageInstallUninstall test function
- [x] 5.2 Add t.Cleanup() for rollback on test failure

## 6. Makefile

- [x] 6.1 Add e2e-test target to Makefile
