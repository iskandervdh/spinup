#!/bin/sh

### This script updates the version number in the .version file

# Get the current version number from the common/.version file
VERSION=`cat common/.version`

# Increment the version number based on the argument passed to the script
if [ "$1" = "major" ]; then
    MAJOR=`echo $VERSION | cut -d. -f1 | cut -c2-`
    MAJOR=`expr $MAJOR + 1`
    VERSION="v$MAJOR.0.0"
elif [ "$1" = "minor" ]; then
    MINOR=`echo $VERSION | cut -d. -f2`
    MINOR=`expr $MINOR + 1`
    VERSION="v`echo $VERSION | cut -d. -f1 | cut -c2-`.${MINOR}.0"
elif [ "$1" = "patch" ]; then
    PATCH=`echo $VERSION | cut -d. -f3`
    PATCH=`expr $PATCH + 1`
    VERSION="v`echo $VERSION | cut -d. -f1 | cut -c2-`.$(echo $VERSION | cut -d. -f2).${PATCH}"
else
    echo "Usage: update-version.sh [major|minor|patch]"
    exit 1
fi

# Write the new version number to the common/.version file
echo $VERSION > common/.version

# Commit the change to the repository
git add common/.version

# Create a new commit with the updated version number
git commit -m "Bumped version number to $VERSION"

# Push the changes to the remote repository
git push

# Tag the commit with the new version number
git tag -a $VERSION -m "Version $VERSION"

# Push the tag to the remote repository
git push origin $VERSION

# Print a message indicating the version number has been updated
echo "Version number updated to $VERSION"
