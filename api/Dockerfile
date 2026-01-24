# Development Dockerfile for OPTel Workout API
# This is for LOCAL DEVELOPMENT ONLY - includes all development tools
# For production, use Dockerfile.prod instead
# Uses Ubuntu-based golang image for better compatibility
FROM golang:1.25 AS dev

# Install development tools
RUN apt-get update && apt-get install -y \
    make \
    git \
    curl \
    bash \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install oapi-codegen
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Install golangci-lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Add Go bin to PATH
ENV PATH="${PATH}:/root/go/bin"

# Set working directory
WORKDIR /workspace

# Default command (keeps container running, can be overridden)
CMD ["tail", "-f", "/dev/null"]
