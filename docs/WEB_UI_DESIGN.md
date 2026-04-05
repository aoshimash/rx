# Web UI Design

This document describes the UI design for the Rx web application.

**Created**: 2026-01-30
**Updated**: 2026-04-05
**Status**: Living document

## Overview

The web UI is designed for powerlifters who plan their training programs and track actual workout results. The primary workflow is: create a Program (template), expand it into a Plan (execution queue), then record Logs against each session.

### Value Proposition

Rx's core value is making training data **easy to input** and **easy to extract**. The UI focuses exclusively on frictionless data entry. Analysis, visualization, and recommendations are out of scope — those are the job of external tools and AI agents.

### Platform Split

| Platform | Primary use | Environment |
|----------|-----------|-------------|
| **PC (Web)** | Program creation, Plan editing, bulk log entry | Home |
| **Mobile** | Real-time log recording during workouts | Gym |

The PC and mobile experiences are designed separately. This document covers the PC (Web) UI.

### Design Principles

1. **Plan-First Approach** — Users create a training program first, then execute sessions from their plan
2. **Session-Based View** — The plan page shows an ordered list of upcoming sessions
3. **Minimal Navigation** — Left sidebar for primary navigation; modals for quick actions
4. **In-Place Editing for Plans** — Plan fields are directly editable by clicking. No mode switching
5. **Explicit Editing for Logs** — Logs are confirmed records. Editing requires an explicit action (Edit button)
6. **Spreadsheet-Like Log Entry** — PC log entry uses a high-density tabular layout for bulk input. No per-set time tracking (that's mobile's job)

## User Journeys

### A. Program Creation

#### A1: Create a Program from Scratch

1. Navigate to `/programs/new` from sidebar or Plan page
2. Enter program name, notes, session structure, exercises, and field values
3. Save as **Draft** (editable, cannot be added to Plan) or **Published** (ready for use)

Programs support multiple load specification formats stored as-is (absolute weight, %1RM, RPE targets, etc.). Rx does not convert between formats — that is the user's or AI agent's responsibility. Programs should indicate the specification format and include conversion formulas when applicable, so external tools (including AI) can process them.

#### A2: Import a Program from External Source

There is no import UI. External data (spreadsheets, PDFs, text) is converted to Rx's API format by AI agents or scripts, then submitted via API. Rx's responsibility is a well-documented API that external tools can target.

#### A3: Duplicate and Modify a Program

1. Open an existing program at `/programs/{id}`
2. Click **Duplicate** → Creates a copy in **Draft** state
3. Edit the copy (rename, adjust weights, swap exercises)
4. Publish when ready

### B. Plan Creation and Editing

#### B1: Expand a Program into a Plan

1. Open a program at `/programs/{id}`
2. Click **Add to Plan** → Confirmation dialog (no session selection)
3. All sessions are added to the Plan in program order
4. Adjust on the Plan page if needed

Session selection modal is intentionally omitted. Users add the full program and remove unwanted sessions from the Plan afterward. This keeps the flow simple and avoids a complex selection UI for programs with many sessions. Duplicate additions are allowed (same program can be expanded multiple times — common for repeating weekly cycles).

#### B2: Manual Plan Edits (Ad-Hoc)

For exceptional cases (injury, schedule change), users can:
- **Add a session** manually from the Plan page
- **Delete** a session from the Plan
- **Reorder** sessions via drag and drop

This is not the primary flow. The target user normally works from a Program.

#### B3: Edit Plan Session Content

Plan session cards display all content inline. Individual fields (exercise name, load, reps, RPE, etc.) are **directly clickable for in-place editing**. There is no "edit mode" toggle — the card is always both a display and an editing surface.

### C. Log Recording

#### C1: Record a Log from a Plan Session (Mobile — Primary)

1. Open the app → Today's session (first in Plan) is shown
2. Plan content is pre-filled as template
3. Tap to mark each set complete → Plan values are recorded as-is
4. Only enter data when deviating from the plan
5. Save → Plan session is consumed

Mobile UI features per-set tap tracking with automatic timestamps, serving as both a set counter and interval timer. This is mobile-only — PC does not track per-set timing.

During recording, users can freely deviate from the plan:
- Add exercises not in the plan
- Replace planned exercises with alternatives
- Reorder exercises

One log = one gym visit (not one plan session).

#### C2: Record a Log without a Plan

For unplanned training (joint sessions, personal training, etc.):
- Entry point: **"Record Log" button on the Logs page** (`/logs`)
- Starts with a blank form
- Important for fatigue management even when off-program

#### C3: Record/Complete a Log from PC

For entering logs after the gym or filling in details:
- Entry point: **Record button on Plan session card** (plan-linked) on the Plan page
- Spreadsheet-like tabular layout for high-density bulk input
- No per-set time tracking

#### C4: Edit a Past Log

1. `/logs` → Click a log → `/logs/{id}`
2. Click **Edit** → `/logs/{id}/edit`
3. Modify and save

Logs are confirmed records. Editing requires an explicit step (navigating to an edit page) to prevent accidental modifications. This is intentionally different from Plan's in-place editing.

### D. Data Review

#### D1–D2: Analysis and Visualization

**Out of scope for initial release.** Rx's core value is data input and output. Analysis (trend charts, plan-vs-actual comparison dashboards, etc.) is the responsibility of external tools and AI agents that consume the API.

The `/logs` page provides a chronological list of past logs and CSV export for external analysis.

## Decisions and Rationale

### What We Do NOT Build (and Why)

| Decision | Rationale |
|----------|-----------|
| No import UI for programs | Non-structured data conversion is an AI agent's job, not Rx's. Rx provides the API target. |
| No session selection modal for Add to Plan | Users normally follow the full program. Removing unwanted sessions after is simpler than selecting from a long checklist. |
| No per-set time tracking on PC | Time tracking requires real-time tap interaction suited for mobile. PC is for bulk entry after the fact. |
| No analysis/visualization (initial) | Core value is data I/O. Analysis is external (AI agents, scripts). Follows "Planning and Analysis are External" principle. |
| No "Record Log" button on Plan page for unplanned sessions | Plan page is for plan-linked workflows. Unplanned log recording lives on the Logs page, keeping each page's purpose clear. |
| No load conversion logic | Conversion formulas vary by methodology and coach. Rx stores values as-is and lets users/AI handle conversion. |

### Key Design Choices (and Why)

| Decision | Rationale |
|----------|-----------|
| In-place field editing for Plan cards | Plan editing is the primary PC activity. Click-to-edit individual fields is faster than a mode toggle. Avoids the "what does card click mean?" ambiguity. |
| Explicit Edit step for Logs | Logs are confirmed records. Accidental edits should be prevented. Different data nature warrants different interaction model. |
| Draft/Published state for Programs | Program creation involves hours of trial and error. Draft state lets users save work-in-progress without it leaking into Plans. |
| Duplicate button for Programs | Common workflow: base next cycle on previous one. Copy → Draft → Edit → Publish. |
| Record button on Plan cards | Log recording from Plan is available on PC (for post-gym entry), but as an explicit button — not the default card click action. |
| One log = one gym visit | Users may deviate from the plan, add exercises, or reorder. The log captures what actually happened, not what was planned. |

## Screen Structure

| Screen | Type | Route | Description |
|--------|------|-------|-------------|
| Dashboard / Plan | Page | `/` | Active plan with upcoming sessions |
| Programs | Page | `/programs` | Program list |
| Program Detail | Page | `/programs/[id]` | View program sessions, add to plan, duplicate |
| Program Editor | Page | `/programs/new`, `/programs/[id]/edit` | Create/edit programs (Draft/Published) |
| Logs | Page | `/logs` | Training log history + Record Log (unplanned) |
| Log Detail | Page | `/logs/[id]` | View individual log |
| Log Editor | Page | `/logs/[id]/edit` | Edit existing log |
| Log Creation | Page | `/logs/new` | Record a workout (from Plan or unplanned) |
| Settings | Page | `/settings` | App settings |

## Layout

The app uses a **left sidebar** layout (`components/layout/Sidebar.tsx`) for primary navigation between Dashboard, Programs, Logs, and Settings.

## Dashboard / Plan View

The main screen shows the user's active Plan — an ordered list of upcoming sessions. Each session card displays all content inline (session name, source program, scheduled date, exercises with full field details).

Actions on session cards:
- Click any field → In-place edit that field
- **Record** button → Navigate to log creation (pre-filled from plan)
- **Delete** button → Confirmation dialog, then remove from plan
- Drag → Reorder sessions

Page-level actions:
- **Add Session** → Create a manual session (for ad-hoc adjustments)

## Program Detail

Displays program sessions with exercises in table format.

Actions:
- **Add to Plan** → Confirmation dialog → All sessions added to Plan
- **Duplicate** → Creates a copy in Draft state, opens editor
- **Edit** → Navigate to program editor
- **Export** → Copy JSON / Download JSON
- **Delete** → Confirmation dialog

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

### State

| State | Description |
|-------|-------------|
| **Draft** | Work in progress. Can be edited. Cannot be added to Plan. |
| **Published** | Ready for use. Can be added to Plan, duplicated, exported. |

### Operations

| Operation | Description |
|-----------|-------------|
| Add Group | Add a group to organize sessions |
| Add Session | Add a session (within a group or ungrouped) |
| Add Exercise | Add an exercise row to a session |
| Delete | Confirmation dialog, then remove |
| Custom Fields | Define program-level and log-level custom fields |

## Logs Page

Displays a chronological list of past training logs.

Actions:
- Click a log → Navigate to log detail
- **Record Log** → Navigate to log creation (blank, no plan link)
- **Export CSV** → Download all logs as CSV

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
│   │   ├── page.tsx                # Log list + Record Log button
│   │   ├── new/page.tsx            # Log creation
│   │   ├── [id]/page.tsx           # Log detail
│   │   └── [id]/edit/page.tsx      # Log editor
│   ├── settings/page.tsx           # Settings
│   └── layout.tsx
├── components/
│   ├── ui/                         # shadcn/ui + shared components
│   ├── layout/
│   │   └── Sidebar.tsx             # Left sidebar navigation
│   ├── plan/
│   │   └── SessionCard.tsx         # Plan session card (in-place editing)
│   ├── programs/
│   │   └── ProgramForm.tsx         # Program editor form
│   ├── dashboard/
│   └── export/
│       └── ExportButton.tsx        # CSV export dropdown
└── lib/
    ├── api/                        # API client modules
    ├── hooks/                      # TanStack Query hooks
    └── utils/                      # Utilities (export)
```

## Related Documents

- [PHILOSOPHY.md](PHILOSOPHY.md) — Why Rx exists and core principles
- [FRONTEND_ARCHITECTURE.md](FRONTEND_ARCHITECTURE.md) — Technology stack and patterns
- [DOMAIN_MODEL.md](DOMAIN_MODEL.md) — Domain model (Program, Plan, Log)
