# PipeOps CLI Documentation

Welcome to the official documentation for PipeOps CLI — a fast, cross-platform command-line tool for authenticating with PipeOps, managing projects, provisioning servers, deploying pipelines, and working with agents for Kubernetes cluster management.

## 🚀 What is PipeOps CLI?

PipeOps CLI provides a unified interface to:

- 🔐 Authenticate securely (OAuth with PKCE)
- 📦 Create, list, and manage PipeOps projects
- 🚀 Provision and manage servers and environments
- 🔧 Create, manage, and deploy CI/CD pipelines
- 🤖 Install and manage PipeOps agents
- 📊 Monitor status with rich terminal output or JSON

## 🎯 Key Features

- Beautiful terminal UX with colors, spinners, and progress indicators
- Comprehensive command set for projects, servers, pipelines, and agents
- Docker-friendly: run the CLI in containers
- CI/CD ready: deterministic commands for automation
- Extensive documentation with examples and references

## 📖 Quick Navigation

=== "Getting Started"

    - [Installation](getting-started/installation.md) — Install PipeOps CLI on your system
    - [Quick Start](getting-started/quick-start.md) — Get up and running in minutes
    - [Configuration](getting-started/configuration.md) — Configure your CLI environment

=== "Commands"

    - [Overview](commands/overview.md) — Complete command reference
    - [Authentication](commands/auth.md) — Login, logout, and auth status
    - [Projects](commands/project.md) — Create, list, and manage projects
    - [Deploy](commands/deploy.md) — Pipelines and deployment workflows
    - [Servers](commands/server.md) — Provisioning and server operations
    - [Agents](commands/agents.md) — Install and manage PipeOps agents

=== "Advanced"

    - [Docker Usage](advanced/docker.md) — Running CLI in containers
    - [CI/CD Integration](advanced/ci-cd.md) — Automation workflows
    - [Troubleshooting](advanced/troubleshooting.md) — Common issues and solutions

## 🏃 Quick Start

Get started with PipeOps CLI in a few steps:

1.  Install the CLI (macOS/Linux):

         curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.sh | sh

    For Windows and alternative methods, see the Installation guide.

2.  Log in:

        pipeops auth login

3.  Check authentication status:

        pipeops auth status

4.  List projects:

        pipeops project list

5.  Get help:

        pipeops --help
        pipeops auth --help
        pipeops project --help

## 🌟 What's New

### Latest Features

- Secure OAuth login with PKCE for device-friendly authentication
- Project lifecycle commands: create, list, configure, deploy
- Agent management for Kubernetes cluster integration
- Agent installation and management across major platforms
- JSON output mode for scripting and automation

### Recent Updates

- Improved error messages and diagnostics
- Enhanced Windows and macOS support
- Better Docker ergonomics and examples
- Expanded command help and documentation

## 📊 Platform Support

| Platform | Architecture | Status |
| -------- | ------------ | ------ |
| Linux    | x86_64       | ✅     |
| Linux    | ARM64        | ✅     |
| Linux    | ARM          | ✅     |
| macOS    | x86_64       | ✅     |
| macOS    | ARM64        | ✅     |
| Windows  | x86_64       | ✅     |
| FreeBSD  | x86_64       | ✅     |

## 🤝 Community & Support

- 📖 Documentation: Comprehensive guides and references
- 🐛 Issues: https://github.com/PipeOpsHQ/pipeops-cli/issues
- 💬 Discussions: https://github.com/PipeOpsHQ/pipeops-cli/discussions
- 🗣️ Discord: https://discord.gg/pipeops
- 📧 Email: support@pipeops.io

## 📄 License

This project is licensed under the [MIT License](reference/license.md).

---

Made with ❤️ by the PipeOps team
