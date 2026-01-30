/**
 * API Client Interface for Web PoC
 * 
 * This file defines the interface for the API client that will be implemented
 * in web/lib/api/. The actual implementation uses ky HTTP client.
 * 
 * Base URL: http://localhost:8080/api/v1
 * Authentication: Bearer token in Authorization header
 */

import type {
  Exercise,
  ExerciseCreate,
  ExerciseListResponse,
  Program,
  ProgramCreate,
  ProgramListResponse,
  Workout,
  WorkoutCreate,
  WorkoutListResponse,
} from './types';

// ============================================================================
// Pagination Parameters
// ============================================================================

interface PaginationParams {
  limit?: number;  // 1-100, default 100
  after?: string;  // Cursor from previous response
}

interface WorkoutFilterParams extends PaginationParams {
  timestamp_from?: string;  // ISO 8601
  timestamp_to?: string;    // ISO 8601
}

// ============================================================================
// Exercise API
// ============================================================================

interface ExercisesApi {
  /**
   * List exercises with pagination
   * GET /exercises
   */
  list(params?: PaginationParams): Promise<ExerciseListResponse>;

  /**
   * Get exercise by ID
   * GET /exercises/{id}
   */
  get(id: string): Promise<Exercise>;

  /**
   * Create new exercise
   * POST /exercises
   */
  create(data: ExerciseCreate): Promise<Exercise>;

  /**
   * Update existing exercise
   * PUT /exercises/{id}
   */
  update(id: string, data: ExerciseCreate): Promise<Exercise>;

  /**
   * Delete exercise
   * DELETE /exercises/{id}
   * Returns 409 Conflict if referenced by workout entries
   */
  delete(id: string): Promise<void>;
}

// ============================================================================
// Program API
// ============================================================================

interface ProgramsApi {
  /**
   * List programs with pagination
   * GET /programs
   */
  list(params?: PaginationParams): Promise<ProgramListResponse>;

  /**
   * Get program by ID (includes full node tree)
   * GET /programs/{id}
   */
  get(id: string): Promise<Program>;

  /**
   * Create new program
   * POST /programs
   */
  create(data: ProgramCreate): Promise<Program>;

  /**
   * Update existing program
   * PUT /programs/{id}
   */
  update(id: string, data: ProgramCreate): Promise<Program>;

  /**
   * Delete program (nodes cascade deleted)
   * DELETE /programs/{id}
   * Returns 409 Conflict if referenced by workouts
   */
  delete(id: string): Promise<void>;
}

// ============================================================================
// Workout API
// ============================================================================

interface WorkoutsApi {
  /**
   * List workouts with pagination and optional date filters
   * GET /workouts
   */
  list(params?: WorkoutFilterParams): Promise<WorkoutListResponse>;

  /**
   * Get workout by ID (includes all entries)
   * GET /workouts/{id}
   */
  get(id: string): Promise<Workout>;

  /**
   * Create new workout
   * POST /workouts
   */
  create(data: WorkoutCreate): Promise<Workout>;

  /**
   * Update existing workout
   * PUT /workouts/{id}
   */
  update(id: string, data: WorkoutCreate): Promise<Workout>;

  /**
   * Delete workout (entries cascade deleted)
   * DELETE /workouts/{id}
   */
  delete(id: string): Promise<void>;
}

// ============================================================================
// Combined API Client
// ============================================================================

export interface ApiClient {
  exercises: ExercisesApi;
  programs: ProgramsApi;
  workouts: WorkoutsApi;
}

// ============================================================================
// Error Types
// ============================================================================

export class ApiError extends Error {
  constructor(
    public code: string,
    message: string,
    public details?: Record<string, unknown>,
    public status?: number
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export class UnauthorizedError extends ApiError {
  constructor(message = 'Authentication required') {
    super('UNAUTHORIZED', message, undefined, 401);
    this.name = 'UnauthorizedError';
  }
}

export class NotFoundError extends ApiError {
  constructor(message = 'Resource not found') {
    super('NOT_FOUND', message, undefined, 404);
    this.name = 'NotFoundError';
  }
}

export class ValidationError extends ApiError {
  constructor(message: string, details?: Record<string, unknown>) {
    super('VALIDATION_ERROR', message, details, 400);
    this.name = 'ValidationError';
  }
}

export class ConflictError extends ApiError {
  constructor(message: string, details?: Record<string, unknown>) {
    super('CONFLICT', message, details, 409);
    this.name = 'ConflictError';
  }
}
