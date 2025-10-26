# Docker Usage

Learn how to use PipeOps CLI in Docker containers for development, CI/CD, and production environments.

## 🐳 Quick Start

### Basic Usage

Run PipeOps CLI in a Docker container:

```bash
# Basic help
docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:latest --help

# Check version
docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:latest --version
```

### With Authentication

```bash
# Mount configuration file
docker run --rm -it \
  -v ~/.pipeops.json:/root/.pipeops.json \
  ghcr.io/pipeopshq/pipeops-cli:latest auth status
```

## 🔧 Configuration

### Environment Variables

Configure PipeOps CLI using environment variables:

```bash
# Run with environment variables
docker run --rm -it \
  -e PIPEOPS_TOKEN="your-token" \
  -e PIPEOPS_API_URL="https://api.pipeops.io" \
  -e PIPEOPS_LOG_LEVEL="debug" \
  ghcr.io/pipeopshq/pipeops-cli:latest auth status
```

### Volume Mounts

Mount configuration and data directories:

```bash
# Mount config and data directories
docker run --rm -it \
  -v ~/.pipeops.json:/root/.pipeops.json \
  -v ~/.pipeops-data:/root/.pipeops-data \
  ghcr.io/pipeopshq/pipeops-cli:latest project list
```

## 📦 Docker Compose

### Basic Setup

```yaml
version: '3.8'
services:
  pipeops-cli:
    image: ghcr.io/pipeopshq/pipeops-cli:latest
    environment:
      - PIPEOPS_TOKEN=${PIPEOPS_TOKEN}
      - PIPEOPS_API_URL=https://api.pipeops.io
    volumes:
      - ~/.pipeops.json:/root/.pipeops.json
    command: ["project", "list"]
```

### Development Environment

```yaml
version: '3.8'
services:
  pipeops-cli:
    image: ghcr.io/pipeopshq/pipeops-cli:latest
    environment:
      - PIPEOPS_API_URL=https://dev-api.pipeops.io
      - PIPEOPS_LOG_LEVEL=debug
      - PIPEOPS_OUTPUT_FORMAT=json
    volumes:
      - ~/.pipeops.json:/root/.pipeops.json
      - ./scripts:/scripts
    working_dir: /scripts
    command: ["sh", "-c", "pipeops project list && pipeops server list"]
```

### CI/CD Pipeline

```yaml
version: '3.8'
services:
  pipeops-deploy:
    image: ghcr.io/pipeopshq/pipeops-cli:latest
    environment:
      - PIPEOPS_TOKEN=${PIPEOPS_SERVICE_TOKEN}
      - PIPEOPS_API_URL=https://api.pipeops.io
      - PIPEOPS_LOG_LEVEL=info
      - PIPEOPS_OUTPUT_FORMAT=json
    volumes:
      - ./deploy:/deploy
    working_dir: /deploy
    command: ["sh", "-c", "pipeops deploy create --name production --image my-app:$BUILD_NUMBER"]
```

## 🚀 Use Cases

### Development

```bash
# Interactive development session
docker run --rm -it \
  -v ~/.pipeops.json:/root/.pipeops.json \
  -v $(pwd):/workspace \
  --workdir /workspace \
  ghcr.io/pipeopshq/pipeops-cli:latest

# Run specific commands
docker run --rm -it \
  -v ~/.pipeops.json:/root/.pipeops.json \
  ghcr.io/pipeopshq/pipeops-cli:latest project create my-project
```

### CI/CD Integration

```bash
# Deploy from CI/CD
docker run --rm \
  -e PIPEOPS_TOKEN="${PIPEOPS_SERVICE_TOKEN}" \
  ghcr.io/pipeopshq/pipeops-cli:latest deploy create \
  --name production \
  --image my-app:${BUILD_NUMBER}
```

### Automation Scripts

```bash
# Automated deployment script
docker run --rm \
  -v ~/.pipeops.json:/root/.pipeops.json \
  -v $(pwd)/scripts:/scripts \
  ghcr.io/pipeopshq/pipeops-cli:latest \
  sh -c "cd /scripts && ./deploy.sh"
```

## 🔒 Security Considerations

### Token Management

```bash
# Use environment variables for tokens
docker run --rm -it \
  -e PIPEOPS_TOKEN="$(cat ~/.pipeops-token)" \
  ghcr.io/pipeopshq/pipeops-cli:latest auth status

# Use Docker secrets (Docker Swarm)
docker service create \
  --secret source=pipeops-token,target=/run/secrets/pipeops-token \
  ghcr.io/pipeopshq/pipeops-cli:latest \
  sh -c "export PIPEOPS_TOKEN=\$(cat /run/secrets/pipeops-token) && pipeops auth status"
```

### Network Security

```bash
# Use custom network
docker network create pipeops-network

docker run --rm -it \
  --network pipeops-network \
  ghcr.io/pipeopshq/pipeops-cli:latest auth status
```

### Resource Limits

```bash
# Set resource limits
docker run --rm -it \
  --memory="512m" \
  --cpus="0.5" \
  ghcr.io/pipeopshq/pipeops-cli:latest project list
```

## 🌐 Multi-Platform Support

### Platform-Specific Images

```bash
# Linux AMD64
docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:latest-linux-amd64 --help

# Linux ARM64
docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:latest-linux-arm64 --help

# Windows (using Windows containers)
docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:latest-windows-amd64 --help
```

### Cross-Platform Builds

```bash
# Build for multiple platforms
docker buildx build --platform linux/amd64,linux/arm64 -t pipeops-cli:latest .
```

## 🔄 Updates and Maintenance

### Image Updates

```bash
# Pull latest image
docker pull ghcr.io/pipeopshq/pipeops-cli:latest

# Update running containers
docker-compose pull
docker-compose up -d
```

### Version Pinning

```bash
# Use specific version
docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:1.0.0 --help

# Use version tags
docker run --rm -it ghcr.io/pipeopshq/pipeops-cli:v1.0.0 --help
```

## 📊 Monitoring and Logging

### Log Collection

```bash
# Collect logs
docker logs container-name

# Follow logs
docker logs -f container-name
```

### Health Checks

```bash
# Add health check
docker run --rm -it \
  --health-cmd="pipeops auth status" \
  --health-interval=30s \
  --health-timeout=10s \
  --health-retries=3 \
  ghcr.io/pipeopshq/pipeops-cli:latest
```

## 🐛 Troubleshooting

### Common Issues

#### Permission Issues

```bash
# Fix permission issues
docker run --rm -it \
  --user $(id -u):$(id -g) \
  -v ~/.pipeops.json:/home/user/.pipeops.json \
  ghcr.io/pipeopshq/pipeops-cli:latest auth status
```

#### Network Issues

```bash
# Debug network connectivity
docker run --rm -it \
  --network host \
  ghcr.io/pipeopshq/pipeops-cli:latest auth status
```

#### Volume Mount Issues

```bash
# Check volume mounts
docker run --rm -it \
  -v ~/.pipeops.json:/root/.pipeops.json:ro \
  ghcr.io/pipeopshq/pipeops-cli:latest auth status
```

### Debug Mode

```bash
# Enable debug mode
docker run --rm -it \
  -e PIPEOPS_LOG_LEVEL=debug \
  ghcr.io/pipeopshq/pipeops-cli:latest auth status
```

## 📚 Related Documentation

- **[Installation](../getting-started/installation.md)** - Installation guide
- **[Configuration](../getting-started/configuration.md)** - Configuration options
- **[CI/CD Integration](ci-cd.md)** - CI/CD workflows
- **[Troubleshooting](troubleshooting.md)** - Common issues
