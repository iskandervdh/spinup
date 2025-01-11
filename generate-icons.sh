#!/bin/bash

SOURCE=images/icon-large.png
SIZES=(1024 512 256 128 64 32 16)
PIXMAPS_DIR=./build/share/pixmaps

mkdir -p ./build/share/icons/hicolor

for size in "${SIZES[@]}"; do
    dir="./build/share/icons/hicolor/${size}x${size}/apps"
    dir_2x="./build/share/icons/hicolor/${size}x${size}@2x/apps"

    # Create the normal icon directory and convert the icon to the correct size
    mkdir -p $dir
    convert $SOURCE -resize "${size}x${size}" $dir/spinup.png

    if [[ "$size" -ne 512 && "$size" -ne 1024 ]]; then
        # Do the same for the 2x icon if the size is not too large
        mkdir -p $dir_2x
        convert $SOURCE -resize "${size}x${size}" $dir_2x/spinup.png
    fi
done

# Convert and copy a 256x256 version to the pixmaps directory
mkdir -p $PIXMAPS_DIR
convert $SOURCE -resize "${size}x${size}" $PIXMAPS_DIR/spinup.png
