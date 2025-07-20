# 🏗️ Build Reference Guide

This guide explains the different build options available for the PipeOps CLI.

## 🚀 Quick Start

```bash
# Default secure build (recommended)
make build

# Development build (faster, includes debug symbols)
make build-dev

# Just run the CLI
make run
```

## 📋 Available Build Targets

### Production Builds (Secure)

| Target | Description | Output | Use Case |
|--------|-------------|--------|----------|
| `make build` | **Default secure production build** | `build/bin/pipeops` | General use, CI/CD, releases |
| `make build-public` | Alias for default build | `build/bin/pipeops` | Public releases (same as default) |
| `make build-race` | Secure build with race detection | `build/bin/pipeops-race` | Testing for race conditions |
| `make build-compressed` | Secure build with UPX compression | `build/bin/pipeops` | Minimal binary size |

### Development Builds

| Target | Description | Output | Use Case |
|--------|-------------|--------|----------|
| `make build-dev` | Development build with debug symbols | `build/bin/pipeops-dev` | Local development, debugging |

### Custom Configuration Builds

| Target | Description | Environment Variables | Use Case |
|--------|-------------|----------------------|----------|
| `make build-secure` | Custom configuration injection | `CLIENT_ID`, `API_URL`, `SCOPES`, `GITHUB_REPO` | Custom deployments |
| `make build-enterprise` | Pre-configured for enterprise | None (hardcoded enterprise settings) | Enterprise deployments |

## 🔧 Build Features

### Default Build (`make build`)
- **Symbol stripping**: `-s -w` flags for smaller binary size
- **Build path removal**: `-trimpath` flag for security
- **Build-time configuration**: Secure defaults injected at compile time
- **Version information**: Git version, build date, commit hash included
- **Production-ready**: Optimized for release and distribution

### Development Build (`make build-dev`)
- **Debug symbols**: Includes full debugging information
- **Faster builds**: No symbol stripping or path trimming
- **Development defaults**: Uses build-time configuration for development
- **Separate binary**: `pipeops-dev` to avoid conflicts

## 🌐 Environment Variables

### For `make build-secure`
```bash
export CLIENT_ID="your_oauth_client_id"
export API_URL="https://your-api.example.com"
export SCOPES="read:user,read:projects,write:projects"
export GITHUB_REPO="your-org/your-repo"
make build-secure
```

### Runtime Configuration
```bash
export PIPEOPS_CLIENT_ID="runtime_client_id"    # Override build-time client ID
export PIPEOPS_API_URL="https://custom-api.com"  # Override build-time API URL
export PIPEOPS_GITHUB_REPO="custom-org/repo"     # Override build-time GitHub repo
export GITHUB_TOKEN="your_github_token"          # For private repository updates
```

## 🔒 Security Features

All production builds include:
- **No hardcoded secrets**: All sensitive data via environment variables
- **Stripped debug symbols**: Smaller binaries, harder to reverse engineer
- **Removed build paths**: No local development paths in binary
- **Build-time injection**: Configuration set during compilation
- **PKCE OAuth**: No client secrets required

## 🛠️ Development Workflow

### For Development
```bash
# Fast iterative development
make build-dev
./build/bin/pipeops-dev --help

# Test with race detection
make build-race
./build/bin/pipeops-race --help
```

### For Testing
```bash
# Build and run
make run

# Test different configurations
CLIENT_ID="test_client" make build-secure
./build/bin/pipeops --help
```

### For Release
```bash
# Production build
make build

# Compressed release
make build-compressed

# Enterprise deployment
make build-enterprise
```

## 🚀 CI/CD Integration

### GitHub Actions Example
```yaml
- name: Build secure release
  run: |
    make build

- name: Build enterprise version
  run: |
    make build-enterprise

- name: Build with custom config
  env:
    CLIENT_ID: ${{ secrets.CLIENT_ID }}
    API_URL: ${{ secrets.API_URL }}
  run: |
    make build-secure
```

### Docker Build
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:latest
COPY --from=builder /app/build/bin/pipeops /usr/local/bin/
```

## 🔍 Binary Verification

### Check Security
```bash
# Check for exposed secrets (should find none)
strings build/bin/pipeops | grep -i "secret\|token\|key\|password"

# Check for build paths (should find none)
strings build/bin/pipeops | grep -E "/Users/|/home/"

# Check binary size
ls -lh build/bin/pipeops
```

### Verify Functionality
```bash
# Test basic functionality
./build/bin/pipeops --help
./build/bin/pipeops version
./build/bin/pipeops auth --help
```

## 📊 Build Comparison

| Feature | Default Build | Development Build | Secure Build |
|---------|---------------|-------------------|--------------|
| Symbol stripping | ✅ | ❌ | ✅ |
| Build path removal | ✅ | ❌ | ✅ |
| Build-time config | ✅ | ✅ | ✅ (custom) |
| Debug symbols | ❌ | ✅ | ❌ |
| Binary size | Small | Large | Small |
| Build speed | Medium | Fast | Medium |
| Security | High | Medium | High |

## 💡 Tips

1. **Use `make build` for most cases** - It's secure and production-ready
2. **Use `make build-dev` for development** - Faster builds with debug info
3. **Use `make build-secure` for custom deployments** - Full configuration control
4. **Check binary security** - Always verify no secrets are exposed
5. **Test functionality** - Ensure the binary works as expected

## 🆘 Troubleshooting

### Build Fails
```bash
# Clean and rebuild
make clean
make build
```

### Binary Too Large
```bash
# Use compressed build
make build-compressed
```

### Missing Dependencies
```bash
# Install UPX for compression
brew install upx  # macOS
apt-get install upx  # Ubuntu
```

---

**Remember**: The default `make build` is secure and production-ready! 🔒✨