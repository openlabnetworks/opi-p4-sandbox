# ============================================================
# P4 BMv2 Software Switch — Docker Build
# ============================================================
# Uses the official p4lang behavioral-model as a base and adds:
# - Pre-compiled sample P4 programs
# - P4Runtime gRPC server (simple_switch_grpc)
# - Convenient entrypoint with default pipeline
# ============================================================

FROM p4lang/behavioral-model:latest AS bmv2-base

# --- Runtime Stage ---
FROM ubuntu:22.04

# Avoid interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    libboost-system1.74.0 \
    libboost-thread1.74.0 \
    libboost-filesystem1.74.0 \
    libboost-program-options1.74.0 \
    libgmp10 \
    libpcap0.8 \
    libnanomsg5 \
    libjsoncpp25 \
    python3 \
    python3-pip \
    iproute2 \
    net-tools \
    tcpdump \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy BMv2 binaries from the base image
COPY --from=bmv2-base /usr/local/bin/simple_switch /usr/local/bin/
COPY --from=bmv2-base /usr/local/bin/simple_switch_grpc /usr/local/bin/ 
COPY --from=bmv2-base /usr/local/lib/ /usr/local/lib/

# Update shared library cache
RUN ldconfig

# Create working directory for P4 programs
RUN mkdir -p /p4/programs /p4/compiled /p4/logs

# Copy sample P4 programs
COPY p4-programs/ /p4/programs/
COPY compiled/ /p4/compiled/
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose P4Runtime gRPC port and Thrift port
EXPOSE 9559 9090

# Health check via simple_switch_grpc thrift API
HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -sf http://localhost:9090/ || pgrep simple_switch || exit 1

WORKDIR /p4

ENTRYPOINT ["/entrypoint.sh"]
CMD ["--device-id", "0", "--log-level", "info"]
