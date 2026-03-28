/**
 * TypeScript Types for Rx API
 *
 * These types are aligned with the OpenAPI specification at api/openapi/openapi.yaml.
 * They should be kept in sync with any API changes.
 */

// ============================================================================
// Log
// ============================================================================

export interface Log {
  id: string;
  program_id?: string;
  session_name?: string;
  performed_at: string;
  started_at?: string;
  finished_at?: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  entries: LogEntry[];
  created_at: string;
  updated_at: string;
}

export interface LogCreate {
  program_id?: string;
  session_name?: string;
  performed_at: string;
  started_at?: string;
  finished_at?: string;
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
  started_at?: string;
  finished_at?: string;
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
  started_at?: string;
  finished_at?: string;
  metadata?: Record<string, unknown>;
}

// ============================================================================
// ProgramTemplate
// ============================================================================

export interface ProgramTemplate {
  id: string;
  name: string;
  description?: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  entries?: ProgramTemplateEntry[];
  weeks?: string;
  days_per_week?: string;
  created_by?: string;
  archived_at?: string;
  created_at: string;
  updated_at: string;
}

export interface ProgramTemplateCreate {
  name: string;
  description?: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  weeks?: string;
  days_per_week?: string;
  entries?: ProgramTemplateEntryCreate[];
}

// ============================================================================
// ProgramTemplateEntry
// ============================================================================

export interface ProgramTemplateEntry {
  id: string;
  program_template_id: string;
  order: number;
  exercise_name: string;
  sets?: number;
  reps?: number;
  rpe?: number;
  percent_1rm?: number;
  notes?: string;
  metadata?: Record<string, unknown>;
}

export interface ProgramTemplateEntryCreate {
  exercise_name: string;
  order: number;
  sets?: number;
  reps?: number;
  rpe?: number;
  percent_1rm?: number;
  notes?: string;
  metadata?: Record<string, unknown>;
}

export interface GenerateProgramRequest {
  name?: string;
  target_weights: Record<string, number>;
  load_increments?: Record<string, number>;
}

// ============================================================================
// Program
// ============================================================================

export type ProgramStatus = 'created' | 'ongoing' | 'completed' | 'cancelled';

export interface Program {
  id: string;
  program_template_id?: string;
  name: string;
  status: ProgramStatus;
  notes?: string;
  metadata?: Record<string, unknown>;
  sessions: ProgramSession[];
  created_at: string;
  updated_at: string;
}

export interface ProgramCreate {
  program_template_id?: string;
  name: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  sessions?: ProgramSessionCreate[];
}

export interface ProgramUpdate {
  name: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  sessions: ProgramSessionCreate[];
}

// ============================================================================
// ProgramSession
// ============================================================================

export interface ProgramSession {
  id: string;
  program_id: string;
  session_name: string;
  order: number;
  date?: string;
  entries: ProgramSessionEntry[];
}

export interface ProgramSessionCreate {
  session_name: string;
  order: number;
  date?: string;
  entries?: ProgramSessionEntryCreate[];
}

// ============================================================================
// ProgramSessionEntry
// ============================================================================

export interface ProgramSessionEntry {
  id: string;
  session_id: string;
  order: number;
  exercise_name: string;
  sets?: number;
  reps?: number;
  load_kg?: number;
  rpe?: number;
  notes?: string;
  metadata?: Record<string, unknown>;
}

export interface ProgramSessionEntryCreate {
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
// LoggedSessions
// ============================================================================

export interface LoggedSession {
  session_name: string;
  log_id: string;
  performed_at: string;
  started_at?: string;
  finished_at?: string;
}

export interface LoggedSessionsResponse {
  sessions: LoggedSession[];
}

// ============================================================================
// Pagination
// ============================================================================

export interface PaginatedResponse<T> {
  data: T[];
  next_cursor?: string | null;
  has_more: boolean;
}

export type LogListResponse = PaginatedResponse<Log>;
export type ProgramTemplateListResponse = PaginatedResponse<ProgramTemplate>;
export type ProgramListResponse = PaginatedResponse<Program>;

// ============================================================================
// Error
// ============================================================================

export interface ApiErrorResponse {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}
