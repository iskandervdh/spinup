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
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/%{_bindir}
cp %{name} $RPM_BUILD_ROOT/%{_bindir}

%clean
rm -rf $RPM_BUILD_ROOT

%files
%{_bindir}/%{name}

%changelog
* Sun Jan 5 2025 Iskander <iskandervdh@gmail.com> - 0.12.0
- First version of spinup as an RPM package
