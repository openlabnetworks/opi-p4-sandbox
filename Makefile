# ============================================================
# OPI & P4 Developer Sandbox — Makefile
# ============================================================
# One-command interface for all sandbox operations.
#
# Quick Start:
#   make start-sandbox   → Pull images and start everything
#   make smoke-test      → Verify the sandbox is working
#   make stop-sandbox    → Tear everything down
#
# Run `make help` to see all available targets.
# ============================================================

.PHONY: help start-sandbox stop-sandbox restart-sandbox status logs \
        smoke-test build-images push-images clean install-tools \
        compile-p4

# --- Configuration ---
COMPOSE       := docker compose
OPI_PORT      ?= 8080
P4_GRPC_PORT  ?= 9559
P4_THRIFT_PORT?= 9090
OPI_URL       := http://localhost:$(OPI_PORT)
P4_PROGRAMS   := src/p4-programs

# Colors for terminal output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
CYAN   := \033[0;36m
BOLD   := \033[1m
RESET  := \033[0m

# ============================================================
# Primary Targets (User-facing)
# ============================================================

## start-sandbox: Pull images and start the full sandbox environment
start-sandbox: _check-docker
	@echo "$(GREEN)$(BOLD)🚀 Starting OPI & P4 Sandbox...$(RESET)"
	@echo ""
	$(COMPOSE) up -d
	@echo ""
	@echo "$(GREEN)$(BOLD)✅ Sandbox is running!$(RESET)"
	@echo ""
	@echo "  OPI API Server : $(CYAN)$(OPI_URL)$(RESET)"
	@echo "  P4 BMv2 gRPC   : $(CYAN)localhost:$(P4_GRPC_PORT)$(RESET)"
	@echo "  P4 BMv2 Thrift : $(CYAN)localhost:$(P4_THRIFT_PORT)$(RESET)"
	@echo ""
	@echo "$(YELLOW)Quick test:$(RESET)"
	@echo "  curl $(OPI_URL)/v1/ports"
	@echo ""

## stop-sandbox: Tear down all sandbox containers and networks
stop-sandbox:
	@echo "$(RED)🛑 Stopping sandbox...$(RESET)"
	$(COMPOSE) down --remove-orphans
	@echo "$(GREEN)✅ Sandbox stopped.$(RESET)"

## restart-sandbox: Stop and start the sandbox
restart-sandbox: stop-sandbox start-sandbox

## status: Show running containers, ports, and health
status:
	@echo "$(BOLD)📊 Sandbox Status$(RESET)"
	@echo "$(CYAN)─────────────────────────────────────────$(RESET)"
	@$(COMPOSE) ps
	@echo ""
	@echo "$(CYAN)─────────────────────────────────────────$(RESET)"
	@echo "$(BOLD)Health Checks:$(RESET)"
	@curl -sf $(OPI_URL)/healthz 2>/dev/null | python3 -m json.tool 2>/dev/null || echo "  OPI Server: $(RED)NOT REACHABLE$(RESET)"

## logs: Tail logs from all containers
logs:
	$(COMPOSE) logs -f --tail=50

## smoke-test: Run end-to-end validation of the sandbox
smoke-test: _check-curl
	@echo "$(BOLD)🧪 Running Smoke Tests...$(RESET)"
	@echo ""
	@echo "$(CYAN)Test 1:$(RESET) Health check..."
	@curl -sf $(OPI_URL)/healthz > /dev/null && \
		echo "  $(GREEN)✅ PASS$(RESET) — Server is healthy" || \
		(echo "  $(RED)❌ FAIL$(RESET) — Server not responding" && exit 1)
	@echo ""
	@echo "$(CYAN)Test 2:$(RESET) List ports..."
	@curl -sf $(OPI_URL)/v1/ports | python3 -m json.tool > /dev/null && \
		echo "  $(GREEN)✅ PASS$(RESET) — ListPorts returned valid JSON" || \
		(echo "  $(RED)❌ FAIL$(RESET) — ListPorts failed" && exit 1)
	@echo ""
	@echo "$(CYAN)Test 3:$(RESET) Create a port..."
	@curl -sf -X POST $(OPI_URL)/v1/ports \
		-H "Content-Type: application/json" \
		-d '{"mac_address":"aa:bb:cc:dd:ee:ff","mtu":9000}' | python3 -m json.tool > /dev/null && \
		echo "  $(GREEN)✅ PASS$(RESET) — CreatePort succeeded" || \
		(echo "  $(RED)❌ FAIL$(RESET) — CreatePort failed" && exit 1)
	@echo ""
	@echo "$(CYAN)Test 4:$(RESET) List pipelines..."
	@curl -sf $(OPI_URL)/v1/pipelines | python3 -m json.tool > /dev/null && \
		echo "  $(GREEN)✅ PASS$(RESET) — ListPipelines returned valid JSON" || \
		(echo "  $(RED)❌ FAIL$(RESET) — ListPipelines failed" && exit 1)
	@echo ""
	@echo "$(CYAN)Test 5:$(RESET) Create a pipeline..."
	@curl -sf -X POST $(OPI_URL)/v1/pipelines \
		-H "Content-Type: application/json" \
		-d '{"p4_program":"custom_program.p4"}' | python3 -m json.tool > /dev/null && \
		echo "  $(GREEN)✅ PASS$(RESET) — CreatePipeline succeeded" || \
		(echo "  $(RED)❌ FAIL$(RESET) — CreatePipeline failed" && exit 1)
	@echo ""
	@echo "$(GREEN)$(BOLD)🎉 All smoke tests passed!$(RESET)"

# ============================================================
# Build Targets (For Contributors)
# ============================================================

## build-images: Build Docker images locally
build-images:
	@echo "$(BOLD)🔨 Building Docker images...$(RESET)"
	$(COMPOSE) build
	@echo "$(GREEN)✅ Images built successfully.$(RESET)"

## push-images: Push images to container registry (CI use)
push-images:
	@echo "$(BOLD)📦 Pushing images to registry...$(RESET)"
	$(COMPOSE) push
	@echo "$(GREEN)✅ Images pushed.$(RESET)"

## compile-p4: Compile a P4 program (usage: make compile-p4 P4=myprogram.p4)
compile-p4:
ifndef P4
	@echo "$(RED)Error: specify P4 file. Usage: make compile-p4 P4=myprogram.p4$(RESET)"
	@exit 1
endif
	@echo "$(BOLD)📝 Compiling P4 program: $(P4)$(RESET)"
	docker run --rm -v $(PWD)/$(P4_PROGRAMS):/p4 p4lang/p4c:latest \
		p4c-bm2-ss --p4v 16 /p4/$(P4) -o /p4/$(basename $(P4) .p4).json
	@echo "$(GREEN)✅ Compiled: $(P4_PROGRAMS)/$(basename $(P4) .p4).json$(RESET)"

# ============================================================
# Utility Targets
# ============================================================

## install-tools: Install development tools (curl, jq, etc.)
install-tools:
	@echo "$(BOLD)🔧 Checking development tools...$(RESET)"
	@command -v docker >/dev/null 2>&1 && echo "  $(GREEN)✅$(RESET) docker" || echo "  $(RED)❌$(RESET) docker — install from https://docs.docker.com/get-docker/"
	@command -v curl >/dev/null 2>&1 && echo "  $(GREEN)✅$(RESET) curl" || echo "  $(RED)❌$(RESET) curl — install with: sudo apt install curl"
	@command -v python3 >/dev/null 2>&1 && echo "  $(GREEN)✅$(RESET) python3" || echo "  $(RED)❌$(RESET) python3 — install with: sudo apt install python3"
	@command -v jq >/dev/null 2>&1 && echo "  $(GREEN)✅$(RESET) jq" || echo "  $(YELLOW)⚠️$(RESET)  jq (optional) — install with: sudo apt install jq"
	@command -v go >/dev/null 2>&1 && echo "  $(GREEN)✅$(RESET) go" || echo "  $(YELLOW)⚠️$(RESET)  go (optional, for development) — install from https://go.dev/dl/"

## clean: Remove all containers, images, and volumes
clean: stop-sandbox
	@echo "$(RED)🧹 Cleaning up...$(RESET)"
	$(COMPOSE) down -v --rmi local
	docker image prune -f
	@echo "$(GREEN)✅ Clean.$(RESET)"

## help: Show this help message
help:
	@echo "$(BOLD)OPI & P4 Developer Sandbox$(RESET)"
	@echo "$(CYAN)═══════════════════════════════════════════$(RESET)"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | \
		awk -F': ' '{printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BOLD)Quick Start:$(RESET)"
	@echo "  make start-sandbox   # Start the sandbox"
	@echo "  make smoke-test      # Verify it works"
	@echo "  make stop-sandbox    # Tear it down"
	@echo ""

# ============================================================
# Internal Targets
# ============================================================

_check-docker:
	@command -v docker >/dev/null 2>&1 || \
		(echo "$(RED)$(BOLD)Error:$(RESET) Docker is not installed." && \
		 echo "Install it from: $(CYAN)https://docs.docker.com/get-docker/$(RESET)" && \
		 exit 1)
	@docker info >/dev/null 2>&1 || \
		(echo "$(RED)$(BOLD)Error:$(RESET) Docker daemon is not running." && \
		 echo "Start it with: $(CYAN)sudo systemctl start docker$(RESET)" && \
		 exit 1)

_check-curl:
	@command -v curl >/dev/null 2>&1 || \
		(echo "$(RED)$(BOLD)Error:$(RESET) curl is not installed." && \
		 echo "Install it with: $(CYAN)sudo apt install curl$(RESET)" && \
		 exit 1)

.DEFAULT_GOAL := help
