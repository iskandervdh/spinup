#!/usr/bin/env bash

SPINUP_VERSION=$(cat ./common/.version | sed 's/^v//')

MAC_OS_DIR="./build/macos"

# Define the path to the binary file
BIN_FILE="./build/bin/Spinup.app/Contents/MacOS/spinup"

# Create the mac os directory
mkdir -p $MAC_OS_DIR

# Copy the contents of the build/unix directory to the MacOS directory
cp -r "./build/unix/." $MAC_OS_DIR

# Copy the binary file to the MacOS directory
mkdir -p "$MAC_OS_DIR/usr/share/spinup/bin"
cp $BIN_FILE "$MAC_OS_DIR/usr/share/spinup/bin"

# Move the sudoers file to the config directory
mv "$MAC_OS_DIR/etc/sudoers.d/spinup" "$MAC_OS_DIR/usr/share/spinup/config/sudoers"

# Create a zip file containing the contents of the MacOS directory
(cd $MAC_OS_DIR && zip -r "../../spinup-${SPINUP_VERSION}-macos.zip" .)
