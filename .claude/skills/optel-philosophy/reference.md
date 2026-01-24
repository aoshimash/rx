# OPTel Philosophy - Detailed Reference

## The "Dumb Backend" Principle (Extended)

### What It Means

OPTel is infrastructure monitoring for the human body. Think of it as Prometheus/Grafana for physical systems, not a fitness app.

### Prohibited Features

The backend MUST NOT:

1. Calculate "health scores" or "wellness indices"
2. Provide "motivation" or "encouragement" messages
3. Make recommendations about training intensity
4. Compare users or create leaderboards
5. Implement gamification (streaks, badges, achievements)
6. Store or process subjective feelings (mood, energy level)

### Permitted Features

The backend MAY:

1. Store raw telemetry data (intensity, volume, duration, timestamp)
2. Provide CRUD operations for workload records
3. Expose time-series data for external analysis
4. Implement filtering and aggregation queries
5. Serve data to AI Agents for analysis (via MCP)

## Infrastructure Metaphor Dictionary

| Fitness Term | Infrastructure Term | Context |
|-------------|---------------------|---------|
| Exercise | Workload | A unit of physical exertion |
| Workout | Task / Job | A scheduled or completed work unit |
| Training Session | Execution Cycle | A period of active work |
| Muscle Group | Subsystem / Unit | A functional component |
| Rep | Iteration | A single cycle |
| Set | Batch | A group of iterations |
| Rest Period | Cooldown / Idle | Non-active interval |
| Personal Record | Peak Output | Maximum recorded value |
| Fatigue | Load / Stress | Accumulated strain |
| Recovery | Reset / Restore | Return to baseline |
| Progress | Delta / Trend | Change over time |
| Goal | Target / Threshold | Desired metric value |

## Comment Examples

### API Handler Comments

```go
// ❌ BAD
// CreateWorkout handles POST requests to create a new workout for the user.

// ✅ GOOD
// CreateWorkload handles POST requests to register a new workload execution.
// Records telemetry data for the specified subsystems.
```

### Domain Model Comments

```go
// ❌ BAD
// Workout represents a user's exercise session with all their lifts.

// ✅ GOOD
// Workload represents a completed unit of physical exertion.
// Contains telemetry data: intensity (RPE), volume (load), duration, and timestamp.
```

### Error Comments

```go
// ❌ BAD
// Returns error if the workout is invalid or user doesn't exist.

// ✅ GOOD
// Returns ErrInvalidInput if telemetry data fails validation.
// Returns ErrNotFound if the specified workload ID does not exist.
```

## API Design Guidelines

### Resource Naming

```
✅ /api/v1/workloads
✅ /api/v1/workloads/{id}/telemetry
✅ /api/v1/programs
✅ /api/v1/programs/{id}/schedule

❌ /api/v1/workouts
❌ /api/v1/exercises
❌ /api/v1/training-sessions
```

### Response Structure

Responses should be data-centric, not user-centric:

```json
// ✅ GOOD
{
  "id": "wl-123",
  "timestamp": "2026-01-24T10:00:00Z",
  "duration_seconds": 3600,
  "intensity_rpe": 7,
  "volume_kg": 5000,
  "subsystems": ["pectoral", "deltoid"]
}

// ❌ BAD
{
  "workout_id": "w-123",
  "user_message": "Great workout! You lifted 5000kg!",
  "achievements": ["first_workout"],
  "streak_count": 5
}
```
