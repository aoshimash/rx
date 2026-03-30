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
// FieldGroup
// ============================================================================

export interface FieldGroup {
  id: string;
  name: string;
  description?: string;
  program_fields: FieldDef[];
  log_fields: FieldDef[];
  created_at: string;
  updated_at: string;
}

export interface FieldGroupCreate {
  name: string;
  description?: string;
  program_fields: FieldDef[];
  log_fields: FieldDef[];
}

export interface FieldGroupUpdate {
  name: string;
  description?: string;
  program_fields: FieldDef[];
  log_fields: FieldDef[];
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
  plan_snapshot?: Record<string, unknown>;
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
  plan_snapshot?: Record<string, unknown>;
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
  sets?: LogSet[];
  notes?: string;
  started_at?: string;
  finished_at?: string;
}

export interface LogEntryCreate {
  exercise_name: string;
  fields?: Record<string, unknown>;
  sets?: LogSetCreate[];
  notes?: string;
  started_at?: string;
  finished_at?: string;
}

// ============================================================================
// LogSet
// ============================================================================

export interface LogSet {
  id: string;
  entry_id: string;
  set_number: number;
  fields: Record<string, unknown>;
  video_object_key?: string;
  notes?: string;
}

export interface LogSetCreate {
  set_number: number;
  fields: Record<string, unknown>;
  video_object_key?: string;
  notes?: string;
}

// ============================================================================
// Program
// ============================================================================

export interface Program {
  id: string;
  name: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  sessions: ProgramSession[];
  created_at: string;
  updated_at: string;
}

export interface ProgramCreate {
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
  field_group_id?: string;
  session_name: string;
  order: number;
  date?: string;
  entries: ProgramSessionEntry[];
}

export interface ProgramSessionCreate {
  session_name: string;
  order: number;
  field_group_id?: string;
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
// Plan
// ============================================================================

export interface Plan {
  id: string;
  name?: string;
  notes?: string;
  sessions: PlanSession[];
  created_at: string;
  updated_at: string;
}

export interface PlanCreate {
  name?: string;
  notes?: string;
  sessions?: PlanSessionCreate[];
}

export interface PlanUpdate {
  name?: string;
  notes?: string;
  sessions?: PlanSessionCreate[];
}

export interface PlanSession {
  id: string;
  plan_id: string;
  field_group_id?: string;
  session_name: string;
  order: number;
  date?: string;
  source_program_id?: string;
  source_session_id?: string;
  entries: PlanSessionEntry[];
}

export interface PlanSessionCreate {
  session_name: string;
  order: number;
  field_group_id?: string;
  date?: string;
  source_program_id?: string;
  source_session_id?: string;
  entries?: PlanSessionEntryCreate[];
}

export interface PlanSessionEntry {
  id: string;
  session_id: string;
  order: number;
  exercise_name: string;
  fields?: Record<string, unknown>;
  notes?: string;
}

export interface PlanSessionEntryCreate {
  exercise_name: string;
  order: number;
  fields?: Record<string, unknown>;
  notes?: string;
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
