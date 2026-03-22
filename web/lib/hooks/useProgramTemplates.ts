import { programTemplatesApi } from '@/lib/api/programTemplates';
import type { GenerateProgramRequest, ProgramTemplateCreate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export function useProgramTemplates(includeArchived = false) {
  return useQuery({
    queryKey: ['program-templates', { includeArchived }],
    queryFn: () => programTemplatesApi.list({ limit: 100, includeArchived }),
  });
}

export function useProgramTemplate(id: string | null) {
  return useQuery({
    queryKey: ['program-templates', id],
    queryFn: () => programTemplatesApi.get(id!),
    enabled: !!id,
  });
}

export function useCreateProgramTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ProgramTemplateCreate) => programTemplatesApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['program-templates'] });
    },
  });
}

export function useDeleteProgramTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => programTemplatesApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['program-templates'] });
    },
  });
}

export function useArchiveProgramTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => programTemplatesApi.archive(id),
    onSuccess: (_data, id) => {
      queryClient.invalidateQueries({ queryKey: ['program-templates'] });
      queryClient.invalidateQueries({ queryKey: ['program-templates', id] });
    },
  });
}

export function useUnarchiveProgramTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => programTemplatesApi.unarchive(id),
    onSuccess: (_data, id) => {
      queryClient.invalidateQueries({ queryKey: ['program-templates'] });
      queryClient.invalidateQueries({ queryKey: ['program-templates', id] });
    },
  });
}

export function useDuplicateProgramTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, name }: { id: string; name?: string }) =>
      programTemplatesApi.duplicate(id, name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['program-templates'] });
    },
  });
}

export function useGenerateProgram() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: GenerateProgramRequest }) =>
      programTemplatesApi.generate(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
    },
  });
}
