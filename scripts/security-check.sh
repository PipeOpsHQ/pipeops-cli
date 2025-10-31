#!/bin/bash
# Security audit script for PipeOps CLI

set -e

echo "🔒 Running Security Checks for PipeOps CLI"
echo "=========================================="
echo ""

# Check for Go vulnerability scanner
if ! command -v govulncheck &> /dev/null; then
    echo "📦 Installing govulncheck..."
    go install golang.org/x/vuln/cmd/govulncheck@latest
fi

# 1. Check for known vulnerabilities
echo "1️⃣  Checking for known vulnerabilities..."
govulncheck ./...
echo "✅ Vulnerability check complete"
echo ""

# 2. Check for hardcoded secrets
echo "2️⃣  Scanning for potential secrets..."
if grep -r "password.*=.*\"" --include="*.go" . 2>/dev/null | grep -v "_test.go" | grep -v "// "; then
    echo "⚠️  Found potential hardcoded passwords"
else
    echo "✅ No hardcoded passwords found"
fi

if grep -r "api.?key.*=.*\"" --include="*.go" . 2>/dev/null | grep -v "_test.go" | grep -v "// "; then
    echo "⚠️  Found potential hardcoded API keys"
else
    echo "✅ No hardcoded API keys found"
fi

if grep -r "secret.*=.*\"" --include="*.go" . 2>/dev/null | grep -v "_test.go" | grep -v "ClientSecret\|generate" | grep -v "// "; then
    echo "⚠️  Found potential hardcoded secrets"
else
    echo "✅ No hardcoded secrets found"
fi
echo ""

# 3. Check for localhost URLs in production code
echo "3️⃣  Checking for localhost references..."
LOCALHOST_COUNT=$(grep -r "localhost" --include="*.go" . 2>/dev/null | grep -v "_test.go" | grep -v "// " | wc -l)
if [ "$LOCALHOST_COUNT" -gt 0 ]; then
    echo "⚠️  Found $LOCALHOST_COUNT localhost references (review for production readiness)"
    grep -n "localhost" --include="*.go" -r . 2>/dev/null | grep -v "_test.go" | grep -v "// "
else
    echo "✅ No localhost references in production code"
fi
echo ""

# 4. Check file permissions for config handling
echo "4️⃣  Verifying secure file permission handling..."
if grep -r "0600\|0400" --include="*.go" internal/config/ internal/auth/; then
    echo "✅ Secure file permissions found in code"
else
    echo "⚠️  No secure file permissions found (files should use 0600 or 0400)"
fi
echo ""

# 5. Check for proper error handling
echo "5️⃣  Checking for os.Exit usage (should be minimal)..."
EXIT_COUNT=$(grep -r "os.Exit" --include="*.go" cmd/ internal/ 2>/dev/null | wc -l)
if [ "$EXIT_COUNT" -gt 5 ]; then
    echo "⚠️  Found $EXIT_COUNT os.Exit calls (consider refactoring to return errors)"
else
    echo "✅ Limited os.Exit usage ($EXIT_COUNT calls)"
fi
echo ""

# 6. Check for SQL injection vulnerabilities
echo "6️⃣  Checking for potential SQL injection..."
if grep -r "fmt.Sprintf.*SELECT\|fmt.Sprintf.*INSERT\|fmt.Sprintf.*UPDATE\|fmt.Sprintf.*DELETE" --include="*.go" . 2>/dev/null; then
    echo "⚠️  Found potential SQL injection vulnerability"
else
    echo "✅ No obvious SQL injection patterns found"
fi
echo ""

# 7. Check for command injection
echo "7️⃣  Checking for potential command injection..."
if grep -r "exec.Command.*fmt.Sprintf\|exec.Command.*+" --include="*.go" . 2>/dev/null | grep -v "_test.go"; then
    echo "⚠️  Found potential command injection vulnerability"
else
    echo "✅ No obvious command injection patterns found"
fi
echo ""

# 8. Verify HTTPS usage
echo "8️⃣  Verifying HTTPS usage..."
if grep -r "http://" --include="*.go" . 2>/dev/null | grep -v "localhost\|127.0.0.1" | grep -v "_test.go" | grep -v "// "; then
    echo "⚠️  Found HTTP URLs (should use HTTPS in production)"
else
    echo "✅ All external URLs use HTTPS"
fi
echo ""

# 9. Check crypto usage
echo "9️⃣  Verifying cryptographic implementations..."
if grep -r "crypto/rand" --include="*.go" internal/auth/; then
    echo "✅ Using crypto/rand for secure random generation"
else
    echo "⚠️  Not using crypto/rand (should use for security-sensitive operations)"
fi

if grep -r "sha256\|sha512" --include="*.go" internal/auth/; then
    echo "✅ Using secure hash functions (SHA-256/SHA-512)"
else
    echo "⚠️  No secure hash functions found"
fi
echo ""

# 10. Check dependencies
echo "🔟 Checking for outdated dependencies..."
go list -m -u all | grep "\[" || echo "✅ All dependencies are up to date"
echo ""

echo "=========================================="
echo "🎉 Security audit complete!"
echo ""
echo "Next steps:"
echo "  - Review any warnings above"
echo "  - Run 'make lint' for code quality checks"
echo "  - Run 'make test' to ensure tests pass"
echo "  - Update dependencies: 'go get -u ./...'"
