# PipeOps CLI 🚀

[![Release](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/release.yml/badge.svg)](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/release.yml)
[![CodeQL Analysis](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/code-analysis.yml/badge.svg)](https://github.com/PipeOpsHQ/pipeops-cli/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![GitHub Release](https://img.shields.io/github/release/PipeOpsHQ/pipeops-cli.svg)](https://github.com/PipeOpsHQ/pipeops-cli/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/pipeops/pipeops-cli.svg)](https://hub.docker.com/r/pipeops/pipeops-cli)

PipeOps CLI is a powerful command-line interface designed to simplify managing cloud-native environments, deploying projects, and interacting with the PipeOps platform. With PipeOps CLI, you can provision servers, deploy applications, manage projects, and monitor your infrastructure seamlessly.

---

## 🚀 Quick Install

### macOS & Linux (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.ps1 | iex
```

### Package Managers

#### Homebrew (macOS/Linux)

```bash
brew tap pipeops/pipeops
brew install pipeops
```

#### Docker

```bash
docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:latest --help
```

#### Go Install

```bash
go install github.com/PipeOpsHQ/pipeops-cli@latest
```

**📋 [Complete Installation Guide](INSTALL.md)** - More installation methods and troubleshooting

---

## ✨ Features

- **🔐 OAuth Authentication**: Secure authentication with PKCE flow
- **📦 Project Management**: Create, manage, and deploy projects
- **🚀 Server Management**: Provision and configure servers across multiple environments
- **🔧 Pipeline Management**: Create, manage, and deploy CI/CD pipelines
- **🤖 Agent Setup**: Install and configure PipeOps agents for various platforms
- **🌐 Cross-Platform Support**: Available for Linux, Windows, macOS, and FreeBSD
- **📊 Status Monitoring**: Real-time status updates and monitoring
- **🎨 Beautiful UI**: Rich terminal interface with colors and progress indicators

---

## 🏃‍♂️ Quick Start

### 1. Authenticate with PipeOps

```bash
pipeops auth login
```

### 2. Check your authentication status

```bash
pipeops auth status
```

### 3. List your projects

```bash
pipeops project list
```

### 4. Get help for any command

```bash
pipeops --help
pipeops auth --help
pipeops project --help
```

---

## 📖 Commands Overview

| Command           | Description                               | Examples                                           |
| ----------------- | ----------------------------------------- | -------------------------------------------------- |
| `pipeops auth`    | Manage authentication and user details    | `pipeops auth login`, `pipeops auth status`        |
| `pipeops project` | Manage, list, and deploy PipeOps projects | `pipeops project list`, `pipeops project create`   |
| `pipeops deploy`  | Manage and deploy CI/CD pipelines         | `pipeops deploy pipeline`, `pipeops deploy status` |
| `pipeops server`  | Manage server-related operations          | `pipeops server deploy`, `pipeops server status`   |
| `pipeops k3s`     | Manage K3s clusters                       | `pipeops k3s install`, `pipeops k3s join`          |
| `pipeops agent`   | Manage PipeOps agents                     | `pipeops agent install`, `pipeops agent status`    |

### Global Flags

- `--help, -h`: Show help for any command
- `--version, -v`: Show version information
- `--json`: Output results in JSON format
- `--quiet, -q`: Suppress non-essential output

---

## 🔧 Configuration

PipeOps CLI stores configuration in `~/.pipeops.json`. This includes:

- Authentication tokens
- User preferences
- Default settings

### Environment Variables

- `PIPEOPS_CONFIG_PATH`: Custom config file location
- `PIPEOPS_API_URL`: Custom API endpoint
- `PIPEOPS_LOG_LEVEL`: Log level (debug, info, warn, error)

---

## 🛠️ Development

### Prerequisites

- [Go](https://golang.org/) 1.23 or later
- [Git](https://git-scm.com/)

### Setup

```bash
# Clone the repository
git clone https://github.com/PipeOpsHQ/pipeops-cli.git
cd pipeops-cli

# Install dependencies
go mod download

# Build the CLI
make build

# Run tests
make test

# Run linter
make lint
```

### Available Make Targets

```bash
make build          # Build the binary
make test           # Run tests
make lint           # Run linter
make clean          # Clean build artifacts
make install        # Install locally
make release        # Create release build
make docker-build   # Build Docker image
make docker-run     # Run in Docker
```

### Project Structure

```
pipeops-cli/
├── cmd/                 # CLI commands
│   ├── auth/           # Authentication commands
│   ├── project/        # Project management commands
│   ├── deploy/         # Deployment commands
│   └── ...
├── internal/           # Internal packages
│   ├── auth/           # Authentication logic
│   ├── client/         # HTTP client
│   ├── config/         # Configuration management
│   └── ...
├── models/             # Data models
├── utils/              # Utility functions
├── .goreleaser.yml     # Release configuration
├── Dockerfile         # Docker image
└── install.sh         # Installation script
```

---

## 🐳 Docker Usage

### Run CLI in Docker

```bash
# Basic usage
docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:latest --help

# With authentication (mount config)
docker run --rm -it \
  -v ~/.pipeops.json:/root/.pipeops.json \
  ghcr.io/pipeopshq/pipeops-cli:latest auth status

# Interactive shell
docker run --rm -it \
  -v ~/.pipeops.json:/root/.pipeops.json \
  --entrypoint /bin/sh \
  ghcr.io/pipeopshq/pipeops-cli:latest
```

### Docker Compose

```yaml
version: "3.8"
services:
  pipeops-cli:
    image: ghcr.io/pipeopshq/pipeops-cli:latest
    volumes:
      - ~/.pipeops.json:/root/.pipeops.json
    command: ["project", "list"]
```

---

## 🌐 Available Platforms

PipeOps CLI supports the following platforms:

| Platform | Architecture   | Status |
| -------- | -------------- | ------ |
| Linux    | x86_64         | ✅     |
| Linux    | ARM64          | ✅     |
| Linux    | ARM            | ✅     |
| macOS    | x86_64 (Intel) | ✅     |
| macOS    | ARM64 (M1/M2)  | ✅     |
| Windows  | x86_64         | ✅     |
| FreeBSD  | x86_64         | ✅     |

---

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork the repository** and create your feature branch
2. **Follow the coding standards** and write tests for new features
3. **Test your changes** with `make test` and `make lint`
4. **Submit a pull request** with a clear description

### Contribution Guidelines

- Follow Go best practices and conventions
- Write clear, commented code
- Include tests for new functionality
- Update documentation as needed
- Be respectful and collaborative

[📋 Detailed Contributing Guide](CONTRIBUTING.md)

---

## 📚 Documentation

- **[Installation Guide](INSTALL.md)** - Comprehensive installation instructions
- **[API Documentation](https://docs.pipeops.io)** - Complete API reference
- **[User Guide](https://docs.pipeops.io/cli)** - Detailed usage instructions
- **[Examples](examples/)** - Usage examples and scripts

---

## 🆘 Support & Community

- **📖 Documentation**: [docs.pipeops.io](https://docs.pipeops.io)
- **🐛 Issues**: [GitHub Issues](https://github.com/PipeOpsHQ/pipeops-cli/issues)
- **💬 Discussions**: [GitHub Discussions](https://github.com/PipeOpsHQ/pipeops-cli/discussions)
- **🗣️ Discord**: [Join our community](https://discord.gg/pipeops)
- **📧 Email**: [support@pipeops.io](mailto:support@pipeops.io)
- **🐦 Twitter**: [@PipeOpsHQ](https://twitter.com/pipeops)

---

## 🔄 Release Process

Releases are automated via GitHub Actions when tags are pushed:

1. **Create a new tag**: `git tag -a v1.0.0 -m "Release v1.0.0"`
2. **Push the tag**: `git push origin v1.0.0`
3. **GitHub Actions** will automatically:
   - Build binaries for all platforms
   - Create GitHub release with binaries
   - Push Docker images to registry
   - Update package managers (Homebrew, AUR, etc.)

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

## 🙏 Acknowledgments

Special thanks to:

- All contributors and users of PipeOps CLI
- The Go community for excellent tools and libraries
- GitHub for providing CI/CD infrastructure
- The open-source community for inspiration and support

---

_Made with ❤️ by the PipeOps team_
