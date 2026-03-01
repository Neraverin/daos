Name:           daos
Version:        %{_version}
Release:        1%{?dist}
Summary:        Deployment and Orchestration Service
License:        MIT
URL:            https://github.com/Neraverin/daos
BuildArch:      x86_64
Source0:        daos-%{version}.tar.gz

%description
DAOS is a service for managing host deployments with Docker Compose packages.

%prep
# Nothing to prep
tar -xzf %{SOURCE0}

%build
# Nothing to build

%install
mkdir -p %{buildroot}/opt/daos
mkdir -p %{buildroot}/etc/daos
cp daemon %{buildroot}/opt/daos/
cp tui %{buildroot}/opt/daos/

%files
%dir /opt/daos
%dir /etc/daos
/opt/daos/daemon
/opt/daos/tui

%changelog
* Sun Mar 01 2026 DAOS Team <team@daos.local> - 0.1.0-1
- Initial RPM release
