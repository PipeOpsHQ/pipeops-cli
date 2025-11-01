# PipeOps CLI

[![Release](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/release.yml/badge.svg)](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/release.yml)
[![CodeQL Analysis](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/code-analysis.yml/badge.svg)](https://github.com/PipeOpsHQ/pipeops-cli/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![GitHub Release](https://img.shields.io/github/release/PipeOpsHQ/pipeops-cli.svg)](https://github.com/PipeOpsHQ/pipeops-cli/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/pipeops/pipeops-cli.svg)](https://hub.docker.com/r/pipeops/pipeops-cli)

PipeOps CLI is a powerful, cross-platform command-line tool that streamlines cloud-native development and deployment workflows. Securely authenticate with OAuth, manage projects and servers, deploy CI/CD pipelines, and monitor infrastructure—all from your terminal with a developer-friendly interface.

## Why PipeOps CLI?

- **Unified Workflow**: Manage your entire development lifecycle from a single CLI
- **Developer-First**: Designed for productivity with intuitive commands and rich output  
- **Secure by Default**: OAuth PKCE authentication with secure credential management
- **Platform Agnostic**: Works consistently across Linux, macOS, Windows, and FreeBSD
- **CI/CD Ready**: Perfect for automation scripts and continuous integration pipelines
- **Real-time Feedback**: Live status updates and progress indicators for long-running operations

---

## Table of Contents

- [PipeOps CLI](#pipeops-cli)
  - [Why PipeOps CLI?](#why-pipeops-cli)
  - [Table of Contents](#table-of-contents)
  - [Quick Install](#quick-install)
    - [macOS \& Linux (recommended)](#macos--linux-recommended)
    - [Windows (PowerShell)](#windows-powershell)
    - [Package managers](#package-managers)
      - [Homebrew (macOS/Linux)](#homebrew-macoslinux)
      - [Docker](#docker)
      - [Go install](#go-install)
  - [Features](#features)
  - [Quick Start](#quick-start)
  - [Commands Overview](#commands-overview)
    - [Global Flags](#global-flags)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
  - [Development](#development)
    - [Prerequisites](#prerequisites)
    - [Setup](#setup)
    - [Make targets](#make-targets)
    - [Project structure](#project-structure)
  - [Docker Usage](#docker-usage)
    - [Run CLI in Docker](#run-cli-in-docker)
    - [Docker Compose](#docker-compose)
  - [Platforms](#platforms)
  - [Contributing](#contributing)
    - [Guidelines](#guidelines)
  - [Documentation](#documentation)
  - [Support \& Community](#support--community)
  - [Release Process](#release-process)
  - [License](#license)

---

## Quick Install

### macOS & Linux (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.ps1 | iex
```

    irm https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.ps1 | iex

### Package managers

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

    go install github.com/PipeOpsHQ/pipeops-cli@latest

📋 For more options and troubleshooting, see our [Installation Guide](docs/getting-started/installation.md)

---

## Features

- **Secure Authentication**: OAuth with PKCE for secure, device-friendly login
- **Project Management**: Complete lifecycle management - create, list, configure, and deploy projects
- **Server Operations**: Server provisioning and environment management
- **CI/CD Integration**: Pipeline management and deployment automation
- **Agent Management**: Install and manage PipeOps agents across platforms
- **Cross-Platform**: Native support for Linux, macOS, Windows, and FreeBSD
- **Developer Experience**: Real-time status updates, rich output formatting, and JSON mode
- **Terminal UX**: Clean interface with progress indicators and color-coded output

---

## Quick Start

1. **Log in to your PipeOps account**

       pipeops auth login

2. **Verify authentication status**

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

## Commands Overview

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
- `--quiet, -q`: Reduce non-essential output

---

## Configuration

PipeOps CLI stores configuration in `~/.pipeops.json`. This includes:

- Authentication tokens
- User preferences
- Default settings

### Environment Variables

- `PIPEOPS_CONFIG_PATH`: Custom config file location
- `PIPEOPS_API_URL`: Custom API endpoint
- `PIPEOPS_LOG_LEVEL`: Log level (debug, info, warn, error)

---

## Development

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

## Docker Usage

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

## Platforms

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

## Contributing

We welcome contributions!

1. Fork the repository and create a feature branch
2. Follow coding standards and add tests
3. Validate with `make test` and `make lint`
4. Open a PR with a clear description

### Guidelines

### Contribution Guidelines

- Follow Go best practices and conventions
- Write clear, commented code
- Include tests for new functionality
- Update documentation as needed
- Be respectful and collaborative

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

---

## Documentation

- **[Installation Guide](docs/getting-started/installation.md)** - Complete installation instructions
- **[Quick Start Guide](docs/getting-started/quick-start.md)** - Get up and running quickly
- **[Commands Reference](docs/commands/overview.md)** - Detailed command documentation
- **[API Reference](https://docs.pipeops.io)** - Complete API documentation
- **[CLI User Guide](https://docs.pipeops.io/cli)** - In-depth CLI usage guide

---

## Support & Community

- **[Documentation](https://docs.pipeops.io)** - Complete guides and references
- **[GitHub Issues](https://github.com/PipeOpsHQ/pipeops-cli/issues)** - Bug reports and feature requests
- **[GitHub Discussions](https://github.com/PipeOpsHQ/pipeops-cli/discussions)** - Community discussions
- **[Discord](https://discord.gg/pipeops)** - Real-time community chat
- **[Email Support](mailto:support@pipeops.io)** - Direct technical support
- **[Twitter](https://twitter.com/PipeOpsHQ)** - Latest news and updates

---

## Release Process

Releases are automated via GitHub Actions when tags are pushed:

1. **Create a release tag**

       git tag -a v1.0.0 -m "Release v1.0.0"

2. **Push the tag**

       git push origin v1.0.0

3. **Automated CI process**
   - Build binaries for all supported platforms
   - Create GitHub release with artifacts and changelog
   - Push Docker images to registry
   - Update package managers (Homebrew, etc.)

---

## License

Special thanks to:

- All contributors and users of PipeOps CLI
- The Go community for excellent tools and libraries
- GitHub for providing CI/CD infrastructure
- The open-source community for inspiration and support

---

_Made with ❤️ by the PipeOps team_
