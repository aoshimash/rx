# Architecture

## Overview

OPTel Workload is an agent-native telemetry backend that treats the human body as a production server. It follows Clean Architecture principles with a Schema-First API approach.

## System Context

```
┌─────────────────────────────────────────────────────────────┐
│                   User's Local Environment                   │
│  ┌─────────────────┐     ┌─────────────────────────────┐    │
│  │   AI Agent      │────▶│  MCP Server (Docker/Python) │    │
│  │  (Cursor等)     │     │  - Tools for workload ops   │    │
│  └─────────────────┘     └──────────────┬──────────────┘    │
└─────────────────────────────────────────┼───────────────────┘
                                          │ HTTP
                                          ▼
┌─────────────────────────────────────────────────────────────┐
│                        Remote Server                         │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              OPTel Workload API (REST)              │    │
│  │              - Pure API, no MCP logic               │    │
│  └──────────────────────────┬──────────────────────────┘    │
│                             │                                │
│                             ▼                                │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                    PostgreSQL                        │    │
│  │                (TimescaleDB for metrics)             │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

**Key Point:** MCP Server runs on the user's local machine, not on the backend server. The backend remains a pure REST API.

## Layer Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Handler Layer                            │
│  - HTTP handlers (chi)                                       │
│  - Request/Response mapping                                  │
│  - Error translation to HTTP status                          │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                             │
│  - Entities (Workload, Program, Telemetry)                   │
│  - Business rules (validation)                               │
│  - Repository interfaces (ports)                             │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                         │
│  - Repository implementations (adapters)                     │
│  - Database access                                           │
│  - External service integrations                             │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
optel-workload/
├── api/                          # REST API (Go)
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # Entry point
│   ├── internal/
│   │   ├── config/               # Configuration
│   │   ├── domain/               # Domain entities & rules
│   │   ├── handler/              # HTTP handlers
│   │   ├── repository/           # Repository interfaces
│   │   └── store/                # Repository implementations
│   │       ├── memory/           # In-memory (Phase 1)
│   │       └── postgres/         # PostgreSQL (future)
│   ├── pkg/
│   │   └── openapi/              # Generated OpenAPI code
│   └── openapi/
│       └── openapi.yaml          # API specification
├── mcp/                          # MCP Server (runs on user's local machine)
├── frontend/                     # Frontend (future)
├── infra/                        # Infrastructure (future)
│   ├── docker/
│   ├── helm/
│   └── terraform/
├── docs/                         # Documentation
└── .claude/skills/               # AI agent skills
```

## Data Flow

### Create Workload

```
1. Client sends POST /api/v1/workloads
2. Handler validates request format (OpenAPI generated)
3. Handler converts to domain entity
4. Domain validates business rules
5. Repository stores the workload
6. Handler returns response
```

### Query Workloads

```
1. Client sends GET /api/v1/workloads?from=...&to=...
2. Handler parses query parameters
3. Repository fetches with filters
4. Handler converts to response format
5. Handler returns response
```

## Key Design Decisions

### 1. Schema-First API

- OpenAPI specification is the source of truth
- Code is generated from the spec using oapi-codegen
- Changes to API always start with spec changes

### 2. Immutable Workloads

- Workloads cannot be modified after creation
- Corrections require creating new records
- Supports audit trail and data integrity

### 3. Repository Pattern

- Domain defines interfaces (ports)
- Implementations are in the store package (adapters)
- Easy to swap implementations (memory → PostgreSQL)

### 4. No Service Layer (Initially)

- For Phase 1, handlers call repositories directly
- Service layer will be added when business logic grows
- Keeps the codebase simple initially

## Technology Choices

| Component | Choice | Rationale |
|-----------|--------|-----------|
| HTTP | chi | Standard net/http compatible, oapi-codegen support |
| OpenAPI | oapi-codegen | Schema-First, generates chi handlers |
| Database | PostgreSQL | Reliable, TimescaleDB extension for time-series |
| Logging | log/slog | Standard library, structured logging |
| Linting | golangci-lint | Comprehensive, strict checks for AI-generated code |

## Future Considerations

### MCP Server

- **Runs on user's local machine** (not on the backend server)
- Distributed as Docker image or Python script
- Communicates with the OPTel API via HTTP
- Provides tools for AI agents to query and record workloads
- Backend remains a pure REST API with no MCP-specific logic

### Scaling

- Horizontal scaling via Kubernetes
- Database connection pooling
- Caching layer if needed (Redis)

### Observability

- OpenTelemetry for tracing
- Prometheus metrics export
- Structured logging for log aggregation
