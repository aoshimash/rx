# Data Model: Define Core Data Structures

**Feature**: 001-data-structure  
**Date**: 2026-01-24

## Entity Overview

```
┌─────────────┐       ┌─────────────┐
│   Exercise  │◄──────│ WorkoutEntry│
└─────────────┘       └──────┬──────┘
                             │ N:1
                             ▼
┌─────────────┐       ┌─────────────┐
│   Program   │◄──────│   Workout   │
└──────┬──────┘       └─────────────┘
       │ 1:N                 │
       ▼                     │ 1:N
┌─────────────┐       ┌─────────────┐
│ ProgramNode │       │TelemetryPt  │
└─────────────┘       └─────────────┘
     │ recursive
     ▼
```

---

## Entity: Exercise

**Description**: A catalog entry representing a canonical exercise.

### Fields

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `id` | UUID | ✅ | Primary key | Unique identifier |
| `name` | string | ✅ | 1-200 chars | Canonical exercise name |
| `description` | string | ❌ | max 2000 chars | Optional description |
| `aliases` | []string | ❌ | - | Alternative names for matching |
| `muscle_groups` | []string | ❌ | - | Target muscle groups |
| `load_increment` | float64 | ❌ | > 0 | Minimum weight increment (kg) |
| `created_at` | timestamp | ✅ | auto | Creation timestamp |
| `updated_at` | timestamp | ✅ | auto | Last update timestamp |

### Validation Rules
- `name` must be unique (case-insensitive)
- `load_increment` defaults to 2.5 if not specified

---

## Entity: Workout

**Description**: A completed training session containing performed entries.

### Fields

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `id` | UUID | ✅ | Primary key | Unique identifier |
| `timestamp` | timestamp | ✅ | ≤ now | When session occurred |
| `session_start` | timestamp | ❌ | ≤ session_end | Precise start time |
| `session_end` | timestamp | ❌ | ≥ session_start | Precise end time |
| `body_weight_kg` | float64 | ❌ | > 0 | Body weight at session |
| `fatigue_level` | int | ❌ | 1-5 | Subjective fatigue scale |
| `sleep_hours` | float64 | ❌ | 0-24 | Sleep before session |
| `condition_notes` | string | ❌ | max 2000 chars | Freeform condition notes |
| `program_node_id` | UUID | ❌ | FK → ProgramNode | Link to planned day/block |
| `program_context` | []string | ❌ | - | Hierarchy path snapshot |
| `notes` | string | ❌ | max 5000 chars | General session notes |
| `created_at` | timestamp | ✅ | auto | Creation timestamp |
| `updated_at` | timestamp | ✅ | auto | Last update timestamp |
| `entries` | []WorkoutEntry | ✅ | ≥ 1 | Performed entries (ordered) |

### Validation Rules
- `timestamp` must not be in the future
- `entries` must have at least one entry
- `session_start` ≤ `session_end` if both provided
- `fatigue_level` must be 1-5 if provided
- `body_weight_kg` rounded to 1 decimal place

---

## Entity: WorkoutEntry

**Description**: A single performed exercise entry within a workout session.

### Fields

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `id` | UUID | ✅ | Primary key | Unique identifier |
| `workout_id` | UUID | ✅ | FK → Workout | Parent workout |
| `order` | int | ✅ | ≥ 0 | Position in workout |
| `exercise_id` | UUID | ✅ | FK → Exercise | Exercise performed |
| `display_name` | string | ❌ | max 200 chars | Override display name |
| `entry_type` | string | ✅ | enum | top/main/backoff/accessory |
| `sets` | int | ✅ | > 0 | Number of sets performed |
| `reps` | int | ✅ | > 0 | Reps per set |
| `load_kg` | float64 | ✅ | ≥ 0 | Weight in kilograms |
| `rpe` | int | ✅ | 1-10 | Rate of Perceived Exertion |
| `entry_start` | timestamp | ❌ | - | When entry started |
| `entry_end` | timestamp | ❌ | - | When entry ended |
| `planned_rest_seconds` | int | ❌ | ≥ 0 | Planned rest between sets |
| `performed_rest_seconds` | int | ❌ | ≥ 0 | Actual rest between sets |
| `per_set_rest_overrides` | []int | ❌ | - | Per-set rest (seconds) |
| `program_node_id` | UUID | ❌ | FK → ProgramNode | Link to prescription |
| `plan_snapshot` | PlanSnapshot | ❌ | - | Snapshot of planned values |
| `notes` | string | ❌ | max 2000 chars | Entry-specific notes |

### Validation Rules
- `entry_type` must be one of: `top`, `main`, `backoff`, `accessory`
- `sets` must be > 0
- `reps` must be > 0
- `load_kg` must be ≥ 0 (0 allowed for bodyweight)
- `load_kg` rounded to 1 decimal place
- `rpe` must be 1-10

---

## Embedded: PlanSnapshot

**Description**: Snapshot of planned values at execution time.

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `program_node_id` | UUID | ❌ | Reference to original prescription |
| `target_sets` | int | ❌ | Planned sets |
| `target_reps` | int | ❌ | Planned reps |
| `target_rpe` | int | ❌ | Planned RPE (1-10) |
| `target_load_kg` | float64 | ❌ | Planned load |
| `percent_1rm` | float64 | ❌ | Planned %1RM (0.0-1.0) |
| `planned_rest_seconds` | int | ❌ | Planned rest |

---

## Entity: Program

**Description**: A training program containing a recursive tree of nodes.

### Fields

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `id` | UUID | ✅ | Primary key | Unique identifier |
| `name` | string | ✅ | 1-200 chars | Program name |
| `description` | string | ❌ | max 2000 chars | Program description |
| `created_at` | timestamp | ✅ | auto | Creation timestamp |
| `updated_at` | timestamp | ✅ | auto | Last update timestamp |
| `root_nodes` | []ProgramNode | ❌ | - | Top-level nodes |

---

## Entity: ProgramNode

**Description**: A node in the program tree (cycle, week, day, block, or exercise prescription).

### Fields

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `id` | UUID | ✅ | Primary key | Unique identifier |
| `program_id` | UUID | ✅ | FK → Program | Parent program |
| `parent_id` | UUID | ❌ | FK → ProgramNode | Parent node (null for root) |
| `name` | string | ✅ | 1-200 chars | Node name |
| `node_type` | string | ✅ | - | User-defined type |
| `order` | int | ✅ | ≥ 0 | Order among siblings |
| `children` | []ProgramNode | ❌ | - | Child nodes |

### Prescription Fields (for leaf nodes)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `exercise_id` | UUID | ❌ | FK → Exercise |
| `target_sets` | int | ❌ | Planned sets |
| `target_reps` | int | ❌ | Planned reps per set |
| `target_rpe` | int | ❌ | Target RPE (1-10) |
| `percent_1rm` | float64 | ❌ | Target %1RM (0.0-1.0) |
| `planned_rest_seconds` | int | ❌ | Rest between sets |
| `muscle_groups` | []string | ❌ | Target muscles |
| `notes` | string | ❌ | Prescription notes |

### Common node_type Values
- `"cycle"` - Training cycle (mesocycle)
- `"week"` - Training week
- `"day"` - Training day
- `"block"` - Exercise group within a day
- `"exercise"` - Leaf node with prescription

---

## Entity: TelemetryPoint

**Description**: A time-series metric data point.

### Fields

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `id` | UUID | ✅ | Primary key | Unique identifier |
| `timestamp` | timestamp | ✅ | - | When metric was recorded |
| `metric_name` | string | ✅ | 1-100 chars | Metric identifier |
| `value` | float64 | ✅ | - | Numeric value |
| `unit` | string | ✅ | 1-50 chars | Unit of measurement |
| `workout_id` | UUID | ❌ | FK → Workout | Optional workout link |
| `created_at` | timestamp | ✅ | auto | Creation timestamp |

### Common Metrics
- `daily_volume_kg` - Total volume per day
- `weekly_volume_kg` - Total volume per week
- `session_duration_minutes` - Workout duration
- `body_weight_kg` - Daily body weight

---

## Indexes (Future PostgreSQL)

```sql
-- Exercise
CREATE UNIQUE INDEX idx_exercise_name_lower ON exercise (LOWER(name));

-- Workout
CREATE INDEX idx_workout_timestamp ON workout (timestamp DESC);
CREATE INDEX idx_workout_program_node ON workout (program_node_id);

-- WorkoutEntry
CREATE INDEX idx_entry_workout ON workout_entry (workout_id, "order");
CREATE INDEX idx_entry_exercise ON workout_entry (exercise_id);

-- ProgramNode
CREATE INDEX idx_node_program ON program_node (program_id);
CREATE INDEX idx_node_parent ON program_node (parent_id, "order");

-- TelemetryPoint
CREATE INDEX idx_telemetry_metric_time ON telemetry_point (metric_name, timestamp DESC);
CREATE INDEX idx_telemetry_workout ON telemetry_point (workout_id);
```
