# Frontend Architecture

This document describes the architecture decisions for Rx frontend applications.

## Overview

Rx has two frontend applications with distinct purposes:

| Application | Purpose | Priority |
|-------------|---------|----------|
| **Web** | Planning, analysis, full management | First |
| **Mobile** | Gym-side workout logging, plan viewing | Second |

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
| **Workout** | Full CRUD | Create + Read (optimized UI) |
| **Exercise** | Full CRUD | Read-only (via workout) |
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
│  │  useWorkouts(), useCreateWorkout(), useExercises()   │    │
│  └───────────────────────┬───────────────────────────────┘    │
│                          │                                     │
│  ┌───────────────────────▼───────────────────────────────┐    │
│  │                   API Client (ky)                       │    │
│  │  workoutsApi.list(), workoutsApi.create(), ...        │    │
│  └───────────────────────┬───────────────────────────────┘    │
└──────────────────────────┼────────────────────────────────────┘
                           │ HTTP (Bearer Token)
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Rx API                         │
│                  http://localhost:8080/api/v1               │
└─────────────────────────────────────────────────────────────┘
```

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

- View cached workout plans offline
- Record workouts offline
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

- [AGENTS.md](../AGENTS.md) - AI agent guidance
- [optel-frontend-standards](../.claude/skills/optel-frontend-standards/) - Coding standards
- [ARCHITECTURE.md](ARCHITECTURE.md) - Overall system architecture
