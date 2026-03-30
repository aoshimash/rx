import type { FieldGroup, FieldGroupCreate, FieldGroupUpdate } from '@/types/api';
import { api } from './client';

export const fieldGroupsApi = {
  async list(): Promise<{ data: FieldGroup[] }> {
    return api.get('field-groups').json<{ data: FieldGroup[] }>();
  },

  async get(id: string): Promise<FieldGroup> {
    return api.get(`field-groups/${id}`).json<FieldGroup>();
  },

  async create(data: FieldGroupCreate): Promise<FieldGroup> {
    return api.post('field-groups', { json: data }).json<FieldGroup>();
  },

  async update(id: string, data: FieldGroupUpdate): Promise<FieldGroup> {
    return api.put(`field-groups/${id}`, { json: data }).json<FieldGroup>();
  },

  async delete(id: string): Promise<void> {
    await api.delete(`field-groups/${id}`);
  },
};
