# Production Dockerfile for RX API
# Multi-stage build: Ubuntu-based builder + distroless runtime
# Distroless provides minimal attack surface (no shell, no package manager)

# Build stage: --platform=$BUILDPLATFORM ensures the builder always runs on the
# host architecture (amd64 in CI), avoiding QEMU emulation. Go cross-compiles
# for the target via TARGETARCH, which Docker buildx sets automatically.
FROM --platform=$BUILDPLATFORM golang:1.26.1 AS builder

ARG TARGETARCH

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary (CGO_ENABLED=0 for distroless compatibility)
# TARGETARCH enables multi-arch builds via Go cross-compilation (no QEMU needed)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o bin/server cmd/server/main.go

# Runtime stage (distroless - no shell, no package manager, minimal attack surface)
FROM gcr.io/distroless/static-debian11:nonroot

# Copy CA certificates for HTTPS (distroless includes them)
# Copy binary from builder
COPY --from=builder /build/bin/server /app/server

# Use nonroot user (distroless provides nonroot user: 65532:65532)
USER nonroot:nonroot

# Expose port
EXPOSE 8080

# Health check (requires /health endpoint to be implemented)
# Note: distroless has no shell, so health check must use binary directly
# Uncomment when health endpoint is available:
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#     CMD ["/app/server", "-health-check"] || exit 1

# Run binary
# Distroless has no shell, so use exec form
ENTRYPOINT ["/app/server"]
