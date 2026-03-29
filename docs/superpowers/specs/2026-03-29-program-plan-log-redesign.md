# Program / Plan / Log Redesign

**Date:** 2026-03-29
**Status:** Draft

## Motivation

The current `Program` entity conflates two concerns: **template** (what to train) and **execution** (tracking progress). This causes:

- Rigid status lifecycle (`created → ongoing → completed/cancelled`) that doesn't match real training flexibility
- "Complete Program" requires all sessions logged, which is impractical (illness, competition prep adjustments)
- No way to reuse the same program across multiple training cycles
- No way to mix sessions from different programs into a single execution plan

This redesign separates these concerns into three distinct entities aligned with how powerlifters and plan-driven trainees actually operate.

## Design Overview

```
Program (persistent template)
├── Hierarchical structure: Groups → Sessions → Entries
├── May use relative prescriptions (%1RM, RPE targets)
├── Reusable across training cycles
└── No lifecycle status

Plan (persistent working document, single per user)
├── Flat ordered list of upcoming sessions
├── Can be created from Program(s) or manually
├── Freely editable: add, remove, reorder sessions
└── No lifecycle constraints (delete when done)

Log (persistent historical record)
├── Actual performance with per-set granularity
├── plan_snapshot: copy of what was planned
└── Independent historical record
```

### Daily User Flow

```
Open app → Plan (top page) → Pick session → Log workout → Done
```

Plan is the primary page. Programs are accessed from within Plan when adding sessions.

## Data Model

### Program (modified)

A reusable training template with hierarchical session organization.

```
Program {
  id: UUID
  name: string (1-200)
  notes?: string (max 5000)
  program_fields?: FieldDef[]    // Field definitions for session entries
  log_fields?: FieldDef[]        // Field definitions for log entries
  groups: ProgramGroup[]
  sessions: ProgramSession[]
  created_at: datetime
  updated_at: datetime
}
```

**Changes from current:**
- **Removed:** `status` (created/ongoing/completed/cancelled)
- **Removed:** `metadata` (redundant with program_fields/log_fields)
- **Added:** `groups` (hierarchical grouping)

### ProgramGroup (new)

Hierarchical grouping for sessions within a program. Supports nesting up to max_depth=2.

```
ProgramGroup {
  id: UUID
  program_id: UUID
  parent_group_id?: UUID         // Self-referential for nesting
  name: string (1-200)
  order: int (0-based)
  notes?: string (max 5000)
}
```

**Depth validation (backend-enforced):**
- depth=0: top-level group (e.g., "Block 1 - Accumulation")
- depth=1: nested group (e.g., "Week 1")
- depth=2 or greater: rejected by API

**Example structure:**
```
Program: "12-Week Meet Prep"
├── Group (depth=0): "Block 1 - Accumulation"
│   ├── Group (depth=1): "Week 1"
│   │   ├── Session: "Day 1 - Squat"
│   │   └── Session: "Day 2 - Bench"
│   └── Group (depth=1): "Week 2 - Deload"
│       └── Session: "Day 1 - Light"
└── Group (depth=0): "Block 2 - Intensification"
    └── ...
```

### ProgramSession (modified)

```
ProgramSession {
  id: UUID
  program_id: UUID
  group_id?: UUID                // Which group this session belongs to
  session_name: string (1-200)
  order: int (0-based)           // Order within its group (or top-level if no group)
  date?: DateOnly
  entries: ProgramSessionEntry[]
}
```

**Changes from current:**
- **Added:** `group_id` (optional reference to ProgramGroup)

### ProgramSessionEntry (unchanged)

```
ProgramSessionEntry {
  id: UUID
  session_id: UUID
  order: int (0-based)
  exercise_name: string (1-200)
  fields?: JSON                  // {sets, reps, weight_kg, rpe_target, percent_1rm, ...}
  notes?: string (max 2000)
}
```

### FieldDef (unchanged)

```
FieldDef {
  name: string
  type: string                   // "number", "text", "select", etc.
  options?: string[]             // Only for type=select
}
```

### Plan (new)

A single working document per user containing upcoming sessions. Freely editable, no lifecycle constraints.

```
Plan {
  id: UUID
  name?: string (max 200)       // Optional label (e.g., "2026 Spring Meet Prep")
  notes?: string (max 5000)
  sessions: PlanSession[]
  created_at: datetime
  updated_at: datetime
}
```

**Constraint:** One Plan per user. API returns 409 if creating a second Plan.

### PlanSession (new)

```
PlanSession {
  id: UUID
  plan_id: UUID
  session_name: string (1-200)
  order: int (0-based)
  date?: DateOnly
  source_program_id?: UUID       // Which Program this came from (for reference)
  source_session_id?: UUID       // Which ProgramSession this came from
  entries: PlanSessionEntry[]
}
```

### PlanSessionEntry (new)

```
PlanSessionEntry {
  id: UUID
  session_id: UUID
  order: int (0-based)
  exercise_name: string (1-200)
  fields?: JSON                  // Execution-ready values (absolute weights filled in)
  notes?: string (max 2000)
}
```

### Log (modified)

```
Log {
  id: UUID
  session_name: string (1-200)
  date: DateOnly
  plan_snapshot?: JSON           // Full copy of PlanSession at time of logging
  program_id?: UUID              // Source Program reference (for analytics)
  metadata?: JSON                // Session-level data (body_weight_kg, fatigue_level, etc.)
  entries: LogEntry[]
  notes?: string (max 5000)
  created_at: datetime
  updated_at: datetime
}
```

**Changes from current:**
- **Added:** `plan_snapshot` (JSON copy of the planned session including all entries)
- **Kept:** `metadata` (session-level data like body weight, fatigue)
- **Kept:** `program_id` (optional, for analytics queries)

**plan_snapshot structure:**
```json
{
  "session_name": "Day 1 - Squat",
  "source_program_id": "uuid",
  "source_session_id": "uuid",
  "entries": [
    {
      "exercise_name": "Back Squat",
      "fields": { "sets": 5, "reps": 3, "weight_kg": 140, "rpe_target": 8 },
      "notes": "80%1RM"
    }
  ]
}
```

### LogEntry (unchanged structure, but with child sets)

```
LogEntry {
  id: UUID
  log_id: UUID
  order: int (0-based)
  exercise_name: string (1-200)
  fields?: JSON                  // Exercise-level summary if needed
  notes?: string (max 2000)
  sets: LogSet[]                 // Per-set actual performance
}
```

### LogSet (new)

Per-set actual performance with optional video and notes.

```
LogSet {
  id: UUID
  entry_id: UUID
  set_number: int (1-based)
  fields: JSON                   // {reps, weight_kg, rpe, ...}
  video_url?: string             // Pre-signed URL reference
  notes?: string (max 2000)
}
```

**Example:**
```json
[
  { "set_number": 1, "fields": { "reps": 3, "weight_kg": 140, "rpe": 7 } },
  { "set_number": 2, "fields": { "reps": 3, "weight_kg": 140, "rpe": 8 }, "video_url": "..." },
  { "set_number": 3, "fields": { "reps": 2, "weight_kg": 140, "rpe": 10 }, "notes": "3rep目潰れた" }
]
```

## API Changes

### Program Endpoints (modified)

| Method | Path | Change |
|--------|------|--------|
| GET | /programs | Remove `status` filter parameter |
| POST | /programs | Remove `status`, `metadata` from request body |
| GET | /programs/{id} | Response includes `groups`, no `status`/`metadata` |
| PUT | /programs/{id} | Remove `status`, `metadata`; include `groups` |
| DELETE | /programs/{id} | Unchanged |
| ~~PATCH~~ | ~~/programs/{id}/status~~ | **Removed** (no status lifecycle) |
| ~~GET~~ | ~~/programs/{id}/logged-sessions~~ | **Removed** (logging is Plan/Log concern) |

### Plan Endpoints (new)

| Method | Path | Description |
|--------|------|-------------|
| GET | /plan | Get the user's Plan (404 if none exists) |
| POST | /plan | Create Plan (409 if already exists) |
| PUT | /plan | Update Plan (name, notes, sessions) |
| DELETE | /plan | Delete Plan |
| POST | /plan/sessions | Add session(s) to Plan |
| PUT | /plan/sessions/{id} | Update a specific PlanSession |
| DELETE | /plan/sessions/{id} | Remove a session from Plan |
| POST | /plan/expand-program/{program_id} | Expand a Program's sessions and **append** to the Plan |

**Note:** `/plan` (singular) rather than `/plans` since there is exactly one Plan per user.

**expand-program behavior:** Appends the Program's sessions to the existing Plan (does not replace). This allows mixing sessions from multiple Programs into one Plan.

### Log Endpoints (modified)

| Method | Path | Change |
|--------|------|--------|
| POST | /logs | Accept `plan_snapshot` in request body |
| GET | /logs/{id} | Response includes `plan_snapshot` and `entries[].sets` |
| PUT | /logs/{id} | Accept `plan_snapshot`, `entries[].sets` |

## Entity Relationships

```
┌──────────────────────┐
│      Program         │
│  (reusable template) │
├──────────────────────┤
│ id, name, notes      │
│ program_fields       │
│ log_fields           │
└──────┬───────────────┘
       │ 1:N
┌──────▼───────────────┐     ┌──────────────────────┐
│   ProgramGroup       │────►│  ProgramGroup         │  (nested, max depth=2)
│ id, name, order      │     │  parent_group_id      │
│ parent_group_id?     │     └──────────────────────┘
└──────────────────────┘
       │ 1:N
┌──────▼───────────────┐
│  ProgramSession      │
│ group_id?, order     │
│ session_name, date?  │
└──────┬───────────────┘
       │ 1:N
┌──────▼───────────────┐
│ ProgramSessionEntry  │
│ exercise_name, fields│
│ order, notes?        │
└──────────────────────┘


┌──────────────────────┐
│        Plan          │  (one per user)
│  name?, notes?       │
└──────┬───────────────┘
       │ 1:N
┌──────▼───────────────┐
│    PlanSession       │
│ session_name, order  │
│ date?                │
│ source_program_id?   │
│ source_session_id?   │
└──────┬───────────────┘
       │ 1:N
┌──────▼───────────────┐
│  PlanSessionEntry    │
│ exercise_name, fields│
│ order, notes?        │
└──────────────────────┘


┌──────────────────────┐
│        Log           │
│ session_name, date   │
│ plan_snapshot? (JSON)│
│ program_id?          │
│ metadata?            │
└──────┬───────────────┘
       │ 1:N
┌──────▼───────────────┐
│      LogEntry        │
│ exercise_name, fields│
│ order, notes?        │
└──────┬───────────────┘
       │ 1:N
┌──────▼───────────────┐
│       LogSet         │  (new)
│ set_number, fields   │
│ video_url?, notes?   │
└──────────────────────┘
```

## UI Navigation (high-level direction)

```
Plan (top page, /)
├── Upcoming sessions list (flat, ordered)
├── [+ Add sessions]
│   ├── From Program → browse templates → expand into Plan
│   └── Manual creation
├── Tap session → Start workout → Log screen
└── [Programs] → template management (create/edit/browse)

Logs (/logs)
├── Historical log list
└── Log detail: planned vs actual comparison

Settings (/settings)
```

**Key change:** Plan replaces Dashboard as the top page. Programs are accessed from within Plan, not as a standalone top-level route.

## Migration Path

### Database Changes

1. **programs table:** Drop `status` column, drop `metadata` column
2. **New tables:** `program_groups`, `plan`, `plan_sessions`, `plan_session_entries`, `log_sets`
3. **logs table:** Add `plan_snapshot` (JSONB, nullable) column
4. **Existing data:** Migrate current Programs (all become templates with no status). Existing Logs retain `program_id` reference.

### API Changes

1. Remove `/programs/{id}/status` endpoint
2. Remove `/programs/{id}/logged-sessions` endpoint
3. Remove `status` and `metadata` from Program schemas
4. Add `groups` to Program schemas
5. Add all Plan endpoints
6. Add `plan_snapshot` to Log schemas
7. Add `sets` to LogEntry schema

### Frontend Changes

1. Remove Program status UI (status badges, Complete/Cancel buttons)
2. Replace Dashboard with Plan as top page
3. Build Plan management UI (session list, add/remove/reorder)
4. Build Program → Plan expansion flow
5. Build LogSet UI (per-set recording with video/notes)
6. Build planned vs actual comparison in Log detail view

## Design Decisions

| Decision | Rationale |
|----------|-----------|
| Program has no status lifecycle | Powerlifters frequently modify, skip, or extend programs mid-cycle. Status enforcement doesn't match reality. |
| Program.metadata removed | Redundant with program_fields/log_fields. YAGNI. |
| ProgramGroup max_depth=2 | Covers Block > Week hierarchy. Deeper nesting creates UI complexity with no practical benefit. Enforced at backend. |
| Plan is singular per user | Plan = "what I'm doing next." Multiple plans would require choosing which is active, adding unnecessary complexity. |
| Plan is persistent but disposable | Users need Plan to survive app restarts and device switches, but Plan has no historical value (that's Log's job). |
| Log stores plan_snapshot as JSON copy | Plan sessions change over time. The snapshot preserves "what was planned" at the moment of logging, enabling planned vs actual analysis. |
| LogSet as first-class entity | Powerlifters need per-set video review and RPE annotation. JSON arrays inside fields would make this data hard to query. |
| Plan sessions are flat (no groups) | Plan is an execution queue, not a design document. Grouping belongs in Program templates. Flat lists are easier to manipulate (insert ad-hoc sessions, reorder). |
