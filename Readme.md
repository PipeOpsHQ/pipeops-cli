# PipeOps CLI 🚀

[![Release](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/release.yml/badge.svg)](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/release.yml)
[![CodeQL Analysis](https://github.com/PipeOpsHQ/pipeops-cli/actions/workflows/code-analysis.yml/badge.svg)](https://github.com/PipeOpsHQ/pipeops-cli/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

PipeOps CLI is a powerful command-line interface designed to simplify managing cloud-native environments, deploying projects, and interacting with the PipeOps platform. With PipeOps CLI, you can provision servers, deploy applications, manage projects, and monitor your infrastructure seamlessly.

---

## Features ✨

- **Server Management**: Provision and configure servers across multiple environments (e.g., K3s, EKS, GKE).
- **Project Deployment**: Deploy your applications directly to servers with ease.
- **Pipeline Management**: Create, manage, and deploy CI/CD pipelines.
- **Agent Setup**: Install and configure PipeOps agents for various platforms.
- **Authentication**: Securely log in and manage your PipeOps account.
- **Cross-Platform Support**: Available for Linux, Windows, and macOS.

---

## Installation 📦

Download the latest version from the [Releases](https://github.com/PipeOpsHQ/pipeops-cli/releases) page or use the following commands:

### Linux / macOS
```bash
curl -LO https://github.com/PipeOpsHQ/pipeops-cli/releases/latest/download/pipeops-cli-linux-amd64
chmod +x pipeops-cli-linux-amd64
sudo mv pipeops-cli-linux-amd64 /usr/local/bin/pipeops
```

### Windows
1. Download `pipeops-cli-windows-amd64.exe` from the [Releases](https://github.com/PipeOpsHQ/pipeops-cli/releases).
2. Add it to your `PATH` for easier access.

---

## Quick Start 🚀

### Authenticate with PipeOps
```bash
pipeops auth login
```

### Deploy a Server
```bash
pipeops server deploy --name my-server --region us-east
```

### Deploy a Project to a Server
```bash
pipeops project deploy --name my-project --server my-server
```

### List Projects
```bash
pipeops project list
```

For a complete list of commands, use:
```bash
pipeops help
```

---

## Commands Overview 📖

| Command               | Description                                      |
|-----------------------|--------------------------------------------------|
| `pipeops auth`        | Manage authentication and user details.          |
| `pipeops server`      | Manage server-related operations (e.g., K3s).    |
| `pipeops project`     | Manage, list, and deploy PipeOps projects.       |
| `pipeops deploy`      | Manage and deploy CI/CD pipelines.               |

---

## Development 🛠️

### Prerequisites
- [Go](https://golang.org/) 1.20 or later
- [Git](https://git-scm.com/)

### Clone the Repository
```bash
git clone https://github.com/PipeOpsHQ/pipeops-cli.git
cd pipeops-cli
```

### Build the CLI
```bash
go build -o pipeops
```

### Run Tests
```bash
go test ./...
```

---

## Contributing 🤝

We welcome contributions! To contribute:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature/my-feature`).
3. Commit your changes (`git commit -m 'Add my feature'`).
4. Push to the branch (`git push origin feature/my-feature`).
5. Open a Pull Request.

Please follow our [Code of Conduct](CODE_OF_CONDUCT.md) and review our [Contributing Guidelines](CONTRIBUTING.md).

---

## License 📜

This project is licensed under the [MIT License](LICENSE).

---

## Support 💬

For questions or support, please open an [issue](https://github.com/PipeOpsHQ/pipeops-cli/issues) or contact us at [support@pipeops.io](mailto:support@pipeops.io).

---

## Acknowledgments 🙌

Special thanks to all contributors and users of PipeOps CLI! 🎉
