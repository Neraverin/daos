## ADDED Requirements

### Requirement: E2E test detects Linux distribution
The test SHALL detect the Linux distribution to determine which package type to use.

#### Scenario: Detect Debian-based system
- **GIVEN** the test runs on a Debian or Ubuntu system
- **WHEN** `detectOS()` is called
- **THEN** it SHALL return "deb"

#### Scenario: Detect RHEL-based system
- **GIVEN** the test runs on a RHEL, CentOS, or Fedora system
- **WHEN** `detectOS()` is called
- **THEN** it SHALL return "rpm"

### Requirement: E2E test requires root privileges
The test SHALL fail gracefully if not run as root.

#### Scenario: Test runs without root
- **WHEN** the test is executed by a non-root user
- **THEN** it SHALL fail with message "This test must be run as root"

### Requirement: Package files exist for installation
The test SHALL verify that package files exist before attempting installation.

#### Scenario: Deb package exists
- **GIVEN** test runs on Debian/Ubuntu
- **WHEN** installation is attempted
- **THEN** file `bin/daos_0.1.0_amd64.deb` SHALL exist

#### Scenario: Rpm package exists
- **GIVEN** test runs on RHEL/CentOS/Fedora
- **WHEN** installation is attempted
- **THEN** file `bin/daos-0.1.0-1.x86_64.rpm` SHALL exist

### Requirement: DAOS package installs successfully
The test SHALL verify that the DAOS package installs correctly.

#### Scenario: Install deb package
- **GIVEN** DAOS deb package exists
- **WHEN** `installDAOS()` is executed
- **THEN** `dpkg -i` SHALL complete without error
- **AND** DAOS service SHALL be started

#### Scenario: Install rpm package
- **GIVEN** DAOS rpm package exists
- **WHEN** `installDAOS()` is executed
- **THEN** `rpm -ivh` SHALL complete without error
- **AND** DAOS service SHALL be started

### Requirement: Installation artifacts are present
After installation, all required files and directories SHALL exist.

#### Scenario: Verify /opt/daos directory
- **GIVEN** DAOS is installed
- **WHEN** checking `/opt/daos/`
- **THEN** directory SHALL exist
- **AND** `daemon` executable SHALL exist
- **AND** `tui` executable SHALL exist

#### Scenario: Verify systemd files
- **GIVEN** DAOS is installed
- **WHEN** checking `/opt/daos/registry/`
- **THEN** `daos-registry.env` SHALL exist

### Requirement: DAOS daemon is running after installation
The test SHALL verify the DAOS daemon is active.

#### Scenario: Check daemon status
- **GIVEN** DAOS is installed
- **WHEN** `systemctl is-active daos` is executed
- **THEN** it SHALL return "active"

### Requirement: Docker Registry container is running
The test SHALL verify the Docker Registry container is running.

#### Scenario: Check registry container
- **GIVEN** DAOS is installed and started
- **WHEN** `docker ps` is executed
- **THEN** output SHALL contain "daos-registry"

### Requirement: DAOS package uninstalls successfully
The test SHALL verify that the DAOS package can be removed cleanly.

#### Scenario: Uninstall deb package
- **GIVEN** DAOS deb package is installed
- **WHEN** `dpkg -r daos` is executed
- **THEN** it SHALL complete without error

#### Scenario: Uninstall rpm package
- **GIVEN** DAOS rpm package is installed
- **WHEN** `rpm -e daos` is executed
- **THEN** it SHALL complete without error

### Requirement: Installation artifacts removed after uninstall
After uninstallation, all files and directories SHALL be removed.

#### Scenario: Verify /opt/daos removed
- **GIVEN** DAOS is uninstalled
- **WHEN** checking `/opt/daos/`
- **THEN** directory SHALL NOT exist

### Requirement: DAOS service removed after uninstall
The test SHALL verify the DAOS service is no longer running.

#### Scenario: Check daemon not running
- **GIVEN** DAOS is uninstalled
- **WHEN** `systemctl is-active daos` is executed
- **THEN** it SHALL return an error (service not found/not active)
