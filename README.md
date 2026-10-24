# 🚀 OpenLanNetworks Sandbox

[![Build & Publish](https://github.com/openlabnetworks/opi-p4-sandbox/actions/workflows/build-publish.yml/badge.svg)](https://github.com/openlabnetworks/opi-p4-sandbox/actions/workflows/build-publish.yml)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/openlabnetworks/opi-p4-sandbox)

**A 1-click sandbox for experimenting with hardware-assisted networking workflows — no DPU or SmartNIC hardware required.**

---

## What is this?

This repository packages an **OPI API emulator** and a **P4 BMv2 software switch** into a pre-built, containerized sandbox that any developer can spin up in under 3 minutes. It lets you:

- 🔌 **Test OPI APIs** — Create and manage virtual network ports and P4 pipelines
- 📝 **Write P4 programs** — Compile custom P4 dataplane programs for BMv2
- 🧪 **Validate architectures** — Test your datapath offloading design against the OPI standard
- ☁️ **Run anywhere** — Locally with Docker, or in the cloud with GitHub Codespaces

## Architecture

```
┌──────────────── Docker Compose ────────────────────┐
│                                                     │
│   ┌────────────────┐       ┌────────────────────┐   │
│   │  OPI Server    │       │  P4 BMv2 Switch    │   │
│   │  REST :8080    │       │  P4Runtime :9559   │   │
│   │                │       │  Thrift    :9090   │   │
│   │  • NetworkSvc  │       │  • simple_switch   │   │
│   │  • PipelineSvc │       │  • External control│   │
│   └────────────────┘       └────────────────────┘   │
│                                                     │
└─────────────────────────────────────────────────────┘
         ▲                           ▲
         │ curl                 P4Runtime client
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

**🎉 You're in!** OpenLanNetworks Sandbox is running and ready for experiments.

---

## 🌐 Interactive Web Testing Dashboard

Prefer a graphical interface? The sandbox includes a built-in, lightweight **Web Testing Dashboard** served directly by the OPI server.

**Open in your browser:** [`http://localhost:8080/ui`](http://localhost:8080/ui) *(or `http://localhost:8080/`)*

Learn about the platform's use cases and architecture at [`http://localhost:8080/about`](http://localhost:8080/about).

### Key Features

- 🎨 **Minimalist Light & Brutalist Design**: High-contrast, Swiss editorial typography, architectural grid canvas, and cobalt hardware passport card.
- 🔌 **Live Port Manager**: View active ports, create new virtual SmartNIC ports with 1-click random MAC generation, and delete ports in real-time.
- 📦 **P4 Pipeline Registry**: Inspect registered pipeline metadata and register P4 programs (`basic_forwarding.p4`, `firewall_drop.p4`).
- ⚡ **Interactive 5-Stage Dataplane Packet Simulator**: Send simulated packets through Parser ➔ Verify Checksum ➔ Ingress Match-Action ➔ Egress ➔ Deparser with animated visual traversal and instant PASS / DROP security verdicts.
- 💻 **Interactive REST API Console**: Fire test requests (`GET /healthz`, `GET /v1/ports`, `POST /v1/ports`) and copy generated `curl` commands.
- 🧪 **One-Click Browser Smoke Test**: Run the complete automated test suite with a single click.

---

## 🧪 Try More

### Create a Virtual Port

```bash
curl -X POST http://localhost:8080/v1/ports \
  -H "Content-Type: application/json" \
  -d '{"mac_address": "de:ad:be:ef:00:01", "mtu": 9000}' \
  | python3 -m json.tool
```

### List and Register P4 Pipelines

```bash
# List all active pipelines
curl http://localhost:8080/v1/pipelines | python3 -m json.tool

# Register a custom pipeline
curl -X POST http://localhost:8080/v1/pipelines \
  -H "Content-Type: application/json" \
  -d '{"p4_program": "firewall_drop.p4"}' \
  | python3 -m json.tool
```

### Compile and Register a P4 Program

The sandbox includes sample P4_16 programs in `src/p4-programs/`:
- `basic_forwarding.p4` — L2/L3 packet routing with IPv4 LPM table lookup and TTL modification
- `firewall_drop.p4` — Security filter that drops all non-IPv4 traffic

```bash
# Compile any P4 program with the containerized p4c compiler
make compile-p4 P4=firewall_drop.p4

# Or compile and register its metadata with the OPI API
bash scripts/load-p4-pipeline.sh firewall_drop.p4
```

> [!IMPORTANT]
> The REST pipeline endpoints are an in-memory metadata emulator and do not load a pipeline into BMv2. Use a P4Runtime client directly against `localhost:9559` to install a compiled program on the switch.

### P4Runtime integration status

The OPI server is not yet a P4Runtime client. Creating or deleting a pipeline
through the REST API changes only the server's in-memory registry; it does not
call BMv2's `SetForwardingPipelineConfig` or remove a configuration from the
switch. Real P4Runtime-backed pipeline loading is planned and tracked in
[issue #6](https://github.com/openlabnetworks/opi-p4-sandbox/issues/6).

### Run Automated Tests

```bash
# 1. Quick smoke test (5 core assertions)
make smoke-test

# 2. Comprehensive 11-step integration test suite
bash scripts/smoke-test.sh

# 3. OPI Server Go unit tests (95% statement coverage)
cd src/opi-server && go test -v ./...
```

---

## 📚 Guided Experiments

| # | Experiment | Difficulty | Description |
|---|-----------|-----------|-------------|
| 01 | [Hello OPI](examples/01-hello-opi/) | Beginner | Your first OPI API interaction |
| 02 | [Custom P4](examples/02-custom-p4/) | Intermediate | Write, compile, and register a P4 program |
| 03 | [Pipeline Management](examples/03-opi-pipeline-mgmt/) | Intermediate | Full pipeline CRUD lifecycle |

---

## ☁️ Zero-Install Cloud Option

Don't want to install anything? Click a badge above to launch the sandbox in:

- **GitHub Codespaces** — Full VS Code in your browser with the sandbox auto-started

The sandbox starts automatically when the workspace opens. Ports are forwarded to your browser.

---

## ☁️ AWS AMI Deployment (Option 1: Pre-Baked AMI with Packer)

For cloud deployment on AWS EC2, the repository includes a complete **HashiCorp Packer** automation template that creates a pre-baked Amazon Machine Image (AMI).

### Why Pre-Baked AMI?
- ⚡ **Zero-Download Boot Time**: Docker, BMv2, and OPI images are 100% pre-built and cached inside the AMI filesystem. Instances boot and become healthy in seconds without downloading layers from the internet.
- 🔄 **Auto-Start on Boot**: Uses systemd (`/etc/systemd/system/opi-sandbox.service`) to start `docker compose up -d` immediately upon EC2 launch.
- 🌐 **Instant Cloud Web Dashboard**: Access `http://<ec2-public-ip>:8080/ui` as soon as the instance launches.

### How to Build the AMI

```bash
# 1. Initialize Packer AWS builder plugins
make ami-init
# or: packer init packer/opi-sandbox.pkr.hcl

# 2. Validate template
make ami-validate

# 3. Build the AMI (requires AWS credentials configured)
make ami-build
# or: packer build -var="aws_region=us-east-1" packer/opi-sandbox.pkr.hcl
```

### Launching the Instance
1. Open the **AWS EC2 Console** ➔ **AMIs** (select "Owned by me").
2. Select `opi-p4-sandbox-v1.0.0` and click **Launch Instance**.
3. Choose instance type `t3.medium` or `t3.large`.
4. Ensure Security Group allows inbound ports:
   - `8080` (OPI REST API & Web Dashboard)
   - `9559` (P4Runtime gRPC)
   - `9090` (BMv2 Thrift CLI)
5. Click **Launch**. Within 60 seconds, visit `http://<ec2-public-ip>:8080/ui`!

---

## 🔧 API Reference

> [!NOTE]
> **REST vs gRPC:** The official OPI APIs are gRPC and protobuf-based. To make this sandbox as accessible as possible without requiring users to install `grpcurl` or compile protobufs, the sandbox server provides a simplified REST/JSON emulator that mirrors the structure and naming conventions of the real OPI APIs. It does not proxy REST pipeline requests to BMv2; program installation is performed separately through P4Runtime.

### Endpoints

| Method | Endpoint | Description |
|--------|---------|-------------|
| `GET` | `/ui` | Interactive Web Testing Dashboard |
| `GET` | `/healthz` | Server health check |
| `GET` | `/` | API info (returns Web UI for browsers, JSON for API/curl) |
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
