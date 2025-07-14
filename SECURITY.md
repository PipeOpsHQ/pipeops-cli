# 🔒 Security Guide for PipeOps CLI

This guide outlines security best practices for building and publishing the PipeOps CLI, including protecting sensitive information from reverse engineering.

## 🛡️ Security Architecture Overview

### Current Security Measures
- **PKCE OAuth Flow**: No client secrets required
- **Local Token Storage**: Tokens stored in user's home directory with restricted permissions
- **Environment-Based Configuration**: Sensitive configuration via environment variables
- **Build-Time Injection**: Configuration injected during compilation

## 🔧 Build-Time Security Configuration

### 1. Basic Secure Build
```bash
# Build with custom configuration
make build-secure \
  CLIENT_ID="your_oauth_client_id" \
  API_URL="https://your-api.example.com" \
  GITHUB_REPO="your-org/your-repo"
```

### 2. Enterprise Build
```bash
# Build for enterprise deployment
make build-enterprise
```

### 3. Public Release Build
```bash
# Build for public distribution (with symbol stripping)
make build-public
```

### 4. Maximum Security Build
```bash
# Build with all security features
make build-stripped
make build-compressed  # Requires UPX
```

## 🔐 Advanced Protection Strategies

### 1. **Symbol Stripping** (Implemented)
```bash
go build -ldflags "-s -w" -trimpath -o bin/pipeops .
```
- `-s`: Strip symbol table
- `-w`: Strip debug information
- `-trimpath`: Remove build path information

### 2. **Binary Obfuscation**
```bash
# Using garble (install: go install mvdan.cc/garble@latest)
garble -literals -tiny build -o bin/pipeops .
```

### 3. **Configuration Encryption**
For highly sensitive deployments:

```go
// Add to internal/config/config.go
import "crypto/aes"

func (c *Config) EncryptSensitiveData(key []byte) error {
    // Encrypt sensitive fields before storage
    // Implementation depends on your security requirements
}
```

### 4. **Runtime Configuration Verification**
```go
// Add integrity checks
func verifyConfigIntegrity() error {
    // Check configuration hasn't been tampered with
    // Verify against known checksums or signatures
}
```

## 🌐 Environment-Based Security

### Production Environment Variables
```bash
# OAuth Configuration
export PIPEOPS_CLIENT_ID="prod_client_id"
export PIPEOPS_API_URL="https://api.pipeops.sh"
export PIPEOPS_SCOPES="read:user,read:projects,write:projects"

# Update Configuration
export PIPEOPS_GITHUB_REPO="PipeOpsHQ/pipeops-cli"
export GITHUB_TOKEN="your_github_token"  # For private repos

# Security Settings
export PIPEOPS_DEBUG="false"  # Never enable in production
```

### Development Environment
```bash
# Development-only settings
export PIPEOPS_CLIENT_ID="dev_client_id"
export PIPEOPS_API_URL="https://dev-api.pipeops.sh"
export PIPEOPS_DEBUG="true"
```

## 🔍 What Gets Protected

### ✅ **Secured Information**
- OAuth client IDs (build-time injection)
- API endpoints (environment variables)
- GitHub repository references (configurable)
- Debug symbols (stripped)
- Build paths (trimmed)

### ❌ **Not Exposed in Binary**
- No OAuth client secrets (using PKCE)
- No private keys
- No user tokens (stored locally)
- No hardcoded credentials

### ⚠️ **Visible in Binary** (by design)
- Public API endpoints
- OAuth flow endpoints
- Command help text
- Error messages

## 🚀 CI/CD Security Pipeline

### GitHub Actions Example
```yaml
name: Secure Build

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.23

      - name: Build with security settings
        env:
          CLIENT_ID: ${{ secrets.OAUTH_CLIENT_ID }}
          API_URL: ${{ secrets.API_URL }}
          GITHUB_REPO: ${{ secrets.GITHUB_REPO }}
        run: |
          make build-secure

      - name: Compress binary
        run: |
          upx --best --ultra-brute bin/pipeops

      - name: Sign binary (optional)
        run: |
          # Sign with your code signing certificate
          # codesign -s "Your Certificate" bin/pipeops
```

## 🔒 Additional Security Measures

### 1. **Code Signing**
```bash
# macOS
codesign -s "Developer ID Application: Your Name" bin/pipeops

# Windows
signtool sign /f certificate.pfx /p password bin/pipeops.exe
```

### 2. **Checksum Verification**
```bash
# Generate checksums for releases
sha256sum bin/pipeops > bin/pipeops.sha256
```

### 3. **Binary Analysis Prevention**
```bash
# Install anti-debugging measures (advanced)
go build -tags="prod" -ldflags="-s -w -X main.debug=false"
```

## 🛠️ Development Security

### Safe Development Practices
1. **Never commit secrets**: Use `.env` files (gitignored)
2. **Use separate OAuth clients**: Dev, staging, prod
3. **Rotate credentials regularly**: Especially for production
4. **Monitor binary size**: Large increases might indicate exposed data

### Testing Security
```bash
# Check for exposed secrets in binary
strings bin/pipeops | grep -i "secret\|token\|key\|password"

# Check for build paths
strings bin/pipeops | grep -i "/Users/\|/home/\|C:\\"

# Verify symbol stripping
objdump -t bin/pipeops | head -20
```

## 📋 Security Checklist

### Before Public Release
- [ ] All secrets moved to environment variables
- [ ] Binary built with stripped symbols (`-s -w`)
- [ ] Build paths removed (`-trimpath`)
- [ ] OAuth client configured for production
- [ ] Update mechanism secured (HTTPS, signed releases)
- [ ] Binary compressed/obfuscated if needed
- [ ] Code signed with valid certificate
- [ ] Checksums generated for verification

### Ongoing Security
- [ ] Regular dependency updates
- [ ] OAuth token rotation
- [ ] Security vulnerability scanning
- [ ] Binary integrity monitoring
- [ ] Access logs monitoring

## 🆘 Incident Response

### If Binary is Compromised
1. **Revoke OAuth client**: Immediately revoke compromised client
2. **Rotate secrets**: Change all configuration values
3. **Release new version**: Build with new security settings
4. **Notify users**: Provide update instructions
5. **Audit logs**: Check for unauthorized access

### Emergency Contacts
- Security team: security@pipeops.sh
- DevOps team: devops@pipeops.sh

---

**Remember**: Security is a continuous process. Regularly review and update these practices as your application evolves.