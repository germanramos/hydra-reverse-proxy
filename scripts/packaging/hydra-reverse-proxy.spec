Name: hydra-reverse-proxy
Version: 0
Release: 1
Summary: Hydra Reverse Proxy
Source0: hydra-reverse-proxy-0.1.tar.gz
License: MIT
Group: custom
URL: https://github.com/innotech/hydra-reverse-proxy
BuildArch: x86_64
BuildRoot: %{_tmppath}/%{name}-buildroot
%description
Reverse proxy service to applications balanced by hydra.
%prep
%setup -q
%build
%install
install -m 0755 -d $RPM_BUILD_ROOT/usr/local/hydra-reverse-proxy
install -m 0755 hydra-reverse-proxy $RPM_BUILD_ROOT/usr/local/hydra-reverse-proxy/hydra-reverse-proxy

install -m 0755 -d $RPM_BUILD_ROOT/etc/init.d
install -m 0755 hydra-reverse-proxy-init.d.sh $RPM_BUILD_ROOT/etc/init.d/hydra-reverse-proxy

install -m 0755 -d $RPM_BUILD_ROOT/etc/hydra-reverse-proxy
install -m 0644 hydra-reverse-proxy.conf $RPM_BUILD_ROOT/etc/hydra-reverse-proxy/hydra-reverse-proxy.conf
%clean
rm -rf $RPM_BUILD_ROOT
%post
echo   You should edit config file /etc/hydra-reverse-proxy/hydra-reverse-proxy.conf
echo   When finished, you may want to run \"update-rc.d hydra-reverse-proxy defaults\"
%files
/usr/local/hydra-reverse-proxy/hydra-reverse-proxy
/etc/init.d/hydra-reverse-proxy
%dir /etc/hydra-reverse-proxy
/etc/hydra-reverse-proxy/hydra-reverse-proxy.conf
/etc/init.d/hydra-reverse-proxy
