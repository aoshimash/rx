import { programsApi } from '@/lib/api/programs';
import type { ConvertProgramToPlanRequest, ProgramCreate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export function usePrograms(includeArchived = false) {
  return useQuery({
    queryKey: ['programs', { includeArchived }],
    queryFn: () => programsApi.list({ limit: 100, includeArchived }),
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

export function useArchiveProgram() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => programsApi.archive(id),
    onSuccess: (_data, id) => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
      queryClient.invalidateQueries({ queryKey: ['programs', id] });
    },
  });
}

export function useUnarchiveProgram() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => programsApi.unarchive(id),
    onSuccess: (_data, id) => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
      queryClient.invalidateQueries({ queryKey: ['programs', id] });
    },
  });
}

export function useDuplicateProgram() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => programsApi.duplicate(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
    },
  });
}

export function useConvertProgramToPlans() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ConvertProgramToPlanRequest) => programsApi.convertToPlans(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plans'] });
      queryClient.invalidateQueries({ queryKey: ['cycles'] });
    },
  });
}
