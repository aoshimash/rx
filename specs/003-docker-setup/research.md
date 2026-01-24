# Research: Docker Development Environment Setup

**Date**: 2026-01-25  
**Branch**: `003-docker-setup`  
**Spec**: [spec.md](./spec.md)

## Summary

This feature is primarily configuration and documentation work. There are no unresolved "NEEDS CLARIFICATION" items in the spec; research is used to record concrete decisions and the rationale to reduce downstream rework.

## Decisions

### Decision: DevContainer is the only development environment

- **Chosen**: Use DevContainer for iterative local development; do not maintain a separate development Dockerfile.
- **Rationale**: DevContainer integrates with the editor, provides a consistent toolchain, and avoids duplicating "dev environment" logic across Dockerfile + compose + docs.
- **Alternatives considered**:
  - Keep a dedicated development Dockerfile and run all dev commands via `docker compose exec`.
  - Install tools locally.

### Decision: Single Dockerfile for production distribution

- **Chosen**: Maintain exactly one Dockerfile for production distribution.
- **Location**: `api/Dockerfile` (clarified in spec).
- **Rationale**: The API is the deployable unit in this repository; keeping the Dockerfile in `api/` aligns the build context with the Go module and reduces path/config churn.
- **Alternatives considered**:
  - Put `Dockerfile` at repo root (requires broader context and more careful `.dockerignore` management).

### Decision: Docker Compose is for local smoke-testing (production-like)

- **Chosen**: Use Docker Compose for local smoke-testing in a production-like way (build with `api/Dockerfile` and run with exposed ports/env).
- **Rationale**: Keeps compose aligned with the distributed artifact and prevents compose from becoming a second dev environment.
- **Alternatives considered**:
  - Use compose for live development (volume mounts, tool installs in container).

## Notes / Follow-ups

- Existing repository files currently include both `api/Dockerfile` (dev) and `api/Dockerfile.prod` (prod). This feature will converge on a single `api/Dockerfile` and remove/replace the redundant one.
- Concrete command examples and step-by-step instructions belong in `quickstart.md` and shared docs, not the feature spec.

