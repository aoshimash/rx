import { api } from './client';
import type {
  Program,
  ProgramCreate,
  ProgramListResponse,
} from '@/types/api';

interface PaginationParams {
  limit?: number;
  after?: string;
}

export const programsApi = {
  /**
   * List programs with pagination
   * GET /programs
   */
  async list(params?: PaginationParams): Promise<ProgramListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);

    return api.get(`programs?${searchParams}`).json<ProgramListResponse>();
  },

  /**
   * Get program by ID (includes full node tree)
   * GET /programs/{id}
   */
  async get(id: string): Promise<Program> {
    return api.get(`programs/${id}`).json<Program>();
  },

  /**
   * Create new program
   * POST /programs
   */
  async create(data: ProgramCreate): Promise<Program> {
    return api.post('programs', { json: data }).json<Program>();
  },

  /**
   * Update existing program
   * PUT /programs/{id}
   */
  async update(id: string, data: ProgramCreate): Promise<Program> {
    return api.put(`programs/${id}`, { json: data }).json<Program>();
  },

  /**
   * Delete program
   * DELETE /programs/{id}
   */
  async delete(id: string): Promise<void> {
    await api.delete(`programs/${id}`);
  },
};
