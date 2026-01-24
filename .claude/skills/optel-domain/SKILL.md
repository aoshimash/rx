---
name: optel-domain
description: Domain models for OPTel project. Covers Workload, Program, and Telemetry entities with validation rules and relationships. Use when implementing entities, designing database schemas, writing API handlers, or when the code references Workload, Program, or Telemetry types.
---

# OPTel Domain Models

## Core Entities

### 1. Workload (The Log)

A completed unit of physical exertion. Immutable once created.

**Required Fields:**
- `id` - Unique identifier
- `timestamp` - When the workload was executed
- `duration` - Duration in seconds
- `intensity` - RPE (Rate of Perceived Exertion), 1-10 scale
- `volume` - Total load in kg or equivalent metric

**Optional Fields:**
- `subsystems` - Affected body subsystems (e.g., "pectoral", "quadriceps")
- `notes` - Freeform telemetry notes
- `program_id` - Reference to parent program (if scheduled)

### 2. Program (The Manifest)

A definition of intended workloads. Think of it as a Job Manifest in K8s.

**Structure:** Recursive Tree
```
Program
└── Phase (e.g., "Hypertrophy Block")
    └── Day (e.g., "Day 1 - Upper")
        └── Block (e.g., "Primary Movement")
            └── Slot (e.g., "Bench Press 4x8")
```

### 3. Telemetry (The Metrics)

Time-series data points derived from Workloads. Used for trend analysis.

## Additional Resources

For detailed schema and relationships, see [reference.md](reference.md).
