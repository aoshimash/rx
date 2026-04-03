import type { FieldGroup, FieldGroupCreate, FieldGroupUpdate } from '@/types/api';
import { api } from './client';

export const fieldGroupsApi = {
  async list(): Promise<{ data: FieldGroup[] }> {
    return api.get('field-groups').json<{ data: FieldGroup[] }>();
  },

  async get(id: string): Promise<FieldGroup> {
    const res = await api.get(`field-groups/${id}`).json<{ fieldGroup: FieldGroup }>();
    return res.fieldGroup;
  },

  async create(data: FieldGroupCreate): Promise<FieldGroup> {
    const res = await api.post('field-groups', { json: data }).json<{ fieldGroup: FieldGroup }>();
    return res.fieldGroup;
  },

  async update(id: string, data: FieldGroupUpdate): Promise<FieldGroup> {
    const res = await api
      .put(`field-groups/${id}`, { json: data })
      .json<{ fieldGroup: FieldGroup }>();
    return res.fieldGroup;
  },

  async delete(id: string): Promise<void> {
    await api.delete(`field-groups/${id}`);
  },
};
