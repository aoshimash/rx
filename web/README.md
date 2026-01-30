# OPTel Workout Web Frontend

Web application for OPTel Workout - training program planning, workout management, and data visualization.

## Overview

This is the web frontend for OPTel Workout, providing:

- **Program Management**: Create and edit training programs
- **Workout Management**: View, create, and edit workout records
- **Exercise Catalog**: Manage exercise definitions
- **Telemetry Visualization**: View time-series training data

## Technology Stack

| Category | Technology |
|----------|------------|
| Framework | Next.js (App Router) |
| Language | TypeScript (no JavaScript) |
| Linter/Formatter | Biome |
| State Management | TanStack Query |
| UI Library | shadcn/ui |
| Styling | Tailwind CSS |
| Forms | React Hook Form + Zod |
| HTTP Client | ky |

## Getting Started

### Prerequisites

- Node.js 20+
- pnpm (recommended) or npm
- OPTel Workout API running at `http://localhost:8080`

### Installation

```bash
cd web
pnpm install
```

### Development

```bash
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Linting and Formatting

```bash
pnpm check       # Run lint and format check
pnpm check:fix   # Fix lint and format issues
```

### Build

```bash
pnpm build
pnpm start
```

## Project Structure

```
web/
├── app/                    # Next.js App Router
│   ├── layout.tsx          # Root layout
│   ├── page.tsx            # Home page
│   ├── workouts/           # Workout pages
│   ├── programs/           # Program pages
│   ├── exercises/          # Exercise pages
│   └── telemetry/          # Telemetry pages
├── components/
│   ├── ui/                 # shadcn/ui components
│   └── features/           # Feature-specific components
├── lib/
│   ├── api/                # API client
│   └── hooks/              # TanStack Query hooks
├── types/                  # TypeScript types
└── schemas/                # Zod validation schemas
```

## Configuration

### Environment Variables

Create `.env.local` for local development:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Authentication

For development, provide a Bearer token via the UI. The token is stored in localStorage.

## Conventions

See [.claude/skills/optel-frontend-standards/](../.claude/skills/optel-frontend-standards/) for coding standards including:

- Directory structure
- Naming conventions
- API communication patterns
- Component separation rules

## Related Documentation

- [Frontend Architecture](../docs/FRONTEND_ARCHITECTURE.md) - Architecture decisions
- [API Documentation](../api/README.md) - Backend API reference
- [OpenAPI Spec](../api/openapi/openapi.yaml) - API specification
