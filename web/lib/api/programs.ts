import type {
  LoggedSessionsResponse,
  Program,
  ProgramCreate,
  ProgramListResponse,
} from '@/types/api';
import { api } from './client';

interface ProgramListParams {
  limit?: number;
  after?: string;
  status?: string;
  program_template_id?: string;
}

export const programsApi = {
  async list(params?: ProgramListParams): Promise<ProgramListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);
    if (params?.status) searchParams.set('status', params.status);
    if (params?.program_template_id)
      searchParams.set('program_template_id', params.program_template_id);

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

  async getLoggedSessions(id: string): Promise<LoggedSessionsResponse> {
    return api.get(`programs/${id}/logged-sessions`).json<LoggedSessionsResponse>();
  },
};
