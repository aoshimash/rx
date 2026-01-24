# OPTel Domain Models - Detailed Reference

## Workout Entity

### Schema Definition

```go
type Workout struct {
    ID          string     `json:"id"`
    Timestamp   time.Time  `json:"timestamp"`
    Duration    int        `json:"duration_seconds"`
    Intensity   int        `json:"intensity_rpe"`      // 1-10 scale
    Volume      float64    `json:"volume_kg"`
    MuscleGroups []string   `json:"muscle_groups,omitempty"`
    Notes       string     `json:"notes,omitempty"`
    ProgramID   *string    `json:"program_id,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
}
```

### Business Rules

1. **Intensity Range** - RPE must be 1-10. Values outside this range are invalid.
2. **Volume** - Must be non-negative. Zero is valid (e.g., mobility work).
3. **Timestamp** - Must not be in the future.

### API Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/workouts | Create a new workout |
| GET | /api/v1/workouts | List workouts (with filters) |
| GET | /api/v1/workouts/{id} | Get a specific workout |
| DELETE | /api/v1/workouts/{id} | Soft-delete a workout |
| PUT | /api/v1/workouts/{id} | Update a workout |
| PATCH | /api/v1/workouts/{id} | Partially update a workout |

## Program Entity

### Schema Definition

```go
type Program struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    Phases      []Phase   `json:"phases"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Phase struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Order       int    `json:"order"`
    Days        []Day  `json:"days"`
}

type Day struct {
    ID     string  `json:"id"`
    Name   string  `json:"name"`
    Order  int     `json:"order"`
    Blocks []Block `json:"blocks"`
}

type Block struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Order int    `json:"order"`
    Slots []Slot `json:"slots"`
}

type Slot struct {
    ID         string   `json:"id"`
    Name       string   `json:"name"`
    Order      int      `json:"order"`
    MuscleGroups []string `json:"muscle_groups"`
    TargetSets int      `json:"target_sets,omitempty"`
    TargetReps int      `json:"target_reps,omitempty"`
    TargetRPE  int      `json:"target_rpe,omitempty"`
}
```

### Tree Structure Rationale

The recursive tree structure supports complex periodization:

1. **Linear Programs** - Single phase with sequential days
2. **Block Periodization** - Multiple phases with different focuses
3. **Undulating Programs** - Varying intensity within weeks
4. **Custom Structures** - Flexible nesting for any methodology

## Telemetry Entity

### Schema Definition

```go
type TelemetryPoint struct {
    Timestamp  time.Time `json:"timestamp"`
    MetricName string    `json:"metric_name"`
    Value      float64   `json:"value"`
    Unit       string    `json:"unit"`
    WorkoutID  string    `json:"workout_id,omitempty"`
}
```

### Derived Metrics

Telemetry points can be derived from Workouts:

| Metric | Derivation | Unit |
|--------|-----------|------|
| `daily_volume` | Sum of volume per day | kg |
| `weekly_volume` | Sum of volume per week | kg |
| `avg_intensity` | Average RPE per period | RPE |
| `frequency` | Workouts per week | count |

## Entity Relationships

```
┌─────────────┐
│   Program   │
└──────┬──────┘
       │ 1:N
       ▼
┌─────────────┐
│   Phase     │
└──────┬──────┘
       │ 1:N
       ▼
┌─────────────┐
│    Day      │
└──────┬──────┘
       │ 1:N
       ▼
┌─────────────┐
│   Block     │
└──────┬──────┘
       │ 1:N
       ▼
┌─────────────┐
│    Slot     │
└─────────────┘

┌─────────────┐      ┌─────────────┐
│   Workout   │─────▶│  Telemetry  │
└─────────────┘ 1:N  └─────────────┘
       │
       │ N:1 (optional)
       ▼
┌─────────────┐
│   Program   │
└─────────────┘
```

## Validation Rules Summary

| Entity | Field | Rule |
|--------|-------|------|
| Workout | intensity_rpe | 1 ≤ value ≤ 10 |
| Workout | volume_kg | value ≥ 0 |
| Workout | duration_seconds | value > 0 |
| Workout | timestamp | value ≤ now |
| Program | phases | at least 1 phase |
| Phase | days | at least 1 day |
| Slot | target_rpe | 1 ≤ value ≤ 10 (if set) |
