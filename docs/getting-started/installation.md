# Installation

This guide covers various methods to install PipeOps CLI on different platforms.

## Quick Install

### macOS & Linux (Recommended)

The easiest way to install PipeOps CLI is using our installation script:

```bash
curl -fsSL https://get.pipeops.dev/cli.sh | bash
```

This script will:
- Detect your operating system and architecture
- Download the appropriate binary
- Install it to `/usr/local/bin/pipeops`
- Verify the installation

### Windows (PowerShell)

For Windows users, use our PowerShell installation script:

```powershell
irm https://get.pipeops.dev/cli.ps1 | iex
```

## Package Managers

### Homebrew (macOS/Linux)

Install using Homebrew:

```bash
# Add the PipeOps tap
brew tap pipeops/pipeops

# Install PipeOps CLI
brew install pipeops
```

### Go Install

If you have Go installed, you can install directly from source:

```bash
go install github.com/PipeOpsHQ/pipeops-cli@latest
```

### Docker

Run PipeOps CLI in a Docker container:

```bash
# Basic usage
docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:latest --help

# With persistent configuration
docker run --rm -it \
  -v ~/.pipeops.json:/root/.pipeops.json \
  ghcr.io/pipeopshq/pipeops-cli:latest auth status
```

## Manual Installation

### Download Binaries

Download the appropriate binary for your platform from the [releases page](https://github.com/PipeOpsHQ/pipeops-cli/releases):

#### Linux
```bash
# Install using installer domain (auto-detects arch)
curl -fsSL https://get.pipeops.dev/cli.sh | bash

# Pin to a specific version
VERSION=v1.2.0 curl -fsSL https://get.pipeops.dev/cli.sh | bash
```

#### macOS
```bash
# Install using installer domain (auto-detects arch)
curl -fsSL https://get.pipeops.dev/cli.sh | bash

# Pin to a specific version
VERSION=v1.2.0 curl -fsSL https://get.pipeops.dev/cli.sh | bash
```

#### Windows
```powershell
# Install using installer domain
irm https://get.pipeops.dev/cli.ps1 | iex

# Pin to a specific version
$env:VERSION = 'v1.2.0'; irm https://get.pipeops.dev/cli.ps1 | iex
```

## Verify Installation

After installation, verify that PipeOps CLI is working correctly:

```bash
# Check version
pipeops --version

# Check help
pipeops --help

# Check authentication status
pipeops auth status
```

## Updating

Keep your PipeOps CLI up to date to access the latest features and security improvements.

### Automatic Updates (Recommended)

PipeOps CLI can update itself automatically:

```bash
# Check for available updates
pipeops update check

# Update to the latest version
pipeops update

# Update to a specific version
pipeops update --version v1.2.0

# View update history
pipeops update history
```

The self-update feature:
- ✅ Automatically detects your platform
- ✅ Downloads the correct binary
- ✅ Preserves your configuration
- ✅ Creates a backup of the current version
- ✅ Verifies the update integrity

### Package Manager Updates

#### Homebrew (macOS/Linux)
```bash
# Update package lists
brew update

# Upgrade PipeOps CLI
brew upgrade pipeops

# Upgrade to a specific version
brew install pipeops@1.2.0
```

#### Go Install
```bash
# Update to latest version
go install github.com/PipeOpsHQ/pipeops-cli@latest

# Install specific version
go install github.com/PipeOpsHQ/pipeops-cli@v1.2.0
```

### Manual Update

For manual updates, download and replace the binary:

```bash
# Method 1: Re-run installer (overwrites existing installation)
curl -fsSL https://get.pipeops.dev/cli.sh | bash

# Method 2: Install specific version
VERSION=v1.2.0 curl -fsSL https://get.pipeops.dev/cli.sh | bash
```

### Update Configuration

Control update behavior through configuration:

```bash
# Disable automatic update checks
pipeops config set skip_update_check true

# Set update channel
pipeops config set update_channel stable  # stable, beta, alpha

# Set update check interval
pipeops config set update_check_interval 24h
```

### Verify Update

After updating, verify the new version:

```bash
# Check version
pipeops --version

# Verify functionality
pipeops auth status
pipeops --help
```

## Uninstalling

Remove PipeOps CLI completely from your system.

### Using Package Managers

#### Homebrew (macOS/Linux)
```bash
# Uninstall PipeOps CLI
brew uninstall pipeops

# Remove tap (optional)
brew untap pipeops/pipeops

# Clean up any remaining files
brew cleanup
```

#### Docker
```bash
# Remove Docker images
docker rmi ghcr.io/pipeopshq/pipeops-cli:latest

# Remove all PipeOps CLI images
docker images | grep pipeops-cli | awk '{print $3}' | xargs docker rmi
```

### Manual Uninstall

#### macOS & Linux

```bash
# Remove binary
sudo rm -f /usr/local/bin/pipeops

# Alternative locations
sudo rm -f /usr/bin/pipeops
rm -f ~/.local/bin/pipeops

# Remove configuration files
rm -f ~/.pipeops.json
rm -rf ~/.pipeops/

# Remove from shell profile (if added manually)
# Edit ~/.bashrc, ~/.zshrc, etc. and remove PipeOps PATH entries
```

#### Windows

```powershell
# Remove executable
Remove-Item "$env:USERPROFILE\bin\pipeops.exe" -Force -ErrorAction SilentlyContinue
Remove-Item "C:\pipeops\pipeops.exe" -Force -ErrorAction SilentlyContinue

# Remove configuration
Remove-Item "$env:USERPROFILE\.pipeops.json" -Force -ErrorAction SilentlyContinue

# Remove from PATH
$currentPath = [Environment]::GetEnvironmentVariable("PATH", [EnvironmentVariableTarget]::User)
$newPath = ($currentPath.Split(';') | Where-Object { $_ -notlike "*pipeops*" }) -join ';'
[Environment]::SetEnvironmentVariable("PATH", $newPath, [EnvironmentVariableTarget]::User)
```

### Cleanup Data and Logs

Remove any remaining data:

```bash
# Remove cache files (Linux/macOS)
rm -rf ~/.cache/pipeops/
rm -rf ~/Library/Caches/pipeops/  # macOS only

# Remove log files
rm -rf ~/.local/share/pipeops/
rm -rf ~/Library/Logs/pipeops/    # macOS only

# Windows cache cleanup
# Remove-Item "$env:LOCALAPPDATA\pipeops" -Recurse -Force -ErrorAction SilentlyContinue
```

### Verify Uninstall

Confirm complete removal:

```bash
# Check if binary is removed
which pipeops
# Should return: pipeops not found

# Check for remaining files
ls -la ~/.*pipeops*
ls -la ~/.local/bin/pipeops*
ls -la /usr/local/bin/pipeops*

# Check environment variables
env | grep PIPEOPS
```

### Troubleshooting Uninstall

If you encounter issues during uninstall:

#### Permission Denied
```bash
# Use sudo for system-wide installations
sudo rm -f /usr/local/bin/pipeops

# Change ownership if needed
sudo chown $(whoami) /usr/local/bin/pipeops
rm /usr/local/bin/pipeops
```

#### Binary in Use
```bash
# Check if PipeOps CLI is running
ps aux | grep pipeops

# Kill any running processes
pkill -f pipeops

# Try uninstall again
sudo rm -f /usr/local/bin/pipeops
```

#### Multiple Installations
```bash
# Find all installations
find /usr -name "pipeops" 2>/dev/null
find ~ -name "pipeops" 2>/dev/null

# Remove all found instances
sudo find /usr -name "pipeops" -delete 2>/dev/null
find ~ -name "pipeops" -delete 2>/dev/null
```

## Troubleshooting

### Permission Issues

If you encounter permission issues during installation:

```bash
# Make sure the install directory is writable
sudo chown -R $(whoami) /usr/local/bin

# Or install to user directory
mkdir -p ~/.local/bin
mv pipeops ~/.local/bin
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Binary Not Found

If the `pipeops` command is not found:

1. **Check if it's in your PATH**:
   ```bash
   which pipeops
   echo $PATH
   ```

2. **Add to PATH** (if needed):
   ```bash
   echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   ```

### Network Issues

If you're behind a corporate firewall or proxy:

```bash
# Set proxy environment variables
export HTTP_PROXY=http://proxy.company.com:8080
export HTTPS_PROXY=http://proxy.company.com:8080

# Then run installation
curl -fsSL https://get.pipeops.dev/cli.sh | bash
```

## System Requirements

### Minimum Requirements

- **Operating System**: Linux, macOS, Windows, or FreeBSD
- **Architecture**: x86_64, ARM64, or ARM
- **Memory**: 50MB available RAM
- **Disk Space**: 20MB for binary and configuration

### Recommended Requirements

- **Operating System**: Latest LTS version
- **Memory**: 100MB+ available RAM
- **Network**: Stable internet connection for authentication and API calls

## Security Considerations

- Binaries are signed and checksums are provided
- Always verify checksums before installation
- Use official installation methods when possible
- Keep your CLI updated to the latest version

## Getting Help

If you encounter issues during installation:

- **Check the [troubleshooting guide](advanced/troubleshooting.md)**
- **Open an [issue on GitHub](https://github.com/PipeOpsHQ/pipeops-cli/issues)**
- **Join our [Discord community](https://discord.gg/pipeops)**
- **Email us at [support@pipeops.io](mailto:support@pipeops.io)**
