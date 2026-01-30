import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { programsApi } from '@/lib/api/programs';
import type { ProgramCreate } from '@/types/api';

/**
 * Fetch all programs
 */
export function usePrograms() {
  return useQuery({
    queryKey: ['programs'],
    queryFn: () => programsApi.list({ limit: 100 }),
  });
}

/**
 * Fetch single program by ID (includes full node tree)
 */
export function useProgram(id: string | null) {
  return useQuery({
    queryKey: ['programs', id],
    queryFn: () => programsApi.get(id!),
    enabled: !!id,
  });
}

/**
 * Create new program
 */
export function useCreateProgram() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: ProgramCreate) => programsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
    },
  });
}

/**
 * Update existing program
 */
export function useUpdateProgram() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ProgramCreate }) =>
      programsApi.update(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
      queryClient.invalidateQueries({ queryKey: ['programs', variables.id] });
    },
  });
}

/**
 * Delete program
 */
export function useDeleteProgram() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (id: string) => programsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
    },
  });
}
