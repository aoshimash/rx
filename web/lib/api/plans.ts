import type { Plan, PlanCreate, PlanListResponse } from '@/types/api';
import { api } from './client';

interface PaginationParams {
  limit?: number;
  after?: string;
  cycle_id?: string;
}

export const plansApi = {
  async list(params?: PaginationParams): Promise<PlanListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);
    if (params?.cycle_id) searchParams.set('cycle_id', params.cycle_id);

    return api.get(`plans?${searchParams}`).json<PlanListResponse>();
  },

  async get(id: string): Promise<Plan> {
    return api.get(`plans/${id}`).json<Plan>();
  },

  async create(data: PlanCreate): Promise<Plan> {
    return api.post('plans', { json: data }).json<Plan>();
  },

  async update(id: string, data: PlanCreate): Promise<Plan> {
    return api.put(`plans/${id}`, { json: data }).json<Plan>();
  },

  async delete(id: string): Promise<void> {
    await api.delete(`plans/${id}`);
  },
};
