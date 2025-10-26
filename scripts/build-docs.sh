#!/bin/bash

# PipeOps CLI Documentation Build Script
# This script builds the MkDocs documentation locally

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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

# Check if Python is installed
check_python() {
    if ! command -v python3 &> /dev/null; then
        print_error "Python 3 is required but not installed."
        exit 1
    fi

    python_version=$(python3 --version 2>&1 | cut -d' ' -f2)
    print_status "Found Python $python_version"
}

# Check if pip is installed
check_pip() {
    if ! command -v pip3 &> /dev/null; then
        print_error "pip3 is required but not installed."
        exit 1
    fi
    print_status "Found pip3"
}

# Install dependencies
install_dependencies() {
    print_status "Installing MkDocs dependencies..."

    pip3 install --upgrade pip
    pip3 install mkdocs-material
    pip3 install mkdocs-git-revision-date-localized-plugin
    pip3 install mkdocs-minify-plugin

    print_success "Dependencies installed successfully"
}

# Get the correct Python path
get_python_path() {
    if command -v python3 &> /dev/null; then
        python3 -c "import sys; print(sys.executable)"
    else
        echo "python3"
    fi
}

# Get the correct pip path
get_pip_path() {
    if command -v pip3 &> /dev/null; then
        pip3 -c "import sys; print(sys.executable)" 2>/dev/null || echo "pip3"
    else
        echo "pip3"
    fi
}

# Build documentation
build_docs() {
    print_status "Building documentation..."

    if [ ! -f "mkdocs.yml" ]; then
        print_error "mkdocs.yml not found. Are you in the project root?"
        exit 1
    fi

    # Try to find mkdocs in the user's local bin
    if [ -f "$HOME/Library/Python/3.9/bin/mkdocs" ]; then
        "$HOME/Library/Python/3.9/bin/mkdocs" build
    elif [ -f "$HOME/.local/bin/mkdocs" ]; then
        "$HOME/.local/bin/mkdocs" build
    elif command -v mkdocs &> /dev/null; then
        mkdocs build
    else
        print_error "mkdocs command not found. Please ensure it's installed and in your PATH."
        exit 1
    fi

    if [ $? -eq 0 ]; then
        print_success "Documentation built successfully"
        print_status "Site generated in ./site directory"
    else
        print_error "Documentation build failed"
        exit 1
    fi
}

# Serve documentation locally
serve_docs() {
    print_status "Starting local documentation server..."
    print_status "Documentation will be available at: http://127.0.0.1:8000"
    print_status "Press Ctrl+C to stop the server"

    # Try to find mkdocs in the user's local bin
    if [ -f "$HOME/Library/Python/3.9/bin/mkdocs" ]; then
        "$HOME/Library/Python/3.9/bin/mkdocs" serve
    elif [ -f "$HOME/.local/bin/mkdocs" ]; then
        "$HOME/.local/bin/mkdocs" serve
    elif command -v mkdocs &> /dev/null; then
        mkdocs serve
    else
        print_error "mkdocs command not found. Please ensure it's installed and in your PATH."
        exit 1
    fi
}

# Clean build artifacts
clean_docs() {
    print_status "Cleaning build artifacts..."

    if [ -d "site" ]; then
        rm -rf site
        print_success "Build artifacts cleaned"
    else
        print_warning "No build artifacts to clean"
    fi
}

# Validate documentation
validate_docs() {
    print_status "Validating documentation..."

    # Check for broken links
    if command -v linkchecker &> /dev/null; then
        print_status "Checking for broken links..."
        linkchecker http://127.0.0.1:8000 --check-extern
    else
        print_warning "linkchecker not installed. Skipping link validation."
        print_status "Install with: pip3 install linkchecker"
    fi

    # Check markdown syntax
    if command -v markdownlint &> /dev/null; then
        print_status "Checking markdown syntax..."
        markdownlint docs/
    else
        print_warning "markdownlint not installed. Skipping markdown validation."
        print_status "Install with: npm install -g markdownlint-cli"
    fi
}

# Show help
show_help() {
    echo "PipeOps CLI Documentation Build Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  build     Build the documentation"
    echo "  serve     Build and serve documentation locally"
    echo "  clean     Clean build artifacts"
    echo "  install   Install dependencies"
    echo "  validate  Validate documentation"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 build          # Build documentation"
    echo "  $0 serve          # Serve documentation locally"
    echo "  $0 clean          # Clean build artifacts"
}

# Main script logic
main() {
    case "${1:-build}" in
        "build")
            check_python
            check_pip
            install_dependencies
            build_docs
            ;;
        "serve")
            check_python
            check_pip
            install_dependencies
            serve_docs
            ;;
        "clean")
            clean_docs
            ;;
        "install")
            check_python
            check_pip
            install_dependencies
            ;;
        "validate")
            validate_docs
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
