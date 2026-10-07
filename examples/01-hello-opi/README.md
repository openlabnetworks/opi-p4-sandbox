# Experiment 01: Hello OPI

Your first interaction with the OPI API. In this experiment, you'll discover the
API, list ports, and create your first virtual network port.

## Prerequisites

The sandbox must be running:
```bash
make start-sandbox
```

## Step 1: Check Server Health

```bash
curl http://localhost:8080/healthz | python3 -m json.tool
```

**Expected output:**
```json
{
    "status": "SERVING",
    "services": {
        "opi.v1.NetworkService": "SERVING",
        "opi.v1.PipelineService": "SERVING"
    }
}
```

## Step 2: Discover the API

```bash
curl http://localhost:8080/ | python3 -m json.tool
```

This shows you the available services and version.

## Step 3: List Existing Ports

```bash
curl http://localhost:8080/v1/ports | python3 -m json.tool
```

You should see the default `port-0` that the sandbox creates on startup.

## Step 4: Create a New Port

```bash
curl -X POST http://localhost:8080/v1/ports \
  -H "Content-Type: application/json" \
  -d '{
    "mac_address": "de:ad:be:ef:00:01",
    "mtu": 9000
  }' | python3 -m json.tool
```

**Expected output:**
```json
{
    "name": "ports/port-1",
    "id": "port-1",
    "mac_address": "de:ad:be:ef:00:01",
    "mtu": 9000,
    "admin_state": "UP",
    "oper_state": "UP",
    "speed": "25G",
    ...
}
```

## Step 5: Get a Specific Port

```bash
curl http://localhost:8080/v1/ports/port-1 | python3 -m json.tool
```

## Step 6: Delete a Port

```bash
curl -X DELETE http://localhost:8080/v1/ports/port-1 -v
```

You should see a `204 No Content` response.

## Step 7: Verify Deletion

```bash
curl http://localhost:8080/v1/ports/port-1
```

You should get a `404 Not Found` error.

---

## 🎉 Congratulations!

You've completed your first OPI API interaction. You now know how to:
- Check server health
- List, create, get, and delete network ports
- Use the OPI resource naming convention (`ports/port-{id}`)

**Next:** Try [Experiment 02 — Custom P4 Pipeline](../02-custom-p4/README.md)
