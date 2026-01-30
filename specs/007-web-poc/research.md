# Research: Web PoC

**Date**: 2026-01-30  
**Feature**: 007-web-poc

## Technology Stack Decisions

### 1. State Management: TanStack Query

**Decision**: Use TanStack Query (React Query) for server state management

**Rationale**:
- Designed for API-centric applications (matches our REST API consumption pattern)
- Automatic caching, background refetching, and stale-while-revalidate
- Built-in loading/error states reduce boilerplate
- Documented in `docs/FRONTEND_ARCHITECTURE.md` as the standard choice

**Alternatives considered**:
- Zustand: Good for client state, but TanStack Query is better for server state
- Redux Toolkit Query: More complex setup, smaller community for this use case
- SWR: Similar to TanStack Query but TanStack has better TypeScript support

### 2. UI Components: shadcn/ui

**Decision**: Use shadcn/ui for base components

**Rationale**:
- Copy-paste model allows full customization without fighting the library
- Built on Radix UI primitives (accessibility built-in)
- Tailwind CSS integration matches existing setup
- Components needed: Accordion, Dialog, Form, Input, Select, Table, Button, AlertDialog, Calendar, Badge, Combobox

**Alternatives considered**:
- Chakra UI: Runtime CSS-in-JS, heavier bundle
- MUI: Opinionated styling, harder to customize
- Headless UI: Good but shadcn/ui provides better starting styles

### 3. HTTP Client: ky

**Decision**: Use ky for HTTP requests

**Rationale**:
- Lightweight wrapper around fetch
- Built-in retry, timeout, and JSON parsing
- Simple API for creating reusable instances with base URL and headers
- Documented in `docs/FRONTEND_ARCHITECTURE.md`

**Alternatives considered**:
- axios: Larger bundle, more features than needed
- Native fetch: Requires more boilerplate for auth headers, error handling

### 4. Form Handling: React Hook Form + Zod

**Decision**: Use React Hook Form with Zod validation

**Rationale**:
- Uncontrolled forms for better performance
- Zod schemas provide runtime validation and TypeScript inference
- Integration via `@hookform/resolvers/zod`
- Documented in `docs/FRONTEND_ARCHITECTURE.md`

**Alternatives considered**:
- Formik: More verbose, controlled by default
- Native forms: Too much boilerplate for complex forms like Program Editor

### 5. Date Handling: date-fns

**Decision**: Use date-fns for date manipulation

**Rationale**:
- Tree-shakeable (only import what you use)
- Functional API (no mutation)
- Needed for schedule generation (skip weekends, avoid consecutive days)

**Alternatives considered**:
- dayjs: Good alternative, date-fns has slightly better tree-shaking
- Temporal API: Not yet stable in all browsers

## API Integration Patterns

### Authentication Flow

1. User navigates to app → app loads without blocking
2. User performs action requiring API (list workouts, save program)
3. If no token in localStorage → show token input modal
4. User enters token → store in localStorage
5. Retry original request with token
6. If 401 response → clear token, show input modal again

### API Client Structure

```typescript
// lib/api/client.ts
import ky from 'ky';

const getToken = () => localStorage.getItem('optel_token');

export const api = ky.create({
  prefixUrl: 'http://localhost:8080/api/v1',
  hooks: {
    beforeRequest: [
      (request) => {
        const token = getToken();
        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`);
        }
      },
    ],
  },
});
```

### TanStack Query Patterns

```typescript
// lib/hooks/usePrograms.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { programsApi } from '@/lib/api/programs';

export const usePrograms = () => {
  return useQuery({
    queryKey: ['programs'],
    queryFn: () => programsApi.list(),
  });
};

export const useCreateProgram = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: programsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
    },
  });
};
```

## UI Implementation Patterns

### Diff Calculation Logic

Diff is calculated client-side by comparing Plan (from ProgramNode) vs Actual (from WorkoutEntry):

```typescript
// lib/utils/diff.ts
type DiffResult = {
  status: 'match' | 'diff' | 'pending';
  differences: string[];
};

function calculateDiff(plan: PlanSnapshot | null, actual: WorkoutEntry | null): DiffResult {
  if (!plan) return { status: 'match', differences: [] }; // unplanned
  if (!actual) return { status: 'pending', differences: [] };
  
  const diffs: string[] = [];
  if (plan.target_sets !== actual.sets) {
    diffs.push(`Sets ${actual.sets - plan.target_sets > 0 ? '+' : ''}${actual.sets - plan.target_sets}`);
  }
  if (plan.target_reps !== actual.reps) {
    diffs.push(`Reps ${actual.reps - plan.target_reps > 0 ? '+' : ''}${actual.reps - plan.target_reps}`);
  }
  if (plan.target_load_kg !== actual.load_kg) {
    diffs.push(`Load ${actual.load_kg - plan.target_load_kg > 0 ? '+' : ''}${actual.load_kg - plan.target_load_kg}kg`);
  }
  // Note: RPE differences are NOT considered as diff per spec
  
  return {
    status: diffs.length > 0 ? 'diff' : 'match',
    differences: diffs,
  };
}
```

### Schedule Generation Logic

```typescript
// lib/utils/schedule.ts
import { addDays, isWeekend } from 'date-fns';

interface ScheduleOptions {
  startDate: Date;
  skipWeekends: boolean;
  avoidConsecutive: boolean;
}

function generateSchedule(
  days: number, // total program days
  options: ScheduleOptions
): Date[] {
  const schedule: Date[] = [];
  let currentDate = options.startDate;
  
  for (let i = 0; i < days; i++) {
    // Skip weekends if enabled
    while (options.skipWeekends && isWeekend(currentDate)) {
      currentDate = addDays(currentDate, 1);
    }
    
    schedule.push(currentDate);
    
    // Move to next day
    currentDate = addDays(currentDate, options.avoidConsecutive ? 2 : 1);
  }
  
  return schedule;
}
```

### CSV Export Format

```typescript
// lib/utils/export.ts
function exportToCSV(data: ExportRow[]): void {
  const BOM = '\uFEFF'; // UTF-8 BOM for Excel
  const headers = [
    'date', 'week', 'day', 'day_name', 'exercise', 'entry_type',
    'plan_sets', 'plan_reps', 'plan_load_kg',
    'actual_sets', 'actual_reps', 'actual_load_kg', 'actual_rpe',
    'diff', 'notes'
  ];
  
  const csv = BOM + [
    headers.join(','),
    ...data.map(row => headers.map(h => escapeCSV(row[h])).join(','))
  ].join('\n');
  
  downloadBlob(csv, 'text/csv;charset=utf-8', `${programName}_${week}.csv`);
}
```

## Dependencies to Add

```json
{
  "dependencies": {
    "@tanstack/react-query": "^5.x",
    "@hookform/resolvers": "^3.x",
    "react-hook-form": "^7.x",
    "zod": "^3.x",
    "ky": "^1.x",
    "date-fns": "^3.x"
  }
}
```

shadcn/ui components will be added via CLI:
```bash
npx shadcn@latest add accordion dialog form input select table button alert-dialog calendar badge command
```

## Open Questions Resolved

| Question | Resolution |
|----------|------------|
| Token storage | localStorage with key `optel_token` |
| API base URL | `http://localhost:8080/api/v1` (configurable via env) |
| Diff calculation | Client-side, Sets/Reps/Load only (not RPE) |
| Schedule storage | Client-side only for PoC (not persisted to API) |
