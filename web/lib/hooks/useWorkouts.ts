import { workoutsApi } from '@/lib/api/workouts';
import type { WorkoutCreate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

interface WorkoutFilters {
  timestamp_from?: string;
  timestamp_to?: string;
}

/**
 * Fetch workouts with optional date filters
 */
export function useWorkouts(filters?: WorkoutFilters) {
  return useQuery({
    queryKey: ['workouts', filters],
    queryFn: () => workoutsApi.list({ ...filters, limit: 100 }),
  });
}

/**
 * Fetch single workout by ID
 */
export function useWorkout(id: string | null) {
  return useQuery({
    queryKey: ['workouts', id],
    queryFn: () => workoutsApi.get(id!),
    enabled: !!id,
  });
}

/**
 * Create new workout
 */
export function useCreateWorkout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: WorkoutCreate) => workoutsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workouts'] });
    },
  });
}

/**
 * Update existing workout
 */
export function useUpdateWorkout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: WorkoutCreate }) => workoutsApi.update(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['workouts'] });
      queryClient.invalidateQueries({ queryKey: ['workouts', variables.id] });
    },
  });
}

/**
 * Delete workout
 */
export function useDeleteWorkout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => workoutsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workouts'] });
    },
  });
}
