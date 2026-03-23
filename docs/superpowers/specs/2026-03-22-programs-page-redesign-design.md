# Programs Page Redesign — Design Spec

## Summary

Redesign the Programs page to visually distinguish programs by status, with ongoing programs prominently displayed and a new status model replacing the current `active / completed / planned` set.

## Status Model

### Statuses

| Status | Meaning | Default |
|---|---|---|
| `created` | Registered, not yet started | Yes (on creation) |
| `ongoing` | In progress | |
| `completed` | All sessions logged and user confirmed | |
| `cancelled` | Stopped mid-way by user | |

Replaces the previous `active / completed / planned` enum.

### Transitions

All transitions are explicit user actions via buttons on the program detail page.

```
created → ongoing      (Start button)
ongoing → completed    (Complete button — enabled only when all sessions have logs)
ongoing → cancelled    (Cancel button — available anytime)
```

No automatic transitions. No backend business logic for transitions (Dumb Backend principle).

### Transition Rules

- **Start**: Available on detail page for `created` programs. Moves to `ongoing`.
- **Complete**: Available on detail page for `ongoing` programs. Button is **disabled** unless every session in the program has at least one associated log. Moves to `completed`.
- **Cancel**: Available on detail page for `ongoing` programs. Moves to `cancelled`. No preconditions.

Both `completed` and `cancelled` are **terminal states**. Cancelled programs cannot be restarted; create a new program instead.

### Transition Validation

The update status handler validates that the requested transition is allowed. Invalid transitions (e.g., `cancelled → ongoing`, `created → completed`) return a 400 error. This is input validation, not business logic, consistent with the Dumb Backend principle.

Allowed transitions:
- `created → ongoing`
- `ongoing → completed`
- `ongoing → cancelled`

### Removal of Auto-Completion Logic

The existing auto-completion logic in `api/internal/handler/log.go` (which automatically transitions programs to `completed` when all sessions are logged) must be **removed**. The `completed` transition is now an explicit user action only.

## Programs List Page Layout

The page is divided into three visual sections, ordered top to bottom:

### 1. Ongoing Section

- **Layout**: Single column, full width
- **Visibility**: Always shown (if any ongoing programs exist)
- **Card style**: Prominent — border highlight, progress bar showing session completion (e.g., "3 / 12 sessions")
- **Section header**: "Ongoing" (small uppercase label)

### 2. Created Section

- **Layout**: 2-column grid
- **Visibility**: Always shown (if any created programs exist)
- **Card style**: Subdued — standard border, no progress bar. Shows program name and session count only.
- **Section header**: "Programs" (small uppercase label)

### 3. Finished Section (completed + cancelled)

- **Layout**: 2-column grid
- **Visibility**: Hidden by default, toggled via "Show finished (N)" / "Hide finished (N)" button
- **Toggle state**: Persisted in localStorage (renamed to `programs:showFinished`). On first load, migrate the old `programs:showCompleted` key to the new key and remove the old one.
- **Card style**: Subdued, reduced opacity. Badge distinguishes `completed` vs `cancelled`.
- **No progress bar**

### Empty State

When no programs exist at all, show the existing empty state with "Create Program" button.

When a section has no programs, that section is simply not rendered (no empty placeholder per section).

## Program Card Variants

### Ongoing Card

- Full-width, single column
- Status badge: `ongoing` (primary/blue variant)
- Progress bar: yes (completed sessions / total sessions)
- Session count text: "3 / 12 sessions"
- Template info: shown if present
- Notes: shown if present

### Created Card

- 2-column grid item
- No status badge (the section heading provides context)
- No progress bar
- Shows: program name, session count (e.g., "8 sessions"), created date

### Finished Card (completed / cancelled)

- 2-column grid item
- Status badge: `completed` (secondary/green) or `cancelled` (destructive/red)
- No progress bar
- Reduced opacity
- Shows: program name, session count, created date

## Program Detail Page Changes

Add status action buttons:

- **created** status: Show "Start Program" button
- **ongoing** status: Show "Complete Program" button (disabled unless all sessions logged) and "Cancel Program" button
- **completed / cancelled** status: No action buttons (terminal states)

## API Changes

### OpenAPI Spec (`api/openapi/openapi.yaml`)

Change the `status` enum from `[active, completed, planned]` to `[created, ongoing, completed, cancelled]`.

Update the status description:
- `created` = registered, not yet started
- `ongoing` = in progress
- `completed` = all sessions logged and confirmed by user
- `cancelled` = stopped mid-way

Update the list endpoint's status query parameter enum accordingly.

Change default status on program creation from `active` to `created`.

### Domain Model (`api/internal/domain/program.go`)

Replace status constants:

```go
const (
    ProgramStatusCreated   ProgramStatus = "created"
    ProgramStatusOngoing   ProgramStatus = "ongoing"
    ProgramStatusCompleted ProgramStatus = "completed"
    ProgramStatusCancelled ProgramStatus = "cancelled"
)
```

### Validation (`api/internal/domain/validation.go`)

Update the status validation to accept the new enum values.

### Handler (`api/internal/handler/program.go`)

Change default status on creation from `active` to `created`.

### Handler (`api/internal/handler/log.go`)

Remove the auto-completion logic that transitions programs to `completed` when all sessions are logged.

### Template Generation (`api/internal/domain/program_template.go`)

Update `GenerateProgram` to set initial status to `ProgramStatusCreated` instead of `ProgramStatusActive`.

### Update Status Handler

Add transition validation: only allow `created→ongoing`, `ongoing→completed`, `ongoing→cancelled`. Return 400 for invalid transitions.

### Generated Code

Run `task generate` after OpenAPI spec changes to regenerate `pkg/openapi/server.gen.go`.

### Documentation (`docs/DOMAIN_MODEL.md`)

Update program status descriptions and remove references to automatic completion.

## Web Changes

### Type Definition (`web/types/api.ts`)

```typescript
export type ProgramStatus = 'created' | 'ongoing' | 'completed' | 'cancelled';
```

### Programs List Page (`web/app/programs/page.tsx`)

- Split programs into three groups: `ongoing`, `created`, `finished` (completed + cancelled)
- Render each group in its own section with appropriate layout
- Rename localStorage key from `programs:showCompleted` to `programs:showFinished`
- Update finished count to include both `completed` and `cancelled`

### ProgramCard (`web/components/programs/ProgramCard.tsx`)

- Accept a `variant` prop or derive display from status
- Ongoing: full-width, progress bar, prominent styling
- Created: compact, no progress bar, no status badge
- Finished: compact, reduced opacity, status badge

### Program Detail Page (`web/app/programs/[id]/page.tsx`)

- Add Start / Complete / Cancel buttons based on current status
- Complete button: disabled unless all sessions have logs (use existing `useLoggedSessions` hook)
- Status change calls the existing program update API endpoint

## Performance Note

Ongoing cards display progress bars using the `useLoggedSessions` hook, which makes a per-program API call. With multiple ongoing programs, this creates N+1 requests. This is acceptable for now since the number of concurrent ongoing programs is expected to be small (typically 1-3). If this becomes a bottleneck, a batch endpoint can be added later.

## Data Migration

Implement as a database migration file:

- `active` → `ongoing`
- `planned` → `created`
- `completed` → `completed` (no change)

Include a rollback migration:

- `ongoing` → `active`
- `created` → `planned`
- `cancelled` → `active` (lossy but acceptable for rollback)
