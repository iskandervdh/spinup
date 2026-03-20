Name:       spinup
Version:    {{version}}
Release:    1%{?dist}
Summary:    Quickly spin up your multi command projects.
License:    MIT
ExclusiveArch:  x86_64

Source0:    %{name}-%{version}.tar.gz
Requires:   dnsmasq, libgtk-3-0, libwebkit2gtk-4.0-dev, nginx

%description
Quickly spin up your multi command projects.

%prep
%setup -q

%install
mkdir -p %{buildroot}/etc/spinup/bin
cp %{_builddir}/%{name}-%{version}/%{name} %{buildroot}/etc/spinup/bin
cp %{_builddir}/%{name}-%{version}/postinst %{buildroot}/etc/spinup/bin/postinst
cp -r %{_builddir}/%{name}-%{version}/etc %{buildroot}
cp -r %{_builddir}/%{name}-%{version}/usr %{buildroot}

%post
sh /etc/spinup/bin/postinst

%clean
rm -rf %{buildroot}

%files
/etc/spinup/bin/%{name}
/etc/spinup/bin/postinst
/etc/spinup/config/dnsmasq.conf
/etc/sudoers.d/spinup
/usr/share/applications/spinup-app.desktop
/usr/share/icons/hicolor/1024x1024/apps/spinup.png
/usr/share/icons/hicolor/128x128/apps/spinup.png
/usr/share/icons/hicolor/128x128@2x/apps/spinup.png
/usr/share/icons/hicolor/16x16/apps/spinup.png
/usr/share/icons/hicolor/16x16@2x/apps/spinup.png
/usr/share/icons/hicolor/256x256/apps/spinup.png
/usr/share/icons/hicolor/256x256@2x/apps/spinup.png
/usr/share/icons/hicolor/32x32/apps/spinup.png
/usr/share/icons/hicolor/32x32@2x/apps/spinup.png
/usr/share/icons/hicolor/512x512/apps/spinup.png
/usr/share/icons/hicolor/64x64/apps/spinup.png
/usr/share/icons/hicolor/64x64@2x/apps/spinup.png
/usr/share/pixmaps/spinup.png

%changelog
* Sun Jan 19 2025 Iskander <iskandervdh@gmail.com> - 0.13.0
- First version of spinup as an RPM package
