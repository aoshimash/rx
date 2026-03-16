# Domain Model

## Three-Stage Training Lifecycle

```
Program → Plan → Log
```

1. **Program** — Reusable training template (no dates, no absolute weights)
2. **Plan** — Concrete training schedule derived from a Program (with dates and calculated weights)
3. **Log** — Actual workout record (what was actually performed)

## Entities

### Program

A reusable, RPE-based training template. Contains no dates and no absolute weights.

```go
type Program struct {
    ID          uuid.UUID       `json:"id"`
    Name        string          `json:"name"`
    Description *string         `json:"description,omitempty"`
    Notes       *string         `json:"notes,omitempty"`
    Metadata    json.RawMessage `json:"metadata,omitempty"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
    Entries     []ProgramEntry  `json:"entries,omitempty"`
}

type ProgramEntry struct {
    ID           uuid.UUID       `json:"id"`
    ProgramID    uuid.UUID       `json:"program_id"`
    Order        int             `json:"order"`
    ExerciseName string          `json:"exercise_name"`
    Sets         *int            `json:"sets,omitempty"`
    Reps         *int            `json:"reps,omitempty"`
    RPE          *int            `json:"rpe,omitempty"`
    Percent1RM   *float64        `json:"percent_1rm,omitempty"`
    Notes        *string         `json:"notes,omitempty"`
    Metadata     json.RawMessage `json:"metadata,omitempty"`
}
```

**Business Rules:**
- `name`: required, 1–200 chars
- `entries`: max 1000
- `rpe`: 1–10 (if set)
- `percent_1rm`: 0.0–1.0 (if set); fraction of 1RM (e.g. 0.80 = 80%)

**API Operations:**

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/programs | Create a new program |
| GET | /api/v1/programs | List programs |
| GET | /api/v1/programs/{id} | Get a specific program |
| PUT | /api/v1/programs/{id} | Update a program |
| DELETE | /api/v1/programs/{id} | Delete a program |
| POST | /api/v1/plans/from-program | Convert a Program into a Plan |

---

### Plan

A concrete training schedule, optionally derived from a Program via `ConvertProgramToPlan()`.

```go
type Plan struct {
    ID          uuid.UUID       `json:"id"`
    ProgramID   *uuid.UUID      `json:"program_id,omitempty"`
    Name        string          `json:"name"`
    Description *string         `json:"description,omitempty"`
    Notes       *string         `json:"notes,omitempty"`
    Metadata    json.RawMessage `json:"metadata,omitempty"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
    Entries     []PlanEntry     `json:"entries,omitempty"`
}

// DateOnly is a date without time, formatted as "2006-01-02" in JSON.
type DateOnly time.Time

type PlanEntry struct {
    ID           uuid.UUID       `json:"id"`
    PlanID       uuid.UUID       `json:"plan_id"`
    Order        int             `json:"order"`
    Date         *DateOnly       `json:"date,omitempty"`
    ExerciseName string          `json:"exercise_name"`
    Sets         *int            `json:"sets,omitempty"`
    Reps         *int            `json:"reps,omitempty"`
    LoadKg       *float64        `json:"load_kg,omitempty"`
    RPE          *int            `json:"rpe,omitempty"`
    Notes        *string         `json:"notes,omitempty"`
    Metadata     json.RawMessage `json:"metadata,omitempty"`
}
```

**Business Rules:**
- `name`: required, 1–200 chars
- `entries`: max 1000
- `load_kg`: >= 0; auto-rounded to 0.1 kg precision on validation (if set)
- `rpe`: 1–10 (if set)
- `date`: date-only (no time component), format `YYYY-MM-DD`

**API Operations:**

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/plans | Create a new plan |
| GET | /api/v1/plans | List plans |
| GET | /api/v1/plans/{id} | Get a specific plan |
| PUT | /api/v1/plans/{id} | Update a plan |
| DELETE | /api/v1/plans/{id} | Delete a plan |
| POST | /api/v1/plans/from-program | Convert Program to Plan |

---

### Log

A record of actual training performed.

```go
type Log struct {
    ID          uuid.UUID       `json:"id"`
    PlanID      *uuid.UUID      `json:"plan_id,omitempty"`
    PerformedAt time.Time       `json:"performed_at"`
    Notes       *string         `json:"notes,omitempty"`
    Metadata    json.RawMessage `json:"metadata,omitempty"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
    Entries     []LogEntry      `json:"entries"`
}

type LogEntry struct {
    ID             uuid.UUID       `json:"id"`
    LogID          uuid.UUID       `json:"log_id"`
    Order          int             `json:"order"`
    ExerciseName   string          `json:"exercise_name"`
    Sets           *int            `json:"sets,omitempty"`
    Reps           *int            `json:"reps,omitempty"`
    LoadKg         *float64        `json:"load_kg,omitempty"`
    RPE            *int            `json:"rpe,omitempty"`
    Notes          *string         `json:"notes,omitempty"`
    VideoObjectKey *string         `json:"video_object_key,omitempty"`
    Metadata       json.RawMessage `json:"metadata,omitempty"`
}
```

**Business Rules:**
- `performed_at`: must not be in the future
- `entries`: at least 1, max 500
- `load_kg`: >= 0; auto-rounded to 0.1 kg precision on validation (if set)
- `rpe`: 1–10 (if set)
- `video_object_key`: 1–500 chars; S3/object storage key for exercise video (if set)

**API Operations:**

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/logs | Create a new log |
| GET | /api/v1/logs | List logs |
| GET | /api/v1/logs/{id} | Get a specific log |
| PUT | /api/v1/logs/{id} | Update a log |
| DELETE | /api/v1/logs/{id} | Delete a log |

---

## Entity Relationships

```
┌─────────────┐
│   Program   │  (template: no dates, no absolute weights)
└──────┬──────┘
       │ ConvertProgramToPlan()
       ▼
┌─────────────┐
│    Plan     │  (schedule: dates + calculated weights)
└──────┬──────┘
       │ 1:N (optional reference)
       ▼
┌─────────────┐
│     Log     │  (record: what was actually performed)
└─────────────┘
```

---

## Program → Plan Conversion

`ConvertProgramToPlan(program, input)` calculates concrete weights from relative intensities:

- If `percent_1rm` is set: `load_kg = percent_1rm * target_weight`, rounded to increment
- If `percent_1rm` is not set: `load_kg = target_weight` (direct copy)
- `LoadIncrements` map allows equipment-aware rounding per exercise

**Input:**

```go
ConvertProgramToPlanInput{
    Name:           "Week 1",                              // optional, defaults to program name
    TargetWeights:  map[string]float64{"Squat": 100.0},   // 1RM or working weight per exercise
    LoadIncrements: map[string]float64{"Squat": 2.5},     // rounding increment per exercise
}
```

---

## Validation Rules Summary

| Entity | Field | Rule |
|--------|-------|------|
| Program | name | required, 1–200 chars |
| Program | entries | max 1000 |
| ProgramEntry | exercise_name | required, 1–200 chars |
| ProgramEntry | sets, reps | > 0 (if set) |
| ProgramEntry | rpe | 1–10 (if set) |
| ProgramEntry | percent_1rm | 0.0–1.0 (if set) |
| Plan | name | required, 1–200 chars |
| Plan | entries | max 1000 |
| PlanEntry | exercise_name | required, 1–200 chars |
| PlanEntry | sets, reps | > 0 (if set) |
| PlanEntry | load_kg | >= 0, rounded to 0.1 kg (if set) |
| PlanEntry | rpe | 1–10 (if set) |
| Log | performed_at | must not be in the future |
| Log | entries | 1–500 (at least one required) |
| LogEntry | exercise_name | required, 1–200 chars |
| LogEntry | sets, reps | > 0 (if set) |
| LogEntry | load_kg | >= 0, rounded to 0.1 kg (if set) |
| LogEntry | rpe | 1–10 (if set) |
| LogEntry | video_object_key | 1–500 chars (if set) |
