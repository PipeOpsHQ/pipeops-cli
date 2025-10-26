# Tailscale Kubernetes CLI Documentation

Welcome to the comprehensive documentation for Tailscale Kubernetes CLI - a powerful command-line interface designed to simplify Tailscale installation, configuration, and management for Kubernetes clusters with Tailscale Funnel support.

## 🚀 What is Tailscale Kubernetes CLI?

Tailscale Kubernetes CLI is a modern, cross-platform command-line tool that provides a unified interface for:

- **🔐 Tailscale Installation**: Automatic Tailscale installation and configuration
- **🌐 Tailscale Funnel**: Easy setup for public port 80 exposure via Tailscale Funnel
- **🚀 Kubernetes Integration**: Native support for k3s, minikube, k3d, and kind clusters
- **🔧 Ingress Management**: Automatic ingress configuration with Tailscale annotations
- **📦 Operator Setup**: Automated Tailscale Kubernetes operator installation
- **🌍 Public Access**: Secure public internet access to your Kubernetes services

## 🎯 Key Features

- **Beautiful Terminal UI**: Rich interface with colors, progress indicators, and intuitive design
- **Comprehensive Command Set**: Commands covering Tailscale installation and Kubernetes integration
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
    - [Agent Management](commands/agents.md) - Install and manage Tailscale agents
    - [K3s](commands/k3s.md) - Kubernetes cluster management
    - [Proxy](commands/proxy.md) - Proxy and tunnel management

=== "Advanced"

    - [Docker Usage](advanced/docker.md) - Running CLI in containers
    - [CI/CD Integration](advanced/ci-cd.md) - Automation workflows
    - [Troubleshooting](advanced/troubleshooting.md) - Common issues and solutions

## 🏃‍♂️ Quick Start

Get started with Tailscale Kubernetes CLI in just a few steps:

1. **Install the CLI**:
   ```bash
   curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-cli/main/install.sh | sh
   ```

2. **Install Tailscale and setup cluster**:
   ```bash
   pipeops agent install tskey-auth-your-key-here
   ```

3. **Check Tailscale status**:
   ```bash
   tailscale status
   ```

4. **Get help**:
   ```bash
   pipeops --help
   ```

## 🌟 What's New

### Latest Features

- **Tailscale Funnel Integration**: Full support for Tailscale Funnel with automatic port 80 exposure
- **Kubernetes Operator**: Automated Tailscale Kubernetes operator installation and configuration
- **Multi-Platform Support**: Native support for k3s, minikube, k3d, and kind clusters
- **Ingress Management**: Automatic ingress configuration with Tailscale annotations
- **Public Access**: Secure public internet access to your Kubernetes services via Tailscale Funnel

### Recent Updates

- Enhanced Tailscale installation with better error handling
- Improved Kubernetes cluster integration
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
