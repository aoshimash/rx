import type {
  Plan,
  PlanCreate,
  PlanSession,
  PlanSessionCreate,
  PlanSessionUpdate,
  PlanUpdate,
} from '@/types/api';
import { api } from './client';

export const plansApi = {
  async get(): Promise<Plan> {
    const res = await api.get('plan').json<{ plan: Plan }>();
    return res.plan;
  },

  async create(data: PlanCreate): Promise<Plan> {
    const res = await api.post('plan', { json: data }).json<{ plan: Plan }>();
    return res.plan;
  },

  async update(data: PlanUpdate): Promise<Plan> {
    const res = await api.put('plan', { json: data }).json<{ plan: Plan }>();
    return res.plan;
  },

  async delete(): Promise<void> {
    await api.delete('plan');
  },

  async addSessions(sessions: PlanSessionCreate[]): Promise<Plan> {
    const res = await api.post('plan/sessions', { json: { sessions } }).json<{ plan: Plan }>();
    return res.plan;
  },

  async updateSession(sessionId: string, data: PlanSessionUpdate): Promise<PlanSession> {
    const res = await api
      .put(`plan/sessions/${sessionId}`, { json: data })
      .json<{ session: PlanSession }>();
    return res.session;
  },

  async deleteSession(sessionId: string): Promise<void> {
    await api.delete(`plan/sessions/${sessionId}`);
  },

  async expandProgram(programId: string): Promise<Plan> {
    const res = await api.post(`plan/expand-program/${programId}`).json<{ plan: Plan }>();
    return res.plan;
  },
};
