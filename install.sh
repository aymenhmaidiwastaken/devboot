#!/bin/sh
# DevBoot installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/aymenhmaidiwastaken/devboot/main/install.sh | sh
set -e

REPO="aymenhmaidiwastaken/devboot"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Get latest version
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Failed to determine latest version"
    exit 1
fi

echo "Installing devboot v${VERSION} (${OS}/${ARCH})..."

# Download
FILENAME="devboot_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

curl -fsSL "$URL" -o "$TMPDIR/$FILENAME"
tar xzf "$TMPDIR/$FILENAME" -C "$TMPDIR"

# Install
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMPDIR/devboot" "$INSTALL_DIR/devboot"
else
    sudo mv "$TMPDIR/devboot" "$INSTALL_DIR/devboot"
fi

chmod +x "$INSTALL_DIR/devboot"

echo "devboot v${VERSION} installed to ${INSTALL_DIR}/devboot"
echo ""
echo "Get started:"
echo "  devboot init        # create devboot.yaml"
echo "  devboot apply       # apply configuration"
echo "  devboot doctor      # check your environment"
