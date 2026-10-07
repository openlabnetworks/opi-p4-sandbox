# 🚀 OPI & P4 Developer Sandbox

[![Build & Publish](https://github.com/openlabnetworks/opi-p4-sandbox/actions/workflows/build-publish.yml/badge.svg)](https://github.com/openlabnetworks/opi-p4-sandbox/actions/workflows/build-publish.yml)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/openlabnetworks/opi-p4-sandbox)
[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/openlabnetworks/opi-p4-sandbox)

**A 1-click sandbox for experimenting with [OPI (Open Programmable Infrastructure)](https://opiproject.org) APIs and [P4](https://p4.org) programmable dataplanes — no hardware required.**

---

## What is this?

This repository packages an **OPI API emulator** and a **P4 BMv2 software switch** into a pre-built, containerized sandbox that any developer can spin up in under 3 minutes. It lets you:

- 🔌 **Test OPI APIs** — Create and manage virtual network ports and P4 pipelines
- 📝 **Write P4 programs** — Compile and load custom P4 dataplane programs
- 🧪 **Validate architectures** — Test your datapath offloading design against the OPI standard
- ☁️ **Run anywhere** — Locally with Docker, or in the cloud with GitHub Codespaces

## Architecture

```
┌──────────────── Docker Compose ────────────────────┐
│                                                     │
│   ┌────────────────┐       ┌────────────────────┐   │
│   │  OPI Server    │       │  P4 BMv2 Switch    │   │
│   │  REST :8080    │◄─────►│  P4Runtime :9559   │   │
│   │                │       │  Thrift    :9090   │   │
│   │  • NetworkSvc  │       │  • simple_switch   │   │
│   │  • PipelineSvc │       │  • P4 compiler     │   │
│   └────────────────┘       └────────────────────┘   │
│                                                     │
└─────────────────────────────────────────────────────┘
         ▲                           ▲
         │ curl / grpcurl            │ P4Runtime
         │                           │
    ┌────┴───────────────────────────┴────┐
    │          Developer Machine          │
    └─────────────────────────────────────┘
```

---

## ⚡ Quick Start (3 Steps)

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) 24+ with Compose v2
- `curl` (pre-installed on most systems)

### Step 1: Clone

```bash
git clone https://github.com/openlabnetworks/opi-p4-sandbox.git
cd opi-p4-sandbox
```

### Step 2: Start the Sandbox

```bash
make start-sandbox
```

This pulls pre-built images and starts the OPI server + P4 switch. Takes ~30 seconds on first run.

### Step 3: Send Your First API Call

```bash
# List all network ports
curl http://localhost:8080/v1/ports | python3 -m json.tool
```

**Expected output:**
```json
{
    "ports": [
        {
            "name": "ports/port-0",
            "id": "port-0",
            "mac_address": "00:1a:2b:3c:4d:5e",
            "mtu": 1500,
            "admin_state": "UP",
            "oper_state": "UP",
            "speed": "100G"
        }
    ]
}
```

**🎉 You're in!** The OPI sandbox is running and ready for experiments.

---

## 🧪 Try More

### Create a Virtual Port

```bash
curl -X POST http://localhost:8080/v1/ports \
  -H "Content-Type: application/json" \
  -d '{"mac_address": "de:ad:be:ef:00:01", "mtu": 9000}' \
  | python3 -m json.tool
```

### List P4 Pipelines

```bash
curl http://localhost:8080/v1/pipelines | python3 -m json.tool
```

### Create a P4 Pipeline

```bash
curl -X POST http://localhost:8080/v1/pipelines \
  -H "Content-Type: application/json" \
  -d '{"p4_program": "basic_forwarding.p4"}' \
  | python3 -m json.tool
```

### Run Smoke Tests

```bash
make smoke-test
```

---

## 📚 Guided Experiments

| # | Experiment | Difficulty | Description |
|---|-----------|-----------|-------------|
| 01 | [Hello OPI](examples/01-hello-opi/) | Beginner | Your first OPI API interaction |
| 02 | [Custom P4](examples/02-custom-p4/) | Intermediate | Write, compile, and load a P4 program |
| 03 | [Pipeline Management](examples/03-opi-pipeline-mgmt/) | Intermediate | Full pipeline CRUD lifecycle |

---

## ☁️ Zero-Install Cloud Option

Don't want to install anything? Click a badge above to launch the sandbox in:

- **GitHub Codespaces** — Full VS Code in your browser with the sandbox auto-started
- **Gitpod** — Alternative cloud IDE with identical setup

The sandbox starts automatically when the workspace opens. Ports are forwarded to your browser.

---

## 🔧 API Reference

> [!NOTE]
> **REST vs gRPC:** The official OPI APIs are gRPC and protobuf-based. To make this sandbox as accessible as possible without requiring users to install `grpcurl` or compile protobufs, the sandbox server provides a simplified REST/JSON emulator that mirrors the structure and naming conventions of the real OPI APIs.

### Endpoints

| Method | Endpoint | Description |
|--------|---------|-------------|
| `GET` | `/healthz` | Server health check |
| `GET` | `/` | API info and version |
| `GET` | `/v1/ports` | List all ports |
| `POST` | `/v1/ports` | Create a new port |
| `GET` | `/v1/ports/{id}` | Get a specific port |
| `DELETE` | `/v1/ports/{id}` | Delete a port |
| `GET` | `/v1/pipelines` | List all pipelines |
| `POST` | `/v1/pipelines` | Create a new pipeline |
| `GET` | `/v1/pipelines/{id}` | Get a specific pipeline |
| `DELETE` | `/v1/pipelines/{id}` | Delete a pipeline |

### Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `OPI_PORT` | `8080` | OPI server REST port |
| `OPI_LOG_LEVEL` | `info` | Log level: debug, info, warn, error |
| `P4_GRPC_PORT` | `9559` | P4 BMv2 P4Runtime port |
| `P4_THRIFT_PORT` | `9090` | P4 BMv2 Thrift port |
| `P4_LOG_LEVEL` | `info` | BMv2 log level |

---

## 🛠 Makefile Commands

```bash
make help              # Show all available commands
make start-sandbox     # Start the full sandbox
make stop-sandbox      # Tear everything down
make restart-sandbox   # Restart
make status            # Check container health
make logs              # Tail container logs
make smoke-test        # Run end-to-end tests
make build-images      # Build images locally
make compile-p4 P4=X   # Compile a P4 program
make clean             # Remove all containers and images
make install-tools     # Check development tool availability
```

---

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## 📖 Documentation

- [Architecture](docs/architecture.md) — System design and component overview
- [Experiments](examples/) — Guided hands-on tutorials

## 📄 License

[Apache License 2.0](LICENSE) — Same as the [OPI Project](https://opiproject.org).

---

## 🔗 References

- [OPI Project](https://opiproject.org) — Open Programmable Infrastructure standard
- [P4 Language](https://p4.org) — Protocol-independent packet processing
- [BMv2](https://github.com/p4lang/behavioral-model) — P4 software switch
- [P4Runtime](https://p4.org/p4-spec/p4runtime/v1/P4Runtime-Spec.html) — P4 control plane API
