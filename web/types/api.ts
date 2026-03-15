/**
 * TypeScript Types for Rx API
 *
 * These types are aligned with the OpenAPI specification at api/openapi/openapi.yaml.
 * They should be kept in sync with any API changes.
 */

// ============================================================================
// Exercise
// ============================================================================

export interface Exercise {
  id: string;
  name: string;
  description?: string;
  aliases?: string[];
  muscle_groups?: string[];
  load_increment?: number;
  created_at: string;
  updated_at: string;
}

export interface ExerciseCreate {
  name: string;
  description?: string;
  aliases?: string[];
  muscle_groups?: string[];
  load_increment?: number;
}

// ============================================================================
// Program
// ============================================================================

export interface Program {
  id: string;
  name: string;
  description?: string;
  entries?: ProgramEntry[];
  created_at: string;
  updated_at: string;
}

export interface ProgramCreate {
  name: string;
  description?: string;
  entries?: ProgramEntryCreate[];
}

// ============================================================================
// ProgramEntry
// ============================================================================

export interface ProgramEntry {
  id: string;
  program_id: string;
  name: string;
  order: number;
  metadata?: Record<string, unknown>;
  exercise_id?: string;
  target_sets?: number;
  target_reps?: number;
  target_rpe?: number;
  percent_1rm?: number;
  planned_rest_seconds?: number;
  muscle_groups?: string[];
  notes?: string;
}

export interface ProgramEntryCreate {
  name: string;
  order: number;
  metadata?: Record<string, unknown>;
  exercise_id?: string;
  target_sets?: number;
  target_reps?: number;
  target_rpe?: number;
  percent_1rm?: number;
  planned_rest_seconds?: number;
  muscle_groups?: string[];
  notes?: string;
}

// ============================================================================
// Workout
// ============================================================================

export interface Workout {
  id: string;
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
  entries: WorkoutEntry[];
  created_at: string;
  updated_at: string;
}

export interface WorkoutCreate {
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
  entries: WorkoutEntryCreate[];
}

// ============================================================================
// WorkoutEntry
// ============================================================================

/**
 * Entry type for workout entries.
 * User-defined values are allowed, common values include: top, main, backoff, accessory.
 */
export type EntryType = string | null;

export interface WorkoutEntry {
  id: string;
  workout_id: string;
  order: number;
  exercise_id: string;
  display_name?: string;
  entry_type?: EntryType;
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
  video_object_key?: string;
}

export interface WorkoutEntryCreate {
  exercise_id: string;
  display_name?: string;
  entry_type?: EntryType;
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

// ============================================================================
// PlanSnapshot
// ============================================================================

export interface PlanSnapshot {
  program_node_id?: string;
  target_sets?: number;
  target_reps?: number;
  target_rpe?: number;
  target_load_kg?: number;
  percent_1rm?: number;
  planned_rest_seconds?: number;
}

// ============================================================================
// Pagination
// ============================================================================

export interface PaginatedResponse<T> {
  data: T[];
  next_cursor?: string | null;
  has_more: boolean;
}

export type ExerciseListResponse = PaginatedResponse<Exercise>;
export type WorkoutListResponse = PaginatedResponse<Workout>;
export type ProgramListResponse = PaginatedResponse<Program>;

// ============================================================================
// Error
// ============================================================================

export interface ApiErrorResponse {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}
