import type {
  ConvertProgramToPlanRequest,
  ConvertProgramToPlansResponse,
  Program,
  ProgramCreate,
  ProgramListResponse,
} from '@/types/api';
import { api } from './client';

interface PaginationParams {
  limit?: number;
  after?: string;
  includeArchived?: boolean;
}

export const programsApi = {
  async list(params?: PaginationParams): Promise<ProgramListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);
    if (params?.includeArchived) searchParams.set('include_archived', 'true');

    return api.get(`programs?${searchParams}`).json<ProgramListResponse>();
  },

  async get(id: string): Promise<Program> {
    return api.get(`programs/${id}`).json<Program>();
  },

  async create(data: ProgramCreate): Promise<Program> {
    return api.post('programs', { json: data }).json<Program>();
  },

  async delete(id: string): Promise<void> {
    await api.delete(`programs/${id}`);
  },

  async archive(id: string): Promise<Program> {
    return api.post(`programs/${id}/archive`).json<Program>();
  },

  async unarchive(id: string): Promise<Program> {
    return api.post(`programs/${id}/unarchive`).json<Program>();
  },

  async duplicate(id: string): Promise<Program> {
    return api.post(`programs/${id}/duplicate`).json<Program>();
  },

  async convertToPlans(data: ConvertProgramToPlanRequest): Promise<ConvertProgramToPlansResponse> {
    return api.post('plans/from-program', { json: data }).json<ConvertProgramToPlansResponse>();
  },
};
