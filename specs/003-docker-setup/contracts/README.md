# Contracts: Docker Development Environment Setup

**Date**: 2026-01-25  
**Branch**: `003-docker-setup`  
**Spec**: [spec.md](../spec.md)

## Scope

This feature introduces **no new HTTP endpoints** and makes **no changes** to the OpenAPI contract.

Instead, the "contracts" for this feature are the expected **file-level deliverables** and their behavior:

- `.devcontainer/` provides the developer environment for iterative local development.
- `api/Dockerfile` is the single Dockerfile for production distribution.
- `docker-compose.yml` enables local smoke-testing using the production container build.

## Non-Goals

- No changes to `api/openapi/openapi.yaml`
- No changes to API request/response schemas
- No new REST resources

