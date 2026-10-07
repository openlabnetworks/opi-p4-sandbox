# Contributing to OPI & P4 Sandbox

Thank you for your interest in contributing! This guide will help you get set up.

## Development Setup

### Prerequisites

- Docker 24+ with Compose v2
- Go 1.22+ (for OPI server development)
- Make

### Clone and Build Locally

```bash
git clone https://github.com/YOUR_ORG/opi-p4-sandbox.git
cd opi-p4-sandbox

# Build images from source
make build-images

# Start with locally built images
make start-sandbox

# Run tests
make smoke-test
```

### Project Structure

```
opi-p4-sandbox/
├── .github/workflows/     # CI/CD pipelines
├── .devcontainer/          # Codespaces / VS Code config
├── docker/
│   ├── opi-server/         # OPI server Dockerfile
│   └── p4-bmv2/            # P4 BMv2 Dockerfile
├── src/
│   ├── opi-server/         # Go source code
│   └── p4-programs/        # Sample P4 programs
├── scripts/                # Utility scripts
├── examples/               # Guided experiments
├── docs/                   # Documentation
├── docker-compose.yml      # Container orchestration
└── Makefile                # Developer interface
```

## Making Changes

### OPI Server (Go)

1. Edit files in `src/opi-server/`
2. Test locally:
   ```bash
   cd src/opi-server
   go test ./...
   go run . --port 8080
   ```
3. Rebuild the Docker image:
   ```bash
   make build-images
   make restart-sandbox
   ```

### P4 Programs

1. Write your P4 program in `src/p4-programs/`
2. Compile it:
   ```bash
   make compile-p4 P4=your_program.p4
   ```

### Documentation

- API docs go in `docs/`
- Experiments go in `examples/XX-name/README.md`

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Ensure tests pass: `make smoke-test`
5. Submit a pull request

## Code Style

- **Go:** Follow standard `gofmt` formatting
- **P4:** Use 4-space indentation, descriptive names
- **Shell:** Pass `shellcheck` validation
- **Markdown:** One sentence per line

## Reporting Issues

Please include:
- OS and Docker version
- Output of `make status`
- Container logs: `make logs`
- Steps to reproduce
