# Implementation Plan: Web PoC

**Branch**: `007-web-poc` | **Date**: 2026-01-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/007-web-poc/spec.md`

## Summary

Build a Proof of Concept web application for OPTel Workout with a **Plan vs Actual** comparison UI. The primary screens are Week View (main), Workout Input (modal), and Program Editor (page). Uses existing Next.js 15 + React 19 setup in `web/` directory, extending with TanStack Query for API state management, shadcn/ui for components, and ky for HTTP client.

## Technical Context

**Language/Version**: TypeScript 5.7+ (strict mode, no JavaScript)  
**Primary Dependencies**: Next.js 15.1 (App Router), React 19, TanStack Query, shadcn/ui, React Hook Form, Zod, ky, date-fns  
**Storage**: Browser localStorage (token persistence only), Backend API (PostgreSQL via REST)  
**Testing**: Vitest (unit), Playwright (E2E) - deferred for PoC  
**Target Platform**: Modern browsers (Chrome, Firefox, Safari, Edge)  
**Project Type**: Web application (frontend consuming existing REST API)  
**Performance Goals**: Page load < 3s, API responses < 1s perceived  
**Constraints**: PoC scope - no login UI, manual token input, no video/telemetry features  
**Scale/Scope**: Single user, ~5 screens (Week View, Program Editor, Program List, Settings)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Dumb Backend**: ✅ PASS - Frontend only consumes existing REST API. No backend changes required. All business logic (diff calculation, schedule generation) is UI-only presentation logic.
- **Domain-Driven Schema-First**: ✅ PASS - Uses existing OpenAPI spec (`api/openapi/openapi.yaml`). Frontend types will be generated or aligned with API contract.
- **Terminology**: ✅ PASS - Uses established terms: Workout, Exercise, Program, Entry (not "health metrics" or similar).
- **Clean Architecture**: ✅ PASS - Frontend follows separation: Pages → Components → Hooks → API Client.
- **Monorepo Structure**: ✅ PASS - All code goes in `web/` directory. Does not modify `api/` or other components.

## Project Structure

### Documentation (this feature)

```text
specs/007-web-poc/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (frontend types)
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (API client interface)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
web/
├── app/                          # Next.js App Router pages
│   ├── page.tsx                  # Week View (main screen)
│   ├── layout.tsx                # Root layout with providers
│   ├── programs/
│   │   ├── page.tsx              # Program List
│   │   ├── new/page.tsx          # Program Editor (create)
│   │   └── [id]/edit/page.tsx    # Program Editor (edit)
│   └── settings/
│       └── page.tsx              # Settings (token input)
├── components/
│   ├── ui/                       # shadcn/ui components
│   ├── week-view/                # Week View feature components
│   │   ├── WeekView.tsx
│   │   ├── DayAccordion.tsx
│   │   ├── ExerciseTable.tsx
│   │   └── StatusBadge.tsx
│   ├── workout-input/            # Workout Input modal components
│   │   ├── WorkoutModal.tsx
│   │   ├── ExerciseInputRow.tsx
│   │   └── AddExerciseButton.tsx
│   ├── program-editor/           # Program Editor components
│   │   ├── ProgramForm.tsx
│   │   ├── WeekAccordion.tsx
│   │   ├── DayAccordion.tsx
│   │   ├── ExerciseTable.tsx
│   │   └── ExerciseCombobox.tsx
│   └── schedule/
│       └── ScheduleModal.tsx
├── lib/
│   ├── api/                      # API client
│   │   ├── client.ts             # ky instance with auth
│   │   ├── programs.ts           # Program API functions
│   │   ├── workouts.ts           # Workout API functions
│   │   └── exercises.ts          # Exercise API functions
│   ├── hooks/                    # TanStack Query hooks
│   │   ├── usePrograms.ts
│   │   ├── useWorkouts.ts
│   │   └── useExercises.ts
│   └── utils/
│       ├── diff.ts               # Plan vs Actual diff calculation
│       ├── schedule.ts           # Schedule generation logic
│       └── export.ts             # CSV export
├── types/                        # TypeScript type definitions
│   └── api.ts                    # API response types (aligned with OpenAPI)
├── schemas/                      # Zod validation schemas
│   └── forms.ts                  # Form validation schemas
└── stores/
    └── auth.ts                   # Token storage (localStorage)
```

**Structure Decision**: Web application structure following Next.js App Router conventions and `docs/FRONTEND_ARCHITECTURE.md` guidelines. Feature components organized by screen (week-view, workout-input, program-editor).

## Complexity Tracking

No constitution violations identified. No complexity justifications required.

---

## Generated Artifacts

| Artifact | Path | Description |
|----------|------|-------------|
| Research | [research.md](./research.md) | Technology decisions and patterns |
| Data Model | [data-model.md](./data-model.md) | Frontend TypeScript types |
| API Client Interface | [contracts/api-client.ts](./contracts/api-client.ts) | API client interface definition |
| API Types | [contracts/types.ts](./contracts/types.ts) | TypeScript types for API |
| Quickstart | [quickstart.md](./quickstart.md) | Setup and development guide |

## Constitution Re-Check (Post-Design)

- **Dumb Backend**: ✅ PASS - No backend changes. Diff calculation, schedule generation are pure frontend presentation logic.
- **Domain-Driven Schema-First**: ✅ PASS - Types in `contracts/types.ts` align with OpenAPI spec.
- **Terminology**: ✅ PASS - All terms (Workout, Exercise, Program, Entry) match API.
- **Clean Architecture**: ✅ PASS - Frontend follows Pages → Components → Hooks → API Client separation.
- **Monorepo Structure**: ✅ PASS - All changes contained in `web/` directory.

## Next Steps

Run `/speckit.tasks` to generate implementation tasks from this plan.
