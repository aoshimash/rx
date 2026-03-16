# Domain Model

## Three-Stage Training Lifecycle

```
Program → Plan → Log
```

1. **Program** — Reusable training template (no dates, no absolute weights)
2. **Plan** — Concrete training schedule derived from a Program (with dates and calculated weights)
3. **Log** — Actual workout record (what was actually performed)

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

For schema details, see `api/internal/domain/`.
