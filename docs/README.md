# PipeOps CLI Documentation

This directory contains the complete documentation for PipeOps CLI, built with MkDocs and Material for MkDocs.

## Structure

```
docs/
├── index.md                    # Homepage
├── getting-started/            # Getting started guides
│   ├── installation.md         # Installation guide
│   ├── quick-start.md          # Quick start guide
│   └── configuration.md        # Configuration guide
├── commands/                   # Command documentation
│   ├── overview.md            # Command overview
│   ├── auth.md                # Authentication commands
│   ├── projects.md            # Project commands
│   ├── deployments.md         # Deployment commands
│   ├── servers.md             # Server commands
│   ├── agents.md              # Agent commands
│   └── k3s.md                 # K3s commands
├── advanced/                   # Advanced topics
│   ├── docker.md              # Docker usage
│   ├── ci-cd.md               # CI/CD integration
│   └── troubleshooting.md     # Troubleshooting
├── development/                # Development guides
│   ├── contributing.md        # Contributing guide
│   ├── building.md            # Building from source
│   └── api.md                 # API reference
└── reference/                  # Reference materials
    ├── changelog.md           # Changelog
    └── license.md             # License
```

## Building Documentation

### Prerequisites

- Python 3.7 or later
- pip3

### Quick Start

```bash
# Build documentation
./scripts/build-docs.sh build

# Serve documentation locally
./scripts/build-docs.sh serve

# Clean build artifacts
./scripts/build-docs.sh clean
```

### Manual Build

```bash
# Install dependencies
pip3 install mkdocs-material
pip3 install mkdocs-git-revision-date-localized-plugin
pip3 install mkdocs-minify-plugin

# Build documentation
mkdocs build

# Serve locally
mkdocs serve
```

## Writing Documentation

### Markdown Guidelines

- Use clear, concise language
- Include code examples
- Add emojis for visual appeal
- Use admonitions for important information
- Include links to related topics

### Code Examples

```bash
# Always include the command prompt
pipeops auth login

# Show expected output
 Successfully authenticated as user@example.com
```

### Admonitions

Use admonitions for important information:

!!! tip "Pro Tip"
    This is a helpful tip for users.

!!! warning "Warning"
    This is a warning about potential issues.

!!! error "Error"
    This indicates an error condition.

### Links

- Use relative links for internal documentation
- Use absolute links for external resources
- Include descriptive link text

## Configuration

The documentation is configured in `mkdocs.yml`:

- **Theme**: Material for MkDocs
- **Plugins**: Search, Git revision dates, Minification
- **Extensions**: Various Markdown extensions for enhanced features

## Deployment

### GitHub Pages

The documentation is automatically deployed to GitHub Pages when changes are pushed to the main branch.

### Manual Deployment

```bash
# Build documentation
mkdocs build

# Deploy to GitHub Pages
mkdocs gh-deploy
```

## Analytics

The documentation includes Google Analytics integration. Update the `G-XXXXXXXXXX` placeholder in `mkdocs.yml` with your Google Analytics ID.

## Search

The documentation includes a built-in search feature powered by MkDocs Material's search plugin.

## Customization

### Theme Customization

The theme can be customized in `mkdocs.yml`:

```yaml
theme:
  name: material
  palette:
    # Light mode
    - media: "(prefers-color-scheme: light)"
      scheme: default
      primary: indigo
      accent: indigo
    # Dark mode
    - media: "(prefers-color-scheme: dark)"
      scheme: slate
      primary: indigo
      accent: indigo
```

### Navigation

Navigation is configured in the `nav` section of `mkdocs.yml`.

## Troubleshooting

### Common Issues

#### Build Failures

```bash
# Check Python version
python3 --version

# Update dependencies
pip3 install --upgrade mkdocs-material

# Clean and rebuild
rm -rf site
mkdocs build
```

#### Missing Dependencies

```bash
# Install all dependencies
pip3 install mkdocs-material mkdocs-git-revision-date-localized-plugin mkdocs-minify-plugin
```

#### Link Validation

```bash
# Install linkchecker
pip3 install linkchecker

# Check links
linkchecker http://127.0.0.1:8000
```

## Resources

- [MkDocs Documentation](https://www.mkdocs.org/)
- [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/)
- [MkDocs Plugins](https://github.com/mkdocs/mkdocs/wiki/MkDocs-Plugins)

## Contributing

When contributing to the documentation:

1. Follow the existing structure and style
2. Include code examples
3. Test your changes locally
4. Update the navigation if adding new pages
5. Submit a pull request

## License

The documentation is licensed under the same license as the main project (MIT License).
