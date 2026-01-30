# Tasks: Web PoC

**Input**: Design documents from `/specs/007-web-poc/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Deferred for PoC (per plan.md)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

This feature uses the OPTel monorepo structure:
- Frontend: `web/` directory
- All paths below are relative to repository root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization, dependency installation, and base configuration

- [X] T001 Install new dependencies: `pnpm add @tanstack/react-query @hookform/resolvers react-hook-form zod ky date-fns` in web/
- [X] T002 Initialize shadcn/ui with `npx shadcn@latest init` in web/
- [X] T003 [P] Add shadcn/ui core components: button, input, form, select, table in web/components/ui/
- [X] T004 [P] Add shadcn/ui layout components: accordion, dialog, alert-dialog in web/components/ui/
- [X] T005 [P] Add shadcn/ui data entry components: calendar, command, popover in web/components/ui/
- [X] T006 [P] Add shadcn/ui display components: badge, card, separator in web/components/ui/
- [X] T007 Create environment configuration file web/.env.local with NEXT_PUBLIC_API_URL

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T008 Create TypeScript API types aligned with OpenAPI spec in web/types/api.ts
- [X] T009 [P] Create Zod validation schemas for forms in web/schemas/forms.ts
- [X] T010 [P] Create auth token store with localStorage in web/stores/auth.ts
- [X] T011 Create API client base with ky and auth headers in web/lib/api/client.ts
- [X] T012 [P] Create exercises API functions in web/lib/api/exercises.ts
- [X] T013 [P] Create programs API functions in web/lib/api/programs.ts
- [X] T014 [P] Create workouts API functions in web/lib/api/workouts.ts
- [X] T015 Create TanStack Query provider component in web/components/providers/QueryProvider.tsx
- [X] T016 Update root layout with QueryProvider in web/app/layout.tsx
- [X] T017 [P] Create useExercises hook with TanStack Query in web/lib/hooks/useExercises.ts
- [X] T018 [P] Create usePrograms hook with TanStack Query in web/lib/hooks/usePrograms.ts
- [X] T019 [P] Create useWorkouts hook with TanStack Query in web/lib/hooks/useWorkouts.ts
- [X] T020 Create diff calculation utility in web/lib/utils/diff.ts

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - View Week with Plan/Actual/Diff (Priority: P1) 🎯 MVP

**Goal**: Display training week with Plan/Actual/Diff comparison, status indicators, and collapsible days

**Independent Test**: Select a program, view a week, verify days are collapsible with Plan/Actual columns and status indicators (✓/≠/○)

### Implementation for User Story 1

- [X] T021 [P] [US1] Create StatusBadge component for ✓/≠/○/📝 indicators in web/components/week-view/StatusBadge.tsx
- [X] T022 [P] [US1] Create ExerciseTable component showing Plan/Actual/Diff columns in web/components/week-view/ExerciseTable.tsx
- [X] T023 [US1] Create DayAccordion component with collapsible day content in web/components/week-view/DayAccordion.tsx
- [X] T024 [US1] Create WeekView component with week navigation (Prev/Next) in web/components/week-view/WeekView.tsx
- [X] T025 [US1] Create UnplannedWorkouts section component in web/components/week-view/UnplannedWorkouts.tsx
- [X] T026 [US1] Implement Week View main page integrating all components in web/app/page.tsx

**Checkpoint**: User Story 1 complete - Week View displays Plan vs Actual with diff indicators

---

## Phase 4: User Story 2 - Record Workout Results (Priority: P1)

**Goal**: Allow users to record actual workout results via modal with pre-populated plan values

**Independent Test**: Open workout input modal, fill in values for exercises, save, verify workout appears in Week View

### Implementation for User Story 2

- [X] T027 [P] [US2] Create ExerciseInputRow component for single exercise entry in web/components/workout-input/ExerciseInputRow.tsx
- [X] T028 [P] [US2] Create AddExerciseButton component for unplanned exercises in web/components/workout-input/AddExerciseButton.tsx
- [X] T029 [US2] Create WorkoutModal component with form, validation, and session notes field in web/components/workout-input/WorkoutModal.tsx
- [X] T030 [US2] Integrate WorkoutModal into Week View with trigger button in web/app/page.tsx
- [X] T031 [US2] Add workout creation mutation with optimistic update in web/lib/hooks/useWorkouts.ts

**Checkpoint**: User Story 2 complete - Users can record workouts and see them in Week View

---

## Phase 5: User Story 3 - Create/Edit Training Program (Priority: P1)

**Goal**: Enable users to create programs with hierarchical structure (Weeks → Days → Exercises)

**Independent Test**: Create a new program with weeks/days/exercises, save, verify it appears in program selection

### Implementation for User Story 3

- [X] T032 [P] [US3] Create ExerciseCombobox component with autocomplete in web/components/program-editor/ExerciseCombobox.tsx
- [X] T033 [P] [US3] Create ExerciseRow component for exercise prescription in web/components/program-editor/ExerciseRow.tsx
- [X] T034 [US3] Create ExerciseTable component for day's exercise list in web/components/program-editor/ExerciseTable.tsx
- [X] T035 [US3] Create DayAccordion component for collapsible day editor in web/components/program-editor/DayAccordion.tsx
- [X] T036 [US3] Create WeekAccordion component for collapsible week editor in web/components/program-editor/WeekAccordion.tsx
- [X] T037 [US3] Create ProgramForm component with full hierarchical editor in web/components/program-editor/ProgramForm.tsx
- [X] T038 [US3] Create Program Editor page (create mode) in web/app/programs/new/page.tsx
- [X] T039 [US3] Create Program Editor page (edit mode) in web/app/programs/[id]/edit/page.tsx
- [X] T040 [US3] Add program creation/update mutations in web/lib/hooks/usePrograms.ts
- [X] T041 [US3] Add delete confirmation dialog to program editor in web/components/program-editor/DeleteConfirmDialog.tsx

**Checkpoint**: User Story 3 complete - Users can create and edit training programs

---

## Phase 6: User Story 4 - Configure Schedule (Priority: P2)

**Goal**: Allow users to assign calendar dates to program days

**Independent Test**: Open Schedule Settings, set start date, verify dates appear in Week View

### Implementation for User Story 4

- [ ] T042 [P] [US4] Create schedule generation utility with skip weekends/avoid consecutive in web/lib/utils/schedule.ts
- [ ] T043 [US4] Create ScheduleModal component with date picker and options in web/components/schedule/ScheduleModal.tsx
- [ ] T044 [US4] Create schedule store for client-side schedule persistence in web/stores/schedule.ts
- [ ] T045 [US4] Integrate ScheduleModal into Week View header in web/app/page.tsx
- [ ] T046 [US4] Update DayAccordion to display scheduled dates in web/components/week-view/DayAccordion.tsx

**Checkpoint**: User Story 4 complete - Schedule dates display in Week View

---

## Phase 7: User Story 5 - Select Program (Priority: P2)

**Goal**: Allow users to select which program to view in Week View

**Independent Test**: View program list, select a program, verify Week View shows that program's weeks

### Implementation for User Story 5

- [ ] T047 [P] [US5] Create ProgramCard component for program list item in web/components/programs/ProgramCard.tsx
- [ ] T048 [US5] Create Program List page in web/app/programs/page.tsx
- [ ] T049 [US5] Create program selection store for current program in web/stores/program.ts
- [ ] T050 [US5] Add program selector dropdown to Week View header in web/app/page.tsx
- [ ] T051 [US5] Update Week View to filter by selected program in web/components/week-view/WeekView.tsx

**Checkpoint**: User Story 5 complete - Users can select and switch programs

---

## Phase 8: User Story 6 - Export Data (Priority: P2)

**Goal**: Allow users to export training data as CSV

**Independent Test**: Click Export, select Current Week, verify CSV downloads with correct columns

### Implementation for User Story 6

- [ ] T052 [P] [US6] Create CSV export utility with UTF-8 BOM in web/lib/utils/export.ts
- [ ] T053 [US6] Create ExportButton component with options dropdown in web/components/export/ExportButton.tsx
- [ ] T054 [US6] Integrate ExportButton into Week View header in web/app/page.tsx

**Checkpoint**: User Story 6 complete - Users can export data to CSV

---

## Phase 9: Authentication & Settings (Cross-cutting)

**Purpose**: Token management UI and settings page

- [ ] T055 [P] Create TokenInputModal component for lazy auth prompt in web/components/auth/TokenInputModal.tsx
- [ ] T056 Create Settings page with token input form in web/app/settings/page.tsx
- [ ] T057 Integrate TokenInputModal with auth:required event in web/app/layout.tsx
- [ ] T058 Add navigation header with links to Programs, Settings in web/components/layout/Header.tsx
- [ ] T059 Update root layout with Header component in web/app/layout.tsx

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T060 [P] Add loading states to all pages with Skeleton components
- [ ] T061 [P] Add error boundaries with user-friendly error messages and retry button for API failures
- [ ] T062 [P] Add empty state components for lists (no programs, no workouts)
- [ ] T063 Ensure responsive design works on tablet/mobile viewports
- [ ] T064 Run Biome lint and format checks across all new files
- [ ] T065 Validate quickstart.md instructions work end-to-end and verify page load < 3 seconds (SC-001)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-8)**: All depend on Foundational phase completion
- **Auth/Settings (Phase 9)**: Depends on Foundational phase
- **Polish (Phase 10)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational - Benefits from US1 but independently testable
- **User Story 3 (P1)**: Can start after Foundational - No dependencies on other stories
- **User Story 4 (P2)**: Can start after Foundational - Integrates with US1 Week View
- **User Story 5 (P2)**: Can start after Foundational - Integrates with US1 Week View
- **User Story 6 (P2)**: Can start after Foundational - Uses US1 Week View data

### Within Each User Story

- Components with [P] can be built in parallel (different files)
- Container components depend on child components
- Hooks/mutations depend on API functions
- Pages integrate all components

### Parallel Opportunities

**Setup Phase (T003-T006)**: All shadcn/ui component additions can run in parallel

**Foundational Phase**:
- T009, T010 can run in parallel
- T012, T013, T014 can run in parallel (API functions)
- T017, T018, T019 can run in parallel (hooks)

**User Story 1**:
- T021, T022 can run in parallel (leaf components)

**User Story 3**:
- T032, T033 can run in parallel (leaf components)

---

## Parallel Example: Foundational Phase

```bash
# Launch API function tasks together:
Task T012: "Create exercises API functions in web/lib/api/exercises.ts"
Task T013: "Create programs API functions in web/lib/api/programs.ts"
Task T014: "Create workouts API functions in web/lib/api/workouts.ts"

# Launch hook tasks together (after API functions):
Task T017: "Create useExercises hook in web/lib/hooks/useExercises.ts"
Task T018: "Create usePrograms hook in web/lib/hooks/usePrograms.ts"
Task T019: "Create useWorkouts hook in web/lib/hooks/useWorkouts.ts"
```

---

## Implementation Strategy

### MVP First (User Stories 1-3 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 - View Week
4. Complete Phase 4: User Story 2 - Record Workout
5. Complete Phase 5: User Story 3 - Create Program
6. **STOP and VALIDATE**: Test all P1 stories independently
7. Add Phase 9: Auth/Settings for token management
8. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add User Story 1 → Week View displays (needs existing data)
3. Add User Story 3 → Can create programs to view
4. Add User Story 2 → Can record workouts (MVP complete!)
5. Add User Story 5 → Can switch programs
6. Add User Story 4 → Schedule dates display
7. Add User Story 6 → Can export data

### Suggested MVP Scope

**Minimum**: User Stories 1, 2, 3 + Auth/Settings
- View Week with Plan/Actual/Diff
- Record workout results
- Create/edit training programs
- Token management

This delivers the core "Plan vs Actual" value proposition.

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Tests deferred for PoC (can be added in future iteration)
