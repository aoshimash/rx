---
name: rx-domain
description: Domain models for Rx project. Covers Workout, Program, and Telemetry entities with validation rules and relationships. Use when implementing entities, designing database schemas, writing API handlers, or when the code references Workout, Program, or Telemetry types.
---

# Rx Domain Models

## Core Entities

### 1. Workout (The Log)

A completed unit of physical exertion.

**Required Fields:**
- `id` - Unique identifier
- `timestamp` - When the workout was executed
- `duration` - Duration in seconds
- `intensity` - RPE (Rate of Perceived Exertion), 1-10 scale
- `volume` - Total load in kg or equivalent metric

**Optional Fields:**
- `muscle_groups` - Affected muscle groups (e.g., "pectoral", "quadriceps")
- `notes` - Freeform notes
- `program_id` - Reference to parent program (if scheduled)

### 2. Program (The Manifest)

A definition of intended workouts. A structured plan for physical training.

**Structure:** Recursive Tree
```
Program
└── Phase (e.g., "Hypertrophy Block")
    └── Day (e.g., "Day 1 - Upper")
        └── Block (e.g., "Primary Movement")
            └── Slot (e.g., "Bench Press 4x8")
```

### 3. Telemetry (The Metrics)

Time-series data points derived from Workouts. Used for trend analysis.

## Additional Resources

For detailed schema and relationships, see [reference.md](reference.md).
