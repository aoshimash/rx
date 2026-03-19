# Domain Model

## Three-Stage Training Lifecycle

```
Program → Plans → Logs
```

1. **Program** — Reusable training template containing multiple sessions (no dates, no absolute weights)
2. **Plan** — A single concrete workout prescription derived from one session of a Program (with date and calculated weights)
3. **Log** — Actual workout record for one training session (what was actually performed)

## Key Concepts

### Session (metadata concept)

A **Session** is a grouping of exercises for one workout within a Program. It exists only as metadata (`metadata.session` name) on ProgramEntries — it is **not** a database entity or a separate API resource.

When a Program is converted to Plans, each session becomes its own Plan.

### Plan = One Workout

A Plan represents **exactly one training session** with:
- A date (when the workout is scheduled)
- A session name (inherited from the Program session)
- Concrete weights (calculated from the Program's relative prescriptions)
- A reference back to the source Program (optional)

Plans derived from the same Program share the same `program_id`, which serves as the grouping mechanism.

### Log and Plan (1:1)

A Log records what was actually performed in one training session. A Log **optionally** references a Plan (`plan_id`), enabling comparison between planned and actual performance. Logs without a `plan_id` represent unplanned/ad-hoc training sessions.

## Entity Relationships

```
┌──────────────┐
│   Program    │  (template: sessions grouped by metadata)
└──────┬───────┘
       │ ConvertProgramToPlans()  — 1 Program → N Plans (one per session)
       ▼
┌──────────────┐
│    Plan      │  (1 workout: date + session_name + calculated weights)
└──────┬───────┘
       │ 1:1 (optional reference)
       ▼
┌──────────────┐
│     Log      │  (record: what was actually performed)
└──────────────┘
```

### Conversion Flow

```
Program (3 sessions: "Upper A", "Lower", "Upper B")
    │
    │ ConvertProgramToPlans(target_weights, dates)
    │
    ├──► Plan 1: "Upper A" on 2026-03-20, with calculated weights
    ├──► Plan 2: "Lower"   on 2026-03-22, with calculated weights
    └──► Plan 3: "Upper B" on 2026-03-24, with calculated weights
```

Conversion parameters (target weights, load increments) are stored in `Plan.metadata.conversion` for traceability.

## Design Decisions

- **Session is not a DB entity** — Grouping by `program_id` is sufficient; no intermediate entity needed
- **Plan grouping via `program_id`** — Plans derived from the same Program are related through this reference, not through a separate grouping entity
- **Log.plan_id is optional** — Supports both planned and ad-hoc training sessions

For schema details, see `api/internal/domain/`.
