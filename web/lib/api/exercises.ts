import type { Exercise, ExerciseCreate, ExerciseListResponse } from '@/types/api';
import { api } from './client';

interface PaginationParams {
  limit?: number;
  after?: string;
}

export const exercisesApi = {
  /**
   * List exercises with pagination
   * GET /exercises
   */
  async list(params?: PaginationParams): Promise<ExerciseListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);

    return api.get(`exercises?${searchParams}`).json<ExerciseListResponse>();
  },

  /**
   * Get exercise by ID
   * GET /exercises/{id}
   */
  async get(id: string): Promise<Exercise> {
    return api.get(`exercises/${id}`).json<Exercise>();
  },

  /**
   * Create new exercise
   * POST /exercises
   */
  async create(data: ExerciseCreate): Promise<Exercise> {
    return api.post('exercises', { json: data }).json<Exercise>();
  },

  /**
   * Update existing exercise
   * PUT /exercises/{id}
   */
  async update(id: string, data: ExerciseCreate): Promise<Exercise> {
    return api.put(`exercises/${id}`, { json: data }).json<Exercise>();
  },

  /**
   * Delete exercise
   * DELETE /exercises/{id}
   */
  async delete(id: string): Promise<void> {
    await api.delete(`exercises/${id}`);
  },
};
