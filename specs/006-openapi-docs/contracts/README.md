# Contracts: OpenAPI Documentation Tool Integration

**Feature Branch**: `006-openapi-docs`  
**Date**: 2026-01-30

## Overview

This feature does not introduce new API contracts. It provides a viewer for the existing OpenAPI specification.

## Existing Contracts (Referenced)

| Contract | Location | Description |
|----------|----------|-------------|
| OpenAPI Specification | `api/openapi/openapi.yaml` | The existing API contract that will be rendered by Scalar |

## Infrastructure Contracts

### Docker Compose Service Definition

The documentation server is defined as a Docker Compose service:

```yaml
# Added to docker-compose.yml
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

## Service Endpoints

| Endpoint | Description |
|----------|-------------|
| `http://localhost:8081/` | Scalar API Reference UI |
| `http://localhost:8081/health` | Health check endpoint |

## Notes

- No new REST API endpoints are added to the OPTel Training API
- The documentation server is a separate service that reads the existing OpenAPI spec
- "Try it" requests go directly from browser to API server (`http://localhost:8080/api/v1`)
