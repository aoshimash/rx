# Implementation Plan: OpenAPI Documentation Tool Integration

**Branch**: `006-openapi-docs` | **Date**: 2026-01-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/006-openapi-docs/spec.md`

## Summary

Integrate Scalar API Reference as an OpenAPI documentation viewer, deployed via Docker Compose. The documentation server reads the existing `api/openapi/openapi.yaml` specification and serves an interactive API reference at `http://localhost:8081`. Developers can view all API endpoints and use the "Try it" feature to test requests against the running API server.

## Technical Context

**Language/Version**: N/A (Docker-based deployment, no application code)  
**Primary Dependencies**: Docker, Docker Compose, Scalar API Reference (`scalarapi/api-reference:latest`)  
**Storage**: N/A (stateless documentation viewer)  
**Testing**: Manual verification, container health check  
**Target Platform**: Local development environment (Linux/macOS/Windows with Docker)  
**Project Type**: Infrastructure/tooling addition to existing monorepo  
**Performance Goals**: Documentation loads within 3 seconds  
**Constraints**: Must not conflict with existing services (API on 8080, PostgreSQL on 5432, MinIO on 9000/9001)  
**Scale/Scope**: Single developer local environment

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with OPTel Workout Constitution principles:

- **Dumb Backend**: ✅ N/A - This feature adds a documentation viewer, not backend logic
- **Domain-Driven Schema-First**: ✅ Compliant - Reads existing OpenAPI spec without modification
- **Terminology**: ✅ Compliant - Uses standard API documentation terminology
- **Clean Architecture**: ✅ N/A - No application code changes, infrastructure only
- **Monorepo Structure**: ✅ Compliant - Adds shared infrastructure service in docker-compose.yml

**Result**: All principles satisfied. No violations to document.

## Project Structure

### Documentation (this feature)

```text
specs/006-openapi-docs/
├── plan.md              # This file
├── research.md          # Scalar Docker integration research
├── data-model.md        # Entity relationships (minimal for this feature)
├── quickstart.md        # Developer quick start guide
├── contracts/           # Infrastructure contracts
│   └── README.md        # Docker Compose service definition
├── checklists/
│   └── requirements.md  # Specification quality checklist
├── spec.md              # Feature specification
└── tasks.md             # Implementation tasks (created by /speckit.tasks)
```

### Source Code (repository root)

```text
# Files modified by this feature
docker-compose.yml       # Add api-docs service
.env.example             # Add API_DOCS_PORT variable
docs/DEVELOPMENT.md      # Add documentation server instructions

# Files read by this feature (no modification)
api/openapi/openapi.yaml # Existing OpenAPI specification
```

**Structure Decision**: Infrastructure-only change. No new source directories. Extends existing docker-compose.yml with a new service.

## Implementation Approach

### Phase 1: Docker Compose Configuration

Add Scalar API Reference service to `docker-compose.yml`:

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

### Phase 2: Environment Configuration

Add to `.env.example`:
```
# API Documentation
API_DOCS_PORT=8081
```

### Phase 3: Documentation Update

Update `docs/DEVELOPMENT.md` with:
- New section for API Documentation
- Startup command: `docker compose up -d api-docs`
- URL: `http://localhost:8081`
- Instructions for "Try it" feature

## Complexity Tracking

> No complexity violations identified. This feature is infrastructure-only with minimal changes.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | - | - |

## Artifacts Generated

| Artifact | Path | Description |
|----------|------|-------------|
| Research | `specs/006-openapi-docs/research.md` | Scalar Docker integration research |
| Data Model | `specs/006-openapi-docs/data-model.md` | Entity relationships |
| Contracts | `specs/006-openapi-docs/contracts/README.md` | Docker Compose service contract |
| Quickstart | `specs/006-openapi-docs/quickstart.md` | Developer guide |

## Next Steps

Run `/speckit.tasks` to generate implementation tasks from this plan.
