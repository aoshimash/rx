import type { Log, LogCreate, LogListResponse } from '@/types/api';
import { api } from './client';

interface PaginationParams {
  limit?: number;
  after?: string;
}

interface LogFilterParams extends PaginationParams {
  performed_at_from?: string;
  performed_at_to?: string;
}

export const logsApi = {
  async list(params?: LogFilterParams): Promise<LogListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);
    if (params?.performed_at_from) searchParams.set('performed_at_from', params.performed_at_from);
    if (params?.performed_at_to) searchParams.set('performed_at_to', params.performed_at_to);

    return api.get(`logs?${searchParams}`).json<LogListResponse>();
  },

  async get(id: string): Promise<Log> {
    return api.get(`logs/${id}`).json<Log>();
  },

  async create(data: LogCreate): Promise<Log> {
    return api.post('logs', { json: data }).json<Log>();
  },

  async update(id: string, data: LogCreate): Promise<Log> {
    return api.put(`logs/${id}`, { json: data }).json<Log>();
  },

  async delete(id: string): Promise<void> {
    await api.delete(`logs/${id}`);
  },
};
