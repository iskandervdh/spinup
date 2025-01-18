#!/usr/bin/env bash

SPINUP_VERSION=$(cat ./common/.version | sed 's/^v//')

MAC_OS_FOLDER="./build/MacOS"

# Define the path to the binary file
BIN_FILE="./build/bin/spinup/Contents/MacOS/spinup"

ls ./build/bin/

# Create a MacOS directory
mkdir -p $MAC_OS_FOLDER

# Copy the contents of the build/unix directory to the MacOS directory
cp -r "./build/unix/" $MAC_OS_FOLDER

# Copy the binary file to the MacOS directory
cp $BIN_FILE "$MAC_OS_FOLDER/usr/share/spinup/bin/spinup"

# Create a zip file containing the contents of the MacOS directory
zip -r "spinup-${SPINUP_VERSION}-macos.zip" $MAC_OS_FOLDER

hdiutil create -volname "Spinup" -srcfolder $MAC_OS_FOLDER -ov -format UDZO "spinup-${SPINUP_VERSION}.dmg"
