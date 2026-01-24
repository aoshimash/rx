---
name: optel-philosophy
description: Core philosophy and constraints for OPTel project. Use when writing code, designing APIs, or making architectural decisions. Enforces "Dumb Backend" principle.
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

## Terminology

Use intuitive, commonly understood physical and fitness terminology. The system records workouts, exercises, and physical exertion data.

**Recommended Terms:**
- `Workout` - A completed unit of physical exertion
- `Exercise` - A specific movement or activity
- `MuscleGroup` - Body parts or muscle groups targeted
- `Metrics` - Measured data (intensity, volume, duration)

**Note:** The "Dumb Backend" principle remains unchanged - no business logic for health calculations, regardless of terminology used.

## Comment Style

Write clear, descriptive comments that explain what the code does.

```go
// ✅ GOOD
// CreateWorkout handles POST requests to create a new workout.
// Records physical exertion data: intensity (RPE), volume, duration, and timestamp.

// ✅ GOOD
// Record peak force output for the pectoral muscle group
```

## Additional Resources

For detailed guidelines, see [reference.md](reference.md).
