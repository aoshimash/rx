# Feature Specification: Web PoC

**Feature Branch**: `007-web-poc`  
**Created**: 2026-01-30  
**Status**: Draft  
**Input**: User description: "Create Web PoC based on docs/WEB_UI_DESIGN"

## Overview

Build a Proof of Concept web application for OPTel Workout based on the UI design in `docs/WEB_UI_DESIGN.md`. The core concept is a **Plan vs Actual** comparison UI that displays training programs alongside recorded workout results, highlighting differences at a glance.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Week with Plan/Actual/Diff (Priority: P1)

A user wants to view their training week showing planned exercises from their program alongside actual recorded results, with clear diff indicators.

**Why this priority**: This is the core value proposition of the UI - showing Plan vs Actual side-by-side with diffs. Without this, the app is just a generic workout logger.

**Independent Test**: Can be fully tested by selecting a program, viewing a week, and verifying that days are collapsible, exercises show Plan/Actual columns, and status indicators (✓/≠/○) display correctly. Delivers immediate value by showing users how their training compares to the plan.

**Acceptance Scenarios**:

1. **Given** a user has a program with scheduled days, **When** they view the Week View, **Then** they see each day with Plan/Actual/Diff columns and appropriate status indicators
2. **Given** a day has matching plan and actual values (Sets/Reps/Load), **When** viewing the day, **Then** it shows a green checkmark (✓)
3. **Given** a day has differences between plan and actual, **When** viewing the day, **Then** it shows a yellow diff indicator (≠) and the Diff column shows specific differences (e.g., "Load +5kg")
4. **Given** a day has a plan but no recorded workout, **When** viewing the day, **Then** it shows a gray pending indicator (○)
5. **Given** days are displayed, **When** the user clicks a day header, **Then** the day collapses/expands to show/hide the exercise table

---

### User Story 2 - Record Workout Results (Priority: P1)

A user wants to record their actual workout results against their planned training day.

**Why this priority**: Recording actual results is essential to populate the Plan vs Actual view. This is equally critical as viewing.

**Independent Test**: Can be fully tested by opening a day's workout input modal, filling in actual values for exercises, and saving. The recorded values appear in the Week View.

**Acceptance Scenarios**:

1. **Given** a user clicks on a day or "+ Record Workout", **When** the modal opens, **Then** planned exercises are pre-populated with Plan values shown as reference
2. **Given** the user enters Sets, Reps, Load, and RPE for each exercise, **When** they click Save, **Then** the workout is recorded and visible in Week View
3. **Given** the user wants to add an unplanned exercise, **When** they click "+ Add Exercise", **Then** they can add exercises not in the original plan
4. **Given** invalid input (e.g., missing required fields), **When** they try to save, **Then** validation errors are clearly displayed

---

### User Story 3 - Create/Edit Training Program (Priority: P1)

A user wants to create a training program with weeks, days, and exercises to establish their training plan.

**Why this priority**: Programs are the foundation of the Plan vs Actual concept. Without programs, there's nothing to compare against.

**Independent Test**: Can be fully tested by creating a new program with at least one week containing days with exercises, saving it, and verifying it appears in program selection.

**Acceptance Scenarios**:

1. **Given** a user navigates to Program Editor, **When** they enter program name and add weeks/days/exercises, **Then** the hierarchical structure is displayed correctly
2. **Given** the user is editing a week, **When** they add a day with exercises, **Then** exercises can be selected from autocomplete and configured with Sets/Reps/Load/RPE
3. **Given** the user saves the program, **When** returning to Week View, **Then** the program is available for selection
4. **Given** existing items, **When** the user clicks delete, **Then** a confirmation dialog appears before deletion

---

### User Story 4 - Configure Schedule (Priority: P2)

A user wants to assign calendar dates to their program days so the Week View shows dates alongside day names.

**Why this priority**: Schedule assignment enhances the Week View but is not essential for basic Plan vs Actual comparison. The system works without dates.

**Independent Test**: Can be fully tested by opening Schedule Settings, setting a start date, and verifying dates appear in Week View.

**Acceptance Scenarios**:

1. **Given** a user opens Schedule Settings, **When** they set a start date, **Then** dates are auto-generated for each program day
2. **Given** the user has options for "skip weekends" or "avoid consecutive days", **When** they enable these, **Then** the schedule adjusts accordingly
3. **Given** the schedule is saved, **When** viewing Week View, **Then** each day shows its assigned date

---

### User Story 5 - Select Program (Priority: P2)

A user wants to select which program to view in the Week View.

**Why this priority**: Program selection is needed when users have multiple programs, but a single-program scenario works without it.

**Independent Test**: Can be fully tested by viewing the program list and selecting a program, which then loads in Week View.

**Acceptance Scenarios**:

1. **Given** multiple programs exist, **When** the user opens program selection, **Then** they see a list of available programs
2. **Given** the user selects a program, **When** returning to Week View, **Then** the selected program's weeks are displayed

---

### User Story 6 - Export Data (Priority: P2)

A user wants to export their training data as CSV for external analysis or backup.

**Why this priority**: Export is a convenience feature that enhances utility but is not core to the Plan vs Actual experience.

**Independent Test**: Can be fully tested by clicking Export and verifying a CSV file is downloaded with correct columns.

**Acceptance Scenarios**:

1. **Given** the user clicks Export, **When** they select "Current Week", **Then** a CSV is downloaded containing that week's data
2. **Given** the exported CSV, **When** opened in a spreadsheet, **Then** columns include date, week, day, exercise, plan values, actual values, and diff

---

### Edge Cases

- What happens when the API server is unreachable? → Display error message with retry option
- How does the system handle expired or invalid authentication tokens? → Prompt user to re-enter their token
- What happens when a program has no scheduled dates? → Week View shows days without date labels; navigation works by week number
- What happens when a workout is recorded for a day not in any program? → Shows in "Unplanned Workouts" section at bottom of Week View
- What happens when an exercise in the plan no longer exists? → Display exercise name from plan_snapshot, show warning indicator

## Requirements *(mandatory)*

### Functional Requirements

#### Week View (Main Screen)
- **FR-001**: System MUST display one week at a time with collapsible/expandable days
- **FR-002**: System MUST show Plan, Actual, and Diff columns for each exercise
- **FR-003**: System MUST calculate diff based on Sets, Reps, and Load only (not RPE)
- **FR-004**: System MUST display status indicators: ✓ (match), ≠ (diff), ○ (pending), 📝 (unplanned)
- **FR-005**: System MUST support week navigation (Previous/Next)
- **FR-006**: System MUST display unplanned workouts in a separate section

#### Workout Input
- **FR-007**: System MUST open a modal for recording workout results
- **FR-008**: System MUST pre-populate planned exercises with Plan values as reference
- **FR-009**: System MUST allow input of Sets, Reps, Load, and RPE for each exercise
- **FR-010**: System MUST allow adding unplanned exercises to a workout
- **FR-011**: System MUST support session-level notes

#### Program Editor
- **FR-012**: System MUST support creating programs with hierarchical structure (Program → Weeks → Days → Exercises)
- **FR-013**: System MUST provide exercise autocomplete from existing exercise catalog
- **FR-014**: System MUST allow configuring exercise prescriptions (Type, Sets, Reps, Load, RPE)
- **FR-015**: System MUST support delete operations with confirmation dialogs
- **FR-016**: System MUST support collapsible Week/Day sections for navigation

#### Schedule Settings
- **FR-017**: System MUST allow setting a program start date
- **FR-018**: System MUST auto-generate schedule from start date
- **FR-019**: System MUST provide options to skip weekends and avoid consecutive days

#### Program Selection
- **FR-020**: System MUST display a list of available programs
- **FR-021**: System MUST allow selecting a program to view in Week View

#### Export
- **FR-022**: System MUST export data as CSV (UTF-8 with BOM for Excel compatibility)
- **FR-023**: System MUST include columns: date, week, day, day_name, exercise, entry_type, plan values, actual values, diff, notes

#### Authentication (PoC)
- **FR-024**: System MUST support Bearer token authentication with manual token input via settings page (no login UI)
- **FR-025**: System MUST persist token in local storage
- **FR-026**: System MUST request token only when an API operation is attempted (lazy authentication), not on initial app load
- **FR-027**: System MUST provide a settings page where users can enter/update their authentication token

### Key Entities

- **Program**: A training plan containing weeks, days, and exercise prescriptions. Has name, description, and hierarchical root_nodes
- **ProgramNode**: A node in the program tree representing Week, Day, or Exercise prescription with type, order, and optional exercise configuration
- **Workout**: A recorded training session linked to a Program Day via program_node_id, containing entries and metadata
- **WorkoutEntry**: A single exercise performance within a workout, storing actual Sets/Reps/Load/RPE and optional plan_snapshot
- **Exercise**: A catalog item representing a type of exercise with name, description, aliases, and muscle groups

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can view their weekly Plan vs Actual comparison within 3 seconds of page load
- **SC-002**: Users can record a full workout session (5 exercises) in under 3 minutes
- **SC-003**: Users can create a new program with 3 weeks × 3 days × 4 exercises in under 10 minutes
- **SC-004**: Diff indicators correctly identify differences between Plan and Actual for 100% of recorded workouts
- **SC-005**: Users can navigate between weeks with a single click
- **SC-006**: 100% of form submissions provide clear validation feedback for invalid inputs
- **SC-007**: Users can complete the token setup process in under 30 seconds
- **SC-008**: Exported CSV opens correctly in Excel with proper column formatting

## Clarifications

### Session 2026-01-30

- Q: When should the app request the authentication token from the user? → A: Lazy request (prompt only when an API operation is attempted)
- Q: How should authentication be handled in this PoC? → A: No login UI; user manually enters token in settings page (developer/test use)

## Assumptions

- The backend API (OPTel Workout API) is running and accessible at `http://localhost:8080/api/v1`
- Bearer token authentication is sufficient for PoC (no login UI needed)
- The existing web project structure (`web/` directory with Next.js, shadcn/ui, TanStack Query) will be extended
- Users will manually provide their authentication token (no OAuth flow in PoC)
- Video upload functionality is out of scope for this PoC
- Telemetry features are out of scope for this PoC
- The UI design in `docs/WEB_UI_DESIGN.md` is the authoritative reference for visual layout
