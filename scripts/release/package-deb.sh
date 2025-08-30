#!/usr/bin/env bash

SPINUP_VERSION=$(cat ./common/.version | sed 's/^v//')

for os_version in "" "-ubuntu24.04"; do
    # Create the directory structure for the .deb package
    mkdir -p deb/spinup-${SPINUP_VERSION}${os_version}/DEBIAN
    mkdir -p deb/spinup-${SPINUP_VERSION}${os_version}/etc/spinup/bin

    # Copy the necessary files to the .deb package directory
    cp build/bin/spinup-${SPINUP_VERSION}${os_version} deb/spinup-${SPINUP_VERSION}${os_version}/etc/spinup/bin/spinup
    cp packaging/DEBIAN/* deb/spinup-${SPINUP_VERSION}${os_version}/DEBIAN
    cp -r packaging/unix/usr deb/spinup-${SPINUP_VERSION}${os_version}

    # Update the control file with the current version number
    echo -e "\nVersion: $SPINUP_VERSION" >> deb/spinup-${SPINUP_VERSION}${os_version}/DEBIAN/control

    # Add the dependencies for the .deb package
    if [ "$os_version" = "" ]; then
        echo "Depends: libgtk-3-0, libwebkit2gtk-4.0-dev, nginx, dnsmasq" >> deb/spinup-${SPINUP_VERSION}${os_version}/DEBIAN/control
    else
        echo "Depends: libgtk-3-0, libwebkit2gtk-4.1-dev, nginx, dnsmasq" >> deb/spinup-${SPINUP_VERSION}${os_version}/DEBIAN/control
    fi

    # Set the permissions for the .deb package files
    sudo chown -R root:root deb/spinup-${SPINUP_VERSION}${os_version}

    # Build the .deb package
    dpkg-deb --build deb/spinup-${SPINUP_VERSION}${os_version}
done
