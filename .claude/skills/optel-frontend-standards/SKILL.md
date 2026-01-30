---
name: optel-frontend-standards
description: Frontend coding standards for OPTel project (Web and Mobile). Covers technology stack, project structure, naming conventions, and API communication patterns. Use when writing React/React Native code, creating new files in web/ or mobile/, or setting up frontend components.
---

# OPTel Frontend Standards

## Technology Stack

### Web

| Component | Choice |
|-----------|--------|
| Framework | Next.js (App Router) |
| Language | TypeScript (no JavaScript) |
| Linter/Formatter | Biome |
| State Management | TanStack Query (React Query) |
| UI Library | shadcn/ui |
| Styling | Tailwind CSS |
| Form Handling | React Hook Form + Zod |
| HTTP Client | ky or fetch |

### Mobile

| Component | Choice |
|-----------|--------|
| Framework | React Native + Expo |
| Language | TypeScript (no JavaScript) |
| Linter/Formatter | Biome |
| State Management | TanStack Query |
| Navigation | Expo Router |

## Project Structure

```
optel-workout/
├── api/              # REST API (Go) - existing
├── web/              # Web frontend (Next.js)
├── mobile/           # Mobile app (React Native + Expo)
└── packages/         # Shared packages (future)
    └── shared/       # Types, API client, validation schemas
```

## Key Principles

1. **Prefer TypeScript** - Use TypeScript whenever possible. JavaScript is acceptable for config files or when tools require it.
2. **React Shared Mental Model** - Web and Mobile share React patterns
3. **Type Safety** - TypeScript everywhere, Zod for runtime validation
4. **API-First** - Use TanStack Query for all API communication
5. **Component Separation** - Clear distinction between Page, Feature, and UI components
6. **Biome Required** - Use Biome for linting and formatting (no ESLint/Prettier)

## Web vs Mobile Scope

| Feature | Web | Mobile |
|---------|-----|--------|
| Program Management (CRUD) | ✅ | ❌ (read-only) |
| Workout Recording | ✅ | ✅ (optimized UI) |
| Workout Viewing | ✅ | ✅ (simple) |
| Exercise Management | ✅ | ❌ |
| Telemetry Visualization | ✅ | ❌ |
| Video Upload | ✅ | ✅ (WiFi only) |
| AI Analysis | 🔜 (future) | ❌ |

## Additional Resources

For detailed guidelines, see [reference.md](reference.md).
