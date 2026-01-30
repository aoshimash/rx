# Quickstart: OPTel Workout Data Structures

**Feature**: 001-data-structure  
**Date**: 2026-01-24

## Overview

OPTel Workout uses 6 core entities to store training data:

| Entity | Purpose |
|--------|---------|
| **Exercise** | Catalog of exercises (id + name) |
| **Workout** | A training session |
| **WorkoutEntry** | A performed exercise within a session |
| **Program** | A training plan |
| **ProgramNode** | A node in the program tree (recursive) |
| **TelemetryPoint** | Time-series metrics |

---

## Quick Examples

### 1. Create an Exercise

```json
POST /api/v1/exercises
{
  "name": "Bench Press",
  "muscle_groups": ["chest", "triceps", "shoulders"],
  "load_increment": 2.5
}
```

### 2. Record a Workout

```json
POST /api/v1/workouts
{
  "timestamp": "2026-01-24T17:00:00Z",
  "body_weight_kg": 74.5,
  "fatigue_level": 2,
  "program_context": ["Cycle 1", "Week 3", "Day 1"],
  "entries": [
    {
      "exercise_id": "550e8400-e29b-41d4-a716-446655440001",
      "entry_type": "top",
      "sets": 1,
      "reps": 1,
      "load_kg": 120.0,
      "rpe": 8
    },
    {
      "exercise_id": "550e8400-e29b-41d4-a716-446655440001",
      "entry_type": "main",
      "sets": 1,
      "reps": 5,
      "load_kg": 112.5,
      "rpe": 7
    },
    {
      "exercise_id": "550e8400-e29b-41d4-a716-446655440001",
      "entry_type": "backoff",
      "sets": 5,
      "reps": 5,
      "load_kg": 100.0,
      "rpe": 6,
      "planned_rest_seconds": 180,
      "performed_rest_seconds": 200
    }
  ]
}
```

### 3. Create a Program

```json
POST /api/v1/programs
{
  "name": "Full Body Hypertrophy 9 Weeks",
  "root_nodes": [
    {
      "name": "Cycle 1",
      "node_type": "cycle",
      "order": 0,
      "children": [
        {
          "name": "Week 1 (RPE6)",
          "node_type": "week",
          "order": 0,
          "children": [
            {
              "name": "Day 1",
              "node_type": "day",
              "order": 0,
              "children": [
                {
                  "name": "Bench Press Top",
                  "node_type": "exercise",
                  "order": 0,
                  "exercise_id": "550e8400-e29b-41d4-a716-446655440001",
                  "target_sets": 1,
                  "target_reps": 1,
                  "target_rpe": 6,
                  "percent_1rm": 0.85
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

### 4. Link Workout to Plan with Snapshot

```json
{
  "timestamp": "2026-01-24T17:00:00Z",
  "program_node_id": "node-day-1-uuid",
  "program_context": ["Cycle 1", "Week 1 (RPE6)", "Day 1"],
  "entries": [
    {
      "exercise_id": "550e8400-e29b-41d4-a716-446655440001",
      "entry_type": "top",
      "sets": 1,
      "reps": 1,
      "load_kg": 115.0,
      "rpe": 7,
      "program_node_id": "node-bp-top-uuid",
      "plan_snapshot": {
        "target_sets": 1,
        "target_reps": 1,
        "target_rpe": 6,
        "target_load_kg": 110.0,
        "percent_1rm": 0.85
      }
    }
  ]
}
```

---

## Key Rules

### Validation

| Field | Rule |
|-------|------|
| `timestamp` | Must not be in the future |
| `load_kg` | ≥ 0, precision 0.1 kg |
| `rpe` | 1-10 |
| `fatigue_level` | 1-5 |
| `sets`, `reps` | > 0 |
| `entry_type` | `top` / `main` / `backoff` / `accessory` |

### Units

| Measurement | Unit | Notes |
|-------------|------|-------|
| Weight/Load | kg | Convert lb → kg before storage |
| Rest | seconds | Integer |
| Time | RFC 3339 | UTC with timezone |
| %1RM | 0.0-1.0 | e.g., 0.75 = 75% |

### Entry Types

| Type | Description |
|------|-------------|
| `top` | Heavy single/double, near-max effort |
| `main` | Primary working sets |
| `backoff` | Volume sets at reduced weight |
| `accessory` | Supplemental exercises |

---

## Fatigue Scale

| Level | Meaning |
|-------|---------|
| 1 | Fully recovered, energetic |
| 2 | Mostly recovered |
| 3 | Normal / baseline |
| 4 | Somewhat fatigued |
| 5 | Very fatigued, low energy |

---

## Program Hierarchy

Programs use a **recursive tree structure** with user-defined node types:

```
Program
└── ProgramNode (cycle)
    └── ProgramNode (week)
        └── ProgramNode (day)
            └── ProgramNode (block)  [optional]
                └── ProgramNode (exercise)  ← leaf with prescription
```

Common `node_type` values:
- `"cycle"` - Mesocycle
- `"week"` - Training week
- `"day"` - Training day
- `"block"` - Group of exercises (e.g., "Upper", "Legs")
- `"exercise"` - Leaf node with prescription

---

## Files Reference

| File | Content |
|------|---------|
| [spec.md](./spec.md) | Feature specification |
| [plan.md](./plan.md) | Implementation plan |
| [research.md](./research.md) | Technical decisions |
| [data-model.md](./data-model.md) | Detailed entity definitions |
| [contracts/openapi-entities.yaml](./contracts/openapi-entities.yaml) | OpenAPI schemas |
