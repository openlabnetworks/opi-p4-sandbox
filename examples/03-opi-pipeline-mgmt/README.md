# Experiment 03: OPI Pipeline Lifecycle Management

Learn the full lifecycle of managing P4 pipelines through the OPI API:
create → inspect → update → delete.

> **Scope:** This experiment demonstrates the OPI server's in-memory pipeline
> registry. It does not install forwarding pipelines in BMv2. P4Runtime-backed
> pipeline management is planned in [issue #6](https://github.com/openlabnetworks/opi-p4-sandbox/issues/6).

## Prerequisites

```bash
make start-sandbox
```

## Step 1: View All Pipelines

```bash
curl -s http://localhost:8080/v1/pipelines | python3 -m json.tool
```

## Step 2: Create Multiple Pipelines

Create pipelines for different use cases:

```bash
# L2 forwarding pipeline
curl -s -X POST http://localhost:8080/v1/pipelines \
  -H "Content-Type: application/json" \
  -d '{"p4_program": "l2_forwarding.p4"}' | python3 -m json.tool

# ACL filtering pipeline
curl -s -X POST http://localhost:8080/v1/pipelines \
  -H "Content-Type: application/json" \
  -d '{"p4_program": "acl_filter.p4"}' | python3 -m json.tool

# Load balancer pipeline
curl -s -X POST http://localhost:8080/v1/pipelines \
  -H "Content-Type: application/json" \
  -d '{"p4_program": "load_balancer.p4"}' | python3 -m json.tool
```

## Step 3: Inspect Each Pipeline

```bash
for id in pipeline-default pipeline-1 pipeline-2 pipeline-3; do
  echo "=== $id ==="
  curl -s http://localhost:8080/v1/pipelines/$id | python3 -m json.tool
  echo ""
done
```

## Step 4: Pipeline Cleanup

Delete pipelines you no longer need:

```bash
# Delete the ACL pipeline
curl -s -X DELETE http://localhost:8080/v1/pipelines/pipeline-2 -w "\nHTTP Status: %{http_code}\n"

# Verify it's gone
curl -s http://localhost:8080/v1/pipelines/pipeline-2
```

## Step 5: Automate with a Script

Create a script that manages the full pipeline lifecycle:

```bash
#!/bin/bash
OPI=http://localhost:8080

# Create
echo "Creating pipeline..."
RESULT=$(curl -s -X POST $OPI/v1/pipelines \
  -H "Content-Type: application/json" \
  -d '{"p4_program": "test.p4"}')
echo "$RESULT" | python3 -m json.tool

# Extract ID
ID=$(echo "$RESULT" | python3 -c "import json,sys; print(json.load(sys.stdin)['id'])")
echo "Created pipeline: $ID"

# Inspect
echo "Inspecting pipeline..."
curl -s $OPI/v1/pipelines/$ID | python3 -m json.tool

# Delete
echo "Deleting pipeline..."
curl -s -X DELETE $OPI/v1/pipelines/$ID -w "HTTP Status: %{http_code}\n"

echo "Done!"
```

---

## 🎉 Congratulations!

You've mastered OPI pipeline lifecycle management:
- Bulk creation of pipelines
- Inspection and enumeration
- Cleanup and deletion
- Scripted automation

These operations update only the OPI metadata registry. Use a separate
P4Runtime client against BMv2 on `localhost:9559` when you need to install a
forwarding pipeline.

**Next:** Explore more advanced scenarios like traffic steering and telemetry.
