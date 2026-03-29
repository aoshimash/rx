# Plan Frontend Design Spec

## Overview

Replace the current Dashboard-centric frontend with a Plan-centric design. The Plan page becomes the top page (`/`), serving as an execution queue of upcoming training sessions. Programs are accessible from a sidebar panel on the Plan page and retain their own detail/edit pages.

## Navigation

Sidebar reduces from 4 items to 3:

| Route | Label | Icon |
|-------|-------|------|
| `/` | Plan | CalendarCheck (or similar) |
| `/logs` | Logs | ClipboardList |
| `/settings` | Settings | Settings |

The Dashboard page (`/`) and its components (`ActivePlans`, `RecentLogs`) are removed. The `/programs` list page is removed; Program detail and edit pages remain at `/programs/:id` and `/programs/:id/edit`.

## Plan Page (`/`)

### Layout: Split 70/30

- **Left (70%)**: Session queue — the user's upcoming sessions, ordered top to bottom
- **Right (30%)**: Program list — compact list of all programs with a "+ New" button

### Session Cards

Each session is rendered as a card showing all exercises in a **compact fixed-width table**:

- Table does NOT stretch to full container width (`width: auto`)
- Exercise name column has a fixed width (e.g., `150px`) so numeric field columns sit immediately to its right
- Field columns are generated dynamically from the session's entries' `fields` keys
- Column headers are the field names as defined by the user

Card header contains:
- **Session name** (bold)
- **Source label**: "from {Program name}" if `source_program_id` is set, or "manual" in italic if not
- **Date badge** (optional): rendered if `PlanSession.date` is set
- **Delete button** (×): removes session from plan with confirmation

### Session Card Interaction

- **Click anywhere on the card** → navigates to Log creation (`/logs/new`) with the session data pre-filled as `plan_snapshot`
- After log is saved, the session is automatically deleted from the Plan queue (via `DELETE /plan/sessions/:session_id`)

### Add Session

A "+ Add Session" button below the session list opens a form/dialog to manually create a session:
- Session name (required)
- Date (optional)
- Exercise entries with free-form fields (no FieldDef schema — user types field names and values directly)

Note: FieldDef-driven column headers only apply to sessions expanded from a Program (which carries `program_fields`). Manual sessions display whatever field keys exist in each entry's `fields` map.

### Program Sidebar (Right 30%)

Compact list of programs:
- Each item shows: program name, session count
- Click → navigates to `/programs/:id`
- "+ New" button → navigates to `/programs/new`

### Empty State

When Plan has no sessions:
- Centered message: "No sessions planned yet"
- Primary CTA: "+ Add Session" button
- Program sidebar still visible (may also be empty with its own guidance text)

## Program Detail Page (`/programs/:id`)

### Existing Features Retained

Program detail page keeps its current functionality (view sessions, exercises, edit, delete).

### New: "Add to Plan" Flow

1. "Add to Plan" button in the page header
2. Opens a session selector: list of all program sessions with checkboxes
3. **All sessions checked by default**
4. User unchecks sessions they don't want
5. "Add to Plan" confirms → calls `POST /plan/sessions` with selected sessions (including `source_program_id` and `source_session_id`)
6. Navigate back to Plan page

This replaces the `POST /plan/expand-program/:program_id` endpoint with a more granular `POST /plan/sessions` call that only includes selected sessions. The expand-program endpoint remains available as a convenience but is not used by this UI.

## Log Creation Flow Changes

### Current Flow

1. Navigate to `/logs/new?programId=X&session=Y`
2. Fetch Program, find matching session, use as template
3. Save Log with `program_id` and `session_name`

### New Flow (from Plan)

1. Click session card on Plan page
2. Navigate to `/logs/new` with PlanSession data passed (via query params or state)
3. `plan_snapshot` is populated from the PlanSession content (exercise names, fields, notes)
4. User fills in actual values (sets with real weights, reps, RPE, etc.)
5. Save Log with:
   - `plan_snapshot`: the session prescription as JSON
   - `program_id`: from `PlanSession.source_program_id` (if present)
   - `session_name`: from `PlanSession.session_name`
   - `entries`: actual performed exercises with sets
6. On successful save, call `DELETE /plan/sessions/:session_id` to remove from queue

### Ad-hoc Logs

Users can still create logs without a Plan session via `/logs/new` directly (from Logs page). These logs will have no `plan_snapshot`.

## Logs Page (`/logs`) Redesign

### Card-Based Layout

Replace the current table/list with cards matching the Plan page aesthetic:

- Each log is a card with:
  - **Header**: session name + source label ("from {Program}" or "ad-hoc") + date
  - **Body**: compact fixed-width table of exercises

### Plan vs Actual Diff

When a log has `plan_snapshot`:
- For each field value, compare plan value vs actual value
- If different: show plan value with strikethrough + actual value in color
  - Green if actual > plan (e.g., lifted more weight)
  - Red if actual < plan
- If same: show value normally

When a log has no `plan_snapshot`: show actual values only (no diff).

### Card Click

Click on a log card → navigates to `/logs/:id` detail page (existing behavior).

## Frontend Infrastructure

### New Files Required

| File | Purpose |
|------|---------|
| `web/types/api.ts` | Add Plan, PlanSession, PlanSessionEntry, PlanCreate, PlanUpdate types |
| `web/lib/api/plans.ts` | Plan API client (getPlan, createPlan, updatePlan, deletePlan, addSessions, updateSession, deleteSession) |
| `web/lib/hooks/usePlans.ts` | TanStack Query hooks for Plan operations |
| `web/app/page.tsx` | Rewrite as Plan page (replace Dashboard) |
| `web/components/plan/SessionCard.tsx` | Session card with dynamic field table |
| `web/components/plan/ProgramSidebar.tsx` | Right sidebar program list |
| `web/components/plan/AddSessionDialog.tsx` | Manual session creation |
| `web/components/plan/EmptyState.tsx` | Empty plan state |
| `web/components/programs/SessionSelector.tsx` | Checkbox session list for "Add to Plan" |
| `web/components/logs/LogCard.tsx` | Card component for Logs page |
| `web/components/logs/PlanDiff.tsx` | Plan vs actual diff display |

### Files to Modify

| File | Change |
|------|--------|
| `web/components/layout/Sidebar.tsx` | Update nav items: Plan / Logs / Settings |
| `web/app/programs/[id]/page.tsx` | Add "Add to Plan" button and session selector |
| `web/app/logs/page.tsx` | Redesign with card layout |
| `web/app/logs/new/page.tsx` | Support plan_snapshot from Plan session |

### Files to Remove

| File | Reason |
|------|--------|
| `web/components/dashboard/ActivePlans.tsx` | Dashboard removed |
| `web/components/dashboard/RecentLogs.tsx` | Dashboard removed |
| `web/app/programs/page.tsx` | Programs list page removed (accessed via Plan sidebar) |

## API Endpoints Used

| Endpoint | Method | Used By |
|----------|--------|---------|
| `GET /plan` | GET | Plan page — fetch session queue |
| `POST /plan` | POST | Plan page — create plan (first session add) |
| `POST /plan/sessions` | POST | "Add to Plan" from Program detail; manual session add |
| `PUT /plan/sessions/:id` | PUT | Edit session (future) |
| `DELETE /plan/sessions/:id` | DELETE | Remove session from queue; auto-remove after log |
| `GET /programs` | GET | Program sidebar |
| `GET /programs/:id` | GET | Program detail page |
| `POST /logs` | POST | Log creation with plan_snapshot |
| `GET /logs` | GET | Logs page |

## Out of Scope

- Drag-and-drop session reordering (can be added later via order field)
- Session editing within Plan (edit session name, entries, etc.)
- Program creation/editing (existing pages unchanged)
- Mobile-specific responsive layout (standard responsive behavior via CSS)
