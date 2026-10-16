#!/bin/bash
# ============================================================
# Compile a P4 Program and Register Its Metadata
# ============================================================
# Usage: ./load-p4-pipeline.sh <p4_file>
# Example: ./load-p4-pipeline.sh basic_forwarding.p4
#
# This helper compiles a P4 source file and optionally registers its name with
# the REST API. It does not install a forwarding pipeline on BMv2.
# ============================================================
set -euo pipefail

P4_FILE="${1:?Usage: $0 <p4_file>}"
P4_DIR="$(cd "$(dirname "$0")/../src/p4-programs" && pwd)"
P4_PATH="${P4_DIR}/${P4_FILE}"
JSON_FILE="${P4_DIR}/$(basename "$P4_FILE" .p4).json"

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
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

echo -e "${CYAN}Step 2:${RESET} Registering pipeline metadata with the OPI API..."
# The REST API is an in-memory emulator. Registration does not contact BMv2.
if curl -sf http://localhost:8080/healthz > /dev/null 2>&1; then
    curl -fsS -X POST http://localhost:8080/v1/pipelines \
        -H "Content-Type: application/json" \
        -d "{\"p4_program\":\"${P4_FILE}\"}" > /dev/null
    echo -e "${GREEN}✅ Registered metadata with OPI server at http://localhost:8080/v1/pipelines${RESET}"
else
    echo -e "${YELLOW}⚠️  OPI server is unavailable; pipeline metadata was not registered.${RESET}"
fi
echo -e "${YELLOW}NOTE: This script does not load the compiled pipeline into BMv2.${RESET}"
echo ""
echo "Compiled JSON: ${JSON_FILE}"
