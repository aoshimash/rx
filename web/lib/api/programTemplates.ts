import type {
  GenerateProgramRequest,
  Program,
  ProgramTemplate,
  ProgramTemplateCreate,
  ProgramTemplateListResponse,
} from '@/types/api';
import { api } from './client';

interface PaginationParams {
  limit?: number;
  after?: string;
  includeArchived?: boolean;
}

export const programTemplatesApi = {
  async list(params?: PaginationParams): Promise<ProgramTemplateListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.after) searchParams.set('after', params.after);
    if (params?.includeArchived) searchParams.set('include_archived', 'true');

    return api.get(`program-templates?${searchParams}`).json<ProgramTemplateListResponse>();
  },

  async get(id: string): Promise<ProgramTemplate> {
    return api.get(`program-templates/${id}`).json<ProgramTemplate>();
  },

  async create(data: ProgramTemplateCreate): Promise<ProgramTemplate> {
    return api.post('program-templates', { json: data }).json<ProgramTemplate>();
  },

  async delete(id: string): Promise<void> {
    await api.delete(`program-templates/${id}`);
  },

  async archive(id: string): Promise<ProgramTemplate> {
    return api.post(`program-templates/${id}/archive`).json<ProgramTemplate>();
  },

  async unarchive(id: string): Promise<ProgramTemplate> {
    return api.post(`program-templates/${id}/unarchive`).json<ProgramTemplate>();
  },

  async duplicate(id: string, name?: string): Promise<ProgramTemplate> {
    const options = name ? { json: { name } } : undefined;
    return api.post(`program-templates/${id}/duplicate`, options).json<ProgramTemplate>();
  },

  async generate(id: string, data: GenerateProgramRequest): Promise<Program> {
    return api.post(`program-templates/${id}/generate`, { json: data }).json<Program>();
  },

  async edit(
    id: string,
    data: ProgramTemplateCreate
  ): Promise<{ template: ProgramTemplate; isNewVersion: boolean }> {
    const response = await api.post(`program-templates/${id}/edit`, { json: data });
    const template = await response.json<ProgramTemplate>();
    return { template, isNewVersion: response.status === 201 };
  },
};
