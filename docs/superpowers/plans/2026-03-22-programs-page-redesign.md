# Programs Page Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Redesign the Programs page with a new 4-status model (created/ongoing/completed/cancelled), sectioned layout, and status transition controls.

**Architecture:** API-first. Update OpenAPI spec → regenerate Go code → update domain/handler/migration → update frontend types → rebuild UI components and pages. TDD for all Go changes.

**Tech Stack:** Go (API), Next.js/React (Web), OpenAPI codegen, PostgreSQL migrations, TanStack Query

**Spec:** `docs/superpowers/specs/2026-03-22-programs-page-redesign-design.md`

---

### Task 1: Update OpenAPI Spec — Status Enum

**Files:**
- Modify: `api/openapi/openapi.yaml:199-204` (list endpoint status filter)
- Modify: `api/openapi/openapi.yaml:835-838` (Program schema status property)

- [ ] **Step 1: Update status enum in Program schema**

In `api/openapi/openapi.yaml`, find the Program schema status property (around line 835-838):

```yaml
# BEFORE
status:
  type: string
  enum: [active, completed, planned]
  description: active = in progress, completed = all sessions have been logged, planned = not yet started
```

Replace with:

```yaml
# AFTER
status:
  type: string
  enum: [created, ongoing, completed, cancelled]
  description: created = registered not yet started, ongoing = in progress, completed = all sessions logged and confirmed, cancelled = stopped mid-way
```

- [ ] **Step 2: Update status filter enum in list endpoint**

Find the list endpoint query parameter (around line 199-204):

```yaml
# BEFORE
- name: status
  in: query
  schema:
    type: string
    enum: [active, completed, planned]
  description: Filter programs by status
```

Replace with:

```yaml
# AFTER
- name: status
  in: query
  schema:
    type: string
    enum: [created, ongoing, completed, cancelled]
  description: Filter programs by status
```

- [ ] **Step 3: Add PATCH endpoint for status update**

Add a new endpoint after `/programs/{id}` delete operation (after line 264):

```yaml
  /programs/{id}/status:
    parameters:
      - $ref: '#/components/parameters/ProgramId'
    patch:
      summary: Update program status
      operationId: updateProgramStatus
      tags: [Program]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [status]
              properties:
                status:
                  type: string
                  enum: [created, ongoing, completed, cancelled]
      responses:
        '200':
          description: Status updated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Program'
        '400':
          $ref: '#/components/responses/ValidationError'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '404':
          $ref: '#/components/responses/NotFound'
```

- [ ] **Step 4: Run code generation**

```bash
cd api && task generate
```

Expected: `pkg/openapi/server.gen.go` is regenerated with updated type definitions. Note: this project uses chi router with manual route registration, not the generated `ServerInterface`. The generated code provides types only.

- [ ] **Step 5: Commit**

```bash
git add api/openapi/openapi.yaml api/pkg/openapi/server.gen.go
git commit -m "feat: update OpenAPI spec with new program status enum and status update endpoint"
```

---

### Task 2: Update Domain Model — Status Constants and Validation

**Files:**
- Modify: `api/internal/domain/program.go:10-17`
- Modify: `api/internal/domain/validation.go:297-302`
- Modify: `api/internal/domain/program_test.go:428-502`

- [ ] **Step 1: Write failing tests for new status validation**

In `api/internal/domain/program_test.go`, update the `TestValidateProgram` function. Change `validProgram()` helper to use `ProgramStatusCreated` and add test cases for new statuses:

```go
// In validProgram(), change:
// Status: ProgramStatusActive,
// to:
// Status: ProgramStatusCreated,

// Add these test cases inside TestValidateProgram:

t.Run("valid with ongoing status", func(t *testing.T) {
    p := validProgram()
    p.Status = ProgramStatusOngoing
    err := ValidateProgram(p)
    assert.NoError(t, err)
})

t.Run("valid with completed status", func(t *testing.T) {
    p := validProgram()
    p.Status = ProgramStatusCompleted
    err := ValidateProgram(p)
    assert.NoError(t, err)
})

t.Run("valid with cancelled status", func(t *testing.T) {
    p := validProgram()
    p.Status = ProgramStatusCancelled
    err := ValidateProgram(p)
    assert.NoError(t, err)
})

t.Run("invalid status", func(t *testing.T) {
    p := validProgram()
    p.Status = "invalid"
    err := ValidateProgram(p)
    assert.Error(t, err)
})

t.Run("old active status is invalid", func(t *testing.T) {
    p := validProgram()
    p.Status = "active"
    err := ValidateProgram(p)
    assert.Error(t, err)
})
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd api && go test ./internal/domain/... -run TestValidateProgram -v
```

Expected: FAIL — `ProgramStatusCreated` and other new constants are not defined.

- [ ] **Step 3: Update status constants in domain model**

In `api/internal/domain/program.go`, replace lines 10-17:

```go
// BEFORE
const (
	ProgramStatusActive    ProgramStatus = "active"
	ProgramStatusCompleted ProgramStatus = "completed"
	ProgramStatusPlanned   ProgramStatus = "planned"
)
```

```go
// AFTER
const (
	ProgramStatusCreated   ProgramStatus = "created"
	ProgramStatusOngoing   ProgramStatus = "ongoing"
	ProgramStatusCompleted ProgramStatus = "completed"
	ProgramStatusCancelled ProgramStatus = "cancelled"
)
```

- [ ] **Step 4: Update validation in validation.go**

In `api/internal/domain/validation.go`, replace lines 297-302:

```go
// BEFORE
if p.Status != ProgramStatusActive && p.Status != ProgramStatusCompleted && p.Status != ProgramStatusPlanned {
    return &ValidationError{
        Field:   "status",
        Message: "status must be 'active', 'completed', or 'planned'",
    }
}
```

```go
// AFTER
if p.Status != ProgramStatusCreated && p.Status != ProgramStatusOngoing && p.Status != ProgramStatusCompleted && p.Status != ProgramStatusCancelled {
    return &ValidationError{
        Field:   "status",
        Message: "status must be 'created', 'ongoing', 'completed', or 'cancelled'",
    }
}
```

- [ ] **Step 5: Update struct comment**

In `api/internal/domain/program.go`, line 45:

```go
// BEFORE
// Status transitions from "active" to "completed" when all sessions have been logged.
// AFTER
// Status transitions: created → ongoing → completed/cancelled. All transitions are explicit user actions.
```

- [ ] **Step 6: Add transition validation function**

Add to `api/internal/domain/program.go` after the status constants:

```go
// ValidateProgramStatusTransition checks if a status transition is allowed.
// Allowed: created→ongoing, ongoing→completed, ongoing→cancelled.
func ValidateProgramStatusTransition(from, to ProgramStatus) error {
	switch {
	case from == ProgramStatusCreated && to == ProgramStatusOngoing:
		return nil
	case from == ProgramStatusOngoing && to == ProgramStatusCompleted:
		return nil
	case from == ProgramStatusOngoing && to == ProgramStatusCancelled:
		return nil
	default:
		return &ValidationError{
			Field:   "status",
			Message: fmt.Sprintf("invalid status transition from '%s' to '%s'", from, to),
		}
	}
}
```

Add `"fmt"` to imports if not already present.

- [ ] **Step 7: Write tests for transition validation**

Add to `api/internal/domain/program_test.go`:

```go
func TestValidateProgramStatusTransition(t *testing.T) {
	t.Run("created to ongoing is valid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusCreated, ProgramStatusOngoing)
		assert.NoError(t, err)
	})

	t.Run("ongoing to completed is valid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusOngoing, ProgramStatusCompleted)
		assert.NoError(t, err)
	})

	t.Run("ongoing to cancelled is valid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusOngoing, ProgramStatusCancelled)
		assert.NoError(t, err)
	})

	t.Run("created to completed is invalid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusCreated, ProgramStatusCompleted)
		assert.Error(t, err)
	})

	t.Run("completed to ongoing is invalid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusCompleted, ProgramStatusOngoing)
		assert.Error(t, err)
	})

	t.Run("cancelled to ongoing is invalid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusCancelled, ProgramStatusOngoing)
		assert.Error(t, err)
	})

	t.Run("cancelled to created is invalid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusCancelled, ProgramStatusCreated)
		assert.Error(t, err)
	})
}
```

- [ ] **Step 8: Run tests to verify all pass**

```bash
cd api && go test ./internal/domain/... -v
```

Expected: ALL PASS

- [ ] **Step 9: Commit**

```bash
git add api/internal/domain/program.go api/internal/domain/validation.go api/internal/domain/program_test.go
git commit -m "feat: update program status enum and add transition validation"
```

---

### Task 3: Update Handlers — Default Status, Auto-Completion Removal, Status Endpoint

**Files:**
- Modify: `api/internal/handler/program.go:68` (creation default)
- Modify: `api/internal/handler/log.go:184-215` (remove auto-completion)
- Modify: `api/internal/domain/program_template.go:194` (template generation default)
- Modify: `api/internal/handler/program_test.go`
- Modify: `api/cmd/server/main.go:145` (register new route)

- [ ] **Step 1: Update existing tests for new default status**

In `api/internal/handler/program_test.go`, update all `assert.Equal(t, "active", ...)` to `assert.Equal(t, "created", ...)`.

Line 86:
```go
// BEFORE
assert.Equal(t, "active", result["status"])
// AFTER
assert.Equal(t, "created", result["status"])
```

Line 188 (status filter test):
```go
// BEFORE
req := httptest.NewRequest(http.MethodGet, "/programs?status=active", nil)
// AFTER
req := httptest.NewRequest(http.MethodGet, "/programs?status=created", nil)
```

- [ ] **Step 2: Write test for status update endpoint**

Add to `api/internal/handler/program_test.go`:

```go
func TestProgramHandler_UpdateProgramStatus(t *testing.T) {
	router, _ := setupProgramTestRouter()

	t.Run("transitions created to ongoing", func(t *testing.T) {
		created := createTestProgram(t, router)
		id := created["id"].(string)

		body := `{"status": "ongoing"}`
		req := httptest.NewRequest(http.MethodPatch, "/programs/"+id+"/status", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.Equal(t, "ongoing", result["status"])
	})

	t.Run("rejects invalid transition created to completed", func(t *testing.T) {
		created := createTestProgram(t, router)
		id := created["id"].(string)

		body := `{"status": "completed"}`
		req := httptest.NewRequest(http.MethodPatch, "/programs/"+id+"/status", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 404 for non-existent program", func(t *testing.T) {
		body := `{"status": "ongoing"}`
		req := httptest.NewRequest(http.MethodPatch, "/programs/00000000-0000-0000-0000-000000000001/status", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
```

- [ ] **Step 3: Run tests to verify they fail**

```bash
cd api && go test ./internal/handler/... -run TestProgramHandler -v
```

Expected: FAIL — default status is still "active", and update endpoint doesn't exist yet.

- [ ] **Step 4: Change default status on creation**

In `api/internal/handler/program.go`, line 68:

```go
// BEFORE
Status:   domain.ProgramStatusActive,
// AFTER
Status:   domain.ProgramStatusCreated,
```

- [ ] **Step 5: Change default status in template generation**

In `api/internal/domain/program_template.go`, line 194:

```go
// BEFORE
Status:            ProgramStatusActive,
// AFTER
Status:            ProgramStatusCreated,
```

- [ ] **Step 6: Remove auto-completion logic**

In `api/internal/handler/log.go`, delete lines 184-215 (the entire auto-completion block):

```go
// DELETE THIS ENTIRE BLOCK:
// Auto-completion: check if all sessions of the referenced program have been logged
if log.ProgramID != nil && h.programRepo != nil {
    // ... (all the way to the closing brace)
}
```

- [ ] **Step 7: Implement UpdateProgramStatus handler**

Add to `api/internal/handler/program.go` after the `ListLoggedSessions` function:

```go
// UpdateProgramStatus handles PATCH /programs/{id}/status
func (h *ProgramHandler) UpdateProgramStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", nil)
		return
	}

	newStatus := domain.ProgramStatus(req.Status)

	program, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	if err := domain.ValidateProgramStatusTransition(program.Status, newStatus); err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if err := h.repo.UpdateStatus(ctx, id, newStatus); err != nil {
		middleware.WriteInternalError(w, "Failed to update program status")
		return
	}

	program.Status = newStatus
	writeJSON(w, http.StatusOK, program)
}
```

- [ ] **Step 8: Register the new route**

In `api/cmd/server/main.go`, add the PATCH route after line 145 (after the existing program routes):

```go
r.Patch("/programs/{id}/status", programHandler.UpdateProgramStatus)
```

The program routes block should now look like:

```go
// Program routes
r.Post("/programs", programHandler.CreateProgram)
r.Get("/programs", programHandler.ListPrograms)
r.Get("/programs/{id}", programHandler.GetProgram)
r.Delete("/programs/{id}", programHandler.DeleteProgram)
r.Get("/programs/{id}/logged-sessions", programHandler.ListLoggedSessions)
r.Patch("/programs/{id}/status", programHandler.UpdateProgramStatus)
```

- [ ] **Step 9: Run tests to verify all pass**

```bash
cd api && go test ./internal/handler/... -v
```

Expected: ALL PASS

- [ ] **Step 10: Commit**

```bash
git add api/internal/handler/program.go api/internal/handler/program_test.go api/internal/handler/log.go api/internal/domain/program_template.go api/cmd/server/main.go
git commit -m "feat: add status update endpoint, change default to created, remove auto-completion"
```

---

### Task 4: Database Migration

**Files:**
- Create: `api/migrations/000013_program_status_redesign.up.sql`
- Create: `api/migrations/000013_program_status_redesign.down.sql`

- [ ] **Step 1: Create up migration**

```sql
-- Migrate existing status values
UPDATE programs SET status = 'ongoing' WHERE status = 'active';
UPDATE programs SET status = 'created' WHERE status = 'planned';

-- Update constraint to new enum values
ALTER TABLE programs DROP CONSTRAINT IF EXISTS programs_status_check;
ALTER TABLE programs ADD CONSTRAINT programs_status_check CHECK (status IN ('created', 'ongoing', 'completed', 'cancelled'));

-- Update default
ALTER TABLE programs ALTER COLUMN status SET DEFAULT 'created';
```

- [ ] **Step 2: Create down migration**

```sql
-- Revert status values
UPDATE programs SET status = 'active' WHERE status IN ('ongoing', 'cancelled');
UPDATE programs SET status = 'planned' WHERE status = 'created';

-- Revert constraint
ALTER TABLE programs DROP CONSTRAINT IF EXISTS programs_status_check;
ALTER TABLE programs ADD CONSTRAINT programs_status_check CHECK (status IN ('active', 'completed', 'planned'));

-- Revert default
ALTER TABLE programs ALTER COLUMN status SET DEFAULT 'active';
```

- [ ] **Step 3: Commit**

```bash
git add api/migrations/000013_program_status_redesign.up.sql api/migrations/000013_program_status_redesign.down.sql
git commit -m "feat: add migration for program status redesign"
```

---

### Task 5: Run API Checks

- [ ] **Step 1: Verify build compiles**

```bash
cd api && go build ./...
```

Expected: No errors. If there are compile errors from old status references (`ProgramStatusActive`, `ProgramStatusPlanned`), fix them — these are likely in handler tests or generated code that references old constants.

- [ ] **Step 2: Run format and lint**

```bash
cd api && task format && task lint
```

Expected: PASS

- [ ] **Step 3: Run all tests**

```bash
cd api && task test
```

Expected: ALL PASS

- [ ] **Step 4: Run full check**

```bash
cd api && task check
```

Expected: ALL PASS

- [ ] **Step 5: Commit any fixes**

If any fixes were needed, commit them:

```bash
git add -A api/
git commit -m "fix: resolve remaining old status references in API"
```

---

### Task 6: Update Web Types and API Client

**Files:**
- Modify: `web/types/api.ts:137`
- Modify: `web/lib/api/programs.ts`
- Modify: `web/lib/hooks/usePrograms.ts`

- [ ] **Step 1: Update ProgramStatus type**

In `web/types/api.ts`, line 137:

```typescript
// BEFORE
export type ProgramStatus = 'active' | 'completed' | 'planned';
// AFTER
export type ProgramStatus = 'created' | 'ongoing' | 'completed' | 'cancelled';
```

- [ ] **Step 2: Add updateStatus API function**

In `web/lib/api/programs.ts`, add to the `programsApi` object:

```typescript
async updateStatus(id: string, status: string): Promise<Program> {
  return api.patch(`programs/${id}/status`, { json: { status } }).json<Program>();
},
```

- [ ] **Step 3: Add useUpdateProgramStatus hook**

In `web/lib/hooks/usePrograms.ts`, add:

```typescript
export function useUpdateProgramStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      programsApi.updateStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
    },
  });
}
```

- [ ] **Step 4: Commit**

```bash
git add web/types/api.ts web/lib/api/programs.ts web/lib/hooks/usePrograms.ts
git commit -m "feat: update web types and add updateStatus API client"
```

---

### Task 7: Redesign ProgramCard Component

**Files:**
- Modify: `web/components/programs/ProgramCard.tsx`

- [ ] **Step 1: Rewrite ProgramCard with status-based variants**

Replace the entire `ProgramCard.tsx` with:

```tsx
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useProgramTemplate } from '@/lib/hooks/useProgramTemplates';
import { useLoggedSessions } from '@/lib/hooks/usePrograms';
import type { Program } from '@/types/api';
import { useRouter } from 'next/navigation';

interface ProgramCardProps {
  program: Program;
}

function TemplateInfo({ program }: { program: Program }) {
  const { data: template } = useProgramTemplate(program.program_template_id ?? null);
  const targetWeights = program.metadata?.target_weights as Record<string, number> | undefined;

  if (!template && !targetWeights) return null;

  return (
    <div className="text-sm text-muted-foreground space-y-0.5">
      {template && <p>Template: {template.name}</p>}
      {targetWeights && (
        <p>
          Targets:{' '}
          {Object.entries(targetWeights)
            .map(([exercise, weight]) => `${exercise}: ${weight}kg`)
            .join(', ')}
        </p>
      )}
    </div>
  );
}

function statusBadgeVariant(
  status: string,
): 'default' | 'secondary' | 'outline' | 'destructive' {
  if (status === 'ongoing') return 'default';
  if (status === 'completed') return 'secondary';
  if (status === 'cancelled') return 'destructive';
  return 'outline';
}

export function OngoingProgramCard({ program }: ProgramCardProps) {
  const router = useRouter();
  const { data: loggedSessions } = useLoggedSessions(program.id);

  const totalSessions = program.sessions.length;
  const loggedSet = new Set(loggedSessions?.sessions ?? []);
  const completedCount =
    totalSessions > 0 ? program.sessions.filter((s) => loggedSet.has(s.session_name)).length : 0;
  const progressPct = totalSessions > 0 ? (completedCount / totalSessions) * 100 : 0;

  return (
    <Card
      className="cursor-pointer transition-colors hover:bg-accent/50"
      onClick={() => router.push(`/programs/${program.id}`)}
    >
      <CardHeader>
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <CardTitle className="text-xl">{program.name}</CardTitle>
            <Badge variant="default">ongoing</Badge>
          </div>
          {program.notes && <p className="text-sm text-muted-foreground">{program.notes}</p>}
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {program.program_template_id && <TemplateInfo program={program} />}
          <div className="space-y-1">
            <div className="flex justify-between text-xs text-muted-foreground">
              <span>
                {completedCount} / {totalSessions} sessions
              </span>
            </div>
            <div className="h-2 w-full rounded-full bg-secondary overflow-hidden">
              <div
                className="h-full rounded-full bg-primary transition-all"
                style={{ width: `${progressPct}%` }}
              />
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export function CreatedProgramCard({ program }: ProgramCardProps) {
  const router = useRouter();

  return (
    <Card
      className="cursor-pointer transition-colors hover:bg-accent/50"
      onClick={() => router.push(`/programs/${program.id}`)}
    >
      <CardHeader className="pb-2">
        <CardTitle className="text-base">{program.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex justify-between text-sm text-muted-foreground">
          <span>{program.sessions.length} sessions</span>
          <span>{new Date(program.created_at).toLocaleDateString()}</span>
        </div>
      </CardContent>
    </Card>
  );
}

export function FinishedProgramCard({ program }: ProgramCardProps) {
  const router = useRouter();

  return (
    <Card
      className="cursor-pointer transition-colors hover:bg-accent/50 opacity-50"
      onClick={() => router.push(`/programs/${program.id}`)}
    >
      <CardHeader className="pb-2">
        <div className="flex items-center gap-2">
          <CardTitle className="text-base">{program.name}</CardTitle>
          <Badge variant={statusBadgeVariant(program.status)}>{program.status}</Badge>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex justify-between text-sm text-muted-foreground">
          <span>{program.sessions.length} sessions</span>
          <span>{new Date(program.created_at).toLocaleDateString()}</span>
        </div>
      </CardContent>
    </Card>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add web/components/programs/ProgramCard.tsx
git commit -m "feat: split ProgramCard into ongoing/created/finished variants"
```

---

### Task 8: Redesign Programs List Page

**Files:**
- Modify: `web/app/programs/page.tsx`

- [ ] **Step 1: Rewrite programs list page with sections**

Replace the program rendering logic in `web/app/programs/page.tsx`. Key changes:

1. Rename localStorage key with migration logic
2. Split programs into 3 groups: ongoing, created, finished
3. Render each group in its own section

```tsx
const SHOW_FINISHED_KEY = 'programs:showFinished';
const OLD_SHOW_COMPLETED_KEY = 'programs:showCompleted';

// In the component, replace the showCompleted state with:
const [showFinished, setShowFinished] = useState(() => {
  if (typeof window === 'undefined') return false;
  // Migrate old key
  const oldValue = localStorage.getItem(OLD_SHOW_COMPLETED_KEY);
  if (oldValue !== null) {
    localStorage.setItem(SHOW_FINISHED_KEY, oldValue);
    localStorage.removeItem(OLD_SHOW_COMPLETED_KEY);
    return oldValue === 'true';
  }
  return localStorage.getItem(SHOW_FINISHED_KEY) === 'true';
});

// Split programs into groups:
const allPrograms = [...(programsData?.data || [])].sort(
  (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
);
const ongoingPrograms = allPrograms.filter((p) => p.status === 'ongoing');
const createdPrograms = allPrograms.filter((p) => p.status === 'created');
const finishedPrograms = allPrograms.filter(
  (p) => p.status === 'completed' || p.status === 'cancelled'
);

const toggleShowFinished = () => {
  const next = !showFinished;
  setShowFinished(next);
  localStorage.setItem(SHOW_FINISHED_KEY, String(next));
};
```

Render sections:

```tsx
{/* Ongoing section */}
{ongoingPrograms.length > 0 && (
  <div className="mb-8">
    <h2 className="text-xs uppercase tracking-wider text-muted-foreground mb-3">Ongoing</h2>
    <div className="space-y-4">
      {ongoingPrograms.map((program) => (
        <OngoingProgramCard key={program.id} program={program} />
      ))}
    </div>
  </div>
)}

{/* Created section */}
{createdPrograms.length > 0 && (
  <div className="mb-8">
    <h2 className="text-xs uppercase tracking-wider text-muted-foreground mb-3">Programs</h2>
    <div className="grid grid-cols-2 gap-4">
      {createdPrograms.map((program) => (
        <CreatedProgramCard key={program.id} program={program} />
      ))}
    </div>
  </div>
)}

{/* Finished section */}
{finishedPrograms.length > 0 && (
  <div className="mb-8">
    <Button variant="ghost" size="sm" onClick={toggleShowFinished}>
      {showFinished
        ? `Hide finished (${finishedPrograms.length})`
        : `Show finished (${finishedPrograms.length})`}
    </Button>
    {showFinished && (
      <div className="grid grid-cols-2 gap-4 mt-3">
        {finishedPrograms.map((program) => (
          <FinishedProgramCard key={program.id} program={program} />
        ))}
      </div>
    )}
  </div>
)}
```

Update imports to use new card components:

```tsx
import { OngoingProgramCard, CreatedProgramCard, FinishedProgramCard } from '@/components/programs/ProgramCard';
```

- [ ] **Step 2: Run biome check**

```bash
cd web && pnpm check
```

Expected: PASS (fix any issues with `pnpm check:fix`)

- [ ] **Step 3: Commit**

```bash
git add web/app/programs/page.tsx
git commit -m "feat: redesign programs list page with sectioned layout"
```

---

### Task 9: Add Status Action Buttons to Detail Page

**Files:**
- Modify: `web/app/programs/[id]/page.tsx`

- [ ] **Step 1: Add status action buttons**

In `web/app/programs/[id]/page.tsx`:

1. Import the new hook:

```tsx
import { useDeleteProgram, useLoggedSessions, useProgram, useUpdateProgramStatus } from '@/lib/hooks/usePrograms';
```

2. Add the mutation hook in the component:

```tsx
const updateStatus = useUpdateProgramStatus();
```

3. Compute whether all sessions are logged (for Complete button):

```tsx
const allSessionsLogged = program.sessions.length > 0 &&
  program.sessions.every((s) => completedSessionSet.has(s.session_name));
```

4. Update the badge variant (line 156):

```tsx
// BEFORE
<Badge variant={program.status === 'active' ? 'default' : 'secondary'}>
// AFTER
<Badge variant={program.status === 'ongoing' ? 'default' : 'secondary'}>
```

5. Add action buttons after the delete button, inside the header flex container:

```tsx
<div className="flex items-center gap-2">
  {program.status === 'created' && (
    <Button
      onClick={() => updateStatus.mutate({ id: programId, status: 'ongoing' })}
      disabled={updateStatus.isPending}
    >
      {updateStatus.isPending ? 'Starting...' : 'Start Program'}
    </Button>
  )}
  {program.status === 'ongoing' && (
    <>
      <Button
        onClick={() => updateStatus.mutate({ id: programId, status: 'completed' })}
        disabled={updateStatus.isPending || !allSessionsLogged}
      >
        {updateStatus.isPending ? 'Completing...' : 'Complete Program'}
      </Button>
      <Button
        variant="outline"
        onClick={() => updateStatus.mutate({ id: programId, status: 'cancelled' })}
        disabled={updateStatus.isPending}
      >
        Cancel Program
      </Button>
    </>
  )}
  <Button variant="outline" onClick={handleDelete} disabled={deleteProgram.isPending}>
    <Trash2 className="h-4 w-4 mr-2" />
    Delete
  </Button>
</div>
```

- [ ] **Step 2: Run biome check**

```bash
cd web && pnpm check
```

Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add web/app/programs/[id]/page.tsx
git commit -m "feat: add status action buttons to program detail page"
```

---

### Task 10: Update Documentation

**Files:**
- Modify: `docs/DOMAIN_MODEL.md`

- [ ] **Step 1: Update program status documentation**

Find and update the status description (around line 29-30):

```markdown
# BEFORE
- Has a `status`: `active` (in progress) or `completed` (finished)

# AFTER
- Has a `status`: `created` (registered, not yet started), `ongoing` (in progress), `completed` (all sessions logged and confirmed), or `cancelled` (stopped mid-way)
```

Find and update the auto-completion description (around line 74-76):

```markdown
# BEFORE
When a Log is created with a `program_id` and `session_name`, the system checks whether all sessions in the referenced Program have been logged. If all sessions have been logged, the Program status transitions automatically from `active` to `completed`.

# AFTER
Program status transitions are explicit user actions: `created` → `ongoing` (Start), `ongoing` → `completed` (Complete, requires all sessions logged), `ongoing` → `cancelled` (Cancel). Both `completed` and `cancelled` are terminal states.
```

- [ ] **Step 2: Commit**

```bash
git add docs/DOMAIN_MODEL.md
git commit -m "docs: update program status model documentation"
```

---

### Task 11: Final Verification

- [ ] **Step 1: Run full API check**

```bash
cd api && task check
```

Expected: ALL PASS

- [ ] **Step 2: Run full web check**

```bash
cd web && pnpm check
```

Expected: ALL PASS

- [ ] **Step 3: Build web**

```bash
cd web && pnpm build
```

Expected: Build succeeds with no type errors.

- [ ] **Step 4: Verify OpenAPI spec and generated code are in sync**

```bash
cd api && task generate && git diff --exit-code pkg/openapi/
```

Expected: No diff (generated code matches spec).
