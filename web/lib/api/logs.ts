import type { Log, LogCreate, LogListResponse } from '@/types/api';
import { api } from './client';

interface PaginationParams {
  limit?: number;
  after?: string;
}

interface LogFilterParams extends PaginationParams {
  program_id?: string;
  performed_at_from?: string;
  performed_at_to?: string;
}

export const logsApi = {
  async list(params?: LogFilterParams): Promise<LogListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);
    if (params?.program_id) searchParams.set('program_id', params.program_id);
    if (params?.performed_at_from) searchParams.set('performed_at_from', params.performed_at_from);
    if (params?.performed_at_to) searchParams.set('performed_at_to', params.performed_at_to);

    return api.get(`logs?${searchParams}`).json<LogListResponse>();
  },

  async get(id: string): Promise<Log> {
    const res = await api.get(`logs/${id}`).json<{ log: Log }>();
    return res.log;
  },

  async create(data: LogCreate): Promise<Log> {
    const res = await api.post('logs', { json: data }).json<{ log: Log }>();
    return res.log;
  },

  async update(id: string, data: LogCreate): Promise<Log> {
    const res = await api.put(`logs/${id}`, { json: data }).json<{ log: Log }>();
    return res.log;
  },

  async delete(id: string): Promise<void> {
    await api.delete(`logs/${id}`);
  },
};
