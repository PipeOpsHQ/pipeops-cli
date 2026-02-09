# PipeOps CLI Agent Guidelines

This document provides instructions and conventions for AI agents working on the PipeOps CLI repository.
Strictly adhere to these guidelines to ensure code quality and consistency.

## 1. Build, Lint, and Test

The project uses a `Makefile` to manage common tasks.

### Build
- **Build binary**: `make build` (Outputs to `build/bin/pipeops` - secure production build)
- **Dev build**: `make build-dev` (Faster, includes debug symbols)
- **Run**: `make run` or `./build/bin/pipeops [command]`
- **Docker**: `make docker-build` (Builds image) and `make docker-run` (Runs container).

### Test
- **Run all tests**: `make test`
- **Run tests with coverage**: `make test-coverage` (Generates HTML report in `build/coverage.html`)
- **Run a single test**:
  ```bash
  # Syntax: go test -v [package_path] -run [TestNameRegex]
  go test -v ./cmd -run TestShouldSkipUpdateCheck
  go test -v ./internal/config -run TestConfigFunctions
  ```
- **Benchmarks**: `make bench`

### Lint & Format
- **Lint**: `make lint` (Uses `golangci-lint`)
- **Format**: `make fmt` (Runs `go fmt ./...`)
- **Check all**: `make check` (Runs fmt, lint, and test)

## 2. Code Structure & Conventions

### Project Structure
- `cmd/`: Cobra command definitions.
  - `root.go`: Root command and persistent flags.
  - `[command].go`: Definition of top-level commands.
  - Commands often use a `register[Name]Subcommands` pattern.
- `internal/`: Private application code.
  - `auth/`: Authentication logic (PKCE, OAuth) and token management.
  - `client/`: HTTP client wrapper with automatic token refresh/auth injection.
  - `config/`: Configuration management via Viper (file: `~/.pipeops/config.yaml`).
  - `updater/`: Self-update logic using GitHub releases.
  - `validation/`: Input validation helpers (e.g., project names, URLs).
  - `terminal/`: Cross-platform terminal utilities (color support check).
- `main.go`: Entry point, calls `cmd.Execute()`.

### Code Style
- **Go Version**: 1.24.0+
- **Formatting**: Strictly follow `gofmt`. Run `make fmt` before committing.
- **Naming**:
  - `PascalCase` for exported types/functions.
  - `camelCase` for internal variables/functions.
  - Acronyms: `ID`, `HTTP`, `URL`, `JSON`, `API` (keep uppercase).
- **Imports**: Group in this order:
  1. Standard library
  2. Third-party packages
  3. Internal packages (`github.com/PipeOpsHQ/pipeops-cli/...`)

### Error Handling
- **Internal Code**:
  - Always check errors: `if err != nil { return fmt.Errorf("context: %w", err) }`.
  - Use `%w` to wrap errors.
- **CLI Commands**:
  - Return errors in `RunE` functions.
  - Do not `os.Exit` inside `internal/` packages.
  - Use `fmt.Fprintf(os.Stderr, ...)` for error messages if handling explicitly.

### CLI Implementation (Cobra)
- **Command Definition**:
  ```go
  var myCmd = &cobra.Command{
      Use:   "my-command",
      Short: "Short description",
      Long:  `Long description...`,
      RunE: func(cmd *cobra.Command, args []string) error {
          return runMyCommand(cmd.Context())
      },
  }
  ```
- **Flags**:
  - `PersistentFlags()`: Inherited by subcommands (e.g., `--json`, `--verbose`).
  - `Flags()`: Command-specific.
- **Context**: Use `cmd.Context()` to propagate context (cancellation, timeouts).
- **Output**:
  - Respect `--json` flag.
  - Use `github.com/fatih/color` for terminal output (e.g., `color.Green(...)`).
  - Use `github.com/olekukonko/tablewriter` for tabular data.
  - Use `github.com/briandowns/spinner` for long-running operations.

### HTTP Client Usage (Internal)
- **Authenticated Client**:
  ```go
  cfg, err := config.Load()
  // Client handles OAuth token injection and refresh automatically
  cli, err := client.NewAuthenticatedClient(cfg)
  if err != nil { ... }
  
  resp, err := cli.Get(ctx, "/api/v1/projects")
  ```
- **Error Handling**: The client wraps errors. Always check them.

### Configuration
- Defined in `internal/config/config.go`.
- Uses `viper` for loading/saving.
- Do not hardcode secrets or URLs. Use build-time variables (ldflags) for defaults (see `Makefile`).

## 3. Testing Guidelines
- **Unit Tests**: Co-locate with source code (e.g., `package_test.go`).
- **Integration Tests**: Use `t.TempDir()` for file system tests.
- **Mocking**:
  - Mock external services (API calls).
  - Mock environment variables:
    ```go
    // Use t.Setenv for Go 1.17+ (preferred over manual restore)
    t.Setenv("VAR_NAME", "value") 
    ```
- **Table-Driven Tests**: Preferred for logic with multiple cases.

## 4. Common Tasks for Agents

### Adding a New Command
1. Create `cmd/[name].go`.
2. Define `var [name]Cmd = &cobra.Command{...}`.
3. In `init()`, add `rootCmd.AddCommand([name]Cmd)`.
4. If it has subcommands, create a `register[Name]Subcommands` function.
5. Add tests in `cmd/[name]_test.go`.

### Adding a Configuration Field
1. Update `Config` struct in `internal/config/config.go`.
2. Update `Load` and `Save` logic if necessary (Viper usually handles struct tags).
3. Add a test case in `internal/config/config_test.go`.

### Debugging
- Use `make build-dev` to build with debug symbols.
- Use `fmt.Printf` temporarily but remove before committing.
- Check logs if the application writes to any log files (default is stdout/stderr).

## 5. Workflow
1. **Explore**: Use `ls -F`, `grep`, `glob` to locate relevant files.
2. **Plan**: Outline changes. Check for existing tests or conventions.
3. **Implement**:
   - Modify/Add files.
   - Run `make fmt` frequently.
4. **Verify**:
   - Run specific tests: `go test -v ./path/to/pkg`.
   - Run linter: `make lint`.
   - Build: `make build`.
   - Manually run the command if applicable: `./build/bin/pipeops [cmd]`.
