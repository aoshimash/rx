# Quickstart: Web PoC

**Date**: 2026-01-30  
**Feature**: 007-web-poc

## Prerequisites

- Node.js 20+ (or use `pnpm` with auto-install)
- pnpm (package manager)
- Backend API running at `http://localhost:8080` (see `api/README.md`)

## Setup

### 1. Install Dependencies

```bash
cd web

# Install existing dependencies
pnpm install

# Add new dependencies for this feature
pnpm add @tanstack/react-query @hookform/resolvers react-hook-form zod ky date-fns

# Add shadcn/ui CLI and init (if not already done)
pnpm add -D @shadcn/ui
npx shadcn@latest init
```

### 2. Add shadcn/ui Components

```bash
cd web

# Core components
npx shadcn@latest add button
npx shadcn@latest add input
npx shadcn@latest add form
npx shadcn@latest add select
npx shadcn@latest add table

# Layout components
npx shadcn@latest add accordion
npx shadcn@latest add dialog
npx shadcn@latest add alert-dialog

# Data entry components
npx shadcn@latest add calendar
npx shadcn@latest add command  # for combobox/autocomplete
npx shadcn@latest add popover

# Display components
npx shadcn@latest add badge
npx shadcn@latest add card
npx shadcn@latest add separator
```

### 3. Environment Configuration

Create `web/.env.local`:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### 4. Start Development

```bash
# Terminal 1: Start backend API
cd api
task run

# Terminal 2: Start web frontend
cd web
pnpm dev
```

Open http://localhost:3000 in your browser.

## Project Structure After Setup

```
web/
├── app/
│   ├── page.tsx                  # Week View (main)
│   ├── layout.tsx                # Root layout with providers
│   ├── programs/
│   │   ├── page.tsx              # Program List
│   │   ├── new/page.tsx          # Program Editor (create)
│   │   └── [id]/edit/page.tsx    # Program Editor (edit)
│   └── settings/
│       └── page.tsx              # Settings (token input)
├── components/
│   ├── ui/                       # shadcn/ui components
│   ├── providers/
│   │   └── QueryProvider.tsx     # TanStack Query provider
│   ├── week-view/
│   ├── workout-input/
│   ├── program-editor/
│   └── schedule/
├── lib/
│   ├── api/
│   │   ├── client.ts
│   │   ├── exercises.ts
│   │   ├── programs.ts
│   │   └── workouts.ts
│   ├── hooks/
│   │   ├── useExercises.ts
│   │   ├── usePrograms.ts
│   │   └── useWorkouts.ts
│   └── utils/
│       ├── diff.ts
│       ├── schedule.ts
│       └── export.ts
├── types/
│   └── api.ts
├── schemas/
│   └── forms.ts
└── stores/
    └── auth.ts
```

## Initial Implementation Steps

### Step 1: Setup TanStack Query Provider

```tsx
// components/providers/QueryProvider.tsx
'use client';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useState } from 'react';

export function QueryProvider({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000, // 1 minute
            retry: 1,
          },
        },
      })
  );

  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
}
```

### Step 2: Update Root Layout

```tsx
// app/layout.tsx
import { QueryProvider } from '@/components/providers/QueryProvider';
import './globals.css';

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <QueryProvider>{children}</QueryProvider>
      </body>
    </html>
  );
}
```

### Step 3: Create API Client

```tsx
// lib/api/client.ts
import ky from 'ky';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

const getToken = () => {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem('optel_token');
};

export const api = ky.create({
  prefixUrl: API_URL,
  hooks: {
    beforeRequest: [
      (request) => {
        const token = getToken();
        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`);
        }
      },
    ],
    afterResponse: [
      async (_request, _options, response) => {
        if (response.status === 401) {
          // Clear invalid token
          localStorage.removeItem('optel_token');
          // Trigger auth modal (implement via global state or event)
          window.dispatchEvent(new CustomEvent('auth:required'));
        }
        return response;
      },
    ],
  },
});
```

## Development Workflow

1. **Backend first**: Ensure API is running and accessible
2. **Token setup**: Use Settings page to enter a valid Bearer token
3. **Test API**: Verify requests work via browser DevTools Network tab
4. **Build features**: Implement screens in priority order (P1 first)

## Testing the Setup

After setup, you should be able to:

1. Open http://localhost:3000
2. See the Week View page (initially empty)
3. Open DevTools → Console, no errors
4. Navigate to /settings to enter a token
5. Navigate to /programs to see program list (empty if no programs exist)

## Common Issues

### CORS Errors

If you see CORS errors, ensure the backend API has CORS enabled for `http://localhost:3000`.

### 401 Unauthorized

The app prompts for a token when API returns 401. Enter any non-empty string for PoC (backend accepts any token).

### Type Errors

Run `pnpm check` to verify TypeScript types are correct. Types in `types/api.ts` should match `api/openapi/openapi.yaml`.
