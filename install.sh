#!/bin/bash
set -e

# Detect OS and architecture
detect_os_arch() {
  OS="$(uname -s)"
  ARCH="$(uname -m)"
  case "$OS" in
  Darwin)
    OS_NAME="darwin"
    ;;
  Linux)
    OS_NAME="linux"
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
  esac
  case "$ARCH" in
  arm64 | aarch64)
    ARCH_NAME="arm64"
    ;;
  x86_64 | amd64)
    ARCH_NAME="amd64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
  esac
}

detect_os_arch

BINARY_URL="https://github.com/jae-labs/flashcards/releases/latest/download/flashcards-${OS_NAME}-${ARCH_NAME}"

# Download binary
curl -fsSL "$BINARY_URL" -o flashcards && chmod +x flashcards

# macOS: remove quarantine attribute
if [ "$OS_NAME" = "darwin" ]; then
  xattr -dr com.apple.quarantine flashcards || true
fi

# Prompt user to add the binary to PATH with copy/paste commands
cat <<'INFO'

Installation complete.

To use "flashcards" from anywhere, move the downloaded binary into a directory that's in your PATH.

For example, you can run the following command:
  sudo mv ./flashcards /usr/local/bin/
INFO
