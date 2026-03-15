import { logsApi } from '@/lib/api/logs';
import type { LogCreate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

interface LogFilters {
  performed_at_from?: string;
  performed_at_to?: string;
}

export function useLogs(filters?: LogFilters) {
  return useQuery({
    queryKey: ['logs', filters],
    queryFn: () => logsApi.list({ ...filters, limit: 100 }),
  });
}

export function useLog(id: string | null) {
  return useQuery({
    queryKey: ['logs', id],
    queryFn: () => logsApi.get(id!),
    enabled: !!id,
  });
}

export function useCreateLog() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: LogCreate) => logsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['logs'] });
    },
  });
}

export function useUpdateLog() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: LogCreate }) => logsApi.update(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['logs'] });
      queryClient.invalidateQueries({ queryKey: ['logs', variables.id] });
    },
  });
}

export function useDeleteLog() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => logsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['logs'] });
    },
  });
}
