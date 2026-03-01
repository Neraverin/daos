%define _version 0.1.0
%define _datadir /opt/daos
%define _configdir /etc/daos
%define _logdir /var/log/daos

Name:           daos-daemon
Version:        %{_version}
Release:        1%{?dist}
Summary:        DAOS - Deployment and Orchestration Service Daemon
License:        Apache-2.0
URL:            https://github.com/Neraverin/daos
BuildArch:      x86_64 aarch64

Requires:       sqlite
Requires:       ansible

%description
DAOS - Deployment and Orchestration Service
REST API server for managing deployments via Ansible.

%prep
# No source package to unpack

%build
# No build needed for binary packages

%install
# Create directories
mkdir -p %{buildroot}%{_datadir}
mkdir -p %{buildroot}%{_configdir}
mkdir -p %{buildroot}%{_logdir}
mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/usr/lib/systemd/system

# Install binaries
install -m 755 bin/daemon %{buildroot}/usr/bin/daos-daemon
install -m 755 bin/tui %{buildroot}/usr/bin/daos-tui

# Install systemd service
install -m 644 packaging/deb/daos-daemon.service %{buildroot}/usr/lib/systemd/system/

# Install default config
cat > %{buildroot}%{_configdir}/daemon.yaml <<EOF
server:
  host: "127.0.0.1"
  port: 8080

database:
  path: "%{_datadir}/daos.db"
EOF

%pre
# Create daemon user if it doesn't exist
getent group daos > /dev/null || groupadd -r daos
getent passwd daos > /dev/null || \
    useradd -r -g daos -s /sbin/nologin -d %{_datadir} daos

%post
# Set ownership
chown -R daos:daos %{_datadir}
chown -R daos:daos %{_logdir}
chown daos:daos %{_configdir}/daemon.yaml

# Enable and start service
systemctl daemon-reload
systemctl enable daos-daemon.service
systemctl start daos-daemon.service || true

%preun
# Stop service on removal
if [ $1 -eq 0 ]; then
    systemctl stop daos-daemon.service || true
    systemctl disable daos-daemon.service || true
fi

%postun
# Reload systemd on upgrade
if [ $1 -ge 1 ]; then
    systemctl daemon-reload
fi
# Clean up on complete removal
if [ $1 -eq 0 ]; then
    rm -rf %{_datadir}
    rm -rf %{_logdir}
    userdel daos || true
    groupdel daos || true
fi

%files
%defattr(-,root,root,-)
%dir %attr(755,daos,daos) %{_datadir}
%dir %attr(755,daos,daos) %{_logdir}
%dir %attr(755,root,root) %{_configdir}
%attr(755,root,root) /usr/bin/daos-daemon
%attr(755,root,root) /usr/bin/daos-tui
%attr(644,root,root) /usr/lib/systemd/system/daos-daemon.service
%attr(600,daos,daos) %{_configdir}/daemon.yaml
