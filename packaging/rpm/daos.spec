Name:           daos
Version:        %{_version}
Release:        1%{?dist}
Summary:        Deployment and Orchestration Service
License:        MIT
URL:            https://github.com/Neraverin/daos
BuildArch:      x86_64
Source0:        daos-%{version}.tar.gz
Requires:       docker-ce, containerd

%description
DAOS is a service for managing host deployments with Docker Compose packages.

%prep
# Nothing to prep
tar -xzf %{SOURCE0}

%build
# Nothing to build

%install
mkdir -p %{buildroot}/opt/daos
mkdir -p %{buildroot}/opt/daos/registry
mkdir -p %{buildroot}/etc/daos
cp daemon %{buildroot}/opt/daos/
cp tui %{buildroot}/opt/daos/
cp systemd/daos.service %{buildroot}/opt/daos/
cp systemd/daos-registry.env %{buildroot}/opt/daos/registry/

%post
# Create directories
mkdir -p /opt/daos/registry
mkdir -p /etc/daos

# Install systemd service
if [ -d /etc/systemd/system ]; then
    cp /opt/daos/daos.service /etc/systemd/system/
    systemctl daemon-reload
fi

# Pull registry image
if command -v docker &> /dev/null; then
    docker pull registry:2 || true
fi

%preun
if [ $1 -eq 0 ]; then
    systemctl stop daos || true
    systemctl disable daos || true
fi

%postun
if [ $1 -eq 0 ]; then
    rm -f /etc/systemd/system/daos.service
    systemctl daemon-reload
fi

%files
%dir /opt/daos
%dir /etc/daos
/opt/daos/daemon
/opt/daos/tui

%changelog
* Sun Mar 01 2026 DAOS Team <team@daos.local> - 0.1.0-1
- Initial RPM release
