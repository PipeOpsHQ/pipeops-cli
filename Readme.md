# PipeOps CLI

[![Release](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/release.yml/badge.svg)](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/release.yml)
[![CodeQL Analysis](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/code-analysis.yml/badge.svg)](https://github.com/PipeOpsHQ/pipeops-cli/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![GitHub Release](https://img.shields.io/github/release/PipeOpsHQ/pipeops-cli.svg)](https://github.com/PipeOpsHQ/pipeops-cli/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/pipeops/pipeops-cli.svg)](https://hub.docker.com/r/pipeops/pipeops-cli)

PipeOps CLI is a fast, cross-platform command-line tool for managing cloud-native projects with PipeOps. Authenticate via OAuth (PKCE), provision and manage servers, deploy pipelines, interact with agents, and monitor status—all from your terminal.

---

## Table of Contents

- 🚀 Quick Install
- ✨ Features
- 🏃 Quick Start
- 📖 Commands Overview
- 🔧 Configuration
- 🛠️ Development
- 🐳 Docker Usage
- 🌐 Platforms
- 🤝 Contributing
- 📚 Documentation
- 🆘 Support & Community
- 🔄 Release Process
- 📄 License
- 🙏 Acknowledgments

---

## 🚀 Quick Install

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

📋 For more options and troubleshooting, see the Complete Installation Guide: INSTALL.md

---

## ✨ Features

- 🔐 OAuth with PKCE for secure, device-friendly login
- 📦 Project lifecycle management: create, list, configure, deploy
- 🚀 Server provisioning and environment management
- 🔧 CI/CD pipeline management and deployment
- 🤖 Agent installation and management across platforms
- 🌐 Cross-platform support: Linux, macOS, Windows, FreeBSD
- 📊 Real-time status, rich output, and JSON mode
- 🎨 Pleasant terminal UX with progress and color

---

## 🏃 Quick Start

1. Log in

   pipeops auth login

2. Check authentication status

   pipeops auth status

3. List projects

   pipeops project list

4. Get help

   pipeops --help
   pipeops auth --help
   pipeops project --help

---

## 📖 Commands Overview

| Command           | Description                            | Examples                                           |
| ----------------- | -------------------------------------- | -------------------------------------------------- |
| `pipeops auth`    | Manage authentication and user details | `pipeops auth login`, `pipeops auth status`        |
| `pipeops project` | Manage, list, and deploy projects      | `pipeops project list`, `pipeops project create`   |
| `pipeops deploy`  | Manage and deploy CI/CD pipelines      | `pipeops deploy pipeline`, `pipeops deploy status` |
| `pipeops server`  | Manage server-related operations       | `pipeops server deploy`, `pipeops server status`   |
| `pipeops k3s`     | Manage K3s clusters                    | `pipeops k3s install`, `pipeops k3s join`          |
| `pipeops agent`   | Install and manage PipeOps agents      | `pipeops agent install`, `pipeops agent status`    |

### Global flags

- `--help, -h`: Show help
- `--version, -v`: Show version
- `--json`: Output JSON
- `--quiet, -q`: Reduce non-essential output

---

## 🔧 Configuration

By default, configuration is stored at `~/.pipeops.json`, including:

- Authentication tokens
- User preferences
- Default settings

### Environment variables

- `PIPEOPS_CONFIG_PATH`: Custom config file location
- `PIPEOPS_API_URL`: Override API endpoint
- `PIPEOPS_LOG_LEVEL`: `debug`, `info`, `warn`, `error`

---

## 🛠️ Development

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

## 🐳 Docker Usage

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

## 🌐 Platforms

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

We welcome contributions!

1. Fork the repository and create a feature branch
2. Follow coding standards and add tests
3. Validate with `make test` and `make lint`
4. Open a PR with a clear description

Guidelines:

- Follow Go best practices
- Write clear, documented code
- Include tests for new functionality
- Update docs as needed
- Be respectful and collaborative

See CONTRIBUTING.md for details.

---

## 📚 Documentation

- Installation Guide: INSTALL.md
- API Reference: https://docs.pipeops.io
- CLI User Guide: https://docs.pipeops.io/cli
- Examples: examples/

---

## 🆘 Support & Community

- Docs: https://docs.pipeops.io
- Issues: https://github.com/PipeOpsHQ/pipeops-cli/issues
- Discussions: https://github.com/PipeOpsHQ/pipeops-cli/discussions
- Discord: https://discord.gg/pipeops
- Email: support@pipeops.io
- Twitter: https://twitter.com/PipeOpsHQ

---

## 🔄 Release Process

Releases are automated via GitHub Actions on tag push:

1. Create a tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
2. Push the tag: `git push origin v1.0.0`
3. CI will:
   - Build binaries for all platforms
   - Create a GitHub release with artifacts
   - Push Docker images
   - Update package managers (e.g., Homebrew)

---

## 📄 License

Licensed under the MIT License. See LICENSE.

---

## 🙏 Acknowledgments

- All PipeOps CLI contributors and users
- The Go community for outstanding tooling
- GitHub for CI/CD infrastructure
- The broader open-source community for inspiration

---

_Made with ❤️ by the PipeOps team_
