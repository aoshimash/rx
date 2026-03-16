# Rx Philosophy

## What Rx Is

Rx is a **plan-driven training management system**. The core loop is:

1. **Plan** — AI creates a training program (via API or CLI)
2. **Execute** — User records actual training sessions (via Web, Mobile, API, or CLI)
3. **Review** — AI or user analyzes the results

This is **not a simple training log**. The plan is a first-class concept, and AI is the primary planner.

## Core Principles

### 1. API-First, Multi-Client

- **API and CLI are required** — AI agents and automation tools interact via API or CLI
- **Web and Mobile are first-class clients** — Full-featured interfaces, not samples or demos
- **No features that only exist in UI** — Every action available in Web/Mobile must also be accessible via API

### 2. AI as Planner

- AI creates and manages Programs (training templates)
- Humans execute and record the actual training sessions as Logs
- The backend stores both Plans and Logs without judging the gap between them

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

## What the Backend Does

### Prohibited Features

The backend MUST NOT:

1. Calculate "health scores" or "wellness indices"
2. Provide "motivation" or "encouragement" messages
3. Make recommendations about training intensity
4. Compare users or create leaderboards
5. Implement gamification (streaks, badges, achievements)

### Permitted Features

The backend MAY:

1. Store and retrieve Programs (training templates)
2. Store and retrieve Plans (concrete training schedules)
3. Store and retrieve Logs (actual training records)
4. Provide CRUD operations for all resources
5. Expose time-series telemetry data for external analysis
6. Implement filtering and aggregation queries
7. Serve data to AI Agents for analysis (via API or CLI)

### Interface Requirements

Rx is a full-stack project. All of the following are first-class interfaces:

- **Web** — Full-featured browser client
- **Mobile** — Full-featured native mobile client
- **REST API** — For AI agents and automation
- **CLI** — For scripting and local automation

Every feature MUST be accessible via API (and ideally CLI). Web and Mobile may provide richer UX on top of that.

## Terminology

Use intuitive, commonly understood physical and fitness terminology.

| Term | Description |
|------|-------------|
| Program | A reusable training template (created typically by AI) |
| Plan | A concrete training schedule derived from a Program |
| Log | A record of actual training performed |
| Exercise | A specific movement or activity |
| Rep | A single repetition of an exercise |
| Set | A group of repetitions |
| RPE | Rate of Perceived Exertion (1–10 scale) |
| 1RM | One-rep max: maximum weight for a single repetition |
| load_kg | Weight used for an exercise, in kilograms |

## Comment Style

Write clear, descriptive comments that explain what the code does.

### API Handler Comments

```go
// ✅ GOOD
// CreateLog handles POST requests to create a new training log.
// Records actual training data: exercises performed, sets, reps, load, and RPE.

// ✅ GOOD
// ListLogs retrieves training log records with optional filtering.
```

### Domain Model Comments

```go
// ✅ GOOD
// Log represents a record of actual training performed.
// Links to the Plan being executed (optional) and contains LogEntry items.

// ✅ GOOD
// Program represents a reusable, RPE-based training template.
// It contains no dates and no absolute weights.
```

### Error Comments

```go
// ✅ GOOD
// Returns ErrInvalidInput if log data fails validation.
// Returns ErrNotFound if the specified log ID does not exist.
```

## API Design Guidelines

### Resource Naming

```
✅ /api/v1/programs
✅ /api/v1/programs/{id}
✅ /api/v1/plans
✅ /api/v1/plans/{id}
✅ /api/v1/plans/from-program
✅ /api/v1/logs
✅ /api/v1/logs/{id}
```

### Response Structure

Responses should be data-centric, not user-centric:

```json
// ✅ GOOD
{
  "id": "...",
  "performed_at": "2026-01-24T10:00:00Z",
  "plan_id": "...",
  "entries": [
    {
      "exercise_name": "Squat",
      "sets": 3,
      "reps": 5,
      "load_kg": 100.0,
      "rpe": 8
    }
  ]
}

// ❌ BAD
{
  "log_id": "...",
  "user_message": "Great session! You squatted 100kg!",
  "achievements": ["personal_record"],
  "streak_count": 5
}
```
