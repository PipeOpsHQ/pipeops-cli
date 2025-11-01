#!/bin/bash

# PipeOps CLI Installation Script
# This script detects your platform and installs the latest version of PipeOps CLI

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
REPO="PipeOpsHQ/pipeops-cli"
BINARY_NAME="pipeops"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${PURPLE}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🚀 PipeOps CLI Installer"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${NC}"
}

# Detect OS and architecture
detect_platform() {
    local os=""
    local arch=""

    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        FreeBSD*)   os="freebsd" ;;
        *)          print_error "Unsupported operating system: $(uname -s)" && exit 1 ;;
    esac

    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)   arch="x86_64" ;;
        arm64|aarch64)  arch="arm64" ;;
        armv7*)         arch="armv6" ;;  # Map armv7 to armv6 (closest available)
        armv6*)         arch="armv6" ;;
        i386|i686)      arch="i386" ;;
        *)              print_error "Unsupported architecture: $(uname -m)" && exit 1 ;;
    esac

    # Format OS name for GitHub releases
    case "$os" in
        linux)   os="Linux" ;;
        darwin)  os="Darwin" ;;
        windows) os="Windows" ;;
        freebsd) os="Freebsd" ;;
    esac

    PLATFORM="${os}"
    ARCH="${arch}"

    print_status "Detected platform: ${PLATFORM} ${ARCH}"
}

# Get the latest release version
get_latest_version() {
    print_status "Fetching latest release information..."

    local api_url="https://api.github.com/repos/${REPO}/releases/latest"

    if command -v curl >/dev/null 2>&1; then
        VERSION=$(curl -s "${api_url}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        VERSION=$(wget -qO- "${api_url}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi

    if [ -z "$VERSION" ]; then
        print_error "Failed to get latest version"
        exit 1
    fi

    print_status "Latest version: ${VERSION}"
}

# Download and install binary
install_binary() {
    local version="${VERSION#v}"  # Remove 'v' prefix if present
    local filename
    local url
    local temp_dir

    # Construct filename based on platform
    if [ "$PLATFORM" = "Windows" ]; then
        filename="${BINARY_NAME}-cli_${PLATFORM}_${ARCH}.zip"
    else
        filename="${BINARY_NAME}-cli_${PLATFORM}_${ARCH}.tar.gz"
    fi

    url="https://github.com/${REPO}/releases/download/${VERSION}/${filename}"

    print_status "Downloading ${filename}..."

    # Create temporary directory
    temp_dir=$(mktemp -d)
    trap "rm -rf $temp_dir" EXIT

    # Download the release
    local download_path="${temp_dir}/${filename}"

    if command -v curl >/dev/null 2>&1; then
        if ! curl -fsSL -o "$download_path" "$url"; then
            print_error "Failed to download from $url"
            exit 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget -q -O "$download_path" "$url"; then
            print_error "Failed to download from $url"
            exit 1
        fi
    fi

    print_status "Extracting binary..."

    # Extract the binary
    cd "$temp_dir"

    if [ "$PLATFORM" = "Windows" ]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "$filename"
        else
            print_error "unzip is required to extract Windows archives"
            exit 1
        fi
    else
        if command -v tar >/dev/null 2>&1; then
            tar -xzf "$filename"
        else
            print_error "tar is required to extract archives"
            exit 1
        fi
    fi

    # Find the binary
    local binary_path=""
    if [ "$PLATFORM" = "Windows" ]; then
        binary_path="${BINARY_NAME}.exe"
    else
        binary_path="${BINARY_NAME}"
    fi

    if [ ! -f "$binary_path" ]; then
        print_error "Binary not found in extracted archive"
        exit 1
    fi

    print_status "Installing to ${INSTALL_DIR}..."

    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        if ! mkdir -p "$INSTALL_DIR" 2>/dev/null; then
            print_warning "Cannot create $INSTALL_DIR. Trying with sudo..."
            sudo mkdir -p "$INSTALL_DIR"
        fi
    fi

    # Install binary
    local install_path="${INSTALL_DIR}/${BINARY_NAME}"
    if [ "$PLATFORM" = "Windows" ]; then
        install_path="${INSTALL_DIR}/${BINARY_NAME}.exe"
    fi

    if ! cp "$binary_path" "$install_path" 2>/dev/null; then
        print_warning "Cannot write to $INSTALL_DIR. Trying with sudo..."
        sudo cp "$binary_path" "$install_path"
    fi

    # Make executable
    if [ "$PLATFORM" != "Windows" ]; then
        if ! chmod +x "$install_path" 2>/dev/null; then
            sudo chmod +x "$install_path"
        fi
    fi

    print_success "PipeOps CLI installed successfully!"
}

# Verify installation
verify_installation() {
    print_status "Verifying installation..."

    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local installed_version
        installed_version=$("$BINARY_NAME" version 2>/dev/null | head -n1 || echo "unknown")
        print_success "Installation verified: $installed_version"
        return 0
    else
        print_warning "Binary installed but not found in PATH."
        print_warning "You may need to add ${INSTALL_DIR} to your PATH."
        print_warning "Add this to your shell profile (.bashrc, .zshrc, etc.):"
        echo "    export PATH=\"${INSTALL_DIR}:\$PATH\""
        return 1
    fi
}

# Print post-installation information
print_completion() {
    echo ""
    echo -e "${GREEN}🎉 Installation Complete!${NC}"
    echo ""
    echo -e "${CYAN}Getting Started:${NC}"
    echo "  1. Authenticate with PipeOps:"
    echo "     ${BINARY_NAME} auth login"
    echo ""
    echo "  2. List your projects:"
    echo "     ${BINARY_NAME} project list"
    echo ""
    echo "  3. Get help:"
    echo "     ${BINARY_NAME} --help"
    echo ""
    echo -e "${CYAN}Documentation:${NC}"
    echo "  🌐 Website: https://pipeops.io"
    echo "  📖 Docs: https://docs.pipeops.io"
    echo "  💬 Discord: https://discord.gg/pipeops"
    echo "  🐙 GitHub: https://github.com/${REPO}"
    echo ""
}

# Main installation function
main() {
    print_header

    # Check if running as root
    if [ "$EUID" -eq 0 ]; then
        print_warning "Running as root. Consider running as a regular user."
    fi

    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --install-dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            --version)
                VERSION="$2"
                shift 2
                ;;
            -h|--help)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --install-dir DIR    Install to specific directory (default: /usr/local/bin)"
                echo "  --version VERSION    Install specific version (default: latest)"
                echo "  -h, --help          Show this help message"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done

    # Check for required tools
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        print_error "Either curl or wget is required for installation"
        exit 1
    fi

    # Detect platform
    detect_platform

    # Get latest version if not specified
    if [ -z "$VERSION" ]; then
        get_latest_version
    fi

    # Install binary
    install_binary

    # Verify installation
    verify_installation

    # Print completion message
    print_completion
}

# Run main function
main "$@"