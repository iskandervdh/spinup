#!/usr/bin/env bash

SPINUP_VERSION=$(cat ./common/.version | sed 's/^v//')

for os_version in "" "-ubuntu24.04"; do
    mkdir -p deb/spinup-${SPINUP_VERSION}${os_version}/DEBIAN
    mkdir -p deb/spinup-${SPINUP_VERSION}${os_version}/usr/share/spinup/bin
    cp build/bin/spinup-${SPINUP_VERSION}${os_version} deb/spinup-${SPINUP_VERSION}${os_version}/usr/share/spinup/bin/spinup
    cp build/DEBIAN/* deb/spinup-${SPINUP_VERSION}${os_version}/DEBIAN
    cp -r build/unix/usr deb/spinup-${SPINUP_VERSION}${os_version}

    echo -e "\nVersion: $SPINUP_VERSION" >> deb/spinup-${SPINUP_VERSION}${os_version}/DEBIAN/control

    if [ "$webkit_version" -eq "40" ]; then
        echo "Depends: dnsmasq, libgtk-3-0, libwebkit2gtk-4.0-dev" >> deb/spinup-${SPINUP_VERSION}${os_version}/DEBIAN/control
    else
        echo "Depends: dnsmasq, libgtk-3-0, libwebkit2gtk-4.1-dev" >> deb/spinup-${SPINUP_VERSION}${os_version}/DEBIAN/control
    fi

    sudo chown -R root:root deb/spinup-${SPINUP_VERSION}${os_version}

    dpkg-deb --build deb/spinup-${SPINUP_VERSION}${os_version}
done
