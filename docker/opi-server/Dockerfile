# ============================================================
# OPI API Server — Multi-stage Docker Build
# ============================================================
# Stage 1: Build the Go binary
# Stage 2: Minimal runtime image
# ============================================================

# --- Builder Stage ---
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /opi-server .

# --- Runtime Stage ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates curl

# Create non-root user
RUN addgroup -S opi && adduser -S opi -G opi

COPY --from=builder /opi-server /usr/local/bin/opi-server

# Health check
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/healthz || exit 1

USER opi

EXPOSE 8080

ENTRYPOINT ["opi-server"]
CMD ["--port", "8080"]
