# Web UI Design

This document describes the UI design for OPTel Workout web application.

**Created**: 2026-01-30  
**Status**: Draft

## Overview

The web UI is designed for powerlifters who plan their training programs and track actual workout results. The primary goal is to display **planned vs actual** data side-by-side, making it easy to see differences (diffs) at a glance.

### Design Principles

1. **Plan-First Approach** - Users create a training program first, then record workouts against it
2. **Weekly View** - One page displays one week of training
3. **Diff Visibility** - Differences between plan and actual are clearly highlighted
4. **Minimal Navigation** - Use modals to reduce page transitions

## Screen Structure

| Screen | Type | Priority | Description |
|--------|------|----------|-------------|
| Week View | Page | P1 | Main screen showing Plan/Actual/Diff |
| Workout Input | Modal | P1 | Record actual workout results |
| Program Editor | Page | P1 | Create and edit training programs |
| Program List | Page/Modal | P2 | Select a program |
| Schedule Settings | Modal | P2 | Assign dates to program days |
| Export | Button → CSV | P2 | Export data as CSV |

## Week View (Main Screen)

The main screen displays one week at a time with collapsible days.

### Wireframe

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Week 3 of "FullBody 9weeks"                   [< Prev] [Next >]            │
│ Jan 27 - Feb 2, 2026                          [📅 Edit Schedule]           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ ▼ Day 1 - Lower Body                                Mon, Jan 27    ✓       │
│ ┌───────────┬──────────────────┬──────────────────┬────────────────────┐   │
│ │ Exercise  │ Plan             │ Actual           │ Diff               │   │
│ ├───────────┼──────────────────┼──────────────────┼────────────────────┤   │
│ │ Squat     │ 5×5 @ 100kg      │ 5×5 @ 100kg      │ —                  │   │
│ │ RDL       │ 3×10 @ 60kg      │ 3×10 @ 65kg      │ Load +5kg          │   │
│ │ Leg Curl  │ 3×12 @ 30kg      │ 3×12 @ 30kg      │ —                  │   │
│ └───────────┴──────────────────┴──────────────────┴────────────────────┘   │
│                                                                             │
│ ▶ Day 2 - Upper Body                                Tue, Jan 28    ≠       │
│                                                                             │
│ ▶ Day 3 - Full Body                                 Not scheduled  ○       │
│                                                                             │
│ ─────────────────────────────────────────────────────────────────────────── │
│ 📝 Unplanned Workouts (1)                                                   │
│   Wed, Jan 29 - Light cardio 30min                                         │
│                                                                             │
│                                              [+ Record Workout]             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Status Indicators

| Status | Icon | Color | Condition |
|--------|------|-------|-----------|
| Match | ✓ | Green | Plan exists, Actual matches (Sets/Reps/Load) |
| Diff | ≠ | Yellow | Plan exists, Actual differs |
| Pending | ○ | Gray | Plan exists, no Actual recorded |
| Unplanned | 📝 | Blue | No Plan, Actual only |

### Diff Calculation

Diff is determined by comparing **Sets**, **Reps**, and **Load** only.  
RPE differences are NOT considered as diff (plan is met even if RPE differs).

### Collapsible Days

- Days are collapsible/expandable by clicking the header
- Collapsed state shows: Day name, scheduled date, and status icon
- Expanded state shows: Full exercise table with Plan/Actual/Diff columns

## Workout Input (Modal)

Opens when clicking a Day or the "+ Record Workout" button.

### Wireframe

```
┌─────────────────────────────────────────────────────────────────┐
│ Day 1 - Lower Body                                    ✕ Close  │
│ Mon, Jan 27, 2026                                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│ ┌─────────┬───────────────┬───────────────────────────────────┐ │
│ │Exercise │ Plan          │ Actual                            │ │
│ ├─────────┼───────────────┼───────────────────────────────────┤ │
│ │ Squat   │ 5×5 @ 100kg   │ [5]×[5] @ [100  ]kg  RPE [7 ]    │ │
│ │ (top)   │ RPE 7         │                                   │ │
│ ├─────────┼───────────────┼───────────────────────────────────┤ │
│ │ RDL     │ 3×10 @ 60kg   │ [3]×[10] @ [65   ]kg  RPE [7 ]   │ │
│ │ (main)  │ RPE 7         │                                   │ │
│ ├─────────┼───────────────┼───────────────────────────────────┤ │
│ │ Leg Curl│ 3×12 @ 30kg   │ [ ]×[ ] @ [     ]kg  RPE [  ]    │ │
│ │ (accessory)             │                                   │ │
│ └─────────┴───────────────┴───────────────────────────────────┘ │
│                                                                 │
│ Notes: [                                                      ] │
│                                                                 │
│ [+ Add Exercise]                            [Cancel] [Save]    │
└─────────────────────────────────────────────────────────────────┘
```

### Features

- Pre-populated with planned exercises
- Plan values shown for reference (read-only)
- Input fields for Actual values (Sets, Reps, Load, RPE)
- Can add unplanned exercises
- Notes field for session-level comments

## Program Editor (Page)

Tree/form-based UI for creating and editing training programs.

### Wireframe

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Program Editor                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Program Name: [FullBody 9weeks                         ]                   │
│ Description:  [Full body training, 3 days per week     ]                   │
│                                                                             │
│ ─────────────────────────────────────────────────────────────────────────── │
│                                                                             │
│ ▼ Week 1                                                        [🗑️ Delete]│
│ │  Target RPE: [6  ]                                                       │
│ │                                                                           │
│ │  ▼ Day 1                                                      [🗑️ Delete]│
│ │  │  Name: [Lower Body        ]                                           │
│ │  │                                                                        │
│ │  │  ┌──────────┬──────┬──────┬────────┬─────┬────────┐                   │
│ │  │  │ Exercise │ Type │ Sets │ Reps   │ Load│ RPE    │                   │
│ │  │  ├──────────┼──────┼──────┼────────┼─────┼────────┤                   │
│ │  │  │ [Squat▼] │[top▼]│ [5 ] │ [5   ] │[100]│ [6   ] │ [🗑️]             │
│ │  │  │ [RDL  ▼] │[main]│ [3 ] │ [10  ] │[60 ]│ [6   ] │ [🗑️]             │
│ │  │  └──────────┴──────┴──────┴────────┴─────┴────────┘                   │
│ │  │  [+ Add Exercise]                                                      │
│ │  │                                                                        │
│ │  ▶ Day 2 - Upper Body (3 exercises)                                      │
│ │  ▶ Day 3 - Full Body (4 exercises)                                       │
│ │  [+ Add Day]                                                              │
│ │                                                                           │
│ ▶ Week 2 (Target RPE: 7)                                                   │
│ ▶ Week 3 (Target RPE: 7)                                                   │
│                                                                             │
│ [+ Add Week]                                                                │
│                                                                             │
│ ─────────────────────────────────────────────────────────────────────────── │
│                                                     [Cancel] [Save]        │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Hierarchy

```
Program
├── Week 1
│   ├── Day 1
│   │   ├── Exercise 1 (top)
│   │   ├── Exercise 2 (main)
│   │   └── Exercise 3 (accessory)
│   ├── Day 2
│   └── Day 3
├── Week 2
└── ...
```

### Operations

| Operation | Description |
|-----------|-------------|
| Add Week | Adds empty week at the end |
| Add Day | Adds empty day within a week |
| Add Exercise | Adds row to exercise table |
| Delete | Confirmation dialog, then remove |
| Collapse/Expand | Click Week/Day header to toggle |
| Drag & Drop | Reorder items (future enhancement) |

### Exercise Selection

- Autocomplete from Exercise catalog
- New exercises are added to catalog automatically

```
┌──────────────────────────────────────┐
│ Exercise: [Sq                      ] │
│           ┌────────────────────────┐ │
│           │ Squat                  │ │
│           │ Sumo Squat             │ │
│           │ Safety Bar Squat       │ │
│           └────────────────────────┘ │
└──────────────────────────────────────┘
```

## Schedule Settings (Modal)

Assign calendar dates to program days.

### Wireframe

```
┌──────────────────────────────────────────────┐
│ Schedule Settings                    ✕ Close │
├──────────────────────────────────────────────┤
│                                              │
│ Program Start Date: [2026-01-27 📅]          │
│                                              │
│ ☑ Skip weekends                             │
│ ☑ Avoid consecutive days                    │
│                                              │
│ Schedule:                                    │
│  Week 1 Day 1 → Mon, Jan 27                 │
│  Week 1 Day 2 → Wed, Jan 29  [Change 📅]    │
│  Week 1 Day 3 → Fri, Jan 31  [Change 📅]    │
│  Week 2 Day 1 → Mon, Feb 3                  │
│  ...                                         │
│                                              │
│ [Auto-generate]              [Cancel] [Save] │
└──────────────────────────────────────────────┘
```

### Features

- Auto-generate schedule from start date
- Options to skip weekends or avoid consecutive days
- Manual override for individual days
- Schedule is optional (program works without dates)

## Export

### Format

- **File Type**: CSV (UTF-8 with BOM for Excel compatibility)
- **Filename**: `{program_name}_{week}.csv` or `{program_name}_all.csv`

### CSV Columns

```csv
date,week,day,day_name,exercise,entry_type,plan_sets,plan_reps,plan_load_kg,actual_sets,actual_reps,actual_load_kg,actual_rpe,diff,notes
2026-01-27,3,1,Lower Body,Squat,top,5,5,100,5,5,100,7,,
2026-01-27,3,1,Lower Body,RDL,main,3,10,60,3,10,65,7,Load +5kg,Felt strong
2026-01-28,3,2,Upper Body,Bench,top,5,5,80,4,5,80,8,Sets -1,Fatigue
2026-01-29,3,,,(unplanned),Cardio,,,,1,1,,5,,Light cardio 30min
```

### Export Options

| Option | Description |
|--------|-------------|
| Current Week | Export currently displayed week |
| All Weeks | Export entire program |
| Date Range | Export specified date range |

## Component Structure

```
web/
├── app/
│   ├── page.tsx                    # Week View (main)
│   ├── programs/
│   │   ├── page.tsx                # Program List
│   │   ├── new/page.tsx            # Program Editor (create)
│   │   └── [id]/edit/page.tsx      # Program Editor (edit)
│   └── layout.tsx
├── components/
│   ├── ui/                         # shadcn/ui components
│   ├── week-view/
│   │   ├── WeekView.tsx
│   │   ├── DayAccordion.tsx
│   │   ├── ExerciseTable.tsx
│   │   └── StatusBadge.tsx
│   ├── workout-input/
│   │   ├── WorkoutModal.tsx
│   │   ├── ExerciseInputRow.tsx
│   │   └── AddExerciseButton.tsx
│   ├── program-editor/
│   │   ├── ProgramForm.tsx
│   │   ├── WeekAccordion.tsx
│   │   ├── DayAccordion.tsx
│   │   ├── ExerciseTable.tsx
│   │   ├── ExerciseRow.tsx
│   │   └── ExerciseCombobox.tsx
│   └── schedule/
│       └── ScheduleModal.tsx
└── lib/
    ├── api/                        # API client
    └── hooks/                      # TanStack Query hooks
```

## UI Components (shadcn/ui)

| Component | Usage |
|-----------|-------|
| Accordion | Collapsible Week/Day sections |
| Dialog | Modal windows |
| Form + Input | Form fields |
| Combobox | Exercise autocomplete |
| Select | Entry type selection |
| Table | Exercise listings |
| Button | Actions |
| AlertDialog | Delete confirmation |
| Calendar | Date picker |
| Badge | Status indicators |

## Data Flow

```
                    ┌─────────────┐
                    │   Program   │
                    │  (計画定義)  │
                    └──────┬──────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
         ▼                 ▼                 ▼
    ┌─────────┐       ┌─────────┐       ┌─────────┐
    │ Week 1  │       │ Week 2  │       │ Week 3  │
    │ Day1,2,3│       │ Day1,2,3│       │ Day1,2,3│
    └────┬────┘       └────┬────┘       └────┬────┘
         │                 │                 │
         ▼                 ▼                 ▼
    ┌─────────┐       ┌─────────┐       ┌─────────┐
    │Workouts │       │Workouts │       │Workouts │
    │  (実績)  │       │  (実績)  │       │  (実績)  │
    └─────────┘       └─────────┘       └─────────┘
```

### Workout ↔ Program Linkage

- `Workout.program_node_id` → References the Day's ProgramNode
- `WorkoutEntry.program_node_id` → References the Exercise prescription
- `WorkoutEntry.plan_snapshot` → Stores plan values at recording time

## Related Documents

- [FRONTEND_ARCHITECTURE.md](FRONTEND_ARCHITECTURE.md) - Technology stack and patterns
- [optel-domain](.claude/skills/optel-domain/) - Domain model definitions
- [specs/001-data-structure/](../specs/001-data-structure/) - Data structure specification
