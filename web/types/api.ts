/**
 * TypeScript Types for Rx API
 *
 * These types are aligned with the OpenAPI specification at api/openapi/openapi.yaml.
 * They should be kept in sync with any API changes.
 */

// ============================================================================
// FieldDef
// ============================================================================

export interface FieldDef {
  name: string;
  type: 'text' | 'number' | 'select';
  options?: string[];
  description?: string;
}

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
  fields?: Record<string, unknown>;
  notes?: string;
  video_object_key?: string;
  started_at?: string;
  finished_at?: string;
}

export interface LogEntryCreate {
  exercise_name: string;
  fields?: Record<string, unknown>;
  notes?: string;
  video_object_key?: string;
  started_at?: string;
  finished_at?: string;
}

// ============================================================================
// Program
// ============================================================================

export type ProgramStatus = 'created' | 'ongoing' | 'completed' | 'cancelled';

export interface Program {
  id: string;
  name: string;
  status: ProgramStatus;
  notes?: string;
  metadata?: Record<string, unknown>;
  program_fields?: FieldDef[];
  log_fields?: FieldDef[];
  sessions: ProgramSession[];
  created_at: string;
  updated_at: string;
}

export interface ProgramCreate {
  name: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  program_fields?: FieldDef[];
  log_fields?: FieldDef[];
  sessions?: ProgramSessionCreate[];
}

export interface ProgramUpdate {
  name: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  program_fields?: FieldDef[];
  log_fields?: FieldDef[];
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
  fields?: Record<string, unknown>;
  notes?: string;
}

export interface ProgramSessionEntryCreate {
  exercise_name: string;
  order: number;
  fields?: Record<string, unknown>;
  notes?: string;
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
export type ProgramListResponse = PaginatedResponse<Program>;

// ============================================================================
// Error
// ============================================================================

export interface ApiErrorResponse {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}
