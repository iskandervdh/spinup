#!/usr/bin/env bash

SPINUP_VERSION=$(cat ./common/.version | sed 's/^v//')

MAC_OS_FOLDER="./build/MacOS"

# Define the path to the binary file
BIN_FILE="./build/bin/Spinup.app/Contents/MacOS/spinup"

# Create a MacOS directory
mkdir -p $MAC_OS_FOLDER

# Copy the contents of the build/unix directory to the MacOS directory
cp -r "./build/unix/." $MAC_OS_FOLDER

# Copy the binary file to the MacOS directory
mkdir -p "$MAC_OS_FOLDER/usr/share/spinup/bin"
cp $BIN_FILE "$MAC_OS_FOLDER/usr/share/spinup/bin"

# Copy postinstall script to the MacOS directory
cp "./build/DEBIAN/postinst" $MAC_OS_FOLDER
mv "$MAC_OS_FOLDER/postinst" "$MAC_OS_FOLDER/post_install.sh"

# Create a zip file containing the contents of the MacOS directory
(cd $MAC_OS_FOLDER && zip -r "../../spinup-${SPINUP_VERSION}-macos.zip" .)
