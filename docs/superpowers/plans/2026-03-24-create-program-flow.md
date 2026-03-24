# Create Program Flow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the single-step "Create Program" dialog with a 3-option flow (From Template, Import, From Scratch), add Program name uniqueness enforcement, and add Export to the Program detail page.

**Architecture:** API-first — add uniqueness constraint and `ExistsByName` check before touching the frontend. Frontend builds 5 focused step-components managed by a single root `CreateProgramDialog`. Import uses Zod for runtime schema validation.

**Tech Stack:** Go (API), Next.js / TypeScript / TanStack Query / Zod / shadcn/ui (frontend)

**Spec:** `docs/superpowers/specs/2026-03-24-create-program-flow-design.md`

---

## File Map

### API (new/modified)
| File | Change |
|---|---|
| `api/migrations/000015_program_name_unique.up.sql` | CREATE UNIQUE INDEX on programs.name |
| `api/migrations/000015_program_name_unique.down.sql` | DROP INDEX |
| `api/internal/repository/program.go` | Add `ExistsByName` to interface |
| `api/internal/store/postgres/program.go` | Implement `ExistsByName` |
| `api/internal/store/memory/program.go` | Implement `ExistsByName` |
| `api/internal/handler/program.go` | Check name conflict → 409 in `CreateProgram` |
| `api/internal/handler/program_test.go` | Add 409 conflict test |

### Frontend (new/modified)
| File | Change |
|---|---|
| `web/components/programs/CreateProgramDialog.tsx` | Root dialog, step state machine |
| `web/components/programs/create-program/ChoiceStep.tsx` | 3-option card selection |
| `web/components/programs/create-program/ScratchStep.tsx` | Extracted from programs/page.tsx |
| `web/components/programs/create-program/TemplateSelectStep.tsx` | Searchable template list |
| `web/components/programs/create-program/TemplateConfigStep.tsx` | 1RM + Increment inputs |
| `web/components/programs/create-program/ImportStep.tsx` | Drag/drop + paste + Zod |
| `web/components/programs/create-program/importSchema.ts` | Zod schema for import format |
| `web/app/programs/page.tsx` | Replace inline dialog with CreateProgramDialog |
| `web/app/programs/[id]/page.tsx` | Add Export button (clipboard + download) |

---

## Task 1: DB migration — Program name uniqueness

**Files:**
- Create: `api/migrations/000015_program_name_unique.up.sql`
- Create: `api/migrations/000015_program_name_unique.down.sql`

- [ ] **Step 1: Write up migration**

```sql
-- api/migrations/000015_program_name_unique.up.sql
CREATE UNIQUE INDEX programs_name_unique ON programs (name);
```

- [ ] **Step 2: Write down migration**

```sql
-- api/migrations/000015_program_name_unique.down.sql
DROP INDEX IF EXISTS programs_name_unique;
```

- [ ] **Step 3: Verify migration files exist and are named correctly**

Run from `api/`:
```bash
ls migrations/ | grep 000015
```
Expected: `000015_program_name_unique.down.sql` and `000015_program_name_unique.up.sql`

- [ ] **Step 4: Commit**

```bash
git add api/migrations/000015_program_name_unique.up.sql api/migrations/000015_program_name_unique.down.sql
git commit -m "chore: add unique index on programs.name"
```

---

## Task 2: Repository interface — ExistsByName

**Files:**
- Modify: `api/internal/repository/program.go`

- [ ] **Step 1: Add `ExistsByName` to the `ProgramRepository` interface**

In `api/internal/repository/program.go`, add to the interface (alongside `ExistsByProgramTemplateID`):

```go
// ExistsByName checks if a Program with the given name already exists.
ExistsByName(ctx context.Context, name string) (bool, error)
```

- [ ] **Step 2: Build to verify interface is unsatisfied**

Run from `api/`:
```bash
go build ./...
```
Expected: compile errors that `ExistsByName` is not implemented by postgres and memory stores.

---

## Task 3: Postgres store — ExistsByName

**Files:**
- Modify: `api/internal/store/postgres/program.go`

- [ ] **Step 1: Write the failing integration test**

In `api/internal/store/postgres/program_test.go` (file already exists with integration tests tagged `//go:build integration`), add:

```go
func TestProgramRepository_ExistsByName(t *testing.T) {
    ctx := context.Background()
    repo := setupTestProgramRepo(t) // use existing test helper

    t.Run("returns false when no program with name", func(t *testing.T) {
        exists, err := repo.ExistsByName(ctx, "NonExistentXYZ")
        require.NoError(t, err)
        assert.False(t, exists)
    })

    t.Run("returns true when program with name exists", func(t *testing.T) {
        p := &domain.Program{Name: "My Program", Status: domain.ProgramStatusCreated}
        require.NoError(t, repo.Create(ctx, p))

        exists, err := repo.ExistsByName(ctx, "My Program")
        require.NoError(t, err)
        assert.True(t, exists)
    })
}
```

- [ ] **Step 2: Implement `ExistsByName` in postgres store**

In `api/internal/store/postgres/program.go`, add (alongside `ExistsByProgramTemplateID`):

```go
func (r *programRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
    var exists bool
    err := r.pool.QueryRow(ctx,
        `SELECT EXISTS(SELECT 1 FROM programs WHERE name = $1)`,
        name,
    ).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("check program name exists: %w", err)
    }
    return exists, nil
}
```

- [ ] **Step 3: Build to verify compile passes**

```bash
go build ./...
```
Expected: only memory store still failing.

---

## Task 4: Memory store — ExistsByName

**Files:**
- Modify: `api/internal/store/memory/program.go`

- [ ] **Step 1: Write the unit test**

In `api/internal/store/memory/program_test.go` (file already exists), add:

```go
func TestProgramMemoryStore_ExistsByName(t *testing.T) {
    store := NewProgramRepository()
    ctx := context.Background()

    t.Run("returns false when empty", func(t *testing.T) {
        exists, err := store.ExistsByName(ctx, "Test")
        require.NoError(t, err)
        assert.False(t, exists)
    })

    t.Run("returns true after create", func(t *testing.T) {
        p := &domain.Program{Name: "Test", Status: domain.ProgramStatusCreated}
        require.NoError(t, store.Create(ctx, p))

        exists, err := store.ExistsByName(ctx, "Test")
        require.NoError(t, err)
        assert.True(t, exists)
    })
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/store/memory/... -run TestProgramMemoryStore_ExistsByName -v
```
Expected: FAIL (method not defined)

- [ ] **Step 3: Implement `ExistsByName` in memory store**

In `api/internal/store/memory/program.go`:

```go
func (s *programStore) ExistsByName(ctx context.Context, name string) (bool, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, p := range s.programs {
        if p.Name == name {
            return true, nil
        }
    }
    return false, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/store/memory/... -run TestProgramMemoryStore_ExistsByName -v
```
Expected: PASS

- [ ] **Step 5: Run full unit tests**

```bash
go test ./...
```
Expected: all pass

- [ ] **Step 6: Commit**

```bash
git add api/internal/repository/program.go \
  api/internal/store/postgres/program.go \
  api/internal/store/postgres/program_test.go \
  api/internal/store/memory/program.go \
  api/internal/store/memory/program_test.go
git commit -m "feat: add ExistsByName to ProgramRepository"
```

---

## Task 5: Handler — 409 on duplicate program name

**Files:**
- Modify: `api/internal/handler/program.go`
- Modify: `api/internal/handler/program_test.go`

- [ ] **Step 1: Write the failing handler test**

Handler tests use the real memory store via `setupProgramTestRouter()` (existing helper). Add a test that creates a program, then tries to create another with the same name:

```go
func TestCreateProgram_ConflictOnDuplicateName(t *testing.T) {
    router, _ := setupProgramTestRouter()

    // First creation — should succeed
    body := `{"name":"Duplicate Program","sessions":[]}`
    req := httptest.NewRequest(http.MethodPost, "/programs", bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer user1")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    require.Equal(t, http.StatusCreated, w.Code)

    // Second creation — same name, should conflict
    req2 := httptest.NewRequest(http.MethodPost, "/programs", bytes.NewBufferString(body))
    req2.Header.Set("Content-Type", "application/json")
    req2.Header.Set("Authorization", "Bearer user1")
    w2 := httptest.NewRecorder()
    router.ServeHTTP(w2, req2)

    assert.Equal(t, http.StatusConflict, w2.Code)
    var resp map[string]interface{}
    require.NoError(t, json.NewDecoder(w2.Body).Decode(&resp))
    assert.Equal(t, "CONFLICT", resp["code"])
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/handler/... -run TestCreateProgram_ConflictOnDuplicateName -v
```
Expected: FAIL (second request returns 201, not 409)

- [ ] **Step 3: Add name conflict check to `CreateProgram` handler**

In `api/internal/handler/program.go`, after decoding the request and before creating the domain object, add:

```go
exists, err := h.repo.ExistsByName(ctx, req.Name)
if err != nil {
    middleware.WriteInternalError(w, "Failed to check program name")
    return
}
if exists {
    middleware.WriteConflictError(w, "A program with this name already exists", map[string]interface{}{
        "field": "name",
    })
    return
}
```

Add this block after `json.NewDecoder(r.Body).Decode(&req)` succeeds and before building the domain `Program` object.

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/handler/... -run TestCreateProgram_ConflictOnDuplicateName -v
```
Expected: PASS

- [ ] **Step 5: Run full check**

```bash
task check
```
Expected: all pass

- [ ] **Step 6: Commit**

```bash
git add api/internal/handler/program.go api/internal/handler/program_test.go
git commit -m "feat: return 409 when creating program with duplicate name"
```

---

## Task 6: Frontend — Zod import schema

**Files:**
- Create: `web/components/programs/create-program/importSchema.ts`

- [ ] **Step 1: Create the Zod schema file**

```typescript
// web/components/programs/create-program/importSchema.ts
import { z } from 'zod';

const programSessionEntryImportSchema = z.object({
  exercise_name: z.string().min(1),
  order: z.number().int().min(0),
  sets: z.number().int().min(1).optional(),
  reps: z.number().int().min(1).optional(),
  load_kg: z.number().positive().optional(),
  rpe: z.number().int().min(1).max(10).optional(),
  notes: z.string().optional(),
});

const programSessionImportSchema = z.object({
  session_name: z.string().min(1),
  order: z.number().int().min(0),
  date: z.string().optional(),
  entries: z.array(programSessionEntryImportSchema).optional(),
});

export const programImportSchema = z.object({
  rx_version: z.literal('1'),
  name: z.string().min(1),
  notes: z.string().optional(),
  sessions: z.array(programSessionImportSchema),
});

export type ProgramImport = z.infer<typeof programImportSchema>;
```

- [ ] **Step 2: Verify TypeScript compiles**

```bash
pnpm build 2>&1 | head -20
```
Expected: no type errors in new file (other errors unrelated to this change are fine)

- [ ] **Step 3: Commit**

```bash
git add web/components/programs/create-program/importSchema.ts
git commit -m "feat: add Zod schema for program import format"
```

---

## Task 7: Frontend — ScratchStep (extract existing form)

**Files:**
- Create: `web/components/programs/create-program/ScratchStep.tsx`

- [ ] **Step 1: Create ScratchStep by extracting the existing dialog body**

The current inline form in `web/app/programs/page.tsx` (lines 27–88) handles name, notes, sessions. Extract it into a self-contained component:

```typescript
// web/components/programs/create-program/ScratchStep.tsx
'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { ProgramSessionCreate } from '@/types/api';
import { Plus, X } from 'lucide-react';
import { useState } from 'react';

interface SessionInput {
  session_name: string;
  order: number;
}

interface ScratchStepProps {
  onBack: () => void;
  onSubmit: (data: { name: string; notes?: string; sessions: ProgramSessionCreate[] }) => void;
  isPending: boolean;
  nameError?: string;
}

export function ScratchStep({ onBack, onSubmit, isPending, nameError }: ScratchStepProps) {
  const [name, setName] = useState('');
  const [notes, setNotes] = useState('');
  const [sessions, setSessions] = useState<SessionInput[]>([{ session_name: '', order: 0 }]);

  const handleAddSession = () => {
    setSessions([...sessions, { session_name: '', order: sessions.length }]);
  };

  const handleRemoveSession = (idx: number) => {
    setSessions(sessions.filter((_, i) => i !== idx).map((s, i) => ({ ...s, order: i })));
  };

  const handleSessionNameChange = (idx: number, value: string) => {
    const updated = [...sessions];
    const s = updated[idx];
    if (s) updated[idx] = { ...s, session_name: value };
    setSessions(updated);
  };

  const handleSubmit = () => {
    const programSessions: ProgramSessionCreate[] = sessions
      .filter((s) => s.session_name.trim())
      .map((s) => ({ session_name: s.session_name.trim(), order: s.order }));
    onSubmit({ name, notes: notes || undefined, sessions: programSessions });
  };

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="scratch-name">Program Name</Label>
        <Input
          id="scratch-name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="e.g., SBD Block 1"
        />
        {nameError && <p className="text-sm text-destructive">{nameError}</p>}
      </div>
      <div className="space-y-2">
        <Label htmlFor="scratch-notes">Notes</Label>
        <Input
          id="scratch-notes"
          value={notes}
          onChange={(e) => setNotes(e.target.value)}
          placeholder="Optional notes"
        />
      </div>
      <div className="space-y-2">
        <Label>Sessions</Label>
        <div className="space-y-2">
          {sessions.map((session, idx) => (
            <div key={idx} className="flex gap-2">
              <Input
                value={session.session_name}
                onChange={(e) => handleSessionNameChange(idx, e.target.value)}
                placeholder={`e.g., Week 1 Day ${idx + 1}`}
              />
              {sessions.length > 1 && (
                <Button variant="ghost" size="sm" onClick={() => handleRemoveSession(idx)}>
                  <X className="h-4 w-4" />
                </Button>
              )}
            </div>
          ))}
        </div>
        <Button variant="outline" size="sm" onClick={handleAddSession} className="w-full">
          <Plus className="h-4 w-4 mr-2" />
          Add Session
        </Button>
      </div>
      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onBack}>
          Back
        </Button>
        <Button onClick={handleSubmit} disabled={isPending || !name.trim()}>
          {isPending ? 'Creating...' : 'Create Program'}
        </Button>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Verify TypeScript compiles**

```bash
pnpm build 2>&1 | head -20
```

- [ ] **Step 3: Commit**

```bash
git add web/components/programs/create-program/ScratchStep.tsx
git commit -m "feat: extract ScratchStep from programs page inline dialog"
```

---

## Task 8: Frontend — ChoiceStep

**Files:**
- Create: `web/components/programs/create-program/ChoiceStep.tsx`

- [ ] **Step 1: Create ChoiceStep**

```typescript
// web/components/programs/create-program/ChoiceStep.tsx
'use client';

import { Button } from '@/components/ui/button';
import { Download, FileJson, Pencil } from 'lucide-react';

export type CreationMethod = 'template' | 'import' | 'scratch';

interface ChoiceStepProps {
  onSelect: (method: CreationMethod) => void;
}

const options: { method: CreationMethod; icon: React.ReactNode; title: string; description: string }[] = [
  {
    method: 'template',
    icon: <Download className="h-6 w-6" />,
    title: 'From Template',
    description: 'Generate from a template with calculated loads',
  },
  {
    method: 'import',
    icon: <FileJson className="h-6 w-6" />,
    title: 'Import',
    description: 'Import from JSON (shared by coach or another device)',
  },
  {
    method: 'scratch',
    icon: <Pencil className="h-6 w-6" />,
    title: 'From Scratch',
    description: 'Create manually with custom sessions',
  },
];

export function ChoiceStep({ onSelect }: ChoiceStepProps) {
  return (
    <div className="grid gap-3">
      {options.map(({ method, icon, title, description }) => (
        <button
          key={method}
          type="button"
          onClick={() => onSelect(method)}
          className="flex items-center gap-4 rounded-lg border p-4 text-left hover:bg-accent transition-colors"
        >
          <div className="text-muted-foreground">{icon}</div>
          <div>
            <p className="font-medium">{title}</p>
            <p className="text-sm text-muted-foreground">{description}</p>
          </div>
        </button>
      ))}
    </div>
  );
}
```

- [ ] **Step 2: Verify TypeScript compiles**

```bash
pnpm build 2>&1 | head -20
```

- [ ] **Step 3: Commit**

```bash
git add web/components/programs/create-program/ChoiceStep.tsx
git commit -m "feat: add ChoiceStep for Create Program dialog"
```

---

## Task 9: Frontend — TemplateSelectStep

**Files:**
- Create: `web/components/programs/create-program/TemplateSelectStep.tsx`

- [ ] **Step 1: Create TemplateSelectStep**

```typescript
// web/components/programs/create-program/TemplateSelectStep.tsx
'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Skeleton } from '@/components/ui/skeleton';
import { useProgramTemplates } from '@/lib/hooks/useProgramTemplates';
import type { ProgramTemplate } from '@/types/api';
import { Search } from 'lucide-react';
import { useState } from 'react';

interface TemplateSelectStepProps {
  onBack: () => void;
  onSelect: (template: ProgramTemplate) => void;
}

export function TemplateSelectStep({ onBack, onSelect }: TemplateSelectStepProps) {
  const { data, isLoading } = useProgramTemplates(false); // exclude archived
  const [query, setQuery] = useState('');

  const templates = (data?.data ?? []).filter((t) =>
    t.name.toLowerCase().includes(query.toLowerCase())
  );

  return (
    <div className="space-y-4">
      <div className="relative">
        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
        <Input
          className="pl-8"
          placeholder="Search templates..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
      </div>

      {isLoading ? (
        <div className="space-y-2">
          <Skeleton className="h-16" />
          <Skeleton className="h-16" />
        </div>
      ) : templates.length === 0 ? (
        <p className="text-sm text-muted-foreground text-center py-6">
          {query ? 'No templates match your search.' : 'No templates yet.'}
        </p>
      ) : (
        <div className="space-y-2 max-h-64 overflow-y-auto">
          {templates.map((template) => {
            const exerciseNames = [
              ...new Set((template.entries ?? []).map((e) => e.exercise_name)),
            ].slice(0, 3);
            const sessionCount = new Set(
              (template.entries ?? []).map((e) => e.metadata?.session as string)
            ).size;
            return (
              <button
                key={template.id}
                type="button"
                onClick={() => onSelect(template)}
                className="w-full flex flex-col gap-1 rounded-lg border p-3 text-left hover:bg-accent transition-colors"
              >
                <p className="font-medium">{template.name}</p>
                <p className="text-xs text-muted-foreground">
                  {sessionCount} session{sessionCount !== 1 ? 's' : ''}
                  {exerciseNames.length > 0 && ` · ${exerciseNames.join(', ')}${(template.entries?.length ?? 0) > 3 ? '...' : ''}`}
                </p>
              </button>
            );
          })}
        </div>
      )}

      <div className="flex justify-start">
        <Button variant="outline" onClick={onBack}>
          Back
        </Button>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Verify TypeScript compiles**

```bash
pnpm build 2>&1 | head -20
```

- [ ] **Step 3: Commit**

```bash
git add web/components/programs/create-program/TemplateSelectStep.tsx
git commit -m "feat: add TemplateSelectStep for Create Program dialog"
```

---

## Task 10: Frontend — TemplateConfigStep

**Files:**
- Create: `web/components/programs/create-program/TemplateConfigStep.tsx`

- [ ] **Step 1: Create TemplateConfigStep**

This step extracts unique exercises from the template and collects 1RM/target weight + increment per exercise. Exercises with `percent_1rm` → "1RM" label. Exercises with only `rpe` → "Weight" label (direct copy).

```typescript
// web/components/programs/create-program/TemplateConfigStep.tsx
'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { GenerateProgramRequest, ProgramTemplate } from '@/types/api';
import { useState } from 'react';

interface ExerciseConfig {
  name: string;
  hasPercentRm: boolean; // true → user enters 1RM; false → user enters direct weight
  weight: string;
  increment: string;
}

function buildExerciseConfigs(template: ProgramTemplate): ExerciseConfig[] {
  const seen = new Map<string, boolean>();
  for (const entry of template.entries ?? []) {
    if (!seen.has(entry.exercise_name)) {
      seen.set(entry.exercise_name, entry.percent_1rm != null);
    } else if (entry.percent_1rm != null) {
      seen.set(entry.exercise_name, true);
    }
  }
  return [...seen.entries()].map(([name, hasPercentRm]) => ({
    name,
    hasPercentRm,
    weight: '',
    increment: '2.5',
  }));
}

interface TemplateConfigStepProps {
  template: ProgramTemplate;
  onBack: () => void;
  onSubmit: (templateId: string, data: GenerateProgramRequest) => void;
  isPending: boolean;
  nameError?: string;
}

export function TemplateConfigStep({
  template,
  onBack,
  onSubmit,
  isPending,
  nameError,
}: TemplateConfigStepProps) {
  const [programName, setProgramName] = useState(template.name);
  const [exercises, setExercises] = useState<ExerciseConfig[]>(() =>
    buildExerciseConfigs(template)
  );

  const updateExercise = (idx: number, field: 'weight' | 'increment', value: string) => {
    setExercises((prev) => prev.map((ex, i) => (i === idx ? { ...ex, [field]: value } : ex)));
  };

  const allFilled = exercises.every((ex) => {
    const w = parseFloat(ex.weight);
    return !isNaN(w) && w > 0;
  });
  const hasExercises = exercises.length > 0;

  const handleSubmit = () => {
    const target_weights: Record<string, number> = {};
    const load_increments: Record<string, number> = {};
    for (const ex of exercises) {
      target_weights[ex.name] = parseFloat(ex.weight);
      const inc = parseFloat(ex.increment);
      if (!isNaN(inc) && inc > 0) load_increments[ex.name] = inc;
    }
    onSubmit(template.id, { name: programName, target_weights, load_increments });
  };

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="template-program-name">Program Name</Label>
        <Input
          id="template-program-name"
          value={programName}
          onChange={(e) => setProgramName(e.target.value)}
        />
        {nameError && <p className="text-sm text-destructive">{nameError}</p>}
      </div>

      {!hasExercises ? (
        <p className="text-sm text-muted-foreground">This template has no exercises.</p>
      ) : (
        <div className="space-y-3 max-h-72 overflow-y-auto">
          {exercises.map((ex, idx) => (
            <div key={ex.name} className="rounded-lg border p-3 space-y-2">
              <div className="flex items-center justify-between">
                <p className="font-medium text-sm">{ex.name}</p>
                <span className="text-xs text-muted-foreground">
                  {ex.hasPercentRm ? '% 1RM' : 'RPE only'}
                </span>
              </div>
              <div className="grid grid-cols-2 gap-2">
                <div className="space-y-1">
                  <Label className="text-xs">{ex.hasPercentRm ? '1RM (kg)' : 'Weight (kg)'}</Label>
                  <Input
                    type="number"
                    min={0.1}
                    step={0.5}
                    value={ex.weight}
                    onChange={(e) => updateExercise(idx, 'weight', e.target.value)}
                    placeholder="e.g., 140"
                  />
                </div>
                <div className="space-y-1">
                  <Label className="text-xs">Increment (kg)</Label>
                  <Input
                    type="number"
                    min={0.5}
                    step={0.5}
                    value={ex.increment}
                    onChange={(e) => updateExercise(idx, 'increment', e.target.value)}
                  />
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onBack}>
          Back
        </Button>
        <Button
          onClick={handleSubmit}
          disabled={isPending || !programName.trim() || !allFilled || !hasExercises}
        >
          {isPending ? 'Generating...' : 'Generate Program'}
        </Button>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Verify TypeScript compiles**

```bash
pnpm build 2>&1 | head -20
```

- [ ] **Step 3: Commit**

```bash
git add web/components/programs/create-program/TemplateConfigStep.tsx
git commit -m "feat: add TemplateConfigStep with 1RM and increment inputs"
```

---

## Task 11: Frontend — ImportStep

**Files:**
- Create: `web/components/programs/create-program/ImportStep.tsx`

- [ ] **Step 1: Create ImportStep**

```typescript
// web/components/programs/create-program/ImportStep.tsx
'use client';

import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { programImportSchema, type ProgramImport } from './importSchema';
import { Upload } from 'lucide-react';
import { useCallback, useRef, useState } from 'react';

interface ImportStepProps {
  onBack: () => void;
  onSubmit: (data: ProgramImport) => void;
  isPending: boolean;
  nameError?: string;
}

export function ImportStep({ onBack, onSubmit, isPending, nameError }: ImportStepProps) {
  const [text, setText] = useState('');
  const [parseError, setParseError] = useState<string | null>(null);
  const [parsed, setParsed] = useState<ProgramImport | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const parseAndStore = (raw: string) => {
    if (!raw.trim()) {
      setParseError(null);
      setParsed(null);
      return;
    }
    try {
      const json = JSON.parse(raw);
      const result = programImportSchema.safeParse(json);
      if (!result.success) {
        const firstError = result.error.errors[0];
        setParseError(`Invalid format: ${firstError?.path.join('.') ?? ''} — ${firstError?.message}`);
        setParsed(null);
        return;
      }
      setParseError(null);
      setParsed(result.data);
    } catch {
      setParseError('Invalid JSON: could not parse');
      setParsed(null);
    }
  };

  const handleTextChange = (value: string) => {
    setText(value);
    parseAndStore(value);
  };

  const loadFile = useCallback((file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      setText(content);
      parseAndStore(content);
    };
    reader.readAsText(file);
  }, []);

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    const file = e.dataTransfer.files[0];
    if (file) loadFile(file);
  };

  const canImport = parsed !== null && !isPending;

  return (
    <div className="space-y-3">
      <div
        onDragOver={(e) => { e.preventDefault(); setIsDragging(true); }}
        onDragLeave={() => setIsDragging(false)}
        onDrop={handleDrop}
        className={`rounded-lg border-2 border-dashed p-6 text-center transition-colors ${
          isDragging ? 'border-primary bg-accent' : 'border-muted-foreground/30'
        }`}
      >
        <Upload className="h-6 w-6 mx-auto mb-2 text-muted-foreground" />
        <p className="text-sm text-muted-foreground">Drop JSON file here or paste below</p>
      </div>

      <Textarea
        value={text}
        onChange={(e) => handleTextChange(e.target.value)}
        placeholder={'{\n  "rx_version": "1",\n  "name": "...",\n  "sessions": []\n}'}
        className="font-mono text-xs h-32"
      />

      {(parseError || nameError) && (
        <p className="text-sm text-destructive">{nameError ?? parseError}</p>
      )}

      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() => fileInputRef.current?.click()}
          type="button"
        >
          Browse file
        </Button>
        <input
          ref={fileInputRef}
          type="file"
          accept=".json"
          className="hidden"
          onChange={(e) => { const f = e.target.files?.[0]; if (f) loadFile(f); }}
        />
      </div>

      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onBack}>
          Back
        </Button>
        <Button
          onClick={() => { if (parsed) onSubmit(parsed); }}
          disabled={!canImport}
        >
          {isPending ? 'Importing...' : 'Import'}
        </Button>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Verify TypeScript compiles**

```bash
pnpm build 2>&1 | head -20
```

- [ ] **Step 3: Commit**

```bash
git add web/components/programs/create-program/ImportStep.tsx
git commit -m "feat: add ImportStep with drag-drop, file browse, and Zod validation"
```

---

## Task 12: Frontend — CreateProgramDialog (root)

**Files:**
- Create: `web/components/programs/CreateProgramDialog.tsx`

- [ ] **Step 1: Create root dialog**

This component owns the step state machine and all mutations.

```typescript
// web/components/programs/CreateProgramDialog.tsx
'use client';

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { useCreateProgram } from '@/lib/hooks/usePrograms';
import { useGenerateProgram } from '@/lib/hooks/useProgramTemplates';
import type { GenerateProgramRequest, ProgramTemplate } from '@/types/api';
import { useState } from 'react';
import { ChoiceStep, type CreationMethod } from './create-program/ChoiceStep';
import { ImportStep } from './create-program/ImportStep';
import { ScratchStep } from './create-program/ScratchStep';
import { TemplateConfigStep } from './create-program/TemplateConfigStep';
import { TemplateSelectStep } from './create-program/TemplateSelectStep';
import type { ProgramImport } from './create-program/importSchema';

type Step =
  | { type: 'choice' }
  | { type: 'template-select' }
  | { type: 'template-config'; template: ProgramTemplate }
  | { type: 'import' }
  | { type: 'scratch' };

interface CreateProgramDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const STEP_TITLES: Record<Step['type'], string> = {
  choice: 'Create Program',
  'template-select': 'Create from Template',
  'template-config': 'Create from Template',
  import: 'Import Program',
  scratch: 'Create from Scratch',
};

const STEP_DESCRIPTIONS: Record<Step['type'], string> = {
  choice: 'How would you like to create it?',
  'template-select': 'Select a template to generate from.',
  'template-config': 'Enter your weights to calculate working loads.',
  import: 'Paste or drop a JSON file exported from this app.',
  scratch: 'Define the program name and session names.',
};

function conflictError(err: unknown): string | undefined {
  if (err && typeof err === 'object' && 'status' in err && (err as { status: number }).status === 409) {
    return 'A program with this name already exists';
  }
  return undefined;
}

export function CreateProgramDialog({ open, onOpenChange }: CreateProgramDialogProps) {
  const [step, setStep] = useState<Step>({ type: 'choice' });
  const [nameError, setNameError] = useState<string | undefined>();
  const createProgram = useCreateProgram();
  const generateProgram = useGenerateProgram();

  const handleOpenChange = (value: boolean) => {
    onOpenChange(value);
    if (!value) {
      setStep({ type: 'choice' });
      setNameError(undefined);
    }
  };

  const handleMethodSelect = (method: CreationMethod) => {
    if (method === 'template') setStep({ type: 'template-select' });
    else if (method === 'import') setStep({ type: 'import' });
    else setStep({ type: 'scratch' });
    setNameError(undefined);
  };

  const handleScratchSubmit = async (data: { name: string; notes?: string; sessions: Parameters<typeof createProgram.mutateAsync>[0]['sessions'] }) => {
    try {
      await createProgram.mutateAsync(data);
      handleOpenChange(false);
    } catch (err) {
      setNameError(conflictError(err) ?? 'Failed to create program');
    }
  };

  const handleGenerateSubmit = async (templateId: string, data: GenerateProgramRequest) => {
    try {
      await generateProgram.mutateAsync({ id: templateId, data });
      handleOpenChange(false);
    } catch (err) {
      setNameError(conflictError(err) ?? 'Failed to generate program');
    }
  };

  const handleImportSubmit = async (importData: ProgramImport) => {
    const { rx_version: _v, ...payload } = importData;
    try {
      await createProgram.mutateAsync(payload);
      handleOpenChange(false);
    } catch (err) {
      setNameError(conflictError(err) ?? 'Failed to import program');
    }
  };

  const isPending = createProgram.isPending || generateProgram.isPending;

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{STEP_TITLES[step.type]}</DialogTitle>
          <DialogDescription>{STEP_DESCRIPTIONS[step.type]}</DialogDescription>
        </DialogHeader>

        {step.type === 'choice' && <ChoiceStep onSelect={handleMethodSelect} />}
        {step.type === 'template-select' && (
          <TemplateSelectStep
            onBack={() => setStep({ type: 'choice' })}
            onSelect={(t) => { setStep({ type: 'template-config', template: t }); setNameError(undefined); }}
          />
        )}
        {step.type === 'template-config' && (
          <TemplateConfigStep
            template={step.template}
            onBack={() => setStep({ type: 'template-select' })}
            onSubmit={handleGenerateSubmit}
            isPending={isPending}
            nameError={nameError}
          />
        )}
        {step.type === 'import' && (
          <ImportStep
            onBack={() => setStep({ type: 'choice' })}
            onSubmit={handleImportSubmit}
            isPending={isPending}
            nameError={nameError}
          />
        )}
        {step.type === 'scratch' && (
          <ScratchStep
            onBack={() => setStep({ type: 'choice' })}
            onSubmit={handleScratchSubmit}
            isPending={isPending}
            nameError={nameError}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}
```

- [ ] **Step 2: Verify TypeScript compiles**

```bash
pnpm build 2>&1 | head -20
```

- [ ] **Step 3: Commit**

```bash
git add web/components/programs/CreateProgramDialog.tsx
git commit -m "feat: add CreateProgramDialog with 3-step flow"
```

---

## Task 13: Frontend — Wire up programs/page.tsx

**Files:**
- Modify: `web/app/programs/page.tsx`

- [ ] **Step 1: Replace inline dialog with CreateProgramDialog**

In `web/app/programs/page.tsx`:

1. Remove all state variables and handlers related to the old inline form (`name`, `notes`, `sessions`, `handleAddSession`, `handleRemoveSession`, `handleSessionNameChange`, `handleSave`, `handleOpenChange`, and the `<Dialog>` JSX block)
2. Remove the `useCreateProgram` import and usage
3. Remove the `SessionInput` interface and `ProgramSessionCreate` import
4. Add the import: `import { CreateProgramDialog } from '@/components/programs/CreateProgramDialog';`
5. Keep only `const [open, setOpen] = useState(false);`
6. Replace `<Dialog ...>...</Dialog>` with: `<CreateProgramDialog open={open} onOpenChange={setOpen} />`

- [ ] **Step 2: Verify pnpm check passes**

```bash
pnpm check
```
Expected: no lint or type errors

- [ ] **Step 3: Manual smoke test**

Start the dev server (`pnpm dev`) and verify:
- "Create Program" button opens the choice dialog
- All 3 flows reach their final step without errors
- From Scratch creates a program and closes the dialog

- [ ] **Step 4: Commit**

```bash
git add web/app/programs/page.tsx
git commit -m "feat: replace inline dialog with CreateProgramDialog in programs page"
```

---

## Task 14: Frontend — Export on Program detail page

**Files:**
- Modify: `web/app/programs/[id]/page.tsx`

- [ ] **Step 1: Add export utility function**

At the top of `web/app/programs/[id]/page.tsx`, add a helper that serializes a Program to the import/export JSON format:

```typescript
import type { Program } from '@/types/api';

function programToExportJson(program: Program): string {
  const payload = {
    rx_version: '1',
    name: program.name,
    ...(program.notes ? { notes: program.notes } : {}),
    sessions: program.sessions.map((s) => ({
      session_name: s.session_name,
      order: s.order,
      ...(s.date ? { date: s.date } : {}),
      entries: s.entries.map((e) => ({
        exercise_name: e.exercise_name,
        order: e.order,
        ...(e.sets != null ? { sets: e.sets } : {}),
        ...(e.reps != null ? { reps: e.reps } : {}),
        ...(e.load_kg != null ? { load_kg: e.load_kg } : {}),
        ...(e.rpe != null ? { rpe: e.rpe } : {}),
        ...(e.notes ? { notes: e.notes } : {}),
      })),
    })),
  };
  return JSON.stringify(payload, null, 2);
}
```

- [ ] **Step 2: Add export buttons to the detail page header**

In `ProgramDetailPage`, add two export actions to the button group in the header alongside the existing status buttons. Use shadcn `DropdownMenu` for the export sub-menu:

```typescript
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Copy, Download, Share2 } from 'lucide-react';
import { useState } from 'react';

// Inside ProgramDetailPage:
const [copied, setCopied] = useState(false);

const handleCopyToClipboard = async () => {
  await navigator.clipboard.writeText(programToExportJson(program));
  setCopied(true);
  setTimeout(() => setCopied(false), 2000);
};

const handleDownload = () => {
  const json = programToExportJson(program);
  const blob = new Blob([json], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${program.name}.json`;
  a.click();
  URL.revokeObjectURL(url);
};
```

Add to the button group in JSX:

```tsx
<DropdownMenu>
  <DropdownMenuTrigger asChild>
    <Button variant="outline" size="sm">
      <Share2 className="h-4 w-4 mr-2" />
      Export
    </Button>
  </DropdownMenuTrigger>
  <DropdownMenuContent align="end">
    <DropdownMenuItem onClick={handleCopyToClipboard}>
      <Copy className="h-4 w-4 mr-2" />
      {copied ? 'Copied!' : 'Copy to clipboard'}
    </DropdownMenuItem>
    <DropdownMenuItem onClick={handleDownload}>
      <Download className="h-4 w-4 mr-2" />
      Download .json
    </DropdownMenuItem>
  </DropdownMenuContent>
</DropdownMenu>
```

- [ ] **Step 3: Verify pnpm check passes**

```bash
pnpm check
```

- [ ] **Step 4: Manual smoke test**

Verify Export button appears on the Program detail page, clipboard copy works, and .json download works.

- [ ] **Step 5: Commit**

```bash
git add web/app/programs/[id]/page.tsx
git commit -m "feat: add export (clipboard + download) to Program detail page"
```

---

## Completion

- [ ] Run `task check` from `api/` — all pass
- [ ] Run `pnpm check` from `web/` — all pass
- [ ] End-to-end smoke test: create program via all 3 flows, export, import back
- [ ] Verify 409 is shown when creating a program with a duplicate name
