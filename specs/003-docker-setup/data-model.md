# Data Model: Docker Development Environment Setup

**Date**: 2026-01-25  
**Branch**: `003-docker-setup`  
**Spec**: [spec.md](./spec.md)

## Purpose

This feature is configuration-centric. The "data model" below describes the **deliverable artifacts** (configuration files) and their required properties/relationships so they can be validated and maintained consistently.

## Entities

### DevContainerConfig

- **Represents**: The DevContainer configuration used for local iterative development.
- **Location**: `.devcontainer/devcontainer.json` (and optional companion files under `.devcontainer/`)
- **Key attributes (conceptual)**:
  - **workspace**: points at the repository workspace in the container
  - **toolchain availability**: required development tools for this repository are available inside the container
  - **editor integration**: recommended editor extensions/settings are included
  - **ports**: ability to develop/run the API locally (port forwarding when needed)
- **Validation rules**:
  - Must enable running the repository standard generate/lint/test workflow inside the container with no local tool installs.
  - Must not require a separate "development Dockerfile" to exist.

### ProductionDockerfile

- **Represents**: The single Dockerfile used for production distribution of the API service.
- **Location**: `api/Dockerfile`
- **Key attributes (conceptual)**:
  - **build**: produces a runnable container image for the API service
  - **security posture**: minimal runtime image, non-root user
  - **size constraints**: optimized image size (see success criteria)
- **Validation rules**:
  - Must be the only maintained Dockerfile for the API (no separate development Dockerfile).
  - Must build a minimal, secure, runnable image for distribution.

### ComposeConfig

- **Represents**: Docker Compose configuration for local smoke-testing.
- **Location**: `docker-compose.yml` (repository root)
- **Key attributes (conceptual)**:
  - **service**: one service representing the API container built from `api/Dockerfile`
  - **ports**: maps the container port to a configurable host port
  - **environment**: supports environment variable configuration
  - **operational commands**: supports start/logs/stop workflows
- **Validation rules**:
  - Must be production-like (build/run the production container image).
  - Must not become the primary iterative development environment.

## Relationships

- `DevContainerConfig` is the primary developer workflow for iterative changes.
- `ComposeConfig` depends on `ProductionDockerfile` as its build source for the API service.
- `ProductionDockerfile` must not depend on `DevContainerConfig` (distribution must be self-contained).

