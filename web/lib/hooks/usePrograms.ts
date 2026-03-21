import { programsApi } from '@/lib/api/programs';
import type { ProgramCreate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export function usePrograms(status?: string) {
  return useQuery({
    queryKey: ['programs', { status }],
    queryFn: () => programsApi.list({ limit: 100, status }),
  });
}

export function useProgramsByTemplateId(programTemplateId: string) {
  return useQuery({
    queryKey: ['programs', { program_template_id: programTemplateId }],
    queryFn: () => programsApi.list({ limit: 100, program_template_id: programTemplateId }),
  });
}

export function useProgram(id: string | null) {
  return useQuery({
    queryKey: ['programs', id],
    queryFn: () => programsApi.get(id!),
    enabled: !!id,
  });
}

export function useCreateProgram() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ProgramCreate) => programsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
    },
  });
}

export function useLoggedSessions(programId: string | null) {
  return useQuery({
    queryKey: ['programs', programId, 'logged-sessions'],
    queryFn: () => programsApi.getLoggedSessions(programId!),
    enabled: !!programId,
  });
}

export function useDeleteProgram() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => programsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
    },
  });
}
