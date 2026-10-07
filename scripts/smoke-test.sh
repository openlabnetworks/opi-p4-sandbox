#!/bin/bash
# ============================================================
# Smoke Test Script for OPI & P4 Sandbox
# ============================================================
# Validates end-to-end functionality of the sandbox.
# Exit code 0 = all tests passed, non-zero = failure.
# ============================================================
set -euo pipefail

OPI_URL="${OPI_URL:-http://localhost:8080}"
PASS=0
FAIL=0

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

run_test() {
    local name="$1"
    local cmd="$2"

    echo -en "${CYAN}[TEST]${RESET} ${name}... "
    if eval "$cmd" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PASS${RESET}"
        PASS=$((PASS + 1))
    else
        echo -e "${RED}❌ FAIL${RESET}"
        FAIL=$((FAIL + 1))
    fi
}

echo -e "${BOLD}🧪 OPI & P4 Sandbox — Smoke Tests${RESET}"
echo -e "${CYAN}═══════════════════════════════════════════${RESET}"
echo ""

# Wait for server to be ready
echo -e "Waiting for OPI server at ${OPI_URL}..."
for i in $(seq 1 30); do
    if curl -sf "${OPI_URL}/healthz" > /dev/null 2>&1; then
        echo -e "${GREEN}Server is ready!${RESET}"
        break
    fi
    if [ "$i" -eq 30 ]; then
        echo -e "${RED}Server did not become ready in 60 seconds.${RESET}"
        exit 1
    fi
    sleep 2
done
echo ""

# --- Tests ---

run_test "Health check returns OK" \
    "curl -sf ${OPI_URL}/healthz | grep -q SERVING"

run_test "API info endpoint" \
    "curl -sf ${OPI_URL}/ | grep -q 'OPI P4 Sandbox'"

run_test "List ports returns JSON array" \
    "curl -sf ${OPI_URL}/v1/ports | python3 -c 'import json,sys; data=json.load(sys.stdin); assert \"ports\" in data'"

run_test "Default port-0 exists" \
    "curl -sf ${OPI_URL}/v1/ports/port-0 | grep -q 'port-0'"

run_test "Create new port" \
    "curl -sf -X POST ${OPI_URL}/v1/ports -H 'Content-Type: application/json' -d '{\"mac_address\":\"aa:bb:cc:dd:ee:ff\",\"mtu\":9000}' | grep -q 'aa:bb:cc:dd:ee:ff'"

run_test "List ports shows 2 ports" \
    "curl -sf ${OPI_URL}/v1/ports | python3 -c 'import json,sys; data=json.load(sys.stdin); assert len(data[\"ports\"]) >= 2'"

run_test "Delete port succeeds" \
    "curl -sf -X DELETE ${OPI_URL}/v1/ports/port-0 -o /dev/null -w '%{http_code}' | grep -q '204'"

run_test "Get deleted port returns 404" \
    "curl -s ${OPI_URL}/v1/ports/port-0 -o /dev/null -w '%{http_code}' | grep -q '404'"

run_test "List pipelines returns JSON array" \
    "curl -sf ${OPI_URL}/v1/pipelines | python3 -c 'import json,sys; data=json.load(sys.stdin); assert \"pipelines\" in data'"

run_test "Default pipeline exists" \
    "curl -sf ${OPI_URL}/v1/pipelines/pipeline-default | grep -q 'pipeline-default'"

run_test "Create new pipeline" \
    "curl -sf -X POST ${OPI_URL}/v1/pipelines -H 'Content-Type: application/json' -d '{\"p4_program\":\"custom.p4\"}' | grep -q 'custom.p4'"

# --- Summary ---
echo ""
echo -e "${CYAN}═══════════════════════════════════════════${RESET}"
TOTAL=$((PASS + FAIL))
echo -e "${BOLD}Results: ${PASS}/${TOTAL} passed${RESET}"

if [ "$FAIL" -gt 0 ]; then
    echo -e "${RED}${BOLD}❌ ${FAIL} test(s) failed!${RESET}"
    exit 1
else
    echo -e "${GREEN}${BOLD}🎉 All tests passed!${RESET}"
    exit 0
fi
