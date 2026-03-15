---
name: rx-philosophy
description: Core philosophy and constraints for Rx project. Use when writing code, designing APIs, or making architectural decisions.
---

# Rx Philosophy

## What Rx Is

Rx is a **plan-driven training management system**. The core loop is:

1. **Plan** — AI creates a training program (via API or CLI)
2. **Execute** — User records actual workouts (via Web, Mobile, API, or CLI)
3. **Review** — AI or user analyzes the results

This is **not a simple training log**. The plan is a first-class concept, and AI is the primary planner.

## Core Principles

### 1. API-First, Multi-Client

- **API and CLI are required** — AI agents and automation tools interact via API or CLI
- **Web and Mobile are first-class clients** — Full-featured interfaces, not samples or demos
- **No features that only exist in UI** — Every action available in Web/Mobile must also be accessible via API

### 2. AI as Planner

- AI creates and manages Programs (training plans)
- Humans execute and record the actual workouts
- The backend stores both plans and actuals without judging the gap between them

### 3. No Opinionated Health Logic

- **No health scores, wellness indices, or motivation messages**
- **No recommendations or comparisons** — the backend stores and retrieves, AI analyzes
- **No gamification** — no streaks, badges, or achievements
- Raw data in, raw data out

### 4. Domain-Driven Schema-First Development

- **Domain models define business rules** — `internal/domain/` contains validation and invariants
- **OpenAPI spec defines the API contract** — single source of truth for HTTP API
- **Code generation from OpenAPI** — Go types and server stubs are generated
- **Handlers bridge the gap** — convert between OpenAPI types and domain models
- **Keep them synchronized** — domain models and OpenAPI specs must stay in sync

**Workflow:**
1. Define domain model with business rules
2. Create OpenAPI spec referencing domain model
3. Generate code from OpenAPI spec
4. Implement handlers that convert between types

## Terminology

Use intuitive, commonly understood physical and fitness terminology.

- `Program` - A training plan created (typically by AI)
- `Workout` - A completed unit of physical exertion (actual execution)
- `Exercise` - A specific movement or activity
- `MuscleGroup` - Body parts or muscle groups targeted

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
