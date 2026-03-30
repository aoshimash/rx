# Web UI Design

This document describes the UI design for the Rx web application.

**Created**: 2026-01-30
**Updated**: 2026-03-30
**Status**: Living document

## Overview

The web UI is designed for powerlifters who plan their training programs and track actual workout results. The primary workflow is: create a Program (template), expand it into a Plan (execution queue), then record Logs against each session.

### Design Principles

1. **Plan-First Approach** — Users create a training program first, then execute sessions from their plan
2. **Session-Based View** — The plan page shows an ordered list of upcoming sessions
3. **Minimal Navigation** — Left sidebar for primary navigation; modals for quick actions
4. **Data Export** — CSV export for external analysis

## Screen Structure

| Screen | Type | Route | Description |
|--------|------|-------|-------------|
| Dashboard / Plan | Page | `/` | Active plan with upcoming sessions |
| Programs | Page | `/programs` | Program list |
| Program Detail | Page | `/programs/[id]` | View program sessions, add to plan |
| Program Editor | Page | `/programs/new`, `/programs/[id]/edit` | Create/edit programs |
| Logs | Page | `/logs` | Training log history |
| Log Detail | Page | `/logs/[id]` | View individual log |
| Log Creation | Page | `/logs/new` | Record a workout |
| Settings | Page | `/settings` | App settings |

## Layout

The app uses a **left sidebar** layout (`components/layout/Sidebar.tsx`) for primary navigation between Dashboard, Programs, Logs, and Settings.

## Dashboard / Plan View

The main screen shows the user's active Plan — an ordered list of upcoming sessions. Each session card shows the session name, source program, scheduled date (optional), and exercise summary.

Actions:
- Click a session card → navigate to log creation (pre-filled from plan)
- Delete a session from the plan
- Add sessions from a program

## Program Editor

A form-based UI for creating and editing training programs with hierarchical group support.

### Hierarchy

```
Program
├── Group (e.g., "Block A")
│   ├── Group (e.g., "Week 1", max 2 levels)
│   │   ├── Session: "Day 1" → Exercises
│   │   └── Session: "Day 2" → Exercises
│   └── Group (e.g., "Week 2")
│       └── ...
└── Ungrouped Session → Exercises
```

### Operations

| Operation | Description |
|-----------|-------------|
| Add Group | Add a group to organize sessions |
| Add Session | Add a session (within a group or ungrouped) |
| Add Exercise | Add an exercise row to a session |
| Delete | Confirmation dialog, then remove |
| Custom Fields | Define program-level and log-level custom fields |

## Export

- **Format**: CSV (UTF-8 with BOM for Excel compatibility)
- **Scope**: All logs
- **Columns**: Date, Time, Session, Exercise, Sets, Reps, Load (kg), Notes

## Component Structure

```
web/
├── app/
│   ├── page.tsx                    # Dashboard (Plan view)
│   ├── programs/
│   │   ├── page.tsx                # Program list
│   │   ├── new/page.tsx            # Program editor (create)
│   │   ├── [id]/page.tsx           # Program detail
│   │   └── [id]/edit/page.tsx      # Program editor (edit)
│   ├── logs/
│   │   ├── page.tsx                # Log list
│   │   ├── new/page.tsx            # Log creation
│   │   └── [id]/page.tsx           # Log detail
│   ├── settings/page.tsx           # Settings
│   └── layout.tsx
├── components/
│   ├── ui/                         # shadcn/ui + shared components
│   ├── layout/
│   │   └── Sidebar.tsx             # Left sidebar navigation
│   ├── plan/
│   │   └── SessionCard.tsx         # Plan session card
│   ├── programs/
│   │   ├── ProgramForm.tsx         # Program editor form
│   │   └── SessionSelector.tsx     # Select sessions to add to plan
│   ├── dashboard/
│   │   └── E1rmChart.tsx           # e1RM trend chart
│   └── export/
│       └── ExportButton.tsx        # CSV export dropdown
└── lib/
    ├── api/                        # API client modules
    ├── hooks/                      # TanStack Query hooks
    └── utils/                      # Utilities (export, e1rm calc)
```

## Related Documents

- [FRONTEND_ARCHITECTURE.md](FRONTEND_ARCHITECTURE.md) — Technology stack and patterns
- [DOMAIN_MODEL.md](DOMAIN_MODEL.md) — Domain model (Program, Plan, Log)
