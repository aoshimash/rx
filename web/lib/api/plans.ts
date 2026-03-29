import type { Plan, PlanCreate, PlanSessionCreate, PlanUpdate } from '@/types/api';
import { api } from './client';

export const plansApi = {
  async get(): Promise<Plan> {
    return api.get('plan').json<Plan>();
  },

  async create(data: PlanCreate): Promise<Plan> {
    return api.post('plan', { json: data }).json<Plan>();
  },

  async update(data: PlanUpdate): Promise<Plan> {
    return api.put('plan', { json: data }).json<Plan>();
  },

  async delete(): Promise<void> {
    await api.delete('plan');
  },

  async addSessions(sessions: PlanSessionCreate[]): Promise<Plan> {
    return api.post('plan/sessions', { json: { sessions } }).json<Plan>();
  },

  async deleteSession(sessionId: string): Promise<void> {
    await api.delete(`plan/sessions/${sessionId}`);
  },

  async expandProgram(programId: string): Promise<Plan> {
    return api.post(`plan/expand-program/${programId}`).json<Plan>();
  },
};
