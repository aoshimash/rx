# Implementation Plan: Docker Development Environment Setup

**Branch**: `003-docker-setup` | **Date**: 2026-01-25 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-docker-setup/spec.md`

**Note**: This plan focuses on configuration and documentation changes (no API contract changes).

## Summary

Provide a DevContainer-first development workflow and converge on a single production Dockerfile for distribution:

- DevContainer is the primary local development environment.
- The API has a single production Dockerfile located at `api/Dockerfile` (no separate development Dockerfile).
- Docker Compose is used for production-like local smoke-testing.

## Technical Context

**Language/Version**: Go 1.25+ (API component)  
**Primary Dependencies**: Docker, Docker Compose, VS Code/Cursor Dev Containers (developer workflow)  
**Storage**: N/A (this feature is environment/configuration only)  
**Testing**: Existing Go test suite + container smoke-test via Docker Compose  
**Target Platform**: Developer workstations (macOS/Linux/Windows) running Linux containers  
**Project Type**: Monorepo (`api/` + repository-level configuration)  
**Performance Goals**: N/A (environment setup); see success criteria in spec for time-based targets  
**Constraints**: DevContainer is the only iterative dev environment; single production Dockerfile at `api/Dockerfile`; Docker Compose is for smoke-testing only  
**Scale/Scope**: Single repository; configuration/documentation changes only (no new endpoints)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with OPTel Workout Constitution principles:

- **Dumb Backend**: No health calculations or business logic in backend - only data storage/retrieval
- **Domain-Driven Schema-First**: Domain models define business logic, OpenAPI spec defines API contract. Code generated from OpenAPI spec.
- **Terminology**: Use intuitive, commonly understood physical and fitness terminology
- **Clean Architecture**: Domain logic separated from infrastructure, repository pattern used
- **Monorepo Structure**: Component independence maintained, shared resources properly managed

If any principle is violated, document justification in Complexity Tracking section below.

## Project Structure

### Documentation (this feature)

```text
specs/003-docker-setup/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── README.md
├── checklists/
│   └── requirements.md
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
optel-workout/
├── .devcontainer/             # DevContainer configuration (added)
│   ├── devcontainer.json
│   ├── README.md
│   └── extensions.json
├── api/
│   └── Dockerfile             # Single production Dockerfile (replaced legacy Dockerfile.prod)
└── docker-compose.yml         # Smoke-testing (builds from ./api/Dockerfile)
```

**Structure Decision**: This is a monorepo configuration change that touches the repository root and the `api/` component; no new components are introduced.

## Complexity Tracking

No constitution violations for this feature.
