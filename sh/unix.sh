#!/bin/bash

set -e

REPO="fahidsarker/tmux_init" 
BINARY_NAME="tmux_init"
INSTALL_DIR="$HOME/.local/bin"

OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Normalize arch
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  i386|i686) ARCH="386" ;;
esac

# Get latest release tag
LATEST=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep tag_name | cut -d '"' -f 4)

FILE="${BINARY_NAME}_${LATEST#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$LATEST/$FILE"

echo "Downloading $FILE..."
echo "URL: $URL"
# mkdir -p "$INSTALL_DIR"
# curl -L "$URL" | tar -xz -C "$INSTALL_DIR"

# echo "✅ Installed to $INSTALL_DIR"

# # PATH help
# if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
#   echo ""
#   echo "👉 Add this to your shell config:"
#   echo "export PATH=\"\$HOME/.local/bin:\$PATH\""
#   echo ""
#   echo "Then run: source ~/.bashrc or source ~/.zshrc"
# fi
