# Data Model: Web PoC

**Date**: 2026-01-30  
**Feature**: 007-web-poc

## Overview

Frontend TypeScript types aligned with the backend OpenAPI specification (`api/openapi/openapi.yaml`). These types are used throughout the web application for type safety.

## Core Entities (from API)

### Exercise

```typescript
interface Exercise {
  id: string;                    // UUID
  name: string;                  // 1-200 chars
  description?: string;          // max 2000 chars
  aliases?: string[];
  muscle_groups?: string[];
  load_increment?: number;       // > 0, kg
  created_at: string;            // ISO 8601
  updated_at: string;            // ISO 8601
}

interface ExerciseCreate {
  name: string;
  description?: string;
  aliases?: string[];
  muscle_groups?: string[];
  load_increment?: number;
}
```

### Program

```typescript
interface Program {
  id: string;                    // UUID
  name: string;                  // 1-200 chars
  description?: string;          // max 2000 chars
  root_nodes?: ProgramNode[];
  created_at: string;
  updated_at: string;
}

interface ProgramCreate {
  name: string;
  description?: string;
  root_nodes?: ProgramNodeCreate[];
}
```

### ProgramNode

Represents hierarchical structure: Week → Day → Exercise prescription

```typescript
interface ProgramNode {
  id: string;                    // UUID
  program_id: string;            // UUID
  parent_id?: string;            // UUID, null for root
  name: string;                  // 1-200 chars
  node_type: string;             // "week" | "day" | "exercise" (1-50 chars)
  order: number;                 // >= 0
  children?: ProgramNode[];
  
  // Exercise prescription fields (when node_type === "exercise")
  exercise_id?: string;          // UUID
  target_sets?: number;          // >= 1
  target_reps?: number;          // >= 1
  target_rpe?: number;           // 1-10
  percent_1rm?: number;          // 0-1 (0.75 = 75%)
  planned_rest_seconds?: number; // >= 0
  muscle_groups?: string[];
  notes?: string;                // max 2000 chars
}

interface ProgramNodeCreate {
  name: string;
  node_type: string;
  order: number;
  children?: ProgramNodeCreate[];
  exercise_id?: string;
  target_sets?: number;
  target_reps?: number;
  target_rpe?: number;
  percent_1rm?: number;
  planned_rest_seconds?: number;
  muscle_groups?: string[];
  notes?: string;
}
```

### Workout

```typescript
interface Workout {
  id: string;                    // UUID
  timestamp: string;             // ISO 8601
  session_start?: string;
  session_end?: string;
  body_weight_kg?: number;       // > 0
  fatigue_level?: number;        // 1-5
  sleep_hours?: number;          // 0-24
  condition_notes?: string;      // max 2000 chars
  program_node_id?: string;      // UUID, link to Day
  program_context?: string[];    // e.g., ["Cycle 1", "Week 3", "Day 2"]
  notes?: string;                // max 5000 chars
  entries: WorkoutEntry[];       // min 1 item
  created_at: string;
  updated_at: string;
}

interface WorkoutCreate {
  timestamp: string;
  session_start?: string;
  session_end?: string;
  body_weight_kg?: number;
  fatigue_level?: number;
  sleep_hours?: number;
  condition_notes?: string;
  program_node_id?: string;
  program_context?: string[];
  notes?: string;
  entries: WorkoutEntryCreate[]; // min 1 item
}
```

### WorkoutEntry

```typescript
type EntryType = 'top' | 'main' | 'backoff' | 'accessory';

interface WorkoutEntry {
  id: string;                    // UUID
  workout_id: string;            // UUID
  order: number;                 // >= 0
  exercise_id: string;           // UUID
  display_name?: string;         // max 200 chars
  entry_type: EntryType;
  sets: number;                  // >= 1
  reps: number;                  // >= 1
  load_kg: number;               // >= 0
  rpe: number;                   // 1-10
  entry_start?: string;
  entry_end?: string;
  planned_rest_seconds?: number;
  performed_rest_seconds?: number;
  per_set_rest_overrides?: number[];
  program_node_id?: string;      // UUID, link to prescription
  plan_snapshot?: PlanSnapshot;
  notes?: string;                // max 2000 chars
  video_object_key?: string;     // max 500 chars (out of scope for PoC)
}

interface WorkoutEntryCreate {
  exercise_id: string;
  display_name?: string;
  entry_type: EntryType;
  sets: number;
  reps: number;
  load_kg: number;
  rpe: number;
  entry_start?: string;
  entry_end?: string;
  planned_rest_seconds?: number;
  performed_rest_seconds?: number;
  per_set_rest_overrides?: number[];
  program_node_id?: string;
  plan_snapshot?: PlanSnapshot;
  notes?: string;
}
```

### PlanSnapshot

Snapshot of planned values at execution time (immutable record of what was planned)

```typescript
interface PlanSnapshot {
  program_node_id?: string;      // UUID
  target_sets?: number;          // >= 1
  target_reps?: number;          // >= 1
  target_rpe?: number;           // 1-10
  target_load_kg?: number;       // >= 0
  percent_1rm?: number;          // 0-1
  planned_rest_seconds?: number; // >= 0
}
```

## API Response Types

### Paginated List Responses

```typescript
interface PaginatedResponse<T> {
  data: T[];
  next_cursor?: string | null;
  has_more: boolean;
}

type ExerciseListResponse = PaginatedResponse<Exercise>;
type WorkoutListResponse = PaginatedResponse<Workout>;
type ProgramListResponse = PaginatedResponse<Program>;
```

### Error Response

```typescript
interface ApiError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}
```

## Frontend-Only Types

### Week View State

```typescript
interface WeekViewState {
  selectedProgramId: string | null;
  currentWeekIndex: number;       // 0-based index into program weeks
  expandedDays: Set<string>;      // Set of day node IDs
}
```

### Day Status (for status indicators)

```typescript
type DayStatus = 'match' | 'diff' | 'pending' | 'unplanned';

interface DayWithStatus {
  day: ProgramNode;               // The day node from program
  workout: Workout | null;        // Recorded workout if exists
  status: DayStatus;
  scheduledDate?: Date;           // From schedule settings
}
```

### Schedule (client-side only for PoC)

```typescript
interface ScheduleSettings {
  programId: string;
  startDate: string;              // ISO 8601 date
  skipWeekends: boolean;
  avoidConsecutive: boolean;
}

interface DaySchedule {
  dayNodeId: string;              // ProgramNode ID
  weekIndex: number;
  dayIndex: number;
  date: string | null;            // ISO 8601 date, null if not scheduled
}
```

### Export Row

```typescript
interface ExportRow {
  date: string;
  week: number;
  day: number;
  day_name: string;
  exercise: string;
  entry_type: EntryType | '';
  plan_sets: number | '';
  plan_reps: number | '';
  plan_load_kg: number | '';
  actual_sets: number | '';
  actual_reps: number | '';
  actual_load_kg: number | '';
  actual_rpe: number | '';
  diff: string;
  notes: string;
}
```

## Validation Rules

### Form Schemas (Zod)

```typescript
// schemas/forms.ts
import { z } from 'zod';

export const exerciseCreateSchema = z.object({
  name: z.string().min(1).max(200),
  description: z.string().max(2000).optional(),
  aliases: z.array(z.string()).optional(),
  muscle_groups: z.array(z.string()).optional(),
  load_increment: z.number().positive().optional(),
});

export const workoutEntrySchema = z.object({
  exercise_id: z.string().uuid(),
  entry_type: z.enum(['top', 'main', 'backoff', 'accessory']),
  sets: z.number().int().min(1),
  reps: z.number().int().min(1),
  load_kg: z.number().min(0),
  rpe: z.number().int().min(1).max(10),
  notes: z.string().max(2000).optional(),
});

export const workoutCreateSchema = z.object({
  timestamp: z.string().datetime(),
  program_node_id: z.string().uuid().optional(),
  notes: z.string().max(5000).optional(),
  entries: z.array(workoutEntrySchema).min(1),
});

export const programNodeSchema: z.ZodType<ProgramNodeCreate> = z.lazy(() =>
  z.object({
    name: z.string().min(1).max(200),
    node_type: z.string().min(1).max(50),
    order: z.number().int().min(0),
    children: z.array(programNodeSchema).optional(),
    exercise_id: z.string().uuid().optional(),
    target_sets: z.number().int().min(1).optional(),
    target_reps: z.number().int().min(1).optional(),
    target_rpe: z.number().int().min(1).max(10).optional(),
    percent_1rm: z.number().min(0).max(1).optional(),
    planned_rest_seconds: z.number().int().min(0).optional(),
    notes: z.string().max(2000).optional(),
  })
);

export const programCreateSchema = z.object({
  name: z.string().min(1).max(200),
  description: z.string().max(2000).optional(),
  root_nodes: z.array(programNodeSchema).optional(),
});
```

## Entity Relationships

```
Program
  └── root_nodes: ProgramNode[] (type: "week")
        └── children: ProgramNode[] (type: "day")
              └── children: ProgramNode[] (type: "exercise")
                    └── exercise_id → Exercise

Workout
  └── program_node_id → ProgramNode (type: "day")
  └── entries: WorkoutEntry[]
        └── exercise_id → Exercise
        └── program_node_id → ProgramNode (type: "exercise")
        └── plan_snapshot: PlanSnapshot (snapshot of prescription at record time)
```
