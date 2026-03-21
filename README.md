# Rx

**Rx** — as in *prescription*. In training, it means the prescribed workout. This system is for trainees who follow a structured program and track their execution against the plan.

## Overview

Rx is a training data backend. **Planning and analysis happen outside Rx** — by humans with spreadsheets, AI agents like Claude Code, or any script you write. Rx does one thing: make it as easy as possible to store and retrieve structured training data.

The Web UI is optimized for frictionless data entry, not for discovery or engagement. Everything in the UI is equally accessible via REST API and CLI — so AI agents, scripts, and automation tools are first-class citizens alongside human users.

### Key Principles

- **Bring Your Own Planning** — Planning and analysis are done externally — by humans, AI agents (e.g., Claude Code), or scripts. Rx stores and serves data; it does not interpret it.
- **Frictionless Data Entry** — The Web UI exists to make entering training data as fast and low-friction as possible.
- **Dumb Backend** — No health scores, no recommendations, no gamification. Raw data in, raw data out.
- **API-First** — Every feature is accessible via REST API. Web and CLI are clients of the same API.
- **Schema-First** — Domain models drive the OpenAPI spec; code is generated from the spec.

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
task check            # generate + format + lint + test
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
