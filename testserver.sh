#!/bin/bash

OPENTTD_URL="https://cdn.openttd.org/openttd-releases/12.2/openttd-12.2-linux-generic-amd64.tar.xz"
OPENGFX_URL="https://cdn.openttd.org/opengfx-releases/7.1/opengfx-7.1-all.zip"

# Warning: this value is also in .gitignore!
SERVER_FOLDER="./_testserver"

if [[ $1 == "-r" ]]; then 
    read -rp "Remove server at $SERVER_FOLDER? [Y/n] " yn
    case $yn in
        [Nn]* ) exit;;
    esac

    echo "Removing server..."
    rm -r $SERVER_FOLDER
    
    echo "Finished! Server successfully removed from $SERVER_FOLDER"

    exit
fi

if [[ -e "$SERVER_FOLDER" ]]; then
    echo "Test server found at $SERVER_FOLDER"
    echo "Starting..."

    $SERVER_FOLDER/openttd -D
else
    echo "Test server not found"
    
    read -rp "Create server at $SERVER_FOLDER? [Y/n] " yn
    case $yn in
        [Nn]* ) echo "Aborting..."; exit;;
    esac
    
    mkdir "$SERVER_FOLDER"
    mkdir "$SERVER_FOLDER/.temp"

    cd "$SERVER_FOLDER/.temp/" || (echo "Error: could not create folder! Aborting..." && exit)

    echo "Downloading OpenTTD..."
    wget -q -O openttd.tar.xz $OPENTTD_URL
    echo "Downloading OpenGFX..."
    wget -q -O opengfx.zip $OPENGFX_URL
    echo "Downloads complete!"

    echo "Extracting OpenTTD..."
    tar --strip-components=1 -C .. -xf openttd.tar.xz > /dev/null
    echo "Extracting OpenGFX..."
    unzip -d ../baseset opengfx.zip > /dev/null
    echo "Extrations complete!"
    
    echo "Cleaning up..."
    cd ..
    rm -r .temp

    echo "Finished! Server successfully created at $SERVER_FOLDER"
    
    read -rp "Start server now? [Y/n] " yn
    case $yn in
        [Nn]* ) exit;;
    esac

    ./openttd -D
fi