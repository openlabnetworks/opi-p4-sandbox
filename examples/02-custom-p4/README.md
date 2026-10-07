# Experiment 02: Custom P4 Pipeline

Learn how to compile a P4 program and manage pipelines through the OPI API.

## Prerequisites

The sandbox must be running:
```bash
make start-sandbox
```

## Step 1: Explore the Sample P4 Program

The sandbox includes a sample P4_16 program at `src/p4-programs/basic_forwarding.p4`.

```bash
cat src/p4-programs/basic_forwarding.p4
```

This program implements:
- **Ethernet parsing** — Extracts source/destination MAC addresses
- **IPv4 parsing** — Extracts IP headers
- **LPM forwarding** — Matches destination IP and forwards to an output port
- **TTL decrement** — Decreases Time-To-Live on each hop
- **Checksum recomputation** — Updates IPv4 header checksum

## Step 2: Compile the P4 Program

Use the Makefile target to compile the P4 program into BMv2 JSON:

```bash
make compile-p4 P4=basic_forwarding.p4
```

This runs the `p4c-bm2-ss` compiler inside a Docker container and produces
`src/p4-programs/basic_forwarding.json`.

## Step 3: List Current Pipelines

```bash
curl http://localhost:8080/v1/pipelines | python3 -m json.tool
```

You should see the default `pipeline-default` that loads on startup.

## Step 4: Create a New Pipeline

```bash
curl -X POST http://localhost:8080/v1/pipelines \
  -H "Content-Type: application/json" \
  -d '{
    "p4_program": "basic_forwarding.p4"
  }' | python3 -m json.tool
```

## Step 5: Check Pipeline Status

```bash
curl http://localhost:8080/v1/pipelines/pipeline-1 | python3 -m json.tool
```

The pipeline status should show `LOADED`.

## Step 6: Write Your Own P4 Program

Create a new file `src/p4-programs/my_program.p4`:

```p4
/* My first P4 program */
#include <core.p4>
#include <v1model.p4>

header ethernet_t {
    bit<48> dstAddr;
    bit<48> srcAddr;
    bit<16> etherType;
}

struct headers { ethernet_t ethernet; }
struct metadata { }

parser MyParser(packet_in pkt, out headers hdr,
                inout metadata meta, inout standard_metadata_t smeta) {
    state start {
        pkt.extract(hdr.ethernet);
        transition accept;
    }
}

control MyIngress(inout headers hdr, inout metadata meta,
                  inout standard_metadata_t smeta) {
    apply {
        // Forward all packets to port 1
        smeta.egress_spec = 1;
    }
}

control MyEgress(inout headers hdr, inout metadata meta,
                 inout standard_metadata_t smeta) { apply { } }

control MyVerifyChecksum(inout headers hdr, inout metadata meta) { apply { } }
control MyComputeChecksum(inout headers hdr, inout metadata meta) { apply { } }

control MyDeparser(packet_out pkt, in headers hdr) {
    apply { pkt.emit(hdr.ethernet); }
}

V1Switch(MyParser(), MyVerifyChecksum(), MyIngress(),
         MyEgress(), MyComputeChecksum(), MyDeparser()) main;
```

Then compile and create a pipeline for it:

```bash
make compile-p4 P4=my_program.p4

curl -X POST http://localhost:8080/v1/pipelines \
  -H "Content-Type: application/json" \
  -d '{"p4_program": "my_program.p4"}' | python3 -m json.tool
```

---

## 🎉 Congratulations!

You've learned how to:
- Read and understand a P4_16 program structure
- Compile P4 programs for BMv2
- Manage pipelines through the OPI API
- Write your own P4 program from scratch

**Next:** Try [Experiment 03 — OPI Pipeline Management](../03-opi-pipeline-mgmt/README.md)
