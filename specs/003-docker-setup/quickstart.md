# Quickstart: Docker Development Environment Setup

**Date**: 2026-01-25  
**Branch**: `003-docker-setup`  
**Spec**: [spec.md](./spec.md)

## Goals

- Develop locally using DevContainer (primary workflow)
- Build a distributable production image using a single Dockerfile (`api/Dockerfile`)
- Smoke-test locally using Docker Compose

## Prerequisites

- Docker + Docker Compose
- VS Code or Cursor + Dev Containers extension (for DevContainer workflow)

## DevContainer (recommended for development)

1. Open the repository folder in VS Code/Cursor.
2. Reopen the folder in DevContainer when prompted.
3. Use the integrated terminal inside the container to run the repository's standard workflows (generate / lint / test).

## Build production image (single Dockerfile)

Build the production image using the single Dockerfile under `api/`.

Example:

```bash
docker build -f api/Dockerfile api -t optel-workout:latest
```

## Smoke-test with Docker Compose

Start the stack:

```bash
docker compose up -d
```

Check logs:

```bash
docker compose logs -f
```

Stop:

```bash
docker compose down
```

## Troubleshooting

- **Port already in use**: Change the host port mapping in `docker-compose.yml`.
- **Docker daemon not running**: Start Docker Desktop / Docker Engine, then retry.
- **DevContainer build fails**: Rebuild the DevContainer from the editor command palette.

