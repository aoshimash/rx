# Rx Philosophy - Detailed Reference

## The "Dumb Backend" Principle (Extended)

### What It Means

Rx is a telemetry backend for physical exertion data. It stores and retrieves workout records without business logic for health calculations or recommendations.

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
2. Provide CRUD operations for workout records
3. Expose time-series data for external analysis
4. Implement filtering and aggregation queries
5. Serve data to AI Agents for analysis (via MCP)

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
