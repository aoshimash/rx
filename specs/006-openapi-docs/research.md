# Research: OpenAPI Documentation Tool Integration

**Feature Branch**: `006-openapi-docs`  
**Date**: 2026-01-30

## Research Tasks

### 1. Scalar Docker Image Configuration

**Question**: How to deploy Scalar API Reference via Docker Compose?

**Decision**: Use official `scalarapi/api-reference` Docker image with OpenAPI file mounting.

**Rationale**: 
- Official Docker image is well-maintained (25.2 MB, updated frequently)
- Supports direct OpenAPI file mounting for automatic discovery
- Includes health check endpoint (`/health`)
- Zero configuration needed when OpenAPI file is mounted

**Alternatives Considered**:
- **npm/npx CLI**: Requires Node.js installation, adds complexity
- **Build custom Docker image**: Unnecessary overhead, official image is sufficient
- **Embed in Go API server**: Would violate separation of concerns, harder to maintain

**Implementation Details**:
- Image: `scalarapi/api-reference:latest`
- Mount: `./api/openapi/openapi.yaml:/app/public/openapi.yaml:ro`
- Port: `8081` (to avoid conflict with API server on `8080`)
- Health check: `curl -f http://localhost:8081/health`

### 2. OpenAPI Specification Path

**Question**: Where is the OpenAPI specification file located?

**Decision**: Use `api/openapi/openapi.yaml`

**Rationale**: 
- Constitution specifies OpenAPI spec location as `api/openapi/openapi.yaml`
- File exists and contains valid OpenAPI 3.0.3 specification
- Already used by oapi-codegen for server stub generation

**Note**: The spec.md incorrectly referenced `api/openapi.yaml`. Actual path is `api/openapi/openapi.yaml`.

### 3. Try It Feature Configuration

**Question**: How to enable "Try it" API testing from documentation?

**Decision**: Configure Scalar to proxy requests to the API server within Docker network.

**Rationale**:
- Scalar supports custom server URLs via configuration
- Docker Compose internal networking allows service-to-service communication
- API server is accessible at `http://api:8080` within Docker network

**Implementation Details**:
- Default spec server URL (`http://localhost:8080/api/v1`) works for browser requests
- No additional proxy configuration needed - browser makes direct requests to API

### 4. Docker Compose Service Configuration

**Question**: What is the optimal Docker Compose configuration?

**Decision**: Add `api-docs` service with dependency on OpenAPI file availability.

**Configuration**:
```yaml
api-docs:
  image: scalarapi/api-reference:latest
  container_name: optel-training-api-docs
  ports:
    - "${API_DOCS_PORT:-8081}:80"
  volumes:
    - ./api/openapi/openapi.yaml:/app/public/openapi.yaml:ro
  healthcheck:
    test: ["CMD", "wget", "-q", "--spider", "http://localhost:80/health"]
    interval: 10s
    timeout: 5s
    retries: 3
  restart: unless-stopped
```

**Rationale**:
- Read-only mount (`:ro`) prevents accidental modifications
- Port 8081 avoids conflict with API server
- Health check ensures container is ready before marking as healthy
- No dependency on API service (documentation works independently)

### 5. Documentation Update Requirements

**Question**: What documentation needs to be updated?

**Decision**: Update `docs/DEVELOPMENT.md` with API documentation startup instructions.

**Content to Add**:
- Command to start documentation server: `docker compose up -d api-docs`
- URL to access documentation: `http://localhost:8081`
- Instructions for using "Try it" feature with running API server

## Open Questions (Resolved)

| Question | Resolution |
|----------|------------|
| Docker image selection | `scalarapi/api-reference:latest` |
| Port number | 8081 (configurable via `API_DOCS_PORT`) |
| OpenAPI file path | `api/openapi/openapi.yaml` |
| Configuration method | File mounting (zero config) |

## References

- [Scalar Docker Documentation](https://scalar.com/products/api-references/integrations/docker)
- [Docker Hub: scalarapi/api-reference](https://hub.docker.com/r/scalarapi/api-reference)
- [Scalar GitHub Repository](https://github.com/scalar/scalar)
