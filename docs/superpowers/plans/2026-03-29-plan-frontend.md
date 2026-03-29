# Plan Frontend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the Dashboard-centric frontend with a Plan-centric design where the Plan page is the top page (`/`), showing a session execution queue with a program sidebar.

**Architecture:** Split layout (70/30) Plan page with session cards on the left and program list on the right. Sessions are rendered as compact fixed-width tables with dynamic field columns. Log creation captures `plan_snapshot` and auto-removes the session from the queue. Logs page redesigned as card-based with plan vs actual diff.

**Tech Stack:** Next.js (App Router), TanStack Query, shadcn/ui, ky (HTTP client), TypeScript

**Spec:** `docs/superpowers/specs/2026-03-29-plan-frontend-design.md`

---

## File Structure

### New Files

| File | Responsibility |
|------|---------------|
| `web/lib/api/plans.ts` | Plan API client (CRUD + sessions) |
| `web/lib/hooks/usePlans.ts` | TanStack Query hooks for Plan |
| `web/components/plan/SessionCard.tsx` | Single session card with dynamic field table |
| `web/components/plan/ProgramSidebar.tsx` | Right sidebar program list |
| `web/components/plan/AddSessionDialog.tsx` | Manual session creation dialog |
| `web/components/plan/EmptyState.tsx` | Empty plan state |
| `web/components/programs/SessionSelector.tsx` | Checkbox session list for "Add to Plan" |
| `web/components/logs/LogCard.tsx` | Card component for Logs page |

### Modified Files

| File | Change |
|------|--------|
| `web/types/api.ts` | Add Plan types + `plan_snapshot` to Log/LogCreate |
| `web/components/layout/Sidebar.tsx` | 3-item nav: Plan / Logs / Settings |
| `web/app/page.tsx` | Rewrite as Plan page |
| `web/app/programs/[id]/page.tsx` | Add "Add to Plan" button + SessionSelector |
| `web/app/logs/page.tsx` | Card-based layout with LogCard |
| `web/app/logs/new/page.tsx` | Support plan_snapshot from Plan session |

### Removed Files

| File | Reason |
|------|--------|
| `web/components/dashboard/ActivePlans.tsx` | Dashboard removed |
| `web/components/dashboard/RecentLogs.tsx` | Dashboard removed |
| `web/app/programs/page.tsx` | Programs list removed (sidebar replaces it) |
| `web/components/logs/LogTable.tsx` | Replaced by LogCard |

---

## Task 1: Add Plan Types to api.ts

**Files:**
- Modify: `web/types/api.ts`

- [ ] **Step 1: Add Plan type definitions**

Add the following after the `ProgramSessionEntryCreate` interface (before the LoggedSessions section):

```typescript
// ============================================================================
// Plan
// ============================================================================

export interface Plan {
  id: string;
  name?: string;
  notes?: string;
  sessions: PlanSession[];
  created_at: string;
  updated_at: string;
}

export interface PlanCreate {
  name?: string;
  notes?: string;
  sessions?: PlanSessionCreate[];
}

export interface PlanUpdate {
  name?: string;
  notes?: string;
  sessions?: PlanSessionCreate[];
}

export interface PlanSession {
  id: string;
  plan_id: string;
  session_name: string;
  order: number;
  date?: string;
  source_program_id?: string;
  source_session_id?: string;
  entries: PlanSessionEntry[];
}

export interface PlanSessionCreate {
  session_name: string;
  order: number;
  date?: string;
  source_program_id?: string;
  source_session_id?: string;
  entries?: PlanSessionEntryCreate[];
}

export interface PlanSessionEntry {
  id: string;
  session_id: string;
  order: number;
  exercise_name: string;
  fields?: Record<string, unknown>;
  notes?: string;
}

export interface PlanSessionEntryCreate {
  exercise_name: string;
  order: number;
  fields?: Record<string, unknown>;
  notes?: string;
}
```

- [ ] **Step 2: Add plan_snapshot to Log and LogCreate**

In the existing `Log` interface, add after `metadata`:

```typescript
  plan_snapshot?: Record<string, unknown>;
```

In the existing `LogCreate` interface, add after `metadata`:

```typescript
  plan_snapshot?: Record<string, unknown>;
```

- [ ] **Step 3: Verify types compile**

Run: `cd web && pnpm build --no-lint 2>&1 | head -20`
Expected: No type errors related to api.ts

- [ ] **Step 4: Commit**

```bash
git add web/types/api.ts
git commit -m "feat: add Plan types and plan_snapshot to Log types"
```

---

## Task 2: Create Plan API Client

**Files:**
- Create: `web/lib/api/plans.ts`

- [ ] **Step 1: Create the Plan API client**

Create `web/lib/api/plans.ts`:

```typescript
import type { Plan, PlanCreate, PlanSessionCreate, PlanUpdate } from '@/types/api';
import { api } from './client';

export const plansApi = {
  async get(): Promise<Plan> {
    return api.get('plan').json<Plan>();
  },

  async create(data: PlanCreate): Promise<Plan> {
    return api.post('plan', { json: data }).json<Plan>();
  },

  async update(data: PlanUpdate): Promise<Plan> {
    return api.put('plan', { json: data }).json<Plan>();
  },

  async delete(): Promise<void> {
    await api.delete('plan');
  },

  async addSessions(sessions: PlanSessionCreate[]): Promise<Plan> {
    return api.post('plan/sessions', { json: { sessions } }).json<Plan>();
  },

  async deleteSession(sessionId: string): Promise<void> {
    await api.delete(`plan/sessions/${sessionId}`);
  },

  async expandProgram(programId: string): Promise<Plan> {
    return api.post(`plan/expand-program/${programId}`).json<Plan>();
  },
};
```

- [ ] **Step 2: Verify it compiles**

Run: `cd web && pnpm build --no-lint 2>&1 | head -20`
Expected: No type errors

- [ ] **Step 3: Commit**

```bash
git add web/lib/api/plans.ts
git commit -m "feat: add Plan API client"
```

---

## Task 3: Create Plan React Query Hooks

**Files:**
- Create: `web/lib/hooks/usePlans.ts`

- [ ] **Step 1: Create hooks file**

Create `web/lib/hooks/usePlans.ts`:

```typescript
import { plansApi } from '@/lib/api/plans';
import type { PlanCreate, PlanSessionCreate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { HTTPError } from 'ky';
import { toast } from 'sonner';

export function usePlan() {
  return useQuery({
    queryKey: ['plan'],
    queryFn: () => plansApi.get(),
    retry: (failureCount, error) => {
      // Don't retry on 404 (no plan yet)
      if (error instanceof HTTPError && error.response.status === 404) return false;
      return failureCount < 1;
    },
  });
}

export function useCreatePlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: PlanCreate) => plansApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan'] });
    },
  });
}

export function useAddPlanSessions() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (sessions: PlanSessionCreate[]) => plansApi.addSessions(sessions),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan'] });
    },
    onError: async (error) => {
      let message = 'Failed to add sessions to plan';
      if (error instanceof HTTPError) {
        try {
          const body = await error.response.json();
          if (body.message) message = body.message;
        } catch {
          // use default
        }
      }
      toast.error(message);
    },
  });
}

export function useDeletePlanSession() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (sessionId: string) => plansApi.deleteSession(sessionId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan'] });
    },
  });
}

export function useExpandProgram() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (programId: string) => plansApi.expandProgram(programId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan'] });
      toast.success('Sessions added to plan');
    },
    onError: async (error) => {
      let message = 'Failed to expand program';
      if (error instanceof HTTPError) {
        try {
          const body = await error.response.json();
          if (body.message) message = body.message;
        } catch {
          // use default
        }
      }
      toast.error(message);
    },
  });
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd web && pnpm build --no-lint 2>&1 | head -20`
Expected: No type errors

- [ ] **Step 3: Commit**

```bash
git add web/lib/hooks/usePlans.ts
git commit -m "feat: add Plan React Query hooks"
```

---

## Task 4: Update Sidebar Navigation

**Files:**
- Modify: `web/components/layout/Sidebar.tsx`

- [ ] **Step 1: Update nav items to Plan / Logs / Settings**

Replace the `navItems` array and update imports in `web/components/layout/Sidebar.tsx`:

```typescript
import { CalendarCheck, ClipboardList, Settings } from 'lucide-react';
```

```typescript
const navItems = [
  { href: '/', label: 'Plan', icon: CalendarCheck },
  { href: '/logs', label: 'Logs', icon: ClipboardList },
  { href: '/settings', label: 'Settings', icon: Settings },
];
```

- [ ] **Step 2: Verify it compiles**

Run: `cd web && pnpm build --no-lint 2>&1 | head -20`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add web/components/layout/Sidebar.tsx
git commit -m "feat: update sidebar nav to Plan / Logs / Settings"
```

---

## Task 5: Create SessionCard Component

**Files:**
- Create: `web/components/plan/SessionCard.tsx`

- [ ] **Step 1: Create SessionCard with dynamic field table**

Create `web/components/plan/SessionCard.tsx`:

```typescript
'use client';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import type { PlanSession } from '@/types/api';
import { CalendarDays, X } from 'lucide-react';

interface SessionCardProps {
  session: PlanSession;
  programName?: string;
  onLog: (session: PlanSession) => void;
  onDelete: (sessionId: string) => void;
}

/** Collect unique field keys across all entries, preserving insertion order. */
function collectFieldKeys(session: PlanSession): string[] {
  const keys: string[] = [];
  const seen = new Set<string>();
  for (const entry of session.entries) {
    if (!entry.fields) continue;
    for (const key of Object.keys(entry.fields)) {
      if (!seen.has(key)) {
        seen.add(key);
        keys.push(key);
      }
    }
  }
  return keys;
}

function formatFieldValue(value: unknown): string {
  if (value == null) return '—';
  if (typeof value === 'boolean') return value ? 'Yes' : 'No';
  return String(value);
}

export function SessionCard({ session, programName, onLog, onDelete }: SessionCardProps) {
  const sortedEntries = [...session.entries].sort((a, b) => a.order - b.order);
  const fieldKeys = collectFieldKeys(session);

  return (
    <Card
      className="cursor-pointer transition-colors hover:border-primary"
      onClick={() => onLog(session)}
    >
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 min-w-0">
            <span className="font-semibold truncate">{session.session_name}</span>
            {session.source_program_id && programName ? (
              <span className="text-xs text-muted-foreground shrink-0">from {programName}</span>
            ) : !session.source_program_id ? (
              <span className="text-xs text-muted-foreground italic shrink-0">manual</span>
            ) : null}
          </div>
          <div className="flex items-center gap-2 shrink-0">
            {session.date && (
              <span className="flex items-center gap-1 text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded-full">
                <CalendarDays className="h-3 w-3" />
                {session.date}
              </span>
            )}
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6 text-muted-foreground hover:text-destructive"
              onClick={(e) => {
                e.stopPropagation();
                onDelete(session.id);
              }}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        {sortedEntries.length === 0 ? (
          <p className="text-sm text-muted-foreground">No exercises</p>
        ) : (
          <table className="text-sm border-collapse">
            <thead>
              <tr className="text-xs text-muted-foreground">
                <th className="text-left font-medium pb-1 pr-4 w-[150px]">Exercise</th>
                {fieldKeys.map((key) => (
                  <th key={key} className="text-left font-medium pb-1 px-2 w-[60px]">
                    {key}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {sortedEntries.map((entry) => (
                <tr key={entry.id} className="text-muted-foreground">
                  <td className="pr-4 py-0.5">{entry.exercise_name}</td>
                  {fieldKeys.map((key) => (
                    <td key={key} className="px-2 py-0.5 tabular-nums">
                      {formatFieldValue(entry.fields?.[key])}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </CardContent>
    </Card>
  );
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd web && pnpm build --no-lint 2>&1 | head -20`
Expected: No type errors

- [ ] **Step 3: Commit**

```bash
git add web/components/plan/SessionCard.tsx
git commit -m "feat: add SessionCard component with dynamic field table"
```

---

## Task 6: Create ProgramSidebar Component

**Files:**
- Create: `web/components/plan/ProgramSidebar.tsx`

- [ ] **Step 1: Create ProgramSidebar**

Create `web/components/plan/ProgramSidebar.tsx`:

```typescript
'use client';

import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { usePrograms } from '@/lib/hooks/usePrograms';
import { Plus } from 'lucide-react';
import Link from 'next/link';

export function ProgramSidebar() {
  const { data, isLoading } = usePrograms();
  const programs = data?.data ?? [];

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <h3 className="font-semibold text-sm">Programs</h3>
        <Button variant="default" size="sm" className="h-7 text-xs" asChild>
          <Link href="/programs/new">
            <Plus className="h-3 w-3 mr-1" />
            New
          </Link>
        </Button>
      </div>

      {isLoading ? (
        <div className="space-y-2">
          <Skeleton className="h-14 w-full" />
          <Skeleton className="h-14 w-full" />
        </div>
      ) : programs.length === 0 ? (
        <p className="text-sm text-muted-foreground">
          No programs yet. Create one to quickly add sessions to your plan.
        </p>
      ) : (
        <div className="space-y-1.5">
          {programs.map((program) => (
            <Link
              key={program.id}
              href={`/programs/${program.id}`}
              className="block rounded-md border p-2.5 text-sm hover:border-primary transition-colors"
            >
              <div className="font-medium truncate">{program.name}</div>
              <div className="text-xs text-muted-foreground mt-0.5">
                {program.sessions.length} session{program.sessions.length !== 1 ? 's' : ''}
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd web && pnpm build --no-lint 2>&1 | head -20`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add web/components/plan/ProgramSidebar.tsx
git commit -m "feat: add ProgramSidebar component"
```

---

## Task 7: Create EmptyState and AddSessionDialog Components

**Files:**
- Create: `web/components/plan/EmptyState.tsx`
- Create: `web/components/plan/AddSessionDialog.tsx`

- [ ] **Step 1: Create EmptyState**

Create `web/components/plan/EmptyState.tsx`:

```typescript
import { Button } from '@/components/ui/button';
import { Plus } from 'lucide-react';

interface EmptyStateProps {
  onAddSession: () => void;
}

export function EmptyState({ onAddSession }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-16">
      <p className="text-lg font-medium mb-1">No sessions planned yet</p>
      <p className="text-sm text-muted-foreground mb-6">Add a session to get started</p>
      <Button onClick={onAddSession}>
        <Plus className="h-4 w-4 mr-2" />
        Add Session
      </Button>
    </div>
  );
}
```

- [ ] **Step 2: Create AddSessionDialog**

Create `web/components/plan/AddSessionDialog.tsx`:

```typescript
'use client';

import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { PlanSessionCreate, PlanSessionEntryCreate } from '@/types/api';
import { Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';

interface AddSessionDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onAdd: (session: PlanSessionCreate) => void;
  nextOrder: number;
}

interface EntryDraft {
  id: string;
  exercise_name: string;
  fields: Record<string, string>;
}

function createEmptyEntry(): EntryDraft {
  return { id: crypto.randomUUID(), exercise_name: '', fields: {} };
}

export function AddSessionDialog({ open, onOpenChange, onAdd, nextOrder }: AddSessionDialogProps) {
  const [sessionName, setSessionName] = useState('');
  const [date, setDate] = useState('');
  const [entries, setEntries] = useState<EntryDraft[]>([createEmptyEntry()]);

  const handleAdd = () => {
    if (!sessionName.trim()) return;

    const validEntries = entries.filter((e) => e.exercise_name.trim() !== '');
    const planEntries: PlanSessionEntryCreate[] = validEntries.map((e, i) => ({
      exercise_name: e.exercise_name.trim(),
      order: i,
      fields:
        Object.keys(e.fields).length > 0
          ? Object.fromEntries(
              Object.entries(e.fields)
                .filter(([, v]) => v.trim() !== '')
                .map(([k, v]) => {
                  const num = Number(v);
                  return [k, Number.isNaN(num) ? v : num];
                })
            )
          : undefined,
    }));

    onAdd({
      session_name: sessionName.trim(),
      order: nextOrder,
      date: date || undefined,
      entries: planEntries.length > 0 ? planEntries : undefined,
    });

    // Reset
    setSessionName('');
    setDate('');
    setEntries([createEmptyEntry()]);
    onOpenChange(false);
  };

  const updateEntry = (id: string, field: string, value: string) => {
    setEntries((prev) =>
      prev.map((e) => (e.id === id ? { ...e, [field]: value } : e))
    );
  };

  const updateEntryField = (id: string, fieldName: string, value: string) => {
    setEntries((prev) =>
      prev.map((e) =>
        e.id === id ? { ...e, fields: { ...e.fields, [fieldName]: value } } : e
      )
    );
  };

  const removeEntry = (id: string) => {
    setEntries((prev) => prev.filter((e) => e.id !== id));
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Add Session</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="session-name">Session Name</Label>
            <Input
              id="session-name"
              placeholder="e.g., Upper Body A"
              value={sessionName}
              onChange={(e) => setSessionName(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="session-date">Date (optional)</Label>
            <Input
              id="session-date"
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label>Exercises</Label>
            {entries.map((entry) => (
              <div key={entry.id} className="flex items-center gap-2">
                <Input
                  placeholder="Exercise name"
                  value={entry.exercise_name}
                  onChange={(e) => updateEntry(entry.id, 'exercise_name', e.target.value)}
                  className="flex-1"
                />
                <Input
                  placeholder="load_kg"
                  value={entry.fields.load_kg ?? ''}
                  onChange={(e) => updateEntryField(entry.id, 'load_kg', e.target.value)}
                  className="w-20"
                />
                <Input
                  placeholder="sets"
                  value={entry.fields.sets ?? ''}
                  onChange={(e) => updateEntryField(entry.id, 'sets', e.target.value)}
                  className="w-16"
                />
                <Input
                  placeholder="reps"
                  value={entry.fields.reps ?? ''}
                  onChange={(e) => updateEntryField(entry.id, 'reps', e.target.value)}
                  className="w-16"
                />
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-8 w-8 shrink-0"
                  onClick={() => removeEntry(entry.id)}
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </Button>
              </div>
            ))}
            <Button
              variant="outline"
              size="sm"
              onClick={() => setEntries((prev) => [...prev, createEmptyEntry()])}
            >
              <Plus className="h-3.5 w-3.5 mr-1" />
              Add Exercise
            </Button>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleAdd} disabled={!sessionName.trim()}>
            Add to Plan
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
```

- [ ] **Step 3: Verify both compile**

Run: `cd web && pnpm build --no-lint 2>&1 | head -20`
Expected: No type errors

- [ ] **Step 4: Commit**

```bash
git add web/components/plan/EmptyState.tsx web/components/plan/AddSessionDialog.tsx
git commit -m "feat: add EmptyState and AddSessionDialog components"
```

---

## Task 8: Build the Plan Page

**Files:**
- Modify: `web/app/page.tsx`

- [ ] **Step 1: Rewrite page.tsx as the Plan page**

Replace the entire contents of `web/app/page.tsx`:

```typescript
'use client';

import { AddSessionDialog } from '@/components/plan/AddSessionDialog';
import { EmptyState } from '@/components/plan/EmptyState';
import { ProgramSidebar } from '@/components/plan/ProgramSidebar';
import { SessionCard } from '@/components/plan/SessionCard';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { useAddPlanSessions, useDeletePlanSession, usePlan } from '@/lib/hooks/usePlans';
import { usePrograms } from '@/lib/hooks/usePrograms';
import type { PlanSession, PlanSessionCreate } from '@/types/api';
import { HTTPError } from 'ky';
import { Plus } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function PlanPage() {
  const router = useRouter();
  const { data: plan, isLoading, error } = usePlan();
  const { data: programsData } = usePrograms();
  const addSessions = useAddPlanSessions();
  const deleteSession = useDeletePlanSession();
  const [addDialogOpen, setAddDialogOpen] = useState(false);

  const noPlan = error instanceof HTTPError && error.response.status === 404;
  const sessions = plan?.sessions?.slice().sort((a, b) => a.order - b.order) ?? [];
  const programMap = new Map((programsData?.data ?? []).map((p) => [p.id, p.name]));

  const handleLog = (session: PlanSession) => {
    const params = new URLSearchParams();
    params.set('planSessionId', session.id);
    if (session.source_program_id) params.set('programId', session.source_program_id);
    params.set('session', session.session_name);
    router.push(`/logs/new?${params}`);
  };

  const handleDelete = async (sessionId: string) => {
    await deleteSession.mutateAsync(sessionId);
  };

  const handleAddSession = async (session: PlanSessionCreate) => {
    await addSessions.mutateAsync([session]);
  };

  if (isLoading) {
    return (
      <main className="flex flex-1">
        <div className="flex-[7] p-6 space-y-3">
          <Skeleton className="h-8 w-[100px]" />
          <Skeleton className="h-[120px] w-full" />
          <Skeleton className="h-[120px] w-full" />
        </div>
        <div className="flex-[3] p-4 border-l space-y-2">
          <Skeleton className="h-6 w-[80px]" />
          <Skeleton className="h-14 w-full" />
          <Skeleton className="h-14 w-full" />
        </div>
      </main>
    );
  }

  const nextOrder = sessions.length > 0 ? Math.max(...sessions.map((s) => s.order)) + 1 : 0;

  return (
    <main className="flex flex-1">
      {/* Sessions (70%) */}
      <div className="flex-[7] p-6">
        <div className="mb-4">
          <h1 className="text-2xl font-bold">Plan</h1>
          {sessions.length > 0 && (
            <p className="text-sm text-muted-foreground">
              {sessions.length} session{sessions.length !== 1 ? 's' : ''} queued
            </p>
          )}
        </div>

        {noPlan || sessions.length === 0 ? (
          <EmptyState onAddSession={() => setAddDialogOpen(true)} />
        ) : (
          <div className="space-y-3">
            {sessions.map((session) => (
              <SessionCard
                key={session.id}
                session={session}
                programName={
                  session.source_program_id
                    ? programMap.get(session.source_program_id)
                    : undefined
                }
                onLog={handleLog}
                onDelete={handleDelete}
              />
            ))}
          </div>
        )}

        {sessions.length > 0 && (
          <Button
            variant="outline"
            className="w-full mt-3 border-dashed"
            onClick={() => setAddDialogOpen(true)}
          >
            <Plus className="h-4 w-4 mr-2" />
            Add Session
          </Button>
        )}

        <AddSessionDialog
          open={addDialogOpen}
          onOpenChange={setAddDialogOpen}
          onAdd={handleAddSession}
          nextOrder={nextOrder}
        />
      </div>

      {/* Programs Sidebar (30%) */}
      <div className="flex-[3] p-4 border-l">
        <ProgramSidebar />
      </div>
    </main>
  );
}
```

- [ ] **Step 2: Remove old dashboard components**

```bash
rm web/components/dashboard/ActivePlans.tsx web/components/dashboard/RecentLogs.tsx
rmdir web/components/dashboard 2>/dev/null || true
```

- [ ] **Step 3: Verify it compiles**

Run: `cd web && pnpm build --no-lint 2>&1 | head -20`
Expected: No errors (dashboard components no longer imported)

- [ ] **Step 4: Commit**

```bash
git add web/app/page.tsx
git add -u web/components/dashboard/
git commit -m "feat: replace Dashboard with Plan page"
```

---

## Task 9: Create SessionSelector and Add "Add to Plan" to Program Detail

**Files:**
- Create: `web/components/programs/SessionSelector.tsx`
- Modify: `web/app/programs/[id]/page.tsx`

- [ ] **Step 1: Create SessionSelector component**

Create `web/components/programs/SessionSelector.tsx`:

```typescript
'use client';

import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import type { PlanSessionCreate, PlanSessionEntryCreate, ProgramSession } from '@/types/api';
import { useState } from 'react';

interface SessionSelectorProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  sessions: ProgramSession[];
  programId: string;
  onConfirm: (sessions: PlanSessionCreate[]) => void;
  isPending: boolean;
}

export function SessionSelector({
  open,
  onOpenChange,
  sessions,
  programId,
  onConfirm,
  isPending,
}: SessionSelectorProps) {
  const sortedSessions = [...sessions].sort((a, b) => a.order - b.order);
  const [selected, setSelected] = useState<Set<string>>(
    () => new Set(sortedSessions.map((s) => s.id))
  );

  const toggleSession = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handleConfirm = () => {
    const planSessions: PlanSessionCreate[] = sortedSessions
      .filter((s) => selected.has(s.id))
      .map((s, i) => ({
        session_name: s.session_name,
        order: i,
        date: s.date || undefined,
        source_program_id: programId,
        source_session_id: s.id,
        entries: s.entries.map(
          (e, j): PlanSessionEntryCreate => ({
            exercise_name: e.exercise_name,
            order: j,
            fields: e.fields || undefined,
            notes: e.notes || undefined,
          })
        ),
      }));

    onConfirm(planSessions);
  };

  // Reset selection when dialog opens
  const handleOpenChange = (isOpen: boolean) => {
    if (isOpen) {
      setSelected(new Set(sortedSessions.map((s) => s.id)));
    }
    onOpenChange(isOpen);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Add to Plan</DialogTitle>
        </DialogHeader>

        <div className="border rounded-md divide-y">
          {sortedSessions.map((session) => (
            <label
              key={session.id}
              className="flex items-center gap-3 px-3 py-2.5 cursor-pointer hover:bg-muted/50"
            >
              <Checkbox
                checked={selected.has(session.id)}
                onCheckedChange={() => toggleSession(session.id)}
              />
              <span className="font-medium text-sm flex-1">{session.session_name}</span>
              <span className="text-xs text-muted-foreground">
                {session.entries.length} exercise{session.entries.length !== 1 ? 's' : ''}
              </span>
            </label>
          ))}
        </div>

        <p className="text-xs text-muted-foreground">
          {selected.size} of {sortedSessions.length} sessions selected
        </p>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleConfirm} disabled={selected.size === 0 || isPending}>
            {isPending ? 'Adding...' : 'Add to Plan'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
```

- [ ] **Step 2: Add "Add to Plan" button to Program detail page**

In `web/app/programs/[id]/page.tsx`, add the imports at the top:

```typescript
import { SessionSelector } from '@/components/programs/SessionSelector';
import { useAddPlanSessions } from '@/lib/hooks/usePlans';
import type { PlanSessionCreate } from '@/types/api';
```

Inside the `ProgramDetailPage` component, add after the existing hooks:

```typescript
  const addPlanSessions = useAddPlanSessions();
  const [selectorOpen, setSelectorOpen] = useState(false);
```

Add the handler:

```typescript
  const handleAddToPlan = async (sessions: PlanSessionCreate[]) => {
    await addPlanSessions.mutateAsync(sessions);
    setSelectorOpen(false);
    router.push('/');
  };
```

In the header button group, add the "Add to Plan" button (before the Edit button):

```typescript
            <Button onClick={() => setSelectorOpen(true)}>
              Add to Plan
            </Button>
```

Add the SessionSelector component at the end of the JSX (before `</main>`):

```typescript
      <SessionSelector
        open={selectorOpen}
        onOpenChange={setSelectorOpen}
        sessions={program.sessions}
        programId={programId}
        onConfirm={handleAddToPlan}
        isPending={addPlanSessions.isPending}
      />
```

Change the delete redirect from `'/programs'` to `'/'`:

```typescript
  const handleDelete = async () => {
    await deleteProgram.mutateAsync(programId);
    router.push('/');
  };
```

- [ ] **Step 3: Verify it compiles**

Run: `cd web && pnpm build --no-lint 2>&1 | head -20`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add web/components/programs/SessionSelector.tsx web/app/programs/\[id\]/page.tsx
git commit -m "feat: add 'Add to Plan' session selector to Program detail"
```

---

## Task 10: Update Log Creation to Support plan_snapshot

**Files:**
- Modify: `web/app/logs/new/page.tsx`

- [ ] **Step 1: Update NewLogPage to handle planSessionId**

Replace the entire contents of `web/app/logs/new/page.tsx`:

```typescript
'use client';

import { LogEntryTable } from '@/components/log-entry-table/LogEntryTable';
import { Skeleton } from '@/components/ui/skeleton';
import { useCreateLog } from '@/lib/hooks/useLogs';
import { useDeletePlanSession, usePlan } from '@/lib/hooks/usePlans';
import { useProgram } from '@/lib/hooks/usePrograms';
import type { LogEntryCreate } from '@/types/api';
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { Suspense } from 'react';

function NewLogPageContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const planSessionId = searchParams.get('planSessionId');
  const programId = searchParams.get('programId');
  const sessionName = searchParams.get('session');

  const { data: plan } = usePlan();
  const { data: program, isLoading: programLoading } = useProgram(programId);
  const createLog = useCreateLog();
  const deletePlanSession = useDeletePlanSession();

  // Find plan session if navigating from Plan page
  const planSession = planSessionId
    ? plan?.sessions.find((s) => s.id === planSessionId)
    : undefined;

  // Fall back to program session if no plan session
  const programSession = program?.sessions.find((s) => s.session_name === sessionName);

  // Use plan session entries, or fall back to program session entries
  const templateEntries = planSession
    ? planSession.entries.slice().sort((a, b) => a.order - b.order)
    : programSession?.entries?.slice().sort((a, b) => a.order - b.order);

  const displayName = planSession?.session_name ?? sessionName;
  const displayProgramName = program?.name;
  const backHref = planSessionId ? '/' : programId ? `/programs/${programId}` : '/logs';

  const handleSave = async (data: {
    entries: LogEntryCreate[];
    notes: string;
    startedAt?: string;
    finishedAt?: string;
  }) => {
    // Build plan_snapshot from plan session if available
    const planSnapshot = planSession
      ? {
          session_name: planSession.session_name,
          entries: planSession.entries.map((e) => ({
            exercise_name: e.exercise_name,
            order: e.order,
            fields: e.fields,
            notes: e.notes,
          })),
        }
      : undefined;

    await createLog.mutateAsync({
      program_id: planSession?.source_program_id ?? programId ?? undefined,
      session_name: displayName ?? undefined,
      performed_at: new Date().toISOString(),
      started_at: data.startedAt,
      finished_at: data.finishedAt,
      notes: data.notes || undefined,
      plan_snapshot: planSnapshot,
      entries: data.entries,
    });

    // Auto-remove session from plan queue
    if (planSessionId) {
      await deletePlanSession.mutateAsync(planSessionId);
    }

    router.push(backHref);
  };

  if (programId && programLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-8 w-[200px]" />
        <Skeleton className="h-[400px] w-full" />
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <Link
        href={backHref}
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-4"
      >
        <ArrowLeft className="h-4 w-4" />
        Back
      </Link>

      <div className="mb-6">
        <h1 className="text-3xl font-bold">Record Log</h1>
        {(displayProgramName || displayName) && (
          <p className="text-muted-foreground mt-1">
            {displayProgramName}
            {displayName && <span> — {displayName}</span>}
          </p>
        )}
      </div>

      <LogEntryTable
        initialEntries={templateEntries}
        onSave={handleSave}
        onCancel={() => router.push(backHref)}
      />
    </main>
  );
}

function NewLogPageFallback() {
  return (
    <main className="container mx-auto p-6 space-y-4">
      <Skeleton className="h-8 w-[200px]" />
      <Skeleton className="h-[400px] w-full" />
    </main>
  );
}

export default function NewLogPage() {
  return (
    <Suspense fallback={<NewLogPageFallback />}>
      <NewLogPageContent />
    </Suspense>
  );
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd web && pnpm build --no-lint 2>&1 | head -20`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add web/app/logs/new/page.tsx
git commit -m "feat: support plan_snapshot and auto-remove in log creation"
```

---

## Task 11: Create LogCard and Redesign Logs Page

**Files:**
- Create: `web/components/logs/LogCard.tsx`
- Modify: `web/app/logs/page.tsx`

- [ ] **Step 1: Create LogCard component**

Create `web/components/logs/LogCard.tsx`:

```typescript
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import type { Log, LogEntry } from '@/types/api';
import Link from 'next/link';

interface LogCardProps {
  log: Log;
  programName?: string;
}

interface SnapshotEntry {
  exercise_name: string;
  fields?: Record<string, unknown>;
}

function parseSnapshot(snapshot: Record<string, unknown> | undefined): SnapshotEntry[] {
  if (!snapshot?.entries || !Array.isArray(snapshot.entries)) return [];
  return snapshot.entries as SnapshotEntry[];
}

/** Collect unique field keys across log entries (and snapshot if present). */
function collectFieldKeys(entries: LogEntry[], snapshotEntries: SnapshotEntry[]): string[] {
  const keys: string[] = [];
  const seen = new Set<string>();
  const addKeys = (fields?: Record<string, unknown>) => {
    if (!fields) return;
    for (const key of Object.keys(fields)) {
      if (!seen.has(key)) {
        seen.add(key);
        keys.push(key);
      }
    }
  };
  for (const e of snapshotEntries) addKeys(e.fields);
  for (const e of entries) addKeys(e.fields);
  return keys;
}

function formatValue(value: unknown): string {
  if (value == null) return '—';
  if (typeof value === 'boolean') return value ? 'Yes' : 'No';
  return String(value);
}

function DiffCell({
  planValue,
  actualValue,
}: {
  planValue: unknown;
  actualValue: unknown;
}) {
  const pStr = formatValue(planValue);
  const aStr = formatValue(actualValue);

  if (pStr === aStr) return <span className="tabular-nums">{aStr}</span>;

  const pNum = typeof planValue === 'number' ? planValue : null;
  const aNum = typeof actualValue === 'number' ? actualValue : null;
  const color =
    pNum != null && aNum != null
      ? aNum > pNum
        ? 'text-green-600'
        : 'text-red-500'
      : '';

  return (
    <span className="tabular-nums">
      <span className="line-through text-muted-foreground mr-1">{pStr}</span>
      <span className={color}>{aStr}</span>
    </span>
  );
}

export function LogCard({ log, programName }: LogCardProps) {
  const performedDate = new Date(log.performed_at).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });

  const snapshotEntries = parseSnapshot(log.plan_snapshot);
  const hasSnapshot = snapshotEntries.length > 0;
  const sortedEntries = [...log.entries].sort((a, b) => a.order - b.order);
  const fieldKeys = collectFieldKeys(sortedEntries, snapshotEntries);

  // Build snapshot lookup by exercise name for diff
  const snapshotByName = new Map(snapshotEntries.map((e) => [e.exercise_name, e]));

  return (
    <Link href={`/logs/${log.id}`} className="block">
      <Card className="cursor-pointer transition-colors hover:border-primary">
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 min-w-0">
              <span className="font-semibold truncate">
                {log.session_name ?? 'Untitled Session'}
              </span>
              {log.program_id && programName ? (
                <span className="text-xs text-muted-foreground shrink-0">from {programName}</span>
              ) : !log.program_id ? (
                <span className="text-xs text-muted-foreground italic shrink-0">ad-hoc</span>
              ) : null}
            </div>
            <span className="text-sm text-muted-foreground shrink-0">{performedDate}</span>
          </div>
        </CardHeader>
        <CardContent className="pt-0">
          {sortedEntries.length === 0 ? (
            <p className="text-sm text-muted-foreground">No entries</p>
          ) : (
            <table className="text-sm border-collapse">
              <thead>
                <tr className="text-xs text-muted-foreground">
                  <th className="text-left font-medium pb-1 pr-4 w-[150px]">Exercise</th>
                  {fieldKeys.map((key) => (
                    <th key={key} className="text-left font-medium pb-1 px-2 w-[60px]">
                      {key}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {sortedEntries.map((entry) => {
                  const snapshot = snapshotByName.get(entry.exercise_name);
                  return (
                    <tr key={entry.id} className="text-muted-foreground">
                      <td className="pr-4 py-0.5">{entry.exercise_name}</td>
                      {fieldKeys.map((key) => (
                        <td key={key} className="px-2 py-0.5">
                          {hasSnapshot && snapshot ? (
                            <DiffCell
                              planValue={snapshot.fields?.[key]}
                              actualValue={entry.fields?.[key]}
                            />
                          ) : (
                            <span className="tabular-nums">
                              {formatValue(entry.fields?.[key])}
                            </span>
                          )}
                        </td>
                      ))}
                    </tr>
                  );
                })}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>
    </Link>
  );
}
```

- [ ] **Step 2: Rewrite Logs page with card layout**

Replace the entire contents of `web/app/logs/page.tsx`:

```typescript
'use client';

import { ExportButton } from '@/components/export/ExportButton';
import { LogCard } from '@/components/logs/LogCard';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { useLogs } from '@/lib/hooks/useLogs';
import { usePrograms } from '@/lib/hooks/usePrograms';
import { Plus } from 'lucide-react';
import Link from 'next/link';

export default function LogsPage() {
  const { data: logsData, isLoading: logsLoading, error: logsError } = useLogs();
  const { data: programsData } = usePrograms();

  const logs = logsData?.data || [];
  const sortedLogs = [...logs].sort(
    (a, b) => new Date(b.performed_at).getTime() - new Date(a.performed_at).getTime()
  );

  const programMap = new Map<string, string>((programsData?.data ?? []).map((p) => [p.id, p.name]));

  if (logsLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <div className="space-y-3">
          <Skeleton className="h-[120px] w-full" />
          <Skeleton className="h-[120px] w-full" />
          <Skeleton className="h-[120px] w-full" />
        </div>
      </main>
    );
  }

  if (logsError) {
    return (
      <main className="container mx-auto p-6">
        <p className="text-destructive">Failed to load logs. Please try again later.</p>
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Logs</h1>
          <p className="text-muted-foreground mt-1">View and record your sessions</p>
        </div>
        <div className="flex items-center gap-2">
          <ExportButton logs={logs} plan={null} />
          <Button asChild>
            <Link href="/logs/new">
              <Plus className="h-4 w-4 mr-2" />
              Record Log
            </Link>
          </Button>
        </div>
      </div>

      {sortedLogs.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">
            No logs yet. Record your first training session.
          </p>
          <Button asChild>
            <Link href="/logs/new">
              <Plus className="h-4 w-4 mr-2" />
              Record Log
            </Link>
          </Button>
        </div>
      ) : (
        <div className="space-y-3">
          {sortedLogs.map((log) => (
            <LogCard
              key={log.id}
              log={log}
              programName={log.program_id ? programMap.get(log.program_id) : undefined}
            />
          ))}
        </div>
      )}
    </main>
  );
}
```

- [ ] **Step 3: Remove old LogTable component**

```bash
rm web/components/logs/LogTable.tsx
```

- [ ] **Step 4: Verify it compiles**

Run: `cd web && pnpm build --no-lint 2>&1 | head -20`
Expected: No errors

- [ ] **Step 5: Commit**

```bash
git add web/components/logs/LogCard.tsx web/app/logs/page.tsx
git add -u web/components/logs/LogTable.tsx
git commit -m "feat: redesign Logs page with card layout and plan vs actual diff"
```

---

## Task 12: Remove Programs List Page and Clean Up

**Files:**
- Remove: `web/app/programs/page.tsx`

- [ ] **Step 1: Remove programs list page**

```bash
rm web/app/programs/page.tsx
```

- [ ] **Step 2: Run lint/format check**

Run: `cd web && pnpm check 2>&1 | tail -30`

If there are lint/format issues, fix them:

Run: `cd web && pnpm check:fix`

- [ ] **Step 3: Verify full build**

Run: `cd web && pnpm build 2>&1 | tail -20`
Expected: Build succeeds with no errors

- [ ] **Step 4: Commit**

```bash
git add -u web/app/programs/page.tsx
git add -A web/
git commit -m "chore: remove programs list page and fix lint issues"
```

---

## Task 13: Final Verification

- [ ] **Step 1: Run full check from web/**

Run: `cd web && pnpm check`
Expected: All checks pass

- [ ] **Step 2: Run full build**

Run: `cd web && pnpm build`
Expected: Build succeeds

- [ ] **Step 3: Verify no broken imports**

Run: `cd web && grep -r "ActivePlans\|RecentLogs\|LogTable\|/programs'" --include='*.tsx' --include='*.ts' .`
Expected: No results (all old references removed). Note: `/programs/` with path continuation (like `/programs/${id}`) is fine.

- [ ] **Step 4: Verify Plan page renders (manual)**

Start dev server and confirm:
1. `/` shows Plan page with split layout
2. Sidebar shows Plan / Logs / Settings
3. Empty state shows "No sessions planned yet" with Add Session button
4. `/logs` shows card-based layout
5. `/programs/:id` has "Add to Plan" button
