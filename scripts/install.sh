#!/bin/bash
set -e

echo "🚀 Detecting system..."

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

[ "$ARCH" = "x86_64" ] && ARCH="amd64"
[ "$ARCH" = "arm64" ] && ARCH="arm64"
[ "$ARCH" = "aarch64" ] && ARCH="arm64"

echo "📦 System: ${OS}/${ARCH}"

echo "🔍 Checking for latest version..."
VERSION=$(curl -s https://api.github.com/repos/ksckaan1/crtui/releases/latest | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | sed 's/v//')
echo "✅ Latest version: ${VERSION}"

URL="https://github.com/ksckaan1/crtui/releases/download/v${VERSION}/crtui_${VERSION}_${OS}_${ARCH}.tar.gz"

echo "📥 Downloading crtui ${VERSION}..."
echo "   URL: ${URL}"
TMPDIR=$(mktemp -d)
curl -fSL "$URL" -o "$TMPDIR/crtui.tar.gz"
echo "✅ Download complete"

echo "📦 Extracting..."
if [ "$OS" = "darwin" ]; then
    tar -xf "$TMPDIR/crtui.tar.gz" -C "$TMPDIR"
else
    tar -xzf "$TMPDIR/crtui.tar.gz" -C "$TMPDIR"
fi

echo "⚙️  Installing to /usr/local/bin/crtui..."
sudo mv -f "$TMPDIR/crtui" /usr/local/bin/crtui
chmod +x /usr/local/bin/crtui

rm -rf "$TMPDIR"

echo ""
echo "✅ crtui ${VERSION} installed successfully!"
echo "Run 'crtui' to get started."
