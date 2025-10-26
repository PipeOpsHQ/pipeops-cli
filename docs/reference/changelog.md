# Changelog

All notable changes to PipeOps CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Enhanced agent installation with intelligent cluster detection
- Support for multiple Kubernetes distributions (k3s, minikube, k3d, kind)
- Monitoring stack integration (Prometheus, Loki, Grafana, OpenCost)
- Worker node joining functionality
- Comprehensive documentation with MkDocs
- GitHub Actions CI/CD for documentation publishing
- Docker support for documentation building

### Changed
- Improved agent installation workflow
- Enhanced error handling and user feedback
- Updated command structure and help text

### Fixed
- Fixed environment variable handling in agent commands
- Improved token validation and error messages

## [1.0.0] - 2024-01-01

### Added
- Initial release of PipeOps CLI
- Authentication system with OAuth and PKCE flow
- Project management commands
- Deployment and pipeline management
- Server management capabilities
- K3s cluster management
- Cross-platform support (Linux, macOS, Windows, FreeBSD)
- Docker container support
- Comprehensive command-line interface
- Configuration management
- Update system
- Proxy functionality
- Status monitoring
- Logging and debugging capabilities

### Features
- **Authentication**: Secure OAuth-based authentication
- **Projects**: Create, manage, and deploy projects
- **Deployments**: Manage CI/CD pipelines and deployments
- **Servers**: Provision and configure servers
- **K3s**: Install and manage K3s clusters
- **Agents**: Install and configure PipeOps agents
- **Cross-Platform**: Support for multiple operating systems
- **Docker**: Containerized deployment support
- **Configuration**: Flexible configuration options
- **Updates**: Automatic update checking and installation

## [0.9.0] - 2023-12-15

### Added
- Beta release with core functionality
- Basic authentication system
- Project management commands
- Initial deployment capabilities
- Server management features
- K3s integration
- Cross-platform builds

### Changed
- Improved command structure
- Enhanced error handling
- Better user experience

### Fixed
- Various bug fixes and improvements
- Performance optimizations

## [0.8.0] - 2023-12-01

### Added
- Alpha release
- Basic CLI structure
- Authentication framework
- Project management foundation
- Initial deployment system

### Changed
- Core architecture improvements
- Command structure updates

### Fixed
- Initial bug fixes
- Stability improvements

## [0.7.0] - 2023-11-15

### Added
- Pre-alpha release
- Basic command structure
- Authentication system
- Project management
- Initial deployment capabilities

### Changed
- Core functionality implementation
- Command structure design

### Fixed
- Initial implementation fixes

## [0.6.0] - 2023-11-01

### Added
- Development release
- Core CLI framework
- Basic authentication
- Project management foundation

### Changed
- Architecture improvements
- Command structure updates

### Fixed
- Development phase fixes

## [0.5.0] - 2023-10-15

### Added
- Early development release
- Basic CLI structure
- Authentication framework
- Project management foundation

### Changed
- Core architecture design
- Command structure implementation

### Fixed
- Development phase improvements

## [0.4.0] - 2023-10-01

### Added
- Initial development release
- Basic CLI framework
- Authentication system
- Project management foundation

### Changed
- Core architecture implementation
- Command structure design

### Fixed
- Initial development fixes

## [0.3.0] - 2023-09-15

### Added
- Pre-development release
- Basic CLI structure
- Authentication framework
- Project management foundation

### Changed
- Core architecture design
- Command structure implementation

### Fixed
- Pre-development improvements

## [0.2.0] - 2023-09-01

### Added
- Early development release
- Basic CLI framework
- Authentication system
- Project management foundation

### Changed
- Core architecture implementation
- Command structure design

### Fixed
- Early development fixes

## [0.1.0] - 2023-08-15

### Added
- Initial development release
- Basic CLI structure
- Authentication framework
- Project management foundation

### Changed
- Core architecture design
- Command structure implementation

### Fixed
- Initial development improvements

## [0.0.1] - 2023-08-01

### Added
- First release
- Basic CLI framework
- Authentication system
- Project management foundation

### Changed
- Core architecture implementation
- Command structure design

### Fixed
- Initial implementation fixes

---

## Version Numbering

We use [Semantic Versioning](https://semver.org/) for version numbering:

- **MAJOR** version when you make incompatible API changes
- **MINOR** version when you add functionality in a backwards compatible manner
- **PATCH** version when you make backwards compatible bug fixes

## Release Schedule

- **Major releases**: Every 6 months
- **Minor releases**: Every month
- **Patch releases**: As needed for bug fixes
- **Security releases**: As needed for security fixes

## Support Policy

- **Current version**: Full support
- **Previous major version**: Security fixes only
- **Older versions**: No support

## Migration Guide

For major version changes, we provide migration guides:

- [v1.0.0 Migration Guide](migration/v1.0.0.md)
- [v0.9.0 Migration Guide](migration/v0.9.0.md)

## Breaking Changes

Breaking changes are documented in the [Breaking Changes](breaking-changes.md) file.

## Deprecation Policy

- Features are deprecated for at least one minor version before removal
- Deprecation warnings are shown in the CLI
- Deprecated features are documented in the changelog

## Security Updates

Security updates are released as needed and are documented in the [Security Advisories](security-advisories.md) file.

---

*This changelog is maintained by the PipeOps team and is updated with each release.*
