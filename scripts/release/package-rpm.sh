#!/usr/bin/env bash

SPINUP_VERSION=$(cat ./common/.version | sed 's/^v//')
SOURCES_DIR=./packaging/rpmbuild/SOURCES

# Add the version number to the spec file
sed -i "s/{{version}}/${SPINUP_VERSION}/g" packaging/rpmbuild/SPECS/spinup.spec

# Create the directory structure for the .rpm package
mkdir -p $SOURCES_DIR/spinup-${SPINUP_VERSION}
cp build/bin/spinup-${SPINUP_VERSION} $SOURCES_DIR/spinup-${SPINUP_VERSION}/
mv $SOURCES_DIR/spinup-${SPINUP_VERSION}/spinup-${SPINUP_VERSION} $SOURCES_DIR/spinup-${SPINUP_VERSION}/spinup

# Add the contents of the build/unix directory to the rpm sources
cp -r packaging/unix/. $SOURCES_DIR/spinup-${SPINUP_VERSION}/

# Add the postinstall script to the rpm sources
cp packaging/DEBIAN/postinst $SOURCES_DIR/spinup-${SPINUP_VERSION}/

# Set the permissions for the .rpm package files
sudo chown -R root:root $SOURCES_DIR/spinup-${SPINUP_VERSION}

# Create the tarball for the .rpm package
cd $SOURCES_DIR
tar --create --gzip --file spinup-${SPINUP_VERSION}.tar.gz spinup-${SPINUP_VERSION}
cd ../../..

# Remove the existing rpmbuild directory and copy the new one
rm -rf ~/rpmbuild
cp -r packaging/rpmbuild ~/

# Build the .rpm package
rpmbuild -bb ~/rpmbuild/SPECS/spinup.spec

# Copy the .rpm package to the current directory
cp ~/rpmbuild/RPMS/x86_64/spinup-${SPINUP_VERSION}-1.x86_64.rpm .
