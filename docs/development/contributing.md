# Contributing

We welcome contributions to PipeOps CLI! This guide will help you get started with contributing to the project.

## 🤝 How to Contribute

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/pipeops-cli.git
cd pipeops-cli

# Add upstream remote
git remote add upstream https://github.com/PipeOpsHQ/pipeops-cli.git
```

### 2. Set Up Development Environment

```bash
# Install Go (1.23 or later)
# Follow instructions at https://golang.org/doc/install

# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test

# Run linter
make lint
```

### 3. Create a Branch

```bash
# Create a new branch for your feature
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/your-bug-fix
```

### 4. Make Changes

- Write clear, commented code
- Follow Go best practices and conventions
- Add tests for new functionality
- Update documentation as needed

### 5. Test Your Changes

```bash
# Run all tests
make test

# Run specific tests
go test ./cmd/agent/...

# Run linter
make lint

# Build and test the binary
make build
./pipeops-cli --help
```

### 6. Commit and Push

```bash
# Add your changes
git add .

# Commit with a descriptive message
git commit -m "feat: add new agent installation feature"

# Push to your fork
git push origin feature/your-feature-name
```

### 7. Create a Pull Request

- Go to your fork on GitHub
- Click "New Pull Request"
- Fill out the pull request template
- Submit the pull request

## 📋 Development Guidelines

### Code Style

- Follow Go conventions and best practices
- Use `gofmt` to format code
- Use `golint` for linting
- Write clear, self-documenting code

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
feat: add new feature
fix: fix a bug
docs: update documentation
style: formatting changes
refactor: code refactoring
test: add or update tests
chore: maintenance tasks
```

### Testing

- Write unit tests for new functionality
- Test edge cases and error conditions
- Ensure tests pass before submitting PR
- Aim for good test coverage

### Documentation

- Update documentation for new features
- Add examples and usage instructions
- Update command help text
- Follow the existing documentation style

## 🏗️ Project Structure

```
pipeops-cli/
├── cmd/                 # CLI commands
│   ├── auth/           # Authentication commands
│   ├── project/        # Project management
│   ├── deploy/         # Deployment commands
│   ├── agent/          # Agent management
│   └── ...
├── internal/           # Internal packages
│   ├── auth/           # Authentication logic
│   ├── client/         # HTTP client
│   ├── config/         # Configuration
│   └── ...
├── models/             # Data models
├── utils/              # Utility functions
├── docs/               # Documentation
├── scripts/            # Build and utility scripts
└── .github/            # GitHub workflows
```

## 🔧 Available Make Targets

```bash
make build          # Build the binary
make test           # Run tests
make lint           # Run linter
make clean          # Clean build artifacts
make install        # Install locally
make release        # Create release build
make docker-build   # Build Docker image
make docker-run     # Run in Docker
make docs-build     # Build documentation
make docs-serve     # Serve documentation locally
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./cmd/agent/...

# Run tests with verbose output
go test -v ./...
```

### Test Structure

- Unit tests go in `*_test.go` files
- Test files should be in the same package
- Use table-driven tests when appropriate
- Mock external dependencies

### Example Test

```go
func TestInstallCommand(t *testing.T) {
    tests := []struct {
        name     string
        args     []string
        expected string
    }{
        {
            name:     "valid token",
            args:     []string{"valid-token"},
            expected: "success",
        },
        {
            name:     "empty token",
            args:     []string{""},
            expected: "error",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 📝 Documentation

### Writing Documentation

- Use clear, concise language
- Include code examples
- Add emojis for visual appeal
- Use admonitions for important information

### Documentation Structure

- Follow the existing structure in `docs/`
- Update navigation in `mkdocs.yml`
- Include links to related topics
- Test documentation locally

### Building Documentation

```bash
# Build documentation
make docs-build

# Serve documentation locally
make docs-serve
```

## 🐛 Bug Reports

When reporting bugs, please include:

- Description of the issue
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, etc.)
- Relevant logs or error messages

## 💡 Feature Requests

When requesting features, please include:

- Description of the feature
- Use case and motivation
- Proposed implementation (if you have ideas)
- Any alternatives considered

## 🔒 Security Issues

For security issues, please:

- **DO NOT** open a public issue
- Email security@pipeops.io
- Include detailed information about the vulnerability
- Allow time for the team to respond

## 📞 Getting Help

If you need help contributing:

- Check existing issues and discussions
- Join our [Discord community](https://discord.gg/pipeops)
- Email us at [support@pipeops.io](mailto:support@pipeops.io)
- Open a discussion on GitHub

## 🎉 Recognition

Contributors will be:

- Listed in the project README
- Mentioned in release notes
- Invited to the contributors' Discord channel
- Given recognition in project documentation

## 📄 License

By contributing to PipeOps CLI, you agree that your contributions will be licensed under the [MIT License](LICENSE).

---

Thank you for contributing to PipeOps CLI! 🚀
