#!/bin/bash

set -e

REPO="fahidsarker/tmux_init"
BINARY_NAME="tinit"
INSTALL_DIR="$HOME/.local/bin/tmux_init"

OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Normalize arch
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  i386|i686) ARCH="386" ;;
  *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest release tag
LATEST=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep tag_name | cut -d '"' -f 4)

FILE="${BINARY_NAME}_${LATEST#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$LATEST/$FILE"

echo "📦 Downloading $FILE..."
mkdir -p "$INSTALL_DIR"
curl -L "$URL" | tar -xz -C "$INSTALL_DIR"

echo "✅ Installed to $INSTALL_DIR"

# PATH update
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  USER_SHELL=$(basename "$SHELL")

  case "$USER_SHELL" in
    bash)
      CONFIG="$HOME/.bashrc"
      ;;
    zsh)
      CONFIG="$HOME/.zshrc"
      ;;
    ksh)
      CONFIG="$HOME/.kshrc"
      ;;
    *)
      CONFIG=""
      ;;
  esac

  if [[ -n "$CONFIG" ]]; then
    echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$CONFIG"
    echo "✅ Added to PATH in $CONFIG"
    echo "👉 Run: source $CONFIG or restart your terminal"
  else
    echo ""
    echo "⚠️ Your shell ($USER_SHELL) is not directly supported for PATH update."
    echo "👉 Please add the following line to your shell config manually:"
    echo ""
    echo "    export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
    echo "Then restart your terminal."
    
    if [[ "$USER_SHELL" == "fish" ]]; then
      echo "🎣 For fish shell users:"
      echo "    set -U fish_user_paths \$HOME/.local/bin \$fish_user_paths"
    fi
  fi
else
  echo "ℹ️ $INSTALL_DIR is already in PATH"
fi
