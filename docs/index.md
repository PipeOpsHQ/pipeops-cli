# PipeOps CLI Documentation

Welcome to the comprehensive documentation for PipeOps CLI - a powerful command-line interface designed to simplify managing cloud-native environments, deploying projects, and interacting with the PipeOps platform.

## 🚀 What is PipeOps CLI?

PipeOps CLI is a modern, cross-platform command-line tool that provides a unified interface for:

- **🔐 Authentication**: Secure OAuth-based authentication with PKCE flow
- **📦 Project Management**: Create, manage, and deploy projects seamlessly
- **🚀 Server Management**: Provision and configure servers across multiple environments
- **🔧 Pipeline Management**: Create, manage, and deploy CI/CD pipelines
- **🤖 Agent Setup**: Install and configure PipeOps agents for various platforms
- **🌐 Cross-Platform Support**: Available for Linux, Windows, macOS, and FreeBSD

## 🎯 Key Features

- **Beautiful Terminal UI**: Rich interface with colors, progress indicators, and intuitive design
- **Comprehensive Command Set**: Over 20 commands covering all aspects of cloud-native management
- **Docker Support**: Run in containers with full functionality
- **CI/CD Ready**: Perfect for automation and integration workflows
- **Extensive Documentation**: Detailed guides, examples, and API references

## 📖 Quick Navigation

=== "Getting Started"

    - [Installation](getting-started/installation.md) - Install PipeOps CLI on your system
    - [Quick Start](getting-started/quick-start.md) - Get up and running in minutes
    - [Configuration](getting-started/configuration.md) - Configure your CLI environment

=== "Commands"

    - [Overview](commands/overview.md) - Complete command reference
    - [Authentication](commands/auth.md) - Login, logout, and user management
    - [Projects](commands/projects.md) - Project creation and management
    - [Deployments](commands/deployments.md) - Deploy applications and pipelines
    - [Agents](commands/agents.md) - Install and manage PipeOps agents
    - [K3s](commands/k3s.md) - Kubernetes cluster management

=== "Advanced"

    - [Docker Usage](advanced/docker.md) - Running CLI in containers
    - [CI/CD Integration](advanced/ci-cd.md) - Automation workflows
    - [Troubleshooting](advanced/troubleshooting.md) - Common issues and solutions

## 🏃‍♂️ Quick Start

Get started with PipeOps CLI in just a few steps:

1. **Install the CLI**:
   ```bash
   curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.sh | sh
   ```

2. **Authenticate**:
   ```bash
   pipeops auth login
   ```

3. **List your projects**:
   ```bash
   pipeops project list
   ```

4. **Get help**:
   ```bash
   pipeops --help
   ```

## 🌟 What's New

### Latest Features

- **Enhanced Agent Installation**: Full support for PipeOps agent installation with intelligent cluster detection
- **Multi-Platform Support**: Native support for k3s, minikube, k3d, and kind clusters
- **Monitoring Integration**: Built-in Prometheus, Loki, Grafana, and OpenCost stack
- **Worker Node Management**: Easy joining of worker nodes to existing clusters

### Recent Updates

- Improved authentication flow with better error handling
- Enhanced project management capabilities
- Better Docker integration and container support
- Comprehensive documentation and examples

## 📊 Platform Support

| Platform | Architecture | Status |
|----------|-------------|---------|
| Linux | x86_64 | ✅ |
| Linux | ARM64 | ✅ |
| Linux | ARM | ✅ |
| macOS | x86_64 (Intel) | ✅ |
| macOS | ARM64 (M1/M2) | ✅ |
| Windows | x86_64 | ✅ |
| FreeBSD | x86_64 | ✅ |

## 🤝 Community & Support

- **📖 Documentation**: Comprehensive guides and API references
- **🐛 Issues**: [GitHub Issues](https://github.com/PipeOpsHQ/pipeops-cli/issues)
- **💬 Discussions**: [GitHub Discussions](https://github.com/PipeOpsHQ/pipeops-cli/discussions)
- **🗣️ Discord**: [Join our community](https://discord.gg/pipeops)
- **📧 Email**: [support@pipeops.io](mailto:support@pipeops.io)

## 📄 License

This project is licensed under the [MIT License](reference/license.md).

---

*Made with ❤️ by the PipeOps team*
