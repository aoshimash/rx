# optel-training

An agent-native telemetry backend to monitor physical exertion. Part of Open Physical Telemetry (OPTel).

## Overview

OPTel (Open Physical Telemetry) is an observability stack for the human physical layer. This repository manages **Workouts** - records of physical exertion.

No fitness tracking, no motivation scores - just data.

## Key Principles

- **"Dumb Backend"** - Stores and retrieves telemetry data. No health calculations.
- **Schema-First** - API defined in OpenAPI, code generated from spec.
- **Agent-Native** - Designed for AI agents (MCP) and analysis tools, not human UIs.

## Quick Start

**Prerequisites:** Docker, Docker Compose (for smoke-testing), VS Code/Cursor with Dev Containers extension (for development)

### Development (Recommended)

Open this repository in VS Code/Cursor and reopen in DevContainer when prompted. All development tools (Go, oapi-codegen, golangci-lint) are automatically available.

```bash
# Inside DevContainer:
cd api
make generate      # Generate code from OpenAPI spec
make lint          # Run linter
make test          # Run tests
make run           # Start server
```

### Smoke-Testing with Docker Compose

For production-like local testing:

```bash
docker compose up -d
curl http://localhost:8080/api/v1/workouts
docker compose logs -f
docker compose down
```

See [Development Guide](docs/DEVELOPMENT.md) for details.

## Documentation

- [AGENTS.md](AGENTS.md) - AI agent guidance
- [Architecture](docs/ARCHITECTURE.md) - System design
- [Development Guide](docs/DEVELOPMENT.md) - Setup and workflow

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25+ |
| HTTP Server | chi |
| API Spec | OpenAPI 3.1 (oapi-codegen) |
| Database | PostgreSQL (future) |
| Infrastructure | Docker, Helm |

## Project Structure

```
optel-training/
├── api/              # REST API (Go)
├── mcp/              # MCP Server (runs on user's local machine)
├── frontend/         # Frontend (future)
├── infra/            # Infrastructure (future)
├── docs/             # Documentation
└── .claude/skills/   # AI agent skills
```

## Versioning

This project follows [Semantic Versioning](https://semver.org/) (SemVer):

- **MAJOR** version: Breaking changes (incompatible API changes, data structure changes)
- **MINOR** version: New features (backward compatible)
- **PATCH** version: Bug fixes (backward compatible)

API versioning uses URL path pattern `/api/v1/` for the current API version.

## License

MIT License - see [LICENSE](LICENSE)
