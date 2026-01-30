import { exercisesApi } from '@/lib/api/exercises';
import type { ExerciseCreate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

/**
 * Fetch all exercises
 */
export function useExercises() {
  return useQuery({
    queryKey: ['exercises'],
    queryFn: () => exercisesApi.list({ limit: 100 }),
  });
}

/**
 * Fetch single exercise by ID
 */
export function useExercise(id: string | null) {
  return useQuery({
    queryKey: ['exercises', id],
    queryFn: () => exercisesApi.get(id!),
    enabled: !!id,
  });
}

/**
 * Create new exercise
 */
export function useCreateExercise() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ExerciseCreate) => exercisesApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['exercises'] });
    },
  });
}

/**
 * Update existing exercise
 */
export function useUpdateExercise() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ExerciseCreate }) =>
      exercisesApi.update(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['exercises'] });
      queryClient.invalidateQueries({ queryKey: ['exercises', variables.id] });
    },
  });
}

/**
 * Delete exercise
 */
export function useDeleteExercise() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => exercisesApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['exercises'] });
    },
  });
}
