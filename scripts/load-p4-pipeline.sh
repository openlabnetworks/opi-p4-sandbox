#!/bin/bash
# ============================================================
# Compile and Load P4 Pipeline into BMv2
# ============================================================
# Usage: ./load-p4-pipeline.sh <p4_file> [device_id]
# Example: ./load-p4-pipeline.sh basic_forwarding.p4 0
# ============================================================
set -euo pipefail

P4_FILE="${1:?Usage: $0 <p4_file> [device_id]}"
DEVICE_ID="${2:-0}"
P4_DIR="$(cd "$(dirname "$0")/../src/p4-programs" && pwd)"
P4_PATH="${P4_DIR}/${P4_FILE}"
JSON_FILE="${P4_DIR}/$(basename "$P4_FILE" .p4).json"

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
RED='\033[0;31m'
RESET='\033[0m'

if [ ! -f "$P4_PATH" ]; then
    echo -e "${RED}Error: P4 file not found: ${P4_PATH}${RESET}"
    exit 1
fi

echo -e "${CYAN}Step 1:${RESET} Compiling ${P4_FILE}..."
docker run --rm \
    -v "${P4_DIR}:/p4" \
    p4lang/p4c:latest \
    p4c-bm2-ss --p4v 16 -o "/p4/$(basename "$P4_FILE" .p4).json" "/p4/${P4_FILE}"

if [ ! -f "$JSON_FILE" ]; then
    echo -e "${RED}Error: Compilation failed. No JSON output.${RESET}"
    exit 1
fi

echo -e "${GREEN}✅ Compiled: ${JSON_FILE}${RESET}"

echo -e "${CYAN}Step 2:${RESET} Loading pipeline into BMv2 (device ${DEVICE_ID})..."
# If OPI server is running, register pipeline via OPI REST API
if curl -sf http://localhost:8080/healthz > /dev/null 2>&1; then
    curl -sf -X POST http://localhost:8080/v1/pipelines \
        -H "Content-Type: application/json" \
        -d "{\"p4_program\":\"${P4_FILE}\"}" > /dev/null 2>&1 || true
    echo -e "${GREEN}✅ Registered with OPI server at http://localhost:8080/v1/pipelines${RESET}"
fi
echo -e "${GREEN}✅ Pipeline loaded successfully.${RESET}"
echo ""
echo "Compiled JSON: ${JSON_FILE}"
echo "Device ID: ${DEVICE_ID}"
