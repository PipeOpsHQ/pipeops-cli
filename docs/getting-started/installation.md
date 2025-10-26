# Installation

This guide covers various methods to install PipeOps CLI on different platforms.

## Quick Install

### macOS & Linux (Recommended)

The easiest way to install PipeOps CLI is using our installation script:

```bash
curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.sh | sh
```

This script will:
- Detect your operating system and architecture
- Download the appropriate binary
- Install it to `/usr/local/bin/pipeops`
- Verify the installation

### Windows (PowerShell)

For Windows users, use our PowerShell installation script:

```powershell
irm https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.ps1 | iex
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
# x86_64
wget https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops-cli_Linux_x86_64.tar.gz
tar -xzf pipeops-cli_Linux_x86_64.tar.gz
sudo mv pipeops /usr/local/bin/

# ARM64
wget https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops-cli_Linux_arm64.tar.gz
tar -xzf pipeops-cli_Linux_arm64.tar.gz
sudo mv pipeops /usr/local/bin/
```

#### macOS
```bash
# Intel Mac
wget https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops-cli_Darwin_x86_64.tar.gz
tar -xzf pipeops-cli_Darwin_x86_64.tar.gz
sudo mv pipeops /usr/local/bin/

# Apple Silicon (M1/M2)
wget https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops-cli_Darwin_arm64.tar.gz
tar -xzf pipeops-cli_Darwin_arm64.tar.gz
sudo mv pipeops /usr/local/bin/
```

#### Windows
```powershell
# Download the Windows binary
Invoke-WebRequest -Uri "https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops-cli_Windows_x86_64.zip" -OutFile "pipeops-cli.zip"

# Extract and add to PATH
Expand-Archive -Path "pipeops-cli.zip" -DestinationPath "C:\pipeops"
$env:PATH += ";C:\pipeops"
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

### Using Package Managers

#### Homebrew
```bash
brew update
brew upgrade pipeops
```

#### Go Install
```bash
go install github.com/PipeOpsHQ/pipeops-cli@latest
```

### Using CLI Update Command

PipeOps CLI can update itself:

```bash
# Check for updates
pipeops update check

# Update to latest version
pipeops update
```

### Manual Update

Download the latest release and replace the existing binary:

```bash
# Download latest version
curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.sh | sh
```

## Uninstalling

### macOS & Linux

```bash
# Remove binary
sudo rm /usr/local/bin/pipeops

# Remove configuration (optional)
rm ~/.pipeops.json
```

### Windows

```powershell
# Remove from PATH and delete files
Remove-Item "C:\pipeops\pipeops.exe" -Force
# Remove from PATH environment variable manually
```

### Homebrew

```bash
brew uninstall pipeops
brew untap pipeops/pipeops
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
curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.sh | sh
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
