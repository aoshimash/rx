import { programsApi } from '@/lib/api/programs';
import type { ProgramCreate, ProgramStatus, ProgramUpdate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export function usePrograms() {
  return useQuery({
    queryKey: ['programs'],
    queryFn: () => programsApi.list({ limit: 100 }),
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

export function useUpdateProgram() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ProgramUpdate }) => programsApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
    },
  });
}

export function useUpdateProgramStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: ProgramStatus }) =>
      programsApi.updateStatus(id, status),
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
