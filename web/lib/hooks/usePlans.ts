import { plansApi } from '@/lib/api/plans';
import type { PlanCreate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export function usePlans() {
  return useQuery({
    queryKey: ['plans'],
    queryFn: () => plansApi.list({ limit: 100 }),
  });
}

export function usePlan(id: string | null) {
  return useQuery({
    queryKey: ['plans', id],
    queryFn: () => plansApi.get(id!),
    enabled: !!id,
  });
}

export function useCreatePlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: PlanCreate) => plansApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plans'] });
    },
  });
}

export function useUpdatePlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: PlanCreate }) => plansApi.update(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['plans'] });
      queryClient.invalidateQueries({ queryKey: ['plans', variables.id] });
    },
  });
}

export function useDeletePlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => plansApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plans'] });
    },
  });
}
