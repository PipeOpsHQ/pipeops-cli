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

### macOS & Linux (recommended)

    curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.sh | sh

### Windows (PowerShell)

    irm https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.ps1 | iex

### Package managers

#### Homebrew (macOS/Linux)

    brew tap pipeops/pipeops
    brew install pipeops

#### Docker

    docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:latest --help

#### Go install

    go install github.com/PipeOpsHQ/pipeops-cli@latest

For more options and troubleshooting, see our [Installation Guide](docs/getting-started/installation.md)

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

       pipeops auth status

3. **List your projects**

       pipeops project list

4. **Get help for any command**

       pipeops --help
       pipeops auth --help
       pipeops project --help

---

## Commands Overview

| Command           | Description                            | Examples                                           |
| ----------------- | -------------------------------------- | -------------------------------------------------- |
| `pipeops auth`    | Manage authentication and user details | `pipeops auth login`, `pipeops auth status`        |
| `pipeops project` | Manage, list, and deploy projects      | `pipeops project list`, `pipeops project create`   |
| `pipeops deploy`  | Manage and deploy CI/CD pipelines      | `pipeops deploy pipeline`, `pipeops deploy status` |
| `pipeops server`  | Manage server-related operations       | `pipeops server list`                              |
| `pipeops agent`   | Install and manage PipeOps agents      | `pipeops agent install`, `pipeops agent status`    |

### Global Flags

- `--help, -h`: Show help information
- `--version, -v`: Display version information
- `--json`: Output results in JSON format
- `--quiet, -q`: Reduce non-essential output

---

## Configuration

Configuration is stored at `~/.pipeops.json` and includes:

- **Authentication tokens**: Secure OAuth credentials
- **User preferences**: Default settings and customizations
- **API settings**: Endpoint configurations

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PIPEOPS_CONFIG_PATH` | Custom config file location | `~/.pipeops.json` |
| `PIPEOPS_API_URL` | Override API endpoint | Default API URL |
| `PIPEOPS_LOG_LEVEL` | Set logging level | `info` |

Supported log levels: `debug`, `info`, `warn`, `error`

---

## Development

### Prerequisites

- Go 1.23+
- Git

### Setup

    # Clone
    git clone https://github.com/PipeOpsHQ/pipeops-cli.git
    cd pipeops-cli

    # Dependencies
    go mod download

    # Build
    make build

    # Test
    make test

    # Lint
    make lint

### Make targets

    make build          # Build the binary
    make test           # Run tests
    make lint           # Run linter
    make clean          # Clean build artifacts
    make install        # Install locally
    make release        # Create release build
    make docker-build   # Build Docker image
    make docker-run     # Run in Docker

### Project structure

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
    ├── Dockerfile          # Docker image
    └── install.sh          # Installation script

---

## Docker Usage

### Run CLI in Docker

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

### Docker Compose

    version: "3.8"
    services:
      pipeops-cli:
        image: ghcr.io/pipeopshq/pipeops-cli:latest
        volumes:
          - ~/.pipeops.json:/root/.pipeops.json
        command: ["project", "list"]

---

## Platforms

| Platform | Architecture   | Status |
| -------- | -------------- | ------ |
| Linux    | x86_64         | Supported |
| Linux    | ARM64          | Supported |
| Linux    | ARM            | Supported |
| macOS    | x86_64 (Intel) | Supported |
| macOS    | ARM64 (M1/M2)  | Supported |
| Windows  | x86_64         | Supported |
| FreeBSD  | x86_64         | Supported |

---

## Contributing

We welcome contributions!

1. Fork the repository and create a feature branch
2. Follow coding standards and add tests
3. Validate with `make test` and `make lint`
4. Open a PR with a clear description

### Guidelines

- **Code Quality**: Follow Go best practices and conventions
- **Documentation**: Write clear, comprehensive documentation
- **Testing**: Include tests for new functionality
- **Compatibility**: Ensure changes work across supported platforms
- **Collaboration**: Be respectful and constructive in discussions

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

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

---

Made with love by the PipeOps team
