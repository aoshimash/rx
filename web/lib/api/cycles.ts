import type { Cycle, CycleListResponse } from '@/types/api';
import { api } from './client';

interface CycleListParams {
  programId?: string;
  limit?: number;
  after?: string;
}

export const cyclesApi = {
  async list(params?: CycleListParams): Promise<CycleListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.programId) searchParams.set('program_id', params.programId);
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);

    return api.get(`cycles?${searchParams}`).json<CycleListResponse>();
  },

  async get(id: string): Promise<Cycle> {
    return api.get(`cycles/${id}`).json<Cycle>();
  },

  async delete(id: string): Promise<void> {
    await api.delete(`cycles/${id}`);
  },
};
