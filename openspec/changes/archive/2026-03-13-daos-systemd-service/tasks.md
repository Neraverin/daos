## 1. Package Dependencies

- [x] 1.1 Update packaging/deb/DEBIAN/control to add docker.io, containerd to Depends
- [x] 1.2 Update packaging/rpm/daos.spec to add docker-ce, containerd to Requires

## 2. Systemd Service

- [x] 2.1 Create packaging/systemd/daos.service with ExecStartPre directives
- [x] 2.2 Add EnvironmentFile directive to load /opt/daos/registry/daos-registry.env
- [x] 2.3 Configure registry port from DAOS_REGISTRY_PORT environment variable

## 3. Environment Configuration

- [x] 3.1 Create packaging/systemd/daos-registry.env with default DAOS_REGISTRY_PORT=5000
- [x] 3.2 Ensure /opt/daos/registry/ directory is created on install

## 4. Package Installation Scripts

- [x] 4.1 Create deb postinst script to install service file and run systemctl daemon-reload
- [x] 4.2 Create rpm %post script to install service file and run systemctl daemon-reload
- [x] 4.3 Add docker pull registry:2 to pre-install/postinstall scripts

## 5. Directory Setup

- [x] 5.1 Update Makefile to include systemd directory in packaging
- [x] 5.2 Update rpm spec to create /opt/daos/registry/ directory
