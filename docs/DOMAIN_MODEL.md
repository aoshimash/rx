# Domain Model

## Three-Tier Training Lifecycle

```
ProgramTemplate → Program → Log
```

1. **ProgramTemplate** — Reusable training blueprint containing sessions with relative prescriptions (RPE / %1RM). No dates, no absolute weights.
2. **Program** — Concrete training program derived from a template (or created manually), containing sessions with scheduled exercises and absolute weights. Tracks training progress with a `status` (active / completed).
3. **Log** — Actual workout record for one training session (what was actually performed).

## Key Concepts

### ProgramTemplate

A **ProgramTemplate** is a reusable blueprint that defines exercise prescriptions using relative intensity (RPE, %1RM). It is the "program design" layer.

- Entries are stored flat with `metadata.session` grouping them into sessions in the frontend
- Can be archived (soft delete), duplicated, and used to generate Programs
- Has no execution semantics — it is purely a template

### Program

A **Program** is a concrete training instance derived from a ProgramTemplate (via the `/generate` endpoint) or created directly. It:

- Contains **ProgramSessions** (named workout days with an order and optional date)
- Each session has **ProgramSessionEntries** (exercises with absolute weights, not relative prescriptions)
- Has a `status`: `created` (registered, not yet started), `ongoing` (in progress), `completed` (all sessions logged and confirmed), or `cancelled` (stopped mid-way)
- Program status transitions are explicit user actions: `created` → `ongoing` (Start), `ongoing` → `completed` (Complete, requires all sessions logged), `ongoing` → `cancelled` (Cancel). Both `completed` and `cancelled` are terminal states.

### Log and Program (many:1)

A Log records what was actually performed in one training session. A Log **optionally** references a Program (`program_id`) and a session name (`session_name`), enabling tracking of which program session was performed. Logs without a `program_id` represent unplanned/ad-hoc training sessions.

## Entity Relationships

```
┌───────────────────┐
│  ProgramTemplate  │  (blueprint: RPE/% prescriptions grouped by session)
└────────┬──────────┘
         │ generate() — template → concrete Program with absolute weights
         ▼
┌───────────────────┐     ┌──────────────────────┐
│     Program       │────►│   ProgramSession      │  (named training day)
│  (status: active/ │     └──────────┬───────────┘
│   completed)      │               │
└────────┬──────────┘     ┌──────────▼───────────┐
         │                │ ProgramSessionEntry   │  (exercise with load_kg)
         │                └──────────────────────┘
         │ 1:N (optional reference)
         ▼
┌───────────────────┐
│       Log         │  (record: what was actually performed)
│  program_id       │
│  session_name     │
└───────────────────┘
```

### Generation Flow

```
ProgramTemplate (3 sessions: "W1D1", "W1D2", "W1D3")
    │
    │ POST /program-templates/{id}/generate
    │   { target_weights: { "Squat": 100, "Bench": 80 }, load_increments: {...} }
    │
    └──► Program (status: active) with sessions:
          ├── ProgramSession "W1D1": entries with calculated load_kg
          ├── ProgramSession "W1D2": entries with calculated load_kg
          └── ProgramSession "W1D3": entries with calculated load_kg
```

### Program Status Transitions

Program status transitions are explicit user actions via the program detail page. The "Complete Program" button is only enabled when all sessions in the program have associated logs. There is no automatic status transition.

## Design Decisions

- **ProgramSession is a DB entity** — Unlike the old metadata-based session grouping, sessions are now proper first-class entities with IDs, ordering, and optional dates
- **Program.status tracks execution progress** — The status transitions automatically; no manual status management required
- **Log.program_id is optional** — Supports both program-linked and ad-hoc training sessions
- **Absolute weights in Program, relative in Template** — Templates use RPE/%1RM, Programs use load_kg; the generate endpoint handles the conversion

For schema details, see `api/internal/domain/`.
