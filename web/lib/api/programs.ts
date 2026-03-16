import type {
  ConvertProgramToPlanRequest,
  Plan,
  Program,
  ProgramCreate,
  ProgramListResponse,
} from '@/types/api';
import { api } from './client';

interface PaginationParams {
  limit?: number;
  after?: string;
}

export const programsApi = {
  async list(params?: PaginationParams): Promise<ProgramListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);

    return api.get(`programs?${searchParams}`).json<ProgramListResponse>();
  },

  async get(id: string): Promise<Program> {
    return api.get(`programs/${id}`).json<Program>();
  },

  async create(data: ProgramCreate): Promise<Program> {
    return api.post('programs', { json: data }).json<Program>();
  },

  async update(id: string, data: ProgramCreate): Promise<Program> {
    return api.put(`programs/${id}`, { json: data }).json<Program>();
  },

  async delete(id: string): Promise<void> {
    await api.delete(`programs/${id}`);
  },

  async convertToPlan(data: ConvertProgramToPlanRequest): Promise<Plan> {
    return api.post('plans/from-program', { json: data }).json<Plan>();
  },
};
