import type { Workout, WorkoutCreate, WorkoutListResponse } from '@/types/api';
import { api } from './client';

interface PaginationParams {
  limit?: number;
  after?: string;
}

interface WorkoutFilterParams extends PaginationParams {
  timestamp_from?: string;
  timestamp_to?: string;
}

export const workoutsApi = {
  /**
   * List workouts with pagination and optional date filters
   * GET /workouts
   */
  async list(params?: WorkoutFilterParams): Promise<WorkoutListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);
    if (params?.timestamp_from) searchParams.set('timestamp_from', params.timestamp_from);
    if (params?.timestamp_to) searchParams.set('timestamp_to', params.timestamp_to);

    return api.get(`workouts?${searchParams}`).json<WorkoutListResponse>();
  },

  /**
   * Get workout by ID (includes all entries)
   * GET /workouts/{id}
   */
  async get(id: string): Promise<Workout> {
    return api.get(`workouts/${id}`).json<Workout>();
  },

  /**
   * Create new workout
   * POST /workouts
   */
  async create(data: WorkoutCreate): Promise<Workout> {
    return api.post('workouts', { json: data }).json<Workout>();
  },

  /**
   * Update existing workout
   * PUT /workouts/{id}
   */
  async update(id: string, data: WorkoutCreate): Promise<Workout> {
    return api.put(`workouts/${id}`, { json: data }).json<Workout>();
  },

  /**
   * Delete workout
   * DELETE /workouts/{id}
   */
  async delete(id: string): Promise<void> {
    await api.delete(`workouts/${id}`);
  },
};
