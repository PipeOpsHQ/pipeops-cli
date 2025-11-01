# Quick Start

Get up and running with PipeOps CLI in just a few minutes. This guide will walk you through the essential steps to start using PipeOps CLI effectively.

## 🚀 Prerequisites

Before you begin, make sure you have:

- PipeOps CLI installed (see [Installation](installation.md))
- A PipeOps account (sign up at [pipeops.io](https://pipeops.io))
- Basic familiarity with command-line interfaces

## 📋 Step 1: Verify Installation

First, let's make sure PipeOps CLI is properly installed:

```bash
# Check version
pipeops --version

# View available commands
pipeops --help
```

You should see output similar to:

```
🚀 PipeOps CLI Version: 1.0.0

Usage:
  pipeops [command]

Available Commands:
  agent       ⚙️ Manage agent-related commands and tasks
  auth        🔐 Authentication and user management
  deploy      🚀 Deploy applications and pipelines
  project     📦 Project management
  server      🖥️ Server management
  k3s         ☸️ K3s cluster management
  help        Help about any command
```

## 🔐 Step 2: Authenticate

Authenticate with your PipeOps account:

```bash
pipeops auth login
```

This will:

1. Open your default web browser
2. Redirect you to the PipeOps authentication page
3. Complete the OAuth flow
4. Save your credentials locally

!!! tip "Authentication Methods"
    PipeOps CLI supports multiple authentication methods:
    - **OAuth (Recommended)**: Secure browser-based authentication
    - **Service Account Token**: For CI/CD and automation
    - **Environment Variables**: For containerized environments

### Verify Authentication

Check your authentication status:

```bash
pipeops auth status
```

You should see your user information and authentication status.

## 📦 Step 3: Explore Your Projects

List your existing projects:

```bash
pipeops project list
```

If you don't have any projects yet, you can create one:

```bash
pipeops project create my-first-project
```

## 🚀 Step 4: Deploy Your First Application

Let's deploy a simple application to see PipeOps CLI in action:

```bash
# Create a new deployment
pipeops deploy create --name hello-world --image nginx:latest

# Check deployment status
pipeops deploy status hello-world

# View deployment logs
pipeops deploy logs hello-world
```

## 🤖 Step 5: Set Up an Agent (Optional)

If you want to manage Kubernetes clusters, install a PipeOps agent:

```bash
# Install agent with automatic cluster detection
pipeops agent install your-token-here

# Or install on existing cluster
pipeops agent install --existing-cluster --cluster-name="my-cluster"
```

## 📊 Step 6: Monitor and Manage

Use these commands to monitor and manage your resources:

```bash
# Check server status
pipeops server list

# View project logs
pipeops project logs my-first-project

# Get system status
pipeops status
```

## 🎯 Common Workflows

### Daily Operations

```bash
# Morning routine - check status
pipeops auth status
pipeops project list
pipeops server list

# Deploy updates
pipeops deploy update my-app --image my-app:v2.0

# Monitor deployments
pipeops deploy logs my-app --follow
```

### Development Workflow

```bash
# Create new project
pipeops project create my-dev-project

# Deploy development environment
pipeops deploy create --name dev-env --image my-app:dev

# Test locally
pipeops proxy start --port 8080

# Clean up
pipeops deploy delete dev-env
```

### CI/CD Integration

```bash
# Authenticate with service account
export PIPEOPS_TOKEN="your-service-account-token"

# Deploy from CI/CD
pipeops deploy create --name production --image my-app:$BUILD_NUMBER

# Verify deployment
pipeops deploy status production
```

## 🔧 Configuration

### Environment Variables

Set up environment variables for automation:

```bash
# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export PIPEOPS_API_URL="https://api.pipeops.io"
export PIPEOPS_LOG_LEVEL="info"
export PIPEOPS_CONFIG_PATH="$HOME/.pipeops.json"
```

### Configuration File

PipeOps CLI stores configuration in `~/.pipeops.json`:

```json
{
  "version": {
    "version": "1.0.0"
  },
  "updates": {
    "last_update_check": "2024-01-01T00:00:00Z",
    "skip_update_check": false
  },
  "service_account_token": "your-token-here"
}
```

## 📚 Next Steps

Now that you're up and running, explore these areas:

- **[Commands Overview](commands/overview.md)**: Complete command reference
- **[Project Management](commands/projects.md)**: Advanced project operations
- **[Deployment Guide](commands/deployments.md)**: Complex deployment scenarios
- **[Agent Management](commands/agents.md)**: Kubernetes cluster management
- **[Docker Usage](advanced/docker.md)**: Containerized workflows
- **[CI/CD Integration](advanced/ci-cd.md)**: Automation and pipelines

## 🆘 Getting Help

If you run into issues:

1. **Check command help**: `pipeops <command> --help`
2. **View logs**: `pipeops logs <resource>`
3. **Check status**: `pipeops status`
4. **Review [troubleshooting guide](advanced/troubleshooting.md)**
5. **Join our [Discord community](https://discord.gg/pipeops)**

## 🎉 Congratulations!

You've successfully set up PipeOps CLI and completed your first operations. You're now ready to:

- ✅ Manage projects and deployments
- ✅ Monitor your infrastructure
- ✅ Set up automated workflows
- ✅ Scale your applications

Happy deploying! 🚀
