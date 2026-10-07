# Architecture

## System Overview

The OPI & P4 Developer Sandbox provides a complete, containerized environment for
experimenting with the Open Programmable Infrastructure (OPI) standard and P4
programmable dataplanes — without requiring physical SmartNIC or DPU hardware.

## Components

```
┌─────────────────────────────────────────────────────────┐
│                  Developer Machine                       │
│                                                         │
│  ┌─── Docker Compose Network (opi-net) ──────────────┐  │
│  │                                                    │  │
│  │   ┌──────────────────┐   ┌──────────────────────┐  │  │
│  │   │  OPI API Server  │   │  P4 BMv2 Switch      │  │  │
│  │   │                  │   │                      │  │  │
│  │   │  REST :8080      │◄─►│  P4Runtime :9559     │  │  │
│  │   │                  │   │  Thrift    :9090     │  │  │
│  │   │  Services:       │   │                      │  │  │
│  │   │  - NetworkSvc    │   │  Features:           │  │  │
│  │   │  - PipelineSvc   │   │  - simple_switch_grpc│  │  │
│  │   │                  │   │  - P4 compiler (p4c) │  │  │
│  │   │  Go binary       │   │  - Virtual interfaces│  │  │
│  │   │  ~30MB image     │   │  - Pipeline loading  │  │  │
│  │   └──────────────────┘   └──────────────────────┘  │  │
│  │                                                    │  │
│  └────────────────────────────────────────────────────┘  │
│                                                         │
│  Developer tools: curl, grpcurl, python3, go            │
└─────────────────────────────────────────────────────────┘
```

## OPI API Server

The OPI API Server is a lightweight Go application that emulates OPI-standard
APIs for infrastructure offloading. It provides:

- **NetworkService** — CRUD operations for virtual network ports
- **PipelineService** — CRUD operations for P4 pipeline configurations

The server uses in-memory storage, so data resets when the container restarts.
This is intentional for a sandbox — it ensures a clean state for each experiment.

### API Design

The API follows OPI resource naming conventions:
- Resources are named with a hierarchical pattern: `{collection}/{resource_id}`
- Example: `ports/port-0`, `pipelines/pipeline-default`
- Standard CRUD operations: List, Get, Create, Delete

## P4 BMv2 Software Switch

The P4 BMv2 (Behavioral Model v2) is a reference P4 software switch developed by
the P4 Language Consortium. In this sandbox it serves as:

- **A P4Runtime target** — Accepts P4Runtime gRPC calls for pipeline management
- **A packet processing engine** — Executes P4 programs on software-emulated ports
- **A development testbed** — Allows rapid iteration on P4 programs

### Virtual Networking

The BMv2 container creates virtual ethernet (veth) pairs for packet I/O:

```
veth0 ←→ veth0_peer     (Port 0)
veth1 ←→ veth1_peer     (Port 1)
veth2 ←→ veth2_peer     (Port 2)
veth3 ←→ veth3_peer     (Port 3)
```

## Data Flow

1. Developer sends API request to OPI server (REST on :8080)
2. OPI server manages port/pipeline resources in memory
3. P4 programs are compiled and loaded into BMv2 via P4Runtime
4. BMv2 processes packets according to the loaded P4 pipeline
5. Forwarding rules are programmed via P4Runtime table entries

## Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| OPI Server | Go 1.22 | API emulation |
| P4 Switch | BMv2 (C++) | Packet processing |
| P4 Compiler | p4c | P4 → BMv2 JSON |
| Networking | veth pairs | Virtual interfaces |
| Orchestration | Docker Compose | Container management |
| CI/CD | GitHub Actions | Build & publish |
| Registry | GHCR | Image distribution |
