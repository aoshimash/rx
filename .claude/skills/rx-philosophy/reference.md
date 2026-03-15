# Rx Philosophy - Detailed Reference

## What the Backend Does

Rx is a plan-driven training management system. The backend stores programs (plans) and workouts (actuals), and exposes them via API and CLI for both human clients and AI agents.

### Prohibited Features

The backend MUST NOT:

1. Calculate "health scores" or "wellness indices"
2. Provide "motivation" or "encouragement" messages
3. Make recommendations about training intensity
4. Compare users or create leaderboards
5. Implement gamification (streaks, badges, achievements)

### Permitted Features

The backend MAY:

1. Store and retrieve Programs (training plans)
2. Store and retrieve Workouts (actual execution records)
3. Provide CRUD operations for all resources
4. Expose time-series telemetry data for external analysis
5. Implement filtering and aggregation queries
6. Serve data to AI Agents for analysis (via API or CLI)

### Interface Requirements

Rx is a full-stack project. All of the following are first-class interfaces:
- **Web** — Full-featured browser client
- **Mobile** — Full-featured native mobile client
- **REST API** — For AI agents and automation
- **CLI** — For scripting and local automation

Every feature MUST be accessible via API (and ideally CLI). Web and Mobile may provide richer UX on top of that.

## Terminology Glossary

Use intuitive, commonly understood terms for physical exertion data:

| Term | Description |
|------|-------------|
| Workout | A completed unit of physical exertion |
| Exercise | A specific movement or activity |
| MuscleGroup | Body parts or muscle groups targeted |
| Metrics | Measured data (intensity, volume, duration) |
| Rep | A single repetition of an exercise |
| Set | A group of repetitions |
| Rest Period | Non-active interval between sets |
| Personal Record | Maximum recorded value |
| Intensity | Rate of Perceived Exertion (RPE) or similar |
| Volume | Total load or work performed |

## Comment Examples

### API Handler Comments

```go
// ✅ GOOD
// CreateWorkout handles POST requests to create a new workout.
// Records physical exertion data: intensity (RPE), volume, duration, and timestamp.

// ✅ GOOD
// ListWorkouts retrieves workout records with optional filtering by date range.
```

### Domain Model Comments

```go
// ✅ GOOD
// Workout represents a completed unit of physical exertion.
// Contains telemetry data: intensity (RPE), volume (load), duration, and timestamp.

// ✅ GOOD
// Exercise represents a specific movement or activity within a workout.
```

### Error Comments

```go
// ✅ GOOD
// Returns ErrInvalidInput if workout data fails validation.
// Returns ErrNotFound if the specified workout ID does not exist.
```

## API Design Guidelines

### Resource Naming

```
✅ /api/v1/workouts
✅ /api/v1/workouts/{id}/metrics
✅ /api/v1/programs
✅ /api/v1/programs/{id}/schedule
```

### Response Structure

Responses should be data-centric, not user-centric:

```json
// ✅ GOOD
{
  "id": "wo-123",
  "timestamp": "2026-01-24T10:00:00Z",
  "duration_seconds": 3600,
  "intensity_rpe": 7,
  "volume_kg": 5000,
  "muscle_groups": ["pectoral", "deltoid"]
}

// ❌ BAD
{
  "workout_id": "w-123",
  "user_message": "Great workout! You lifted 5000kg!",
  "achievements": ["first_workout"],
  "streak_count": 5
}
```
