import { programsApi } from '@/lib/api/programs';
import type { ProgramCreate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export function usePrograms(status?: string) {
  return useQuery({
    queryKey: ['programs', { status }],
    queryFn: () => programsApi.list({ limit: 100, status }),
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

export function useDeleteProgram() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => programsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
    },
  });
}
