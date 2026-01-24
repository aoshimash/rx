---
name: optel-philosophy
description: Core philosophy and constraints for OPTel project. Use when writing code, designing APIs, or making architectural decisions. Enforces "Dumb Backend" principle and infrastructure terminology.
---

# OPTel Philosophy

## Core Principles

### 1. "Dumb Backend" Principle

- **No business logic for "health"** - No health scores, motivation calculations, or wellness metrics
- **Stateless & Headless** - No UI components. Pure API.
- **Agent-Native** - Primary consumers are AI Agents (via MCP) or analysis tools, not human UIs

### 2. Domain-Driven Schema-First Development

- **Domain models define business logic** - Domain entities (`internal/domain/`) contain business rules, validation, and invariants
- **OpenAPI spec defines API contract** - OpenAPI specification defines the HTTP API contract for clients
- **Code generation from OpenAPI** - Go types and server stubs are generated from the OpenAPI spec
- **Handlers bridge the gap** - HTTP handlers convert between OpenAPI types and domain models
- **Keep them synchronized** - Domain models and OpenAPI specs must stay in sync

**Workflow:**
1. Define domain model with business rules
2. Create OpenAPI spec referencing domain model
3. Generate code from OpenAPI spec
4. Implement handlers that convert between types

### 3. Data Integrity

- **Immutable Logs** - Past workload logs are "committed transactions." Never mutable.
- All changes create new records, not updates to existing ones

## Terminology Rules

Use **Infrastructure/Physics** terminology. Never use fitness terminology.

| ❌ Forbidden | ✅ Required |
|-------------|------------|
| Exercise | Workload |
| Workout | Task |
| Muscle | Unit / Subsystem |
| Sweat | Output |
| Fitness | Telemetry |
| Training | Stress |

## Comment Style

Write comments from the perspective of an SRE monitoring a system.

```go
// ❌ BAD
// Track user's bench press

// ✅ GOOD  
// Record peak force output for the pectoral subsystem
```

## Additional Resources

For detailed guidelines, see [reference.md](reference.md).
