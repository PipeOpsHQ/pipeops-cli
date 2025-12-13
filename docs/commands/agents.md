# Agent Commands

Comprehensive guide to PipeOps agent installation and management commands.

## Overview

PipeOps agents enable you to manage Kubernetes clusters seamlessly. The agent commands provide functionality to install, configure, and manage PipeOps agents across different Kubernetes distributions.

## Available Commands

| Command | Description | Usage |
|---------|-------------|-------|
| `install` | Install PipeOps agent and Kubernetes cluster | `pipeops agent install [token]` |
| `join` | Join worker node to existing cluster | `pipeops agent join <server-url> <token>` |
| `info` | Show cluster information and join commands | `pipeops agent info` |

## Installation

### Basic Installation

Install PipeOps agent with automatic cluster detection:

```bash
# Install with token as argument
pipeops agent install your-token-here

# Install using environment variables
export PIPEOPS_TOKEN="your-token"
export CLUSTER_NAME="my-cluster"
pipeops agent install
```

### Windows Notes

- The bootstrap installer (`pipeops agent install` without `--existing-cluster`) runs a Bash script. On Windows, run it from WSL2 or Git Bash.
- If you already have a Kubernetes cluster and `kubectl` is configured, use `pipeops agent install --existing-cluster --cluster-name="my-cluster"`.

### Installation Options

#### Cluster Type Selection

Choose your Kubernetes distribution:

```bash
# Specify cluster type
pipeops agent install --cluster-type k3s
pipeops agent install --cluster-type minikube
pipeops agent install --cluster-type k3d
pipeops agent install --cluster-type kind
pipeops agent install --cluster-type auto  # Default
```

#### Cluster Naming

Set a custom name for your cluster:

```bash
# Set cluster name
pipeops agent install --cluster-name "production-cluster"
```

#### Monitoring Stack

Control monitoring stack installation:

```bash
# Install without monitoring stack
pipeops agent install --no-monitoring

# Install with monitoring (default)
pipeops agent install
```

#### Existing Cluster Installation

Install agent on an existing Kubernetes cluster:

```bash
# Install on existing cluster
pipeops agent install --existing-cluster --cluster-name="my-existing-cluster"
```

### Installation Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--cluster-name` | Name for the cluster | `pipeops-cluster` |
| `--cluster-type` | Kubernetes distribution (k3s\|minikube\|k3d\|kind\|auto) | `auto` |
| `--existing-cluster` | Install agent on existing Kubernetes cluster | `false` |
| `--no-monitoring` | Skip installation of monitoring stack | `false` |
| `--update` | Update the agent to the latest version | `false` |
| `--uninstall` | Uninstall the agent and monitoring stack | `false` |

## Joining Worker Nodes

### Basic Join

Join a worker node to an existing cluster:

```bash
# Join with server URL and token
pipeops agent join https://192.168.1.100:6443 abc123def456
```

### Environment Variables

Use environment variables for joining:

```bash
# Set environment variables
export K3S_URL="https://192.168.1.100:6443"
export K3S_TOKEN="abc123def456"

# Join using environment variables
pipeops agent join
```

## Cluster Information

### View Cluster Info

Get cluster information and join commands:

```bash
# Show cluster information
pipeops agent info
```

This command displays:
- Server URL for joining worker nodes
- Join token for worker nodes
- Cluster status and configuration
- Connection details

## Environment Variables

The agent commands support various environment variables:

| Variable | Description | Required |
|----------|-------------|----------|
| `PIPEOPS_TOKEN` | Control plane token | Yes |
| `CLUSTER_NAME` | Cluster name | No |
| `CLUSTER_TYPE` | Kubernetes distribution | No |
| `K3S_URL` | Server URL for joining | Yes (for join) |
| `K3S_TOKEN` | Token for joining | Yes (for join) |
| `INSTALL_MONITORING` | Enable/disable monitoring | No |

## Usage Examples

### Complete Installation Workflow

```bash
# 1. Install agent with monitoring
export PIPEOPS_TOKEN="your-token"
export CLUSTER_NAME="production-cluster"
pipeops agent install

# 2. Get cluster information
pipeops agent info

# 3. Join additional worker nodes
export K3S_URL="https://192.168.1.100:6443"
export K3S_TOKEN="abc123def456"
pipeops agent join
```

### Development Environment Setup

```bash
# Install lightweight cluster for development
pipeops agent install \
  --cluster-name "dev-cluster" \
  --cluster-type k3s \
  --no-monitoring
```

### Production Environment Setup

```bash
# Install production cluster with monitoring
pipeops agent install \
  --cluster-name "production-cluster" \
  --cluster-type auto \
  --cluster-name "prod-k8s"
```

### Existing Cluster Integration

```bash
# Install agent on existing cluster
pipeops agent install \
  --existing-cluster \
  --cluster-name "existing-cluster" \
  --no-monitoring
```

## Update and Maintenance

### Update Agent

Update the agent to the latest version:

```bash
# Update agent
pipeops agent install --update
```

### Uninstall Agent

Remove the agent and monitoring stack:

```bash
# Uninstall agent
pipeops agent install --uninstall
```

## Verification

After installation, verify the setup:

```bash
# Check agent pods
kubectl get pods -n pipeops-system

# Check monitoring pods
kubectl get pods -n pipeops-monitoring

# Check cluster nodes
kubectl get nodes

# Check agent logs
kubectl logs deployment/pipeops-agent -n pipeops-system
```

## Troubleshooting

### Common Issues

#### Installation Fails

```bash
# Check token validity
pipeops auth status

# Verify network connectivity
curl -I https://api.pipeops.io

# Check system requirements
pipeops agent install --help
```

#### Join Fails

```bash
# Verify server URL
curl -k https://192.168.1.100:6443/version

# Check token validity
echo $K3S_TOKEN

# Verify network connectivity
ping 192.168.1.100
```

#### Agent Not Starting

```bash
# Check agent logs
kubectl logs deployment/pipeops-agent -n pipeops-system

# Check agent status
kubectl get pods -n pipeops-system

# Restart agent
kubectl rollout restart deployment/pipeops-agent -n pipeops-system
```

### Debug Mode

Enable debug logging:

```bash
# Set debug log level
export PIPEOPS_LOG_LEVEL=debug

# Run installation with debug
pipeops agent install --verbose
```

## Security Considerations

### Token Management

- Store tokens securely
- Use environment variables in production
- Rotate tokens regularly
- Never commit tokens to version control

### Network Security

- Use HTTPS for all connections
- Configure firewall rules appropriately
- Consider VPN for remote access
- Monitor network traffic

## Related Documentation

- **[Installation Guide](../getting-started/installation.md)** - CLI installation
- **[Quick Start](../getting-started/quick-start.md)** - Getting started
- **[Troubleshooting](../advanced/troubleshooting.md)** - Common issues
- **[Docker Usage](../advanced/docker.md)** - Container usage

## Support

If you encounter issues:

1. **Check the [troubleshooting guide](../advanced/troubleshooting.md)**
2. **Open an [issue on GitHub](https://github.com/PipeOpsHQ/pipeops-cli/issues)**
3. **Join our [Discord community](https://discord.gg/pipeops)**
4. **Email us at [support@pipeops.io](mailto:support@pipeops.io)**
