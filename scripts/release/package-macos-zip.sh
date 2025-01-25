#!/usr/bin/env bash

SPINUP_VERSION=$(cat ./common/.version | sed 's/^v//')

MAC_OS_DIR="./packaging/macos"
SPINUP_SHARE_DIR="/usr/local/share/spinup"

# Define the path to the binary file
BIN_FILE="./build/bin/Spinup.app/Contents/MacOS/spinup"

# Create the mac os directory
mkdir -p $MAC_OS_DIR

# Create the MacOS directory structure
mkdir -p "$MAC_OS_DIR/etc"
mkdir -p "$MAC_OS_DIR$SPINUP_SHARE_DIR"

# Copy the contents of the build/unix directory to the MacOS directory
cp -r "./packaging/unix/etc/." $MAC_OS_DIR/etc
cp -r "./packaging/unix/usr/share/spinup/." $MAC_OS_DIR$SPINUP_SHARE_DIR

# Copy the binary file to the MacOS directory
mkdir -p "$MAC_OS_DIR$SPINUP_SHARE_DIR/bin"
cp $BIN_FILE "$MAC_OS_DIR$SPINUP_SHARE_DIR/bin"

# Create a zip file containing the contents of the MacOS directory
(cd $MAC_OS_DIR && zip -r "../../spinup-${SPINUP_VERSION}-macos.zip" .)
