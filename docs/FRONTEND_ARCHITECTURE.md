# Frontend Architecture

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

Both Web and Mobile use React (Next.js for Web, React Native/Expo for Mobile), enabling shared knowledge, patterns, and potential code sharing.

### 3. AI-Friendly Development

Technology choices prioritize large communities, simple configuration, type safety, and consistent patterns.

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

## Project Structure

```
rx/
├── api/        # Go REST API
├── web/        # Next.js (app/, components/, lib/, types/, schemas/)
├── mobile/     # React Native + Expo (future)
└── packages/   # Shared types, schemas, API client (future)
```

## Data Flow

```
Page → Feature Component → TanStack Query Hook → API Client (ky) → Rx API
```

Components never call the API directly. All server state goes through TanStack Query.

## Coding Standards

### Component Separation

| Type | Location | Role |
|------|----------|------|
| **Page** | `app/**/page.tsx` | Data fetching, layout composition |
| **Feature Component** | `components/features/` | Feature-specific UI + logic |
| **UI Component** | `components/ui/` | Reusable presentational |

### Naming Conventions

- **Files**: kebab-case (`program-list.tsx`, `use-programs.ts`)
- **Components**: PascalCase (`ProgramList`, `LogEntryCard`)
- **Props types**: `ComponentNameProps`
- **Variables/functions**: camelCase
- **Event handlers**: `handle` + action (`handleSubmit`, `handleProgramCreate`)
- **Booleans**: `is`/`has`/`should` prefix (`isLoading`, `hasError`)
- **API objects**: `{resource}Api` with `list`, `get`, `create`, `update`, `delete` methods

### TypeScript Policy

- TypeScript everywhere (`.ts`/`.tsx`). JavaScript only for config files.
- No `any` — use `unknown` or proper types
- No `// @ts-ignore` — fix the type issue

### Linting

Use Biome for linting and formatting (no ESLint/Prettier). Run `pnpm check` before committing.

## Authentication

**Current (PoC):** Bearer token stored in localStorage, manually provided by user.

**Future:** OAuth2 / JWT with proper login UI and multi-user support.

## Offline Support (Mobile, Future)

- View cached plans offline
- Record training sessions offline, sync when reconnected
- Video upload only on WiFi

## AI Integration (Future)

1. **Phase 1 (Current)**: External agents via MCP Server
2. **Phase 2**: Web-embedded AI chat/analysis
3. **Phase 3**: AI-powered recommendations

## Related Documents

- [DOMAIN_MODEL.md](DOMAIN_MODEL.md) — Domain model (Program, Plan, Log)
- [ARCHITECTURE.md](ARCHITECTURE.md) — Overall system architecture
- [PHILOSOPHY.md](PHILOSOPHY.md) — Core product philosophy
