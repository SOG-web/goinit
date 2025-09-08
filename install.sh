#!/bin/bash

# GoInit CLI Installer
# Installs the GoInit CLI tool from GitHub releases

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $OS in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            echo -e "${RED}Unsupported OS: $OS${NC}"
            exit 1
            ;;
    esac

    case $ARCH in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    echo -e "${BLUE}Fetching latest version...${NC}"
    LATEST_VERSION=$(curl -s https://api.github.com/repos/SOG-web/goinit/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$LATEST_VERSION" ]; then
        echo -e "${RED}Failed to fetch latest version${NC}"
        exit 1
    fi
    echo -e "${GREEN}Latest version: $LATEST_VERSION${NC}"
}

# Download and install
install_binary() {
    BINARY_NAME="goinit-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi

    DOWNLOAD_URL="https://github.com/SOG-web/goinit/releases/download/${LATEST_VERSION}/${BINARY_NAME}.tar.gz"
    if [ "$OS" = "windows" ]; then
        DOWNLOAD_URL="https://github.com/SOG-web/goinit/releases/download/${LATEST_VERSION}/${BINARY_NAME}.zip"
    fi

    echo -e "${BLUE}Downloading $BINARY_NAME...${NC}"
    if ! curl -L -o "/tmp/${BINARY_NAME}.tar.gz" "$DOWNLOAD_URL"; then
        echo -e "${RED}Failed to download binary${NC}"
        exit 1
    fi

    echo -e "${BLUE}Installing...${NC}"
    if [ "$OS" = "windows" ]; then
        unzip -q "/tmp/${BINARY_NAME}.zip" -d /tmp/
        mv "/tmp/${BINARY_NAME}" "$INSTALL_DIR/goinit.exe"
    else
        tar -xzf "/tmp/${BINARY_NAME}.tar.gz" -C /tmp/
        mv "/tmp/${BINARY_NAME}" "$INSTALL_DIR/goinit"
        chmod +x "$INSTALL_DIR/goinit"
    fi

    # Cleanup
    rm -f "/tmp/${BINARY_NAME}.tar.gz"
    rm -f "/tmp/${BINARY_NAME}"

    echo -e "${GREEN}✅ GoInit installed successfully!${NC}"
    echo -e "${YELLOW}📖 Run 'goinit --help' to get started${NC}"
}

# Main installation process
main() {
    echo -e "${BLUE}🚀 GoInit CLI Installer${NC}"
    echo -e "${BLUE}========================${NC}"

    # Detect platform
    detect_platform
    echo -e "${BLUE}Detected platform: ${OS}/${ARCH}${NC}"

    # Determine install directory
    if [ "$OS" = "windows" ]; then
        INSTALL_DIR="$HOME/AppData/Local/Microsoft/WindowsApps"
    else
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
        export PATH="$INSTALL_DIR:$PATH"
    fi

    # Get latest version
    get_latest_version

    # Install binary
    install_binary

    # Verify installation
    if command -v goinit &> /dev/null; then
        echo -e "${GREEN}✅ Installation verified!${NC}"
        goinit --version
    else
        echo -e "${YELLOW}⚠️  Please add $INSTALL_DIR to your PATH${NC}"
        echo -e "${YELLOW}   export PATH=\"$INSTALL_DIR:\$PATH\"${NC}"
    fi
}

# Run main function
main "$@"