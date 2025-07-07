# 📦 PipeOps CLI Installation Guide

Multiple ways to install PipeOps CLI on your system.

## 🚀 Quick Install (Recommended)

### macOS & Linux
```bash
# One-line installer
curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.sh | sh
```

### Windows (PowerShell)
```powershell
# One-line installer
irm https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.ps1 | iex
```

## 📋 Installation Methods

### 1. Package Managers

#### Homebrew (macOS/Linux)
```bash
# Add the PipeOps tap
brew tap pipeops/pipeops

# Install PipeOps CLI
brew install pipeops
```

#### APT (Ubuntu/Debian)
```bash
# Download and install .deb package
curl -LO https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops_Linux_x86_64.deb
sudo dpkg -i pipeops_Linux_x86_64.deb
```

#### YUM/DNF (RHEL/CentOS/Fedora)
```bash
# Download and install .rpm package
curl -LO https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops_Linux_x86_64.rpm
sudo rpm -i pipeops_Linux_x86_64.rpm
```

#### APK (Alpine Linux)
```bash
# Download and install .apk package
curl -LO https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops_Linux_x86_64.apk
sudo apk add --allow-untrusted pipeops_Linux_x86_64.apk
```

#### AUR (Arch Linux)
```bash
# Using yay
yay -S pipeops-cli

# Using makepkg
git clone https://aur.archlinux.org/pipeops-cli.git
cd pipeops-cli
makepkg -si
```

### 2. Go Install
```bash
# Install directly from source
go install github.com/PipeOpsHQ/pipeops-cli@latest
```

### 3. Docker

#### Run directly
```bash
# Run PipeOps CLI in Docker
docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:latest --help

# With authentication (mount config)
docker run --rm -it -v ~/.pipeops.json:/root/.pipeops.json ghcr.io/pipeopshq/pipeops-cli:latest auth status
```

#### Docker Compose
```yaml
version: '3.8'
services:
  pipeops-cli:
    image: ghcr.io/pipeopshq/pipeops-cli:latest
    volumes:
      - ~/.pipeops.json:/root/.pipeops.json
    command: ["--help"]
```

### 4. Manual Download

Download the appropriate binary for your platform from the [releases page](https://github.com/PipeOpsHQ/pipeops-cli/releases/latest):

#### Linux
```bash
# x86_64
curl -LO https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops_Linux_x86_64.tar.gz
tar -xzf pipeops_Linux_x86_64.tar.gz
sudo mv pipeops /usr/local/bin/

# ARM64
curl -LO https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops_Linux_arm64.tar.gz
tar -xzf pipeops_Linux_arm64.tar.gz
sudo mv pipeops /usr/local/bin/
```

#### macOS
```bash
# Intel Macs
curl -LO https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops_Darwin_x86_64.tar.gz
tar -xzf pipeops_Darwin_x86_64.tar.gz
sudo mv pipeops /usr/local/bin/

# Apple Silicon (M1/M2)
curl -LO https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops_Darwin_arm64.tar.gz
tar -xzf pipeops_Darwin_arm64.tar.gz
sudo mv pipeops /usr/local/bin/
```

#### Windows
```powershell
# Download and extract
Invoke-WebRequest -Uri "https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops_Windows_x86_64.zip" -OutFile "pipeops.zip"
Expand-Archive -Path "pipeops.zip" -DestinationPath "."
Move-Item "pipeops.exe" "C:\Windows\System32\"
```

## ✅ Verify Installation

After installation, verify that PipeOps CLI is working:

```bash
# Check version
pipeops version

# Show help
pipeops --help

# Check authentication status
pipeops auth status
```

## 🏃‍♂️ Getting Started

1. **Authenticate with PipeOps:**
   ```bash
   pipeops auth login
   ```

2. **List your projects:**
   ```bash
   pipeops project list
   ```

3. **Get help for any command:**
   ```bash
   pipeops --help
   pipeops auth --help
   pipeops project --help
   ```

## 🔧 Troubleshooting

### Command not found
If you get "command not found" after installation:

1. **Check PATH:** Make sure the installation directory is in your PATH
   ```bash
   echo $PATH
   ```

2. **Add to PATH:** Add the installation directory to your shell profile
   ```bash
   # For bash
   echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc

   # For zsh
   echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
   source ~/.zshrc
   ```

### Permission denied
If you get permission errors:

1. **Make executable:**
   ```bash
   chmod +x /usr/local/bin/pipeops
   ```

2. **Check ownership:**
   ```bash
   ls -la /usr/local/bin/pipeops
   ```

### Update to latest version
```bash
# Using package managers (they handle updates)
brew upgrade pipeops
sudo apt update && sudo apt upgrade pipeops

# Manual update - re-run the installer
curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.sh | sh
```

## 🏗️ Build from Source

If you prefer to build from source:

```bash
# Clone the repository
git clone https://github.com/PipeOpsHQ/pipeops-cli.git
cd pipeops-cli

# Build the binary
go build -o pipeops .

# Install to system
sudo mv pipeops /usr/local/bin/
```

## 📚 Platform-Specific Notes

### macOS
- **Gatekeeper:** You may need to allow the binary in Security & Privacy settings
- **Homebrew:** Recommended method for macOS users
- **Universal Binary:** Works on both Intel and Apple Silicon Macs

### Linux
- **Package Managers:** Use your distribution's package manager when available
- **Dependencies:** Most dependencies are statically linked
- **Permissions:** You may need sudo for system-wide installation

### Windows
- **PowerShell:** Requires PowerShell 3.0 or later
- **Antivirus:** Some antivirus software may flag the binary
- **PATH:** Windows installation script automatically updates PATH

### FreeBSD
```bash
# Manual installation
curl -LO https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops_Freebsd_x86_64.tar.gz
tar -xzf pipeops_Freebsd_x86_64.tar.gz
sudo mv pipeops /usr/local/bin/
```

## 🔍 Available Architectures

PipeOps CLI is available for the following platforms:

| OS | Architecture | Binary |
|---|---|---|
| Linux | x86_64 | ✅ |
| Linux | ARM64 | ✅ |
| Linux | ARM | ✅ |
| Linux | i386 | ✅ |
| macOS | x86_64 | ✅ |
| macOS | ARM64 (M1/M2) | ✅ |
| Windows | x86_64 | ✅ |
| Windows | i386 | ✅ |
| FreeBSD | x86_64 | ✅ |
| FreeBSD | ARM64 | ✅ |

## 🆘 Support

- **Documentation:** [docs.pipeops.io](https://docs.pipeops.io)
- **Issues:** [GitHub Issues](https://github.com/PipeOpsHQ/pipeops-cli/issues)
- **Discussions:** [GitHub Discussions](https://github.com/PipeOpsHQ/pipeops-cli/discussions)
- **Discord:** [Join our community](https://discord.gg/pipeops)
- **Email:** [support@pipeops.io](mailto:support@pipeops.io)

---

*Made with ❤️ by the PipeOps team*