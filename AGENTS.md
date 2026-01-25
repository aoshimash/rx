# AGENTS.md

This file provides guidance for AI coding agents working on this repository.

## Project Overview

**OPTel (Open Physical Telemetry)** - An observability stack for the human physical layer.

This repository (`optel-training`) is the **Output/Training Management** component. It manages "Workouts" (physical exertion records) as an agent-native telemetry backend.

## Key Principles

1. **"Dumb Backend"** - No business logic for "health." Strictly stores and retrieves telemetry data.
2. **Domain-Driven Schema-First** - Domain models define business logic, OpenAPI spec defines API contract. Code is generated from OpenAPI spec.
3. **Use Latest Stable Versions** - Always use the latest stable version of containers and tools, but specify explicit version numbers (not `latest` tags) to ensure reproducible builds and consistent development environments.

For details, see `.claude/skills/optel-philosophy/`.

## Project Structure

```
optel-training/
├── api/                  # REST API (Go)
├── mcp/                  # MCP Server (runs on user's local machine)
├── frontend/             # Frontend (future)
├── infra/                # Terraform/Helm (future)
├── docs/                 # Documentation
└── .claude/skills/       # AI agent skills
```

## Skills Reference

| Skill | Description |
|-------|-------------|
| [optel-philosophy](.claude/skills/optel-philosophy/) | Core philosophy and constraints |
| [optel-domain](.claude/skills/optel-domain/) | Domain models (Workout, Program, Telemetry) |
| [optel-go-standards](.claude/skills/optel-go-standards/) | Go coding standards |

## Quick Reference

- **Language**: Go 1.25+
- **HTTP Server**: chi
- **OpenAPI**: oapi-codegen (Domain-Driven Schema-First)
- **Linter**: golangci-lint (strict)
- **Logging**: log/slog
- **Testing**: standard testing package

## Documentation

- [Architecture](docs/ARCHITECTURE.md)
- [Development Guide](docs/DEVELOPMENT.md)
