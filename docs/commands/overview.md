# Commands Overview

This page provides a comprehensive overview of all PipeOps CLI commands, their usage, and examples.

## 📋 Command Structure

PipeOps CLI follows a hierarchical command structure:

```
pipeops <command> <subcommand> [flags] [arguments]
```

## 🔐 Authentication Commands

Manage authentication and user details.

### `pipeops auth login`

Authenticate with PipeOps using OAuth.

```bash
# Interactive login
pipeops auth login

# Login with specific provider
pipeops auth login --provider github
```

### `pipeops auth logout`

Log out from PipeOps.

```bash
pipeops auth logout
```

### `pipeops auth status`

Check authentication status and user information.

```bash
pipeops auth status
```

### `pipeops auth me`

Display current user information.

```bash
pipeops auth me
```

## 📦 Project Commands

Manage PipeOps projects.

### `pipeops project list`

List all projects.

```bash
# List all projects
pipeops project list

# List with JSON output
pipeops project list --json
```

### `pipeops project create`

Create a new project.

```bash
# Create project with name
pipeops project create my-project

# Create with description
pipeops project create my-project --description "My awesome project"
```

### `pipeops project logs`

View project logs.

```bash
# View recent logs
pipeops project logs my-project

# Follow logs in real-time
pipeops project logs my-project --follow

# View logs with timestamps
pipeops project logs my-project --timestamps
```

## 🚀 Deployment Commands

Manage deployments and pipelines.

### `pipeops deploy pipeline`

Manage deployment pipelines.

```bash
# Create new pipeline
pipeops deploy pipeline create --name my-pipeline

# List pipelines
pipeops deploy pipeline list

# Get pipeline status
pipeops deploy pipeline status my-pipeline
```

### `pipeops deploy create`

Create a new deployment.

```bash
# Deploy with image
pipeops deploy create --name my-app --image nginx:latest

# Deploy with environment variables
pipeops deploy create --name my-app --image my-app:latest --env KEY=value
```

### `pipeops deploy status`

Check deployment status.

```bash
# Check specific deployment
pipeops deploy status my-app

# Check all deployments
pipeops deploy status --all
```

### `pipeops deploy logs`

View deployment logs.

```bash
# View logs
pipeops deploy logs my-app

# Follow logs
pipeops deploy logs my-app --follow
```

## 🖥️ Server Commands

Manage servers and infrastructure.

### `pipeops server list`

List all servers.

```bash
# List servers
pipeops server list

# List with details
pipeops server list --verbose
```

### `pipeops server deploy`

Deploy to a server.

```bash
# Deploy application
pipeops server deploy --name my-app --server server-1
```

## ⚙️ Agent Commands

Manage PipeOps agents for Kubernetes clusters.

### `pipeops agent install`

Install PipeOps agent and Kubernetes cluster.

```bash
# Install with token
pipeops agent install your-token-here

# Install using environment variables
export PIPEOPS_TOKEN="your-token"
export CLUSTER_NAME="my-cluster"
pipeops agent install

# Install on existing cluster
pipeops agent install --existing-cluster --cluster-name="my-existing-cluster"

# Install without monitoring
pipeops agent install --no-monitoring

# Update agent
pipeops agent install --update

# Uninstall agent
pipeops agent install --uninstall
```

### `pipeops agent join`

Join worker node to existing cluster.

```bash
# Join with server URL and token
pipeops agent join https://192.168.1.100:6443 abc123def456

# Join using environment variables
export K3S_URL="https://192.168.1.100:6443"
export K3S_TOKEN="abc123def456"
pipeops agent join
```

### `pipeops agent info`

Show cluster information and join commands.

```bash
pipeops agent info
```

## ☸️ K3s Commands

Manage K3s clusters.

### `pipeops k3s install`

Install K3s server.

```bash
# Install K3s
pipeops k3s install

# Install with specific version
pipeops k3s install --version v1.28.0
```

### `pipeops k3s join`

Join worker node to K3s cluster.

```bash
# Join worker node
pipeops k3s join --server https://192.168.1.100:6443 --token abc123
```

### `pipeops k3s restart`

Restart K3s service.

```bash
# Restart K3s
pipeops k3s restart
```

### `pipeops k3s kill`

Stop K3s service.

```bash
# Stop K3s
pipeops k3s kill
```

## 🔧 Utility Commands

### `pipeops status`

Show overall system status.

```bash
# Check system status
pipeops status

# Check with verbose output
pipeops status --verbose
```

### `pipeops version`

Show version information.

```bash
# Show version
pipeops version

# Show version with build info
pipeops version --verbose
```

### `pipeops update`

Update PipeOps CLI.

```bash
# Check for updates
pipeops update check

# Update to latest version
pipeops update
```

### `pipeops proxy`

Manage proxy connections.

```bash
# Start proxy
pipeops proxy start --port 8080

# Stop proxy
pipeops proxy stop
```

## 🌐 Global Flags

All commands support these global flags:

| Flag | Description | Example |
|------|-------------|---------|
| `--help, -h` | Show help for command | `pipeops auth --help` |
| `--version, -v` | Show version information | `pipeops --version` |
| `--json` | Output in JSON format | `pipeops project list --json` |
| `--verbose` | Enable verbose output | `pipeops status --verbose` |
| `--quiet, -q` | Suppress non-essential output | `pipeops deploy create --quiet` |
| `--config` | Use custom config file | `pipeops --config ~/.pipeops-custom.json` |

## 📝 Command Examples

### Daily Workflow

```bash
# Morning routine
pipeops auth status
pipeops project list
pipeops server list

# Deploy updates
pipeops deploy create --name my-app --image my-app:v2.0
pipeops deploy status my-app

# Monitor
pipeops deploy logs my-app --follow
```

### Development Workflow

```bash
# Create development environment
pipeops project create dev-project
pipeops deploy create --name dev-env --image my-app:dev

# Test locally
pipeops proxy start --port 8080

# Clean up
pipeops deploy delete dev-env
```

### CI/CD Pipeline

```bash
# Authenticate
export PIPEOPS_TOKEN="your-token"

# Deploy
pipeops deploy create --name production --image my-app:$BUILD_NUMBER

# Verify
pipeops deploy status production
```

## 🔍 Getting Help

For detailed help on any command:

```bash
# General help
pipeops --help

# Command-specific help
pipeops auth --help
pipeops project create --help

# Subcommand help
pipeops deploy pipeline --help
```

## 📚 Related Documentation

- **[Authentication Commands](auth.md)** - Detailed auth command reference
- **[Project Commands](projects.md)** - Project management guide
- **[Deployment Commands](deployments.md)** - Deployment operations
- **[Agent Commands](agents.md)** - Agent installation and management
- **[K3s Commands](k3s.md)** - K3s cluster management
