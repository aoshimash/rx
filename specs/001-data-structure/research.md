# Research: Define Core Data Structures

**Feature**: 001-data-structure  
**Date**: 2026-01-24

## Overview

This document captures technical decisions and research findings for the core data structure definitions.

---

## Decision 1: ID Generation Strategy

**Decision**: Use UUIDs (v4) for all entity identifiers.

**Rationale**:
- Globally unique without coordination (important for future distributed scenarios)
- No sequential guessing (minor security benefit)
- Standard Go support via `google/uuid` package
- Compatible with PostgreSQL `uuid` type

**Alternatives Considered**:
- Auto-increment integers: Rejected—exposes record counts, requires DB coordination
- ULIDs: Considered—sortable, but UUID is more widely supported and sortability is not critical here

---

## Decision 2: Recursive ProgramNode Implementation

**Decision**: Use adjacency list pattern with `parent_id` reference.

**Rationale**:
- Simple to implement in Go structs and SQL
- Easy to query single nodes or immediate children
- Sufficient for typical program depths (3-5 levels: cycle → week → day → block → exercise)

**Alternatives Considered**:
- Nested sets: Rejected—complex to maintain on updates
- Materialized path: Considered—good for path queries, but adjacency list is simpler for this use case
- Closure table: Rejected—overkill for shallow hierarchies

**Implementation Notes**:
```go
type ProgramNode struct {
    ID        uuid.UUID
    ParentID  *uuid.UUID  // nil for root nodes
    Name      string
    NodeType  string      // "cycle", "week", "day", "block", etc.
    Order     int         // ordering among siblings
    Children  []ProgramNode  // populated by query, not stored
    // ... prescription fields for leaf nodes
}
```

---

## Decision 3: Rest/Interval Representation

**Decision**: Store rest as duration in seconds (integer).

**Rationale**:
- Simple numeric comparison and arithmetic
- Sufficient precision (1 second) for training rest periods
- Easy to format for display (mm:ss)

**Alternatives Considered**:
- ISO 8601 duration strings: Rejected—harder to compare and compute
- Float minutes: Rejected—less intuitive than seconds

**Implementation Notes**:
- `planned_rest_seconds: *int` (nullable, nil = not specified)
- `performed_rest_seconds: *int` (nullable, nil = not recorded)
- Per-set overrides stored as JSON array: `[90, 120, 120]` (seconds per set)

---

## Decision 4: Plan Snapshot Storage

**Decision**: Store plan snapshot as embedded struct on WorkoutEntry.

**Rationale**:
- Self-contained: no need to join to Program to see what was planned
- Immutable snapshot: Program edits don't affect historical records
- Query simplicity: single table contains both planned and performed

**Alternatives Considered**:
- Separate snapshot table: Rejected—adds join complexity
- Store only reference, reconstruct from Program history: Rejected—requires Program versioning

**Implementation Notes**:
```go
type PlanSnapshot struct {
    ProgramNodeID  *uuid.UUID  // reference to original prescription
    TargetSets     *int
    TargetReps     *int
    TargetRPE      *int        // 1-10
    TargetLoad     *float64    // kg
    Percent1RM     *float64    // 0.0-1.0 (e.g., 0.75 = 75%)
    PlannedRest    *int        // seconds
}
```

---

## Decision 5: Entry Type (entry_type) Values

**Decision**: Use lowercase string enum: `"top"`, `"main"`, `"backoff"`, `"accessory"`.

**Rationale**:
- Human-readable in JSON and logs
- Extensible: new types can be added without schema changes
- Matches CSV source terminology (after normalization)

**Alternatives Considered**:
- Integer codes: Rejected—requires lookup table for readability
- Uppercase: Rejected—inconsistent with typical JSON conventions

---

## Decision 6: Timestamp Handling

**Decision**: Store all timestamps as UTC with timezone info (RFC 3339).

**Rationale**:
- Unambiguous time representation
- Standard format supported by Go `time.Time` and PostgreSQL `timestamptz`
- Client can convert to local time for display

**Implementation Notes**:
- `timestamp`: required, when session occurred (date/time)
- `session_start`, `session_end`: optional, precise start/end times
- `entry_start`, `entry_end`: optional, per-entry timing

---

## Decision 7: Load Precision Implementation

**Decision**: Store load as `float64` with application-level rounding to 1 decimal place.

**Rationale**:
- Go's `float64` provides sufficient precision for training loads
- Rounding enforced at validation/input layer, not storage
- Allows future precision increase if needed

**Implementation Notes**:
```go
func RoundLoad(kg float64) float64 {
    return math.Round(kg*10) / 10
}
```

---

## Decision 8: Fatigue Level Scale

**Decision**: Integer 1-5 with defined meanings.

| Level | Meaning |
|-------|---------|
| 1 | Fully recovered, energetic |
| 2 | Mostly recovered |
| 3 | Normal / baseline |
| 4 | Somewhat fatigued |
| 5 | Very fatigued, low energy |

**Rationale**:
- Simple enough to assess quickly before training
- Sufficient granularity for tracking trends
- Matches common subjective scales

---

## Open Items

None. All NEEDS CLARIFICATION items from Technical Context have been resolved.
