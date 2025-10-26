# Configuration

Learn how to configure PipeOps CLI for your specific needs and environment.

## 📁 Configuration File

PipeOps CLI stores configuration in `~/.pipeops.json` by default. You can specify a custom location using the `--config` flag or `PIPEOPS_CONFIG_PATH` environment variable.

### Configuration Structure

```json
{
  "version": {
    "version": "1.0.0"
  },
  "updates": {
    "last_update_check": "2024-01-01T00:00:00Z",
    "skip_update_check": false
  },
  "service_account_token": "your-token-here",
  "api_url": "https://api.pipeops.io",
  "log_level": "info",
  "output_format": "text"
}
```

## 🌍 Environment Variables

Configure PipeOps CLI using environment variables:

### Authentication

```bash
# Service account token
export PIPEOPS_TOKEN="your-service-account-token"

# API endpoint
export PIPEOPS_API_URL="https://api.pipeops.io"
```

### Configuration

```bash
# Custom config file location
export PIPEOPS_CONFIG_PATH="$HOME/.pipeops-custom.json"

# Log level
export PIPEOPS_LOG_LEVEL="debug"  # debug, info, warn, error

# Output format
export PIPEOPS_OUTPUT_FORMAT="json"  # text, json
```

### Agent Configuration

```bash
# Cluster settings
export CLUSTER_NAME="my-cluster"
export CLUSTER_TYPE="k3s"  # k3s, minikube, k3d, kind, auto

# Monitoring
export INSTALL_MONITORING="true"  # true, false

# K3s settings
export K3S_URL="https://192.168.1.100:6443"
export K3S_TOKEN="your-k3s-token"
```

## 🔧 Configuration Options

### API Configuration

```bash
# Set custom API endpoint
pipeops --config ~/.pipeops.json config set api_url "https://custom-api.pipeops.io"

# Set timeout
pipeops --config ~/.pipeops.json config set timeout "30s"
```

### Logging Configuration

```bash
# Set log level
pipeops --config ~/.pipeops.json config set log_level "debug"

# Set log format
pipeops --config ~/.pipeops.json config set log_format "json"
```

### Update Configuration

```bash
# Disable automatic update checks
pipeops --config ~/.pipeops.json config set skip_update_check true

# Set update check interval
pipeops --config ~/.pipeops.json config set update_check_interval "24h"
```

## 🐳 Docker Configuration

### Environment Variables

```bash
# Run with custom configuration
docker run --rm -it \
  -e PIPEOPS_TOKEN="your-token" \
  -e PIPEOPS_API_URL="https://api.pipeops.io" \
  -e PIPEOPS_LOG_LEVEL="debug" \
  ghcr.io/pipeopshq/pipeops-cli:latest auth status
```

### Docker Compose

```yaml
version: '3.8'
services:
  pipeops-cli:
    image: ghcr.io/pipeopshq/pipeops-cli:latest
    environment:
      - PIPEOPS_TOKEN=${PIPEOPS_TOKEN}
      - PIPEOPS_API_URL=https://api.pipeops.io
      - PIPEOPS_LOG_LEVEL=info
    volumes:
      - ~/.pipeops.json:/root/.pipeops.json
    command: ["project", "list"]
```

## 🔒 Security Configuration

### Token Management

```bash
# Store token securely
export PIPEOPS_TOKEN="$(cat ~/.pipeops-token)"

# Use keychain (macOS)
security add-generic-password -a "pipeops" -s "pipeops-token" -w "your-token"

# Retrieve from keychain
export PIPEOPS_TOKEN="$(security find-generic-password -a "pipeops" -s "pipeops-token" -w)"
```

### TLS Configuration

```bash
# Skip TLS verification (not recommended for production)
export PIPEOPS_TLS_INSECURE_SKIP_VERIFY="true"

# Custom CA certificate
export PIPEOPS_TLS_CA_CERT="/path/to/ca.crt"
```

## 📊 Output Configuration

### JSON Output

```bash
# Global JSON output
pipeops --json project list

# Command-specific JSON output
pipeops project list --json
```

### Verbose Output

```bash
# Enable verbose output
pipeops --verbose project list

# Quiet mode
pipeops --quiet project list
```

### Custom Output Format

```bash
# Set default output format
export PIPEOPS_OUTPUT_FORMAT="json"

# Override per command
pipeops project list --format text
```

## 🔄 Update Configuration

### Automatic Updates

```bash
# Enable automatic updates
pipeops config set auto_update true

# Set update channel
pipeops config set update_channel "stable"  # stable, beta, alpha
```

### Manual Updates

```bash
# Check for updates
pipeops update check

# Update to latest version
pipeops update

# Update to specific version
pipeops update --version 1.0.0
```

## 🌐 Proxy Configuration

### HTTP Proxy

```bash
# Set HTTP proxy
export HTTP_PROXY="http://proxy.company.com:8080"
export HTTPS_PROXY="http://proxy.company.com:8080"

# Set proxy authentication
export HTTP_PROXY="http://user:pass@proxy.company.com:8080"
```

### SOCKS Proxy

```bash
# Set SOCKS proxy
export ALL_PROXY="socks5://proxy.company.com:1080"
```

## 🏢 Enterprise Configuration

### Custom Endpoints

```bash
# Set custom API endpoint
export PIPEOPS_API_URL="https://api.company.pipeops.io"

# Set custom authentication endpoint
export PIPEOPS_AUTH_URL="https://auth.company.pipeops.io"
```

### Corporate Certificates

```bash
# Set custom CA certificate
export PIPEOPS_TLS_CA_CERT="/etc/ssl/certs/company-ca.crt"

# Set client certificate
export PIPEOPS_TLS_CLIENT_CERT="/etc/ssl/certs/client.crt"
export PIPEOPS_TLS_CLIENT_KEY="/etc/ssl/certs/client.key"
```

## 🔍 Configuration Validation

### Validate Configuration

```bash
# Validate configuration file
pipeops config validate

# Check configuration
pipeops config check
```

### Test Connection

```bash
# Test API connection
pipeops config test-connection

# Test authentication
pipeops auth status
```

## 📝 Configuration Examples

### Development Environment

```bash
# Development configuration
export PIPEOPS_API_URL="https://dev-api.pipeops.io"
export PIPEOPS_LOG_LEVEL="debug"
export PIPEOPS_OUTPUT_FORMAT="json"
export PIPEOPS_SKIP_UPDATE_CHECK="true"
```

### Production Environment

```bash
# Production configuration
export PIPEOPS_API_URL="https://api.pipeops.io"
export PIPEOPS_LOG_LEVEL="warn"
export PIPEOPS_OUTPUT_FORMAT="text"
export PIPEOPS_AUTO_UPDATE="true"
```

### CI/CD Environment

```bash
# CI/CD configuration
export PIPEOPS_TOKEN="${PIPEOPS_SERVICE_TOKEN}"
export PIPEOPS_API_URL="https://api.pipeops.io"
export PIPEOPS_LOG_LEVEL="error"
export PIPEOPS_OUTPUT_FORMAT="json"
export PIPEOPS_SKIP_UPDATE_CHECK="true"
```

## 🐛 Troubleshooting Configuration

### Common Issues

#### Configuration Not Loading

```bash
# Check config file location
echo $PIPEOPS_CONFIG_PATH

# Check config file permissions
ls -la ~/.pipeops.json

# Validate config file
pipeops config validate
```

#### Environment Variables Not Working

```bash
# Check environment variables
env | grep PIPEOPS

# Test configuration
pipeops config check
```

#### Authentication Issues

```bash
# Check token
echo $PIPEOPS_TOKEN

# Test authentication
pipeops auth status

# Re-authenticate
pipeops auth login
```

## 📚 Related Documentation

- **[Installation](installation.md)** - Installation guide
- **[Quick Start](quick-start.md)** - Getting started
- **[Authentication Commands](../commands/auth.md)** - Authentication management
- **[Troubleshooting](../advanced/troubleshooting.md)** - Common issues
