# Rx Frontend Standards - Detailed Reference

## Web Directory Structure

```
web/
├── app/                        # Next.js App Router (routing)
│   ├── layout.tsx              # Root layout
│   ├── page.tsx                # Home page
│   ├── workouts/
│   │   ├── page.tsx            # Workout list
│   │   └── [id]/
│   │       └── page.tsx        # Workout detail
│   ├── programs/
│   │   ├── page.tsx            # Program list
│   │   └── [id]/
│   │       └── page.tsx        # Program detail/edit
│   ├── exercises/
│   │   └── page.tsx            # Exercise list
│   └── telemetry/
│       └── page.tsx            # Telemetry dashboard
├── components/
│   ├── ui/                     # shadcn/ui components
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   └── ...
│   └── features/               # Feature-specific components
│       ├── workout/
│       │   ├── workout-list.tsx
│       │   ├── workout-form.tsx
│       │   └── workout-entry-card.tsx
│       ├── program/
│       │   ├── program-list.tsx
│       │   ├── program-form.tsx
│       │   └── program-node-tree.tsx
│       └── exercise/
│           ├── exercise-list.tsx
│           └── exercise-form.tsx
├── lib/
│   ├── api/                    # API client
│   │   ├── client.ts           # Base HTTP client
│   │   ├── workouts.ts         # Workout API functions
│   │   ├── programs.ts         # Program API functions
│   │   ├── exercises.ts        # Exercise API functions
│   │   └── telemetry.ts        # Telemetry API functions
│   ├── hooks/                  # Custom React hooks
│   │   ├── use-workouts.ts
│   │   ├── use-programs.ts
│   │   └── use-exercises.ts
│   └── utils/                  # Utility functions
│       ├── date.ts             # Date formatting (date-fns)
│       └── validation.ts       # Common validation helpers
├── types/                      # TypeScript type definitions
│   ├── workout.ts
│   ├── program.ts
│   ├── exercise.ts
│   └── api.ts                  # API response types
└── schemas/                    # Zod validation schemas
    ├── workout.ts
    ├── program.ts
    └── exercise.ts
```

## Component Separation Rules

| Type | Location | Role | Example |
|------|----------|------|---------|
| **Page** | `app/**/page.tsx` | Data fetching, layout composition | `app/workouts/page.tsx` |
| **Feature Component** | `components/features/` | Feature-specific UI + logic | `WorkoutList`, `WorkoutForm` |
| **UI Component** | `components/ui/` | Reusable presentational | `Button`, `Card`, `Input` |

### Page Component Example

```tsx
// app/workouts/page.tsx
import { WorkoutList } from '@/components/features/workout/workout-list';

export default function WorkoutsPage() {
  return (
    <div className="container py-8">
      <h1 className="text-2xl font-bold mb-6">Workouts</h1>
      <WorkoutList />
    </div>
  );
}
```

### Feature Component Example

```tsx
// components/features/workout/workout-list.tsx
'use client';

import { useWorkouts } from '@/lib/hooks/use-workouts';
import { WorkoutCard } from './workout-card';
import { Button } from '@/components/ui/button';

export function WorkoutList() {
  const { data, isLoading, error } = useWorkouts();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div className="space-y-4">
      {data?.items.map((workout) => (
        <WorkoutCard key={workout.id} workout={workout} />
      ))}
    </div>
  );
}
```

## Naming Conventions

### Files

```
workout-list.tsx        # Components: kebab-case
use-workouts.ts         # Hooks: use-*.ts
workout.ts              # Types/schemas: singular, kebab-case
```

### Components and Variables

```typescript
// Component names: PascalCase
export function WorkoutList() { ... }
export function WorkoutEntryCard() { ... }

// Props types: ComponentNameProps
type WorkoutListProps = { ... }
type WorkoutEntryCardProps = { ... }

// Variables and functions: camelCase
const workoutList = [];
const handleSubmit = () => {};

// Event handlers: handle + Action
const handleClick = () => {};
const handleSubmit = () => {};
const handleWorkoutCreate = () => {};

// Boolean variables: is/has/should prefix
const isLoading = true;
const hasError = false;
const shouldRefresh = true;
```

### API Functions

```typescript
// lib/api/workouts.ts
export const workoutsApi = {
  list: (params?: WorkoutListParams) => ...,
  get: (id: string) => ...,
  create: (data: WorkoutCreate) => ...,
  update: (id: string, data: WorkoutUpdate) => ...,
  delete: (id: string) => ...,
};
```

## API Communication Pattern

### HTTP Client Setup

```typescript
// lib/api/client.ts
import ky from 'ky';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const apiClient = ky.create({
  prefixUrl: API_BASE_URL,
  hooks: {
    beforeRequest: [
      (request) => {
        const token = localStorage.getItem('auth_token');
        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`);
        }
      },
    ],
  },
});
```

### API Functions

```typescript
// lib/api/workouts.ts
import { apiClient } from './client';
import type { Workout, WorkoutCreate, WorkoutListResponse } from '@/types/workout';

export const workoutsApi = {
  list: async (params?: { limit?: number; after?: string; timestamp_from?: string; timestamp_to?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', String(params.limit));
    if (params?.after) searchParams.set('after', params.after);
    if (params?.timestamp_from) searchParams.set('timestamp_from', params.timestamp_from);
    if (params?.timestamp_to) searchParams.set('timestamp_to', params.timestamp_to);
    
    return apiClient.get('workouts', { searchParams }).json<WorkoutListResponse>();
  },

  get: async (id: string) => {
    return apiClient.get(`workouts/${id}`).json<Workout>();
  },

  create: async (data: WorkoutCreate) => {
    return apiClient.post('workouts', { json: data }).json<Workout>();
  },

  update: async (id: string, data: WorkoutCreate) => {
    return apiClient.put(`workouts/${id}`, { json: data }).json<Workout>();
  },

  delete: async (id: string) => {
    await apiClient.delete(`workouts/${id}`);
  },
};
```

### TanStack Query Hooks

```typescript
// lib/hooks/use-workouts.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { workoutsApi } from '@/lib/api/workouts';
import type { WorkoutCreate } from '@/types/workout';

// Query keys
export const workoutKeys = {
  all: ['workouts'] as const,
  lists: () => [...workoutKeys.all, 'list'] as const,
  list: (params?: object) => [...workoutKeys.lists(), params] as const,
  details: () => [...workoutKeys.all, 'detail'] as const,
  detail: (id: string) => [...workoutKeys.details(), id] as const,
};

// List workouts
export function useWorkouts(params?: { limit?: number; after?: string }) {
  return useQuery({
    queryKey: workoutKeys.list(params),
    queryFn: () => workoutsApi.list(params),
  });
}

// Get single workout
export function useWorkout(id: string) {
  return useQuery({
    queryKey: workoutKeys.detail(id),
    queryFn: () => workoutsApi.get(id),
    enabled: !!id,
  });
}

// Create workout
export function useCreateWorkout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: WorkoutCreate) => workoutsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workoutKeys.lists() });
    },
  });
}

// Update workout
export function useUpdateWorkout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: WorkoutCreate }) => 
      workoutsApi.update(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: workoutKeys.lists() });
      queryClient.invalidateQueries({ queryKey: workoutKeys.detail(id) });
    },
  });
}

// Delete workout
export function useDeleteWorkout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => workoutsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workoutKeys.lists() });
    },
  });
}
```

## Form Handling Pattern

### Zod Schema

```typescript
// schemas/workout.ts
import { z } from 'zod';

export const workoutEntrySchema = z.object({
  exercise_id: z.string().uuid(),
  entry_type: z.enum(['warmup', 'main', 'top', 'backoff']),
  sets: z.number().int().min(1),
  reps: z.number().int().min(1),
  load_kg: z.number().min(0),
  rpe: z.number().min(1).max(10),
  notes: z.string().optional(),
});

export const workoutCreateSchema = z.object({
  timestamp: z.string().datetime(),
  body_weight_kg: z.number().positive().optional(),
  fatigue_level: z.number().int().min(1).max(5).optional(),
  sleep_hours: z.number().min(0).max(24).optional(),
  condition_notes: z.string().optional(),
  notes: z.string().optional(),
  entries: z.array(workoutEntrySchema).min(1),
});

export type WorkoutCreate = z.infer<typeof workoutCreateSchema>;
export type WorkoutEntry = z.infer<typeof workoutEntrySchema>;
```

### React Hook Form Integration

```tsx
// components/features/workout/workout-form.tsx
'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { workoutCreateSchema, type WorkoutCreate } from '@/schemas/workout';
import { useCreateWorkout } from '@/lib/hooks/use-workouts';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export function WorkoutForm() {
  const createWorkout = useCreateWorkout();

  const form = useForm<WorkoutCreate>({
    resolver: zodResolver(workoutCreateSchema),
    defaultValues: {
      timestamp: new Date().toISOString(),
      entries: [],
    },
  });

  const handleSubmit = form.handleSubmit((data) => {
    createWorkout.mutate(data, {
      onSuccess: () => {
        form.reset();
      },
    });
  });

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Form fields */}
      <Button type="submit" disabled={createWorkout.isPending}>
        {createWorkout.isPending ? 'Creating...' : 'Create Workout'}
      </Button>
    </form>
  );
}
```

## Error Handling Pattern

### Global Error Boundary

```tsx
// app/error.tsx
'use client';

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      <h2 className="text-2xl font-bold mb-4">Something went wrong!</h2>
      <p className="text-gray-600 mb-4">{error.message}</p>
      <button
        onClick={reset}
        className="px-4 py-2 bg-blue-500 text-white rounded"
      >
        Try again
      </button>
    </div>
  );
}
```

### API Error Handling

```typescript
// lib/api/client.ts
import ky, { HTTPError } from 'ky';

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string,
    public details?: unknown
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export const apiClient = ky.create({
  prefixUrl: API_BASE_URL,
  hooks: {
    afterResponse: [
      async (_request, _options, response) => {
        if (!response.ok) {
          const body = await response.json().catch(() => ({}));
          throw new ApiError(
            body.message || 'An error occurred',
            response.status,
            body.code,
            body.details
          );
        }
      },
    ],
  },
});
```

## Authentication Token Handling

### Token Storage (Simple approach for PoC)

```typescript
// lib/auth.ts
const TOKEN_KEY = 'auth_token';

export const auth = {
  getToken: () => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(TOKEN_KEY);
  },

  setToken: (token: string) => {
    localStorage.setItem(TOKEN_KEY, token);
  },

  removeToken: () => {
    localStorage.removeItem(TOKEN_KEY);
  },

  isAuthenticated: () => {
    return !!auth.getToken();
  },
};
```

### Auth Provider (optional)

```tsx
// components/providers/auth-provider.tsx
'use client';

import { createContext, useContext, useState, useEffect } from 'react';
import { auth } from '@/lib/auth';

type AuthContextType = {
  isAuthenticated: boolean;
  token: string | null;
  setToken: (token: string) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setTokenState] = useState<string | null>(null);

  useEffect(() => {
    setTokenState(auth.getToken());
  }, []);

  const setToken = (newToken: string) => {
    auth.setToken(newToken);
    setTokenState(newToken);
  };

  const logout = () => {
    auth.removeToken();
    setTokenState(null);
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated: !!token,
        token,
        setToken,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
```

## Date Handling

Use `date-fns` for date manipulation and formatting:

```typescript
// lib/utils/date.ts
import { format, parseISO, formatDistanceToNow } from 'date-fns';
import { ja } from 'date-fns/locale';

export function formatDate(isoString: string): string {
  return format(parseISO(isoString), 'yyyy/MM/dd', { locale: ja });
}

export function formatDateTime(isoString: string): string {
  return format(parseISO(isoString), 'yyyy/MM/dd HH:mm', { locale: ja });
}

export function formatRelative(isoString: string): string {
  return formatDistanceToNow(parseISO(isoString), { addSuffix: true, locale: ja });
}
```

## TypeScript Policy

**Prefer TypeScript**: Use TypeScript (`.ts`/`.tsx`) whenever possible. JavaScript is acceptable for configuration files or when tools require it.

### Guidelines

- ✅ Use `.ts` for non-React files
- ✅ Use `.tsx` for React component files
- ✅ Use `.js`/`.mjs` for config files if needed (e.g., `postcss.config.mjs`)
- ⚠️ Avoid `// @ts-ignore` and `// @ts-nocheck` (fix the type issue instead)
- ⚠️ Avoid `any` type (use `unknown` or proper types)

### Strict TypeScript Configuration

```json
// tsconfig.json (key settings)
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true
  }
}
```

## Linting and Formatting with Biome

**Biome is required** for all frontend projects. It replaces ESLint and Prettier with a single, fast tool.

### Installation

```bash
pnpm add -D @biomejs/biome
```

### Biome Configuration

```json
// biome.json
{
  "$schema": "https://biomejs.dev/schemas/1.9.4/schema.json",
  "organizeImports": {
    "enabled": true
  },
  "linter": {
    "enabled": true,
    "rules": {
      "recommended": true,
      "complexity": {
        "noExcessiveCognitiveComplexity": "warn"
      },
      "correctness": {
        "noUnusedImports": "error",
        "noUnusedVariables": "error",
        "useExhaustiveDependencies": "warn"
      },
      "style": {
        "noNonNullAssertion": "warn",
        "useConst": "error",
        "useImportType": "error"
      },
      "suspicious": {
        "noExplicitAny": "error"
      }
    }
  },
  "formatter": {
    "enabled": true,
    "indentStyle": "space",
    "indentWidth": 2,
    "lineWidth": 100
  },
  "javascript": {
    "formatter": {
      "quoteStyle": "single",
      "semicolons": "always",
      "trailingCommas": "es5"
    }
  },
  "files": {
    "ignore": [
      "node_modules",
      ".next",
      "dist",
      "*.gen.ts"
    ]
  }
}
```

### Package.json Scripts

```json
{
  "scripts": {
    "lint": "biome lint .",
    "lint:fix": "biome lint --write .",
    "format": "biome format --write .",
    "check": "biome check .",
    "check:fix": "biome check --write ."
  }
}
```

### Pre-commit Check

Run before committing:

```bash
pnpm check
```

This runs both linting and formatting checks in a single command.
