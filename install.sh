#!/usr/bin/env bash
# CrunchyCleaner Install Script
# Usage: curl -sSL https://raw.githubusercontent.com/Knuspii/CrunchyCleaner/main/install.sh | sudo bash

set -e

REPO="Knuspii/CrunchyCleaner"
BINARY_NAME="crunchycleaner"
INSTALL_PATH="/usr/local/bin/$BINARY_NAME"

echo "Installing CrunchyCleaner..."

# Detect OS and Architecture
ARCH=$(uname -m)

if [ "$ARCH" == "aarch64" ] || [ "$ARCH" == "arm64" ]; then
    ARCH="-arm64"
else
    ARCH=""
fi

URL="https://github.com/$REPO/releases/latest/download/${BINARY_NAME}${ARCH}"

echo "Downloading $BINARY_NAME"
curl -L "$URL" -o "$BINARY_NAME"

echo "Setting permissions..."
chmod +x "$BINARY_NAME"

echo "Installing to $INSTALL_PATH..."
sudo install -m 755 "$BINARY_NAME" "$INSTALL_PATH"
sudo rm "$BINARY_NAME"

echo ""
echo "Done! You can now use 'CrunchyCleaner'"
echo "Type: 'crunchycleaner -h'"
