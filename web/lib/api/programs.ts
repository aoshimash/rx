import type { Program, ProgramCreate, ProgramListResponse, ProgramUpdate } from '@/types/api';
import { api } from './client';

interface ProgramListParams {
  limit?: number;
  after?: string;
}

export const programsApi = {
  async list(params?: ProgramListParams): Promise<ProgramListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);

    return api.get(`programs?${searchParams}`).json<ProgramListResponse>();
  },

  async get(id: string): Promise<Program> {
    const res = await api.get(`programs/${id}`).json<{ program: Program }>();
    return res.program;
  },

  async create(data: ProgramCreate): Promise<Program> {
    const res = await api.post('programs', { json: data }).json<{ program: Program }>();
    return res.program;
  },

  async update(id: string, data: ProgramUpdate): Promise<Program> {
    const res = await api.put(`programs/${id}`, { json: data }).json<{ program: Program }>();
    return res.program;
  },

  async delete(id: string): Promise<void> {
    await api.delete(`programs/${id}`);
  },
};
