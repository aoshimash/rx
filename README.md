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

**Prerequisites:** Docker, Docker Compose, Make, curl

```bash
cd api
make docker-up      # Start development container
make generate       # Generate code from OpenAPI spec
make run            # Start server (runs in container)
# In another terminal:
curl http://localhost:8080/api/v1/workouts
```

All Go development tools run inside Docker. See [Development Guide](docs/DEVELOPMENT.md) for details.

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

## License

MIT License - see [LICENSE](LICENSE)
