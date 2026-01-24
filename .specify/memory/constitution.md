<!--
Sync Impact Report:
Version change: 2.0.0 (removed immutable logs principle)
Modified principles: Removed IV. Data Integrity - Immutable Logs
Added sections: N/A
Removed sections: IV. Data Integrity - Immutable Logs
Templates requiring updates:
  ✅ plan-template.md - Constitution Check section updated (removed immutable records)
  ✅ spec-template.md - Schema-First principle aligns with requirements
  ✅ tasks-template.md - Task organization aligns with principles (includes monorepo path conventions)
  ⚠ pending - checklist-template.md (no constitution references found)
  ⚠ pending - agent-file-template.md (no constitution references found)
Skills integration:
  ✅ optel-philosophy - Updated (removed immutable logs)
  ✅ optel-domain - Updated (removed immutable constraint)
  ✅ optel-go-standards - Referenced in Testing, Code Quality, Error Handling, Compliance Review
Follow-up TODOs: None
-->

# OPTel Workout Constitution

## Core Principles

### I. Dumb Backend Principle (NON-NEGOTIABLE)

The backend MUST store and retrieve telemetry data only. No business logic for "health" calculations, motivation scores, or wellness metrics. The system is stateless, headless, and agent-native - primary consumers are AI Agents (via MCP) or analysis tools, not human UIs.

**Rationale**: This ensures the system remains a pure telemetry backend, allowing analysis and interpretation to happen at the consumer level. It prevents coupling between data storage and business logic, maintaining flexibility for future analysis tools.

**Implementation Details**: For prohibited/permitted features list and detailed examples, see `.claude/skills/optel-philosophy/reference.md`.

### II. Domain-Driven Schema-First Development (NON-NEGOTIABLE)

Development follows a **Domain-Driven Schema-First** approach:
- **Domain models define business logic** - Domain entities (`internal/domain/`) contain business rules, validation, and invariants
- **OpenAPI spec defines API contract** - OpenAPI specification defines the HTTP API contract for clients
- **Code generation from OpenAPI** - Go types and server stubs are generated from the OpenAPI spec using oapi-codegen
- **Handlers bridge the gap** - HTTP handlers convert between OpenAPI types and domain models
- **Synchronization required** - Domain models and OpenAPI specs must be kept in sync

**Rationale**: This approach combines the benefits of domain modeling (clear business logic, testability) with schema-first API development (explicit contracts, code generation, documentation). It enables code generation while maintaining domain-driven design principles.

**Implementation Details**: 
- Domain models are defined in `internal/domain/` with business rules and validation
- OpenAPI spec is defined in `api/openapi/openapi.yaml` referencing domain models
- Code is generated from OpenAPI spec using `oapi-codegen`
- Handlers convert between OpenAPI types and domain models
- API versioning uses URL path pattern `/api/v1/`
- For domain model details (Workout, Program, Telemetry), see `.claude/skills/optel-domain/`
- For Go code generation patterns, see `.claude/skills/optel-go-standards/`
- For development workflow, see `docs/DEVELOPMENT.md`

### III. Terminology - Use Intuitive Physical Terms

Use intuitive, commonly understood physical and fitness terminology. The system records workouts, exercises, and physical exertion data. Terminology should be clear and accessible to developers and users.

**Rationale**: Intuitive terminology improves developer understanding and reduces cognitive overhead. The system remains a "dumb backend" that stores data without business logic, regardless of terminology choices.

**Implementation Details**: See `.claude/skills/optel-philosophy/reference.md` for terminology glossary and API design guidelines (resource naming, response structure).

### IV. Clean Architecture & Repository Pattern

Domain entities and business rules reside in the domain layer. Repository interfaces (ports) are defined in the domain layer. Implementations (adapters) are in the infrastructure layer. Handlers map HTTP requests/responses and delegate to domain logic.

**Rationale**: This separation of concerns enables testability, maintainability, and the ability to swap implementations (e.g., memory store → PostgreSQL) without changing domain logic.

### V. Monorepo Structure & Component Independence

This repository is a monorepo containing multiple components (api, mcp, frontend, mobile, infra). Each component MUST be independently deployable, testable, and maintainable. Components communicate via well-defined contracts (OpenAPI for API, MCP protocol for MCP server). Shared resources (OpenAPI specifications, documentation) are managed at the repository root level.

**Component Structure:**
- `api/` - REST API (Go) - Backend service
- `mcp/` - MCP Server (Python/Docker) - Runs on user's local machine
- `frontend/` - Web frontend (future)
- `mobile/` - Mobile applications (future)
- `infra/` - Infrastructure as Code (Terraform/Helm) (future)
- `docs/` - Shared documentation
- `.claude/skills/` - AI agent skills (shared)

**Rationale**: Monorepo structure enables coordinated development across components while maintaining clear boundaries. Component independence ensures that changes to one component don't break others, and each can evolve at its own pace with appropriate technology choices.

## Development Standards

### Code Generation

- OpenAPI specifications MUST be defined before implementation
- Generated code MUST NOT be manually edited (regenerate after spec changes)
- All API endpoints MUST be defined in `api/openapi/openapi.yaml`
- Components consuming the API (frontend, mobile, mcp) SHOULD generate client code from the same OpenAPI spec

### Testing Requirements

- Each component MUST have its own test suite
- API component: Unit tests for domain logic, integration tests for repositories, contract tests for endpoints
  - Table-driven tests are mandatory for API component (see `.claude/skills/optel-go-standards/`)
- Other components: Component-appropriate testing (e.g., unit tests, integration tests, E2E tests)
- Tests MUST be runnable in containerized environments where applicable
- Cross-component integration tests SHOULD be defined when components interact

### Code Quality

- Each component MUST have appropriate linting/formatting tools configured
- API component: golangci-lint with strict configuration (see `.claude/skills/optel-go-standards/` for details)
- Other components: Component-appropriate linting tools (e.g., ESLint for frontend, pylint for Python)
- All code MUST pass linting before commit
- Structured logging preferred across all components

### Error Handling

- Domain errors MUST be defined in the domain layer
- HTTP handlers translate domain errors to appropriate HTTP status codes
- Errors MUST be logged with structured logging
- For detailed error handling patterns, see `.claude/skills/optel-go-standards/`

## Technology Constraints

### Component-Specific Technologies

Each component in the monorepo MAY use different technologies appropriate to its purpose:

- **api/**: Go 1.25+, chi, oapi-codegen, golangci-lint, log/slog
- **mcp/**: Python (or Docker), MCP protocol
- **frontend/**: Technology TBD (React, Vue, etc.)
- **mobile/**: Technology TBD (React Native, Flutter, native, etc.)
- **infra/**: Terraform, Helm, Kubernetes manifests

### Shared Resources

- **OpenAPI Specification**: Defined in `api/openapi/openapi.yaml`, serves as contract for API consumers (frontend, mobile, mcp)
- **Documentation**: Shared in `docs/` directory
- **AI Agent Skills**: Shared in `.claude/skills/` directory

### Development Environment

- Component-specific development tools run in Docker containers where applicable
- Each component SHOULD have its own Makefile or build scripts
- Docker Compose for local development across components
- No requirement for local language installations (containerized development)

### Future Technologies (Planned)

- **Database**: PostgreSQL with TimescaleDB extension (for api/)
- **Infrastructure**: Kubernetes, Helm charts (in infra/)
- **Observability**: OpenTelemetry, Prometheus metrics (shared across components)

## Governance

### Amendment Process

1. Proposed amendments MUST be documented with rationale
2. Amendments affecting core principles (I-III) require explicit approval
3. Version MUST be incremented per semantic versioning:
   - **MAJOR**: Backward incompatible changes, principle removals/redefinitions
   - **MINOR**: New principles, materially expanded guidance
   - **PATCH**: Clarifications, wording fixes, non-semantic refinements
4. Constitution changes MUST be reflected in dependent templates and documentation

### Compliance Review

- All PRs/reviews MUST verify compliance with constitution principles
- Constitution Check section in implementation plans MUST be validated
- Complexity violations MUST be justified in plan.md Complexity Tracking section
- For detailed implementation guidelines, refer to AGENTS.md and `.claude/skills/`:
  - **optel-philosophy**: Dumb Backend principle details, terminology glossary, API design guidelines
  - **optel-domain**: Domain models (Workout, Program, Telemetry), validation rules, entity relationships
  - **optel-go-standards**: Go project structure, error handling patterns, testing standards, linting configuration

### Supremacy

This constitution supersedes all other practices and guidelines. When conflicts arise, constitution principles take precedence. Development decisions MUST align with these principles or provide explicit justification for exceptions.

**Version**: 2.0.0 | **Ratified**: 2026-01-24 | **Last Amended**: 2026-01-24
