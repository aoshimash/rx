# Domain Model

## Two-Tier Training Lifecycle

```
Program → Log
```

1. **Program** — Concrete training program containing sessions with scheduled exercises and absolute weights. Tracks training progress with a `status` (created / ongoing / completed / cancelled). Can be created manually or imported from JSON.
2. **Log** — Actual workout record for one training session (what was actually performed).

## Key Concepts

### Program

A **Program** is a concrete training plan. It:

- Contains **ProgramSessions** (named workout days with an order and optional date)
- Each session has **ProgramSessionEntries** (exercises with absolute weights)
- Has a `status`: `created` (registered, not yet started), `ongoing` (in progress), `completed` (all sessions logged and confirmed), or `cancelled` (stopped mid-way)
- Program status transitions are explicit user actions: `created` → `ongoing` (Start), `ongoing` → `completed` (Complete, requires all sessions logged), `ongoing` → `cancelled` (Cancel). Both `completed` and `cancelled` are terminal states.

### Log and Program (many:1)

A Log records what was actually performed in one training session. A Log **optionally** references a Program (`program_id`) and a session name (`session_name`), enabling tracking of which program session was performed. Logs without a `program_id` represent unplanned/ad-hoc training sessions.

## Entity Relationships

```
┌───────────────────┐     ┌──────────────────────┐
│     Program       │────►│   ProgramSession      │  (named training day)
│  (status: created/│     └──────────┬───────────┘
│   ongoing/        │               │
│   completed/      │     ┌──────────▼───────────┐
│   cancelled)      │     │ ProgramSessionEntry   │  (exercise with load_kg)
└────────┬──────────┘     └──────────────────────┘
         │ 1:N (optional reference)
         ▼
┌───────────────────┐
│       Log         │  (record: what was actually performed)
│  program_id       │
│  session_name     │
└───────────────────┘
```

### Program Status Transitions

Program status transitions are explicit user actions via the program detail page. The "Complete Program" button is only enabled when all sessions in the program have associated logs. There is no automatic status transition.

## Design Decisions

- **ProgramSession is a DB entity** — Sessions are first-class entities with IDs, ordering, and optional dates
- **Program.status tracks execution progress** — Status transitions are explicit user actions
- **Log.program_id is optional** — Supports both program-linked and ad-hoc training sessions
- **Flexible metadata** — Domain-specific data like RPE, tempo, rest time is stored in the `metadata` JSON field rather than dedicated columns, following the "Dumb Backend" philosophy

For schema details, see `api/internal/domain/`.
