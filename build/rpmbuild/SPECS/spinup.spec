Name:       spinup
Version:    {{version}}
Release:    1%{?dist}
Summary:    Quickly spin up your multi command projects.
License:    MIT
BuildArch:  x86_64

Source0:    %{name}-%{version}.tar.gz
Requires:   bash

%description
Quickly spin up your multi command projects.

%prep
%setup -q

%install
mkdir -p %{buildroot}/usr/share/spinup/bin
cp %{_builddir}/%{name}-%{version}/%{name} %{buildroot}/usr/share/spinup/bin

%clean
rm -rf %{buildroot}

%files
/usr/share/spinup/bin/%{name}

%changelog
* Sun Jan 19 2025 Iskander <iskandervdh@gmail.com> - 0.13.0
- First version of spinup as an RPM package
