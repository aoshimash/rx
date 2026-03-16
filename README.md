# Rx

**Rx** — as in *prescription*. In training, it means the prescribed workout. This system is for trainees who follow a structured program and track their execution against the plan.

## Overview

Rx is a plan-driven workout tracking system for powerlifters. Plans are the manifest; logs are their execution results. No fitness tracking, no motivation scores — just data.

### Key Principles

- **Plan-Driven** — Plans define what to do; logs record what actually happened. Compare intent vs execution.
- **Dumb Backend** — Stores and retrieves data. No health scores, no motivation, no gamification.
- **Schema-First** — Domain models drive OpenAPI spec, code generated from spec.
- **Full-Stack, Multi-Interface** — Web, Mobile, REST API, and CLI are all first-class interfaces.
- **Agent-Native** — Designed for AI agents (MCP) and analysis tools, not just human UIs.

## Tech Stack

| Component | Technology |
|-----------|------------|
| API | Go 1.25+, chi, oapi-codegen |
| Web | Next.js 16 (App Router), TypeScript, shadcn/ui, Tailwind CSS 4, TanStack Query |
| Database | PostgreSQL 18 (TimescaleDB) |
| Object Storage | S3 / Cloudflare R2 (video uploads) |
| API Docs | Scalar |
| Infrastructure | Docker, Docker Compose |
| Tooling | aqua, go-task, pnpm, Biome |

## Quick Start

**Prerequisites:** Docker, Docker Compose, [aqua](https://aquaproj.github.io/)

```bash
# Install development tools
aqua install

# Start PostgreSQL
docker compose up -d postgres

# API
cd api
task generate   # Generate code from OpenAPI spec
task lint        # Run linter
task test        # Run tests
task run         # Start API server (localhost:8080)

# Web
cd web
pnpm install
pnpm dev         # Start dev server (localhost:3000)
```

### Smoke-Testing with Docker Compose

```bash
# Start all services (API + PostgreSQL + MinIO + API docs)
docker compose up -d

# Test the API
curl http://localhost:8080/api/v1/plans

# Browse API docs
open http://localhost:8081

# Stop
docker compose down
```

Set `HOST_PORT` to change the API port (e.g., `HOST_PORT=9090 docker compose up -d`).

> **Note:** Docker Compose is for smoke-testing only. Use aqua-managed tools on your host for development.

## Project Structure

```
rx/
├── api/                # REST API (Go)
│   ├── cmd/server/     # Entry point
│   ├── internal/       # Domain, handlers, repositories
│   ├── migrations/     # Database migrations
│   ├── openapi/        # OpenAPI spec (source of truth)
│   └── pkg/openapi/    # Generated code
├── web/                # Next.js web app
│   ├── app/            # App Router pages
│   ├── components/     # UI components (shadcn/ui)
│   ├── lib/            # API client, hooks (TanStack Query)
│   └── stores/         # Zustand stores
├── specs/              # Phase-based specifications
├── docs/               # Architecture & design docs
├── scripts/            # Utility scripts
└── .claude/            # AI agent skills (gh-cli)
```

## API Endpoints

All endpoints are under `/api/v1/`. See the [OpenAPI spec](api/openapi/openapi.yaml) or run `docker compose up -d api-docs` to browse interactively.

| Resource | Endpoints |
|----------|-----------|
| Plans | CRUD — list (with pagination), create, get, update, delete |
| Logs | CRUD — list (with pagination & date filtering), create, get, update, delete |
| Videos | Generate presigned upload/download URLs |

**Authentication:** Bearer token. MVP accepts any non-empty string.

## Development

### Task Commands (API)

```bash
task generate         # OpenAPI → Go code generation
task lint             # golangci-lint
task test             # Tests with race detection
task format           # Format code
task check            # format + lint + test
task validate-openapi # Validate OpenAPI spec
task run              # Dev server
task build            # Build binary
```

### pnpm Scripts (Web)

```bash
pnpm dev      # Dev server
pnpm build    # Production build
pnpm lint     # Biome lint
pnpm format   # Biome format
pnpm check    # lint + format check
```

### Pre-commit Hooks

Run `./scripts/setup-githooks.sh` to enable automatic format, lint, and test checks before each commit.

## Documentation

- [Philosophy](docs/PHILOSOPHY.md) — Core principles and constraints
- [Domain Model](docs/DOMAIN_MODEL.md) — Program / Plan / Log lifecycle
- [Go Standards](docs/GO_STANDARDS.md) — Go coding standards
- [Frontend Architecture](docs/FRONTEND_ARCHITECTURE.md) — Web/Mobile architecture and coding standards
- [Architecture](docs/ARCHITECTURE.md) — System design and layer structure
- [Development Guide](docs/DEVELOPMENT.md) — Detailed setup and workflow
- [Web UI Design](docs/WEB_UI_DESIGN.md) — Screen designs and interaction patterns
- [AGENTS.md](AGENTS.md) — AI agent guidance

## License

MIT License — see [LICENSE](LICENSE)
