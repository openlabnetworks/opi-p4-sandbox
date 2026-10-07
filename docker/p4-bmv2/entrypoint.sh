#!/bin/bash
# ============================================================
# P4 BMv2 Entrypoint Script
# ============================================================
# Starts simple_switch_grpc with a default or user-supplied P4 pipeline.
# If a compiled JSON pipeline exists, it is loaded automatically.
# ============================================================
set -e

COMPILED_DIR="/p4/compiled"
LOG_DIR="/p4/logs"
DEFAULT_PIPELINE="${COMPILED_DIR}/basic_forwarding.json"
DEVICE_ID="${P4_DEVICE_ID:-0}"
GRPC_PORT="${P4_GRPC_PORT:-9559}"
THRIFT_PORT="${P4_THRIFT_PORT:-9090}"
LOG_LEVEL="${P4_LOG_LEVEL:-info}"

echo "╔══════════════════════════════════════════════════╗"
echo "║      P4 BMv2 Software Switch — Sandbox           ║"
echo "╠══════════════════════════════════════════════════╣"
echo "║  Device ID  : ${DEVICE_ID}                                ║"
echo "║  gRPC Port  : ${GRPC_PORT}                             ║"
echo "║  Thrift Port: ${THRIFT_PORT}                             ║"
echo "║  Log Level  : ${LOG_LEVEL}                           ║"
echo "╚══════════════════════════════════════════════════╝"

# Create virtual interfaces if they don't exist
create_veth_pairs() {
    echo "[BMv2] Creating virtual ethernet pairs..."
    for i in 0 1 2 3; do
        if ! ip link show "veth${i}" &>/dev/null; then
            ip link add "veth${i}" type veth peer name "veth${i}_peer" 2>/dev/null || true
            ip link set "veth${i}" up 2>/dev/null || true
            ip link set "veth${i}_peer" up 2>/dev/null || true
            echo "[BMv2] Created veth${i} <-> veth${i}_peer"
        fi
    done
}

# Try to create veth pairs (may fail without NET_ADMIN, which is fine)
create_veth_pairs 2>/dev/null || echo "[BMv2] Skipping veth creation (no NET_ADMIN capability)"

# Build the command
CMD="simple_switch_grpc"
CMD_ARGS=""

# Add pipeline if available
if [ -f "${DEFAULT_PIPELINE}" ]; then
    echo "[BMv2] Loading default pipeline: ${DEFAULT_PIPELINE}"
    CMD_ARGS="${CMD_ARGS} ${DEFAULT_PIPELINE}"
else
    echo "[BMv2] No default pipeline found. Starting without a loaded program."
    echo "[BMv2] Use P4Runtime to load a pipeline at runtime."
fi

# Add interfaces if veths exist
for i in 0 1 2 3; do
    if ip link show "veth${i}" &>/dev/null; then
        CMD_ARGS="${CMD_ARGS} -i ${i}@veth${i}"
    fi
done

# Add standard flags
CMD_ARGS="${CMD_ARGS} --device-id ${DEVICE_ID}"
CMD_ARGS="${CMD_ARGS} --log-file ${LOG_DIR}/bmv2.log"
CMD_ARGS="${CMD_ARGS} --log-level ${LOG_LEVEL}"
CMD_ARGS="${CMD_ARGS} --thrift-port ${THRIFT_PORT}"
CMD_ARGS="${CMD_ARGS} -- --grpc-server-addr 0.0.0.0:${GRPC_PORT}"

# Append any extra args passed to the container
CMD_ARGS="${CMD_ARGS} $@"

echo "[BMv2] Starting: ${CMD} ${CMD_ARGS}"
exec ${CMD} ${CMD_ARGS}
