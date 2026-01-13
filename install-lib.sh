#!/bin/bash
# install-lib.sh - Install SochDB native library for Go SDK
#
# This script installs libsochdb_storage to system paths so the Go SDK
# can find it automatically without manual CGO_LDFLAGS setup.

set -e

# Determine OS
OS="$(uname -s)"
ARCH="$(uname -m)"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔧 SochDB Native Library Installer"
echo "=================================="
echo ""

# Check if running as root for system install
NEED_SUDO=""
if [ "$EUID" -ne 0 ]; then 
    NEED_SUDO="sudo"
    echo -e "${YELLOW}Note: System installation requires sudo privileges${NC}"
fi

# Find the native library
SOCHDB_ROOT="${SOCHDB_ROOT:-}"
if [ -z "$SOCHDB_ROOT" ]; then
    # Try to find it relative to this script
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    if [ -f "$SCRIPT_DIR/../sochdb/target/debug/libsochdb_storage.dylib" ]; then
        SOCHDB_ROOT="$SCRIPT_DIR/../sochdb"
    elif [ -f "$SCRIPT_DIR/../sochdb/target/release/libsochdb_storage.dylib" ]; then
        SOCHDB_ROOT="$SCRIPT_DIR/../sochdb"
    elif [ -f "$SCRIPT_DIR/../sochdb/target/debug/libsochdb_storage.so" ]; then
        SOCHDB_ROOT="$SCRIPT_DIR/../sochdb"
    else
        echo -e "${RED}Error: Cannot find libsochdb_storage${NC}"
        echo "Please set SOCHDB_ROOT environment variable to the sochdb repository path"
        echo "Example: export SOCHDB_ROOT=/path/to/sochdb"
        exit 1
    fi
fi

# Determine build type (prefer release, fallback to debug)
BUILD_TYPE="release"
if [ ! -d "$SOCHDB_ROOT/target/release" ]; then
    BUILD_TYPE="debug"
    echo -e "${YELLOW}Using debug build (release not found)${NC}"
fi

# Set library source and destination based on OS
case "$OS" in
    Darwin)
        LIB_SRC="$SOCHDB_ROOT/target/$BUILD_TYPE/libsochdb_storage.dylib"
        LIB_DEST="/usr/local/lib/libsochdb_storage.dylib"
        PKG_CONFIG_DEST="/usr/local/lib/pkgconfig/libsochdb_storage.pc"
        ;;
    Linux)
        LIB_SRC="$SOCHDB_ROOT/target/$BUILD_TYPE/libsochdb_storage.so"
        LIB_DEST="/usr/local/lib/libsochdb_storage.so"
        PKG_CONFIG_DEST="/usr/local/lib/pkgconfig/libsochdb_storage.pc"
        ;;
    MINGW*|MSYS*|CYGWIN*)
        echo -e "${RED}Windows installation not yet supported${NC}"
        echo "Please add the library path to your PATH environment variable manually"
        exit 1
        ;;
    *)
        echo -e "${RED}Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

# Check if library exists
if [ ! -f "$LIB_SRC" ]; then
    echo -e "${RED}Error: Library not found at $LIB_SRC${NC}"
    echo "Please build SochDB first:"
    echo "  cd $SOCHDB_ROOT"
    echo "  cargo build --release"
    exit 1
fi

echo "📦 Found library: $LIB_SRC"
echo "📍 Install destination: $LIB_DEST"
echo ""

# Create directories if needed
$NEED_SUDO mkdir -p "$(dirname "$LIB_DEST")"
$NEED_SUDO mkdir -p "$(dirname "$PKG_CONFIG_DEST")"

# Copy library
echo "📋 Copying library..."
$NEED_SUDO cp "$LIB_SRC" "$LIB_DEST"
$NEED_SUDO chmod 755 "$LIB_DEST"

# Create pkg-config file
echo "📝 Creating pkg-config file..."
cat > /tmp/libsochdb_storage.pc <<EOF
prefix=/usr/local
exec_prefix=\${prefix}
libdir=\${exec_prefix}/lib
includedir=\${prefix}/include

Name: libsochdb_storage
Description: SochDB Native Storage Library
Version: 0.4.0
Libs: -L\${libdir} -lsochdb_storage
Cflags: -I\${includedir}
EOF

$NEED_SUDO mv /tmp/libsochdb_storage.pc "$PKG_CONFIG_DEST"

# Update library cache on Linux
if [ "$OS" = "Linux" ]; then
    echo "🔄 Updating library cache..."
    $NEED_SUDO ldconfig
fi

echo ""
echo -e "${GREEN}✅ Installation complete!${NC}"
echo ""
echo "The Go SDK will now work without manual CGO_LDFLAGS setup:"
echo "  go get github.com/sochdb/sochdb-go@v0.4.0"
echo "  go run your-app.go"
echo ""
echo "To uninstall:"
echo "  $NEED_SUDO rm $LIB_DEST"
echo "  $NEED_SUDO rm $PKG_CONFIG_DEST"
