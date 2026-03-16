# Frontend Architecture

This document describes the architecture decisions and coding standards for Rx frontend applications.

## Overview

Rx has two frontend applications with distinct purposes:

| Application | Purpose | Priority |
|-------------|---------|----------|
| **Web** | Planning, analysis, full management | First |
| **Mobile** | Gym-side training logging, plan viewing | Second |

## Design Principles

### 1. Web is a Superset of Mobile

Everything possible on Mobile should also be possible on Web. Mobile has a focused, optimized UI for gym use.

### 2. React Shared Mental Model

Both Web and Mobile use React (React/Next.js for Web, React Native/Expo for Mobile), enabling:
- Shared knowledge and patterns
- Potential code sharing (types, API clients, validation)
- Consistent development experience

### 3. AI-Friendly Development

Technology choices prioritize:
- Large community and documentation (better AI training data)
- Simple configuration (less boilerplate)
- Type safety (fewer AI mistakes)
- Consistent patterns (predictable code generation)

## Feature Scope

| Feature | Web | Mobile |
|---------|-----|--------|
| **Program** | Full CRUD | Read-only |
| **Plan** | Full CRUD | Read-only |
| **Log** | Full CRUD | Create + Read (optimized UI) |
| **Exercise** | Full CRUD | Read-only (via log) |
| **Telemetry** | View + Visualize | Not available |
| **Video** | Upload/Download | Upload (WiFi only) |
| **AI Analysis** | Future | Not planned |

## Technology Stack

### Web

| Category | Technology | Rationale |
|----------|------------|-----------|
| Framework | Next.js (App Router) | React official recommendation, largest community |
| Language | TypeScript | Type safety, AI-friendly (**no JavaScript**) |
| Linter/Formatter | Biome | Single tool, fast, simple config |
| State Management | TanStack Query | API-focused, automatic caching |
| UI Library | shadcn/ui | Copy-paste components, easy customization |
| Styling | Tailwind CSS | Utility-first, consistent patterns |
| Forms | React Hook Form + Zod | Type-safe validation |
| HTTP Client | ky | Simple, lightweight |
| Date Handling | date-fns | Tree-shakeable, functional |

### Mobile

| Category | Technology | Rationale |
|----------|------------|-----------|
| Framework | React Native + Expo | Simplified build, easier AI iteration |
| Language | TypeScript | Shared with Web (**no JavaScript**) |
| Linter/Formatter | Biome | Consistent with Web |
| State Management | TanStack Query | Consistent with Web |
| Navigation | Expo Router | File-based routing like Next.js |

### Shared (Future)

```
packages/
└── shared/
    ├── types/          # TypeScript type definitions
    ├── schemas/        # Zod validation schemas
    └── api/            # API client functions
```

## Project Structure

```
rx/
├── api/                    # Go REST API (existing)
├── web/                    # Next.js Web application
│   ├── app/                # App Router pages
│   ├── components/
│   │   ├── ui/             # shadcn/ui components
│   │   └── features/       # Feature-specific components
│   ├── lib/
│   │   ├── api/            # API client
│   │   └── hooks/          # TanStack Query hooks
│   ├── types/              # Type definitions
│   └── schemas/            # Zod schemas
├── mobile/                 # React Native + Expo (future)
├── packages/               # Shared packages (future)
│   └── shared/
└── docs/                   # Documentation
```

### Web Directory Structure (detailed)

```
web/
├── app/                        # Next.js App Router (routing)
│   ├── layout.tsx              # Root layout
│   ├── page.tsx                # Home page
│   ├── programs/
│   │   ├── page.tsx            # Program list
│   │   └── [id]/
│   │       └── page.tsx        # Program detail/edit
│   ├── plans/
│   │   ├── page.tsx            # Plan list
│   │   └── [id]/
│   │       └── page.tsx        # Plan detail
│   └── logs/
│       ├── page.tsx            # Log list
│       └── [id]/
│           └── page.tsx        # Log detail
├── components/
│   ├── ui/                     # shadcn/ui components
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   └── ...
│   └── features/               # Feature-specific components
│       ├── program/
│       │   ├── program-list.tsx
│       │   └── program-form.tsx
│       ├── plan/
│       │   ├── plan-list.tsx
│       │   └── plan-form.tsx
│       └── log/
│           ├── log-list.tsx
│           └── log-form.tsx
├── lib/
│   ├── api/                    # API client
│   │   ├── client.ts           # Base HTTP client
│   │   ├── programs.ts         # Program API functions
│   │   ├── plans.ts            # Plan API functions
│   │   └── logs.ts             # Log API functions
│   ├── hooks/                  # Custom React hooks
│   │   ├── use-programs.ts
│   │   ├── use-plans.ts
│   │   └── use-logs.ts
│   └── utils/                  # Utility functions
│       └── date.ts             # Date formatting (date-fns)
├── types/                      # TypeScript type definitions
│   ├── program.ts
│   ├── plan.ts
│   ├── log.ts
│   └── api.ts                  # API response types
└── schemas/                    # Zod validation schemas
    ├── program.ts
    ├── plan.ts
    └── log.ts
```

## Data Flow

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                    Components                         │    │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐             │    │
│  │  │  Page   │  │ Feature │  │   UI    │             │    │
│  │  │         │──│Component│──│Component│             │    │
│  │  └────┬────┘  └────┬────┘  └─────────┘             │    │
│  │       │            │                                 │    │
│  └───────┼────────────┼─────────────────────────────────┘    │
│          │            │                                       │
│  ┌───────▼────────────▼─────────────────────────────────┐    │
│  │              TanStack Query Hooks                      │    │
│  │  usePrograms(), usePlans(), useLogs(), ...            │    │
│  └───────────────────────┬───────────────────────────────┘    │
│                          │                                     │
│  ┌───────────────────────▼───────────────────────────────┐    │
│  │                   API Client (ky)                       │    │
│  │  programsApi.list(), logsApi.create(), ...             │    │
│  └───────────────────────┬───────────────────────────────┘    │
└──────────────────────────┼────────────────────────────────────┘
                           │ HTTP (Bearer Token)
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Rx API                                    │
│                  http://localhost:8080/api/v1                │
└─────────────────────────────────────────────────────────────┘
```

## Coding Standards

### Component Separation Rules

| Type | Location | Role | Example |
|------|----------|------|---------|
| **Page** | `app/**/page.tsx` | Data fetching, layout composition | `app/programs/page.tsx` |
| **Feature Component** | `components/features/` | Feature-specific UI + logic | `ProgramList`, `LogForm` |
| **UI Component** | `components/ui/` | Reusable presentational | `Button`, `Card`, `Input` |

#### Page Component Example

```tsx
// app/programs/page.tsx
import { ProgramList } from '@/components/features/program/program-list';

export default function ProgramsPage() {
  return (
    <div className="container py-8">
      <h1 className="text-2xl font-bold mb-6">Programs</h1>
      <ProgramList />
    </div>
  );
}
```

#### Feature Component Example

```tsx
// components/features/program/program-list.tsx
'use client';

import { usePrograms } from '@/lib/hooks/use-programs';
import { ProgramCard } from './program-card';

export function ProgramList() {
  const { data, isLoading, error } = usePrograms();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div className="space-y-4">
      {data?.items.map((program) => (
        <ProgramCard key={program.id} program={program} />
      ))}
    </div>
  );
}
```

### Naming Conventions

#### Files

```
program-list.tsx        # Components: kebab-case
use-programs.ts         # Hooks: use-*.ts
program.ts              # Types/schemas: singular, kebab-case
```

#### Components and Variables

```typescript
// Component names: PascalCase
export function ProgramList() { ... }
export function LogEntryCard() { ... }

// Props types: ComponentNameProps
type ProgramListProps = { ... }
type LogEntryCardProps = { ... }

// Variables and functions: camelCase
const programList = [];
const handleSubmit = () => {};

// Event handlers: handle + Action
const handleClick = () => {};
const handleSubmit = () => {};
const handleProgramCreate = () => {};

// Boolean variables: is/has/should prefix
const isLoading = true;
const hasError = false;
const shouldRefresh = true;
```

#### API Functions

```typescript
// lib/api/programs.ts
export const programsApi = {
  list: (params?: ProgramListParams) => ...,
  get: (id: string) => ...,
  create: (data: ProgramCreate) => ...,
  update: (id: string, data: ProgramUpdate) => ...,
  delete: (id: string) => ...,
};
```

### API Communication Pattern

#### HTTP Client Setup

```typescript
// lib/api/client.ts
import ky from 'ky';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

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

#### TanStack Query Hooks

```typescript
// lib/hooks/use-programs.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { programsApi } from '@/lib/api/programs';
import type { ProgramCreate } from '@/types/program';

// Query keys
export const programKeys = {
  all: ['programs'] as const,
  lists: () => [...programKeys.all, 'list'] as const,
  list: (params?: object) => [...programKeys.lists(), params] as const,
  details: () => [...programKeys.all, 'detail'] as const,
  detail: (id: string) => [...programKeys.details(), id] as const,
};

export function usePrograms(params?: { limit?: number }) {
  return useQuery({
    queryKey: programKeys.list(params),
    queryFn: () => programsApi.list(params),
  });
}

export function useProgram(id: string) {
  return useQuery({
    queryKey: programKeys.detail(id),
    queryFn: () => programsApi.get(id),
    enabled: !!id,
  });
}

export function useCreateProgram() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: ProgramCreate) => programsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: programKeys.lists() });
    },
  });
}
```

### Form Handling Pattern

```typescript
// schemas/log.ts
import { z } from 'zod';

export const logEntrySchema = z.object({
  exercise_name: z.string().min(1).max(200),
  order: z.number().int().min(0),
  sets: z.number().int().min(1).optional(),
  reps: z.number().int().min(1).optional(),
  load_kg: z.number().min(0).optional(),
  rpe: z.number().int().min(1).max(10).optional(),
  notes: z.string().max(2000).optional(),
});

export const logCreateSchema = z.object({
  performed_at: z.string().datetime(),
  plan_id: z.string().uuid().optional(),
  notes: z.string().max(5000).optional(),
  entries: z.array(logEntrySchema).min(1).max(500),
});

export type LogCreate = z.infer<typeof logCreateSchema>;
```

```tsx
// components/features/log/log-form.tsx
'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { logCreateSchema, type LogCreate } from '@/schemas/log';
import { useCreateLog } from '@/lib/hooks/use-logs';
import { Button } from '@/components/ui/button';

export function LogForm() {
  const createLog = useCreateLog();

  const form = useForm<LogCreate>({
    resolver: zodResolver(logCreateSchema),
    defaultValues: {
      performed_at: new Date().toISOString(),
      entries: [],
    },
  });

  const handleSubmit = form.handleSubmit((data) => {
    createLog.mutate(data, {
      onSuccess: () => form.reset(),
    });
  });

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Form fields */}
      <Button type="submit" disabled={createLog.isPending}>
        {createLog.isPending ? 'Saving...' : 'Save Log'}
      </Button>
    </form>
  );
}
```

### Date Handling

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

### TypeScript Policy

**Prefer TypeScript**: Use TypeScript (`.ts`/`.tsx`) whenever possible. JavaScript is acceptable for configuration files or when tools require it.

- ✅ Use `.ts` for non-React files
- ✅ Use `.tsx` for React component files
- ✅ Use `.js`/`.mjs` for config files if needed (e.g., `postcss.config.mjs`)
- ⚠️ Avoid `// @ts-ignore` and `// @ts-nocheck` (fix the type issue instead)
- ⚠️ Avoid `any` type (use `unknown` or proper types)

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

### Linting and Formatting with Biome

**Biome is required** for all frontend projects. It replaces ESLint and Prettier with a single, fast tool.

```bash
pnpm add -D @biomejs/biome
```

```json
// biome.json
{
  "$schema": "https://biomejs.dev/schemas/1.9.4/schema.json",
  "organizeImports": { "enabled": true },
  "linter": {
    "enabled": true,
    "rules": {
      "recommended": true,
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
      "suspicious": { "noExplicitAny": "error" }
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
    "ignore": ["node_modules", ".next", "dist", "*.gen.ts"]
  }
}
```

Package.json scripts:

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

Run `pnpm check` before committing.

## Authentication

### Current (PoC)

- Bearer token authentication
- Token stored in localStorage
- Token manually provided by user (no login UI)
- Backend accepts any non-empty token for MVP

### Future

- OAuth2 / JWT authentication
- Proper login/logout UI
- Token refresh mechanism
- Multi-user support (coach/athlete sharing)

## Offline Support (Mobile)

### Requirements

- View cached training plans offline
- Record training sessions offline
- Sync when connection available
- Video upload only on WiFi

### Implementation (Future)

- TanStack Query persistence
- Optimistic updates
- Background sync with Expo

## AI Integration (Future)

AI analysis will be integrated in phases:

1. **Phase 1 (Current)**: External agents (Cursor + MCP Server)
2. **Phase 2**: Web-embedded AI chat/analysis
3. **Phase 3**: AI-powered recommendations

## Related Documents

- [DOMAIN_MODEL.md](DOMAIN_MODEL.md) — Domain model reference (Program, Plan, Log)
- [ARCHITECTURE.md](ARCHITECTURE.md) — Overall system architecture
- [PHILOSOPHY.md](PHILOSOPHY.md) — Core product philosophy and principles
