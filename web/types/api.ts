/**
 * TypeScript Types for Rx API
 *
 * These types are aligned with the OpenAPI specification at api/openapi/openapi.yaml.
 * They should be kept in sync with any API changes.
 */

// ============================================================================
// Plan
// ============================================================================

export interface Plan {
  id: string;
  name: string;
  program_id?: string;
  date?: string;
  session_name?: string;
  description?: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  entries?: PlanEntry[];
  created_at: string;
  updated_at: string;
}

export interface PlanCreate {
  name: string;
  date?: string;
  session_name?: string;
  description?: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  entries?: PlanEntryCreate[];
}

// ============================================================================
// PlanEntry
// ============================================================================

export interface PlanEntry {
  id: string;
  plan_id: string;
  exercise_name: string;
  order: number;
  sets?: number;
  reps?: number;
  load_kg?: number;
  rpe?: number;
  notes?: string;
  metadata?: Record<string, unknown>;
}

export interface PlanEntryCreate {
  exercise_name: string;
  order: number;
  sets?: number;
  reps?: number;
  load_kg?: number;
  rpe?: number;
  notes?: string;
  metadata?: Record<string, unknown>;
}

// ============================================================================
// Log
// ============================================================================

export interface Log {
  id: string;
  plan_id?: string;
  performed_at: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  entries: LogEntry[];
  created_at: string;
  updated_at: string;
}

export interface LogCreate {
  plan_id?: string;
  performed_at: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  entries: LogEntryCreate[];
}

// ============================================================================
// LogEntry
// ============================================================================

export interface LogEntry {
  id: string;
  log_id: string;
  exercise_name: string;
  order: number;
  sets?: number;
  reps?: number;
  load_kg?: number;
  rpe?: number;
  notes?: string;
  video_object_key?: string;
  metadata?: Record<string, unknown>;
}

export interface LogEntryCreate {
  exercise_name: string;
  sets?: number;
  reps?: number;
  load_kg?: number;
  rpe?: number;
  notes?: string;
  video_object_key?: string;
  metadata?: Record<string, unknown>;
}

// ============================================================================
// Program
// ============================================================================

export interface Program {
  id: string;
  name: string;
  description?: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  entries?: ProgramEntry[];
  created_at: string;
  updated_at: string;
}

export interface ProgramCreate {
  name: string;
  description?: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  entries?: ProgramEntryCreate[];
}

// ============================================================================
// ProgramEntry
// ============================================================================

export interface ProgramEntry {
  id: string;
  program_id: string;
  order: number;
  exercise_name: string;
  sets?: number;
  reps?: number;
  rpe?: number;
  percent_1rm?: number;
  notes?: string;
  metadata?: Record<string, unknown>;
}

export interface ProgramEntryCreate {
  exercise_name: string;
  order: number;
  sets?: number;
  reps?: number;
  rpe?: number;
  percent_1rm?: number;
  notes?: string;
  metadata?: Record<string, unknown>;
}

export interface ConvertProgramToPlanRequest {
  program_id: string;
  name?: string;
  target_weights: Record<string, number>;
  load_increments?: Record<string, number>;
}

// ============================================================================
// Pagination
// ============================================================================

export interface PaginatedResponse<T> {
  data: T[];
  next_cursor?: string | null;
  has_more: boolean;
}

export type PlanListResponse = PaginatedResponse<Plan>;
export type LogListResponse = PaginatedResponse<Log>;
export type ProgramListResponse = PaginatedResponse<Program>;

// ============================================================================
// Error
// ============================================================================

export interface ApiErrorResponse {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}
