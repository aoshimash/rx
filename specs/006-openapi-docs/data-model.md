# Data Model: OpenAPI Documentation Tool Integration

**Feature Branch**: `006-openapi-docs`  
**Date**: 2026-01-30

## Overview

This feature does not introduce new domain entities or data models. It provides a documentation viewer for the existing OpenAPI specification.

## Entities (Read-Only)

### OpenAPI Specification File

| Attribute | Type | Description |
|-----------|------|-------------|
| Path | string | `api/openapi/openapi.yaml` |
| Format | OpenAPI 3.0.3 | Industry-standard API specification format |
| Access Mode | Read-only | Mounted as read-only volume in Docker |

**Notes**:
- The OpenAPI specification is the source of truth for the API contract
- It is already managed by the existing development workflow (spec → codegen → implementation)
- This feature only reads the specification; it does not modify it

## Configuration Entities

### Docker Compose Service Configuration

| Attribute | Type | Default | Description |
|-----------|------|---------|-------------|
| API_DOCS_PORT | integer | 8081 | Host port for documentation server |
| Container name | string | optel-training-api-docs | Docker container identifier |

## Relationships

```
┌─────────────────────────┐
│  OpenAPI Specification  │
│  (api/openapi/openapi.yaml)
└───────────┬─────────────┘
            │ read-only mount
            ▼
┌─────────────────────────┐
│  Scalar API Reference   │
│  (Docker container)     │
└───────────┬─────────────┘
            │ serves
            ▼
┌─────────────────────────┐
│  Developer Browser      │
│  (http://localhost:8081)│
└───────────┬─────────────┘
            │ "Try it" requests
            ▼
┌─────────────────────────┐
│  OPTel Training API     │
│  (http://localhost:8080)│
└─────────────────────────┘
```

## State Transitions

Not applicable - this feature provides a stateless documentation viewer.

## Validation Rules

Not applicable - no new data validation is introduced by this feature.
