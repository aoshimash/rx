import { programsApi } from '@/lib/api/programs';
import type { ProgramCreate, ProgramUpdate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { HTTPError } from 'ky';
import { toast } from 'sonner';

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

export function useUpdateProgram() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ProgramUpdate }) => programsApi.update(id, data),
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

const statusLabels: Record<string, string> = {
  ongoing: 'started',
  completed: 'completed',
  cancelled: 'cancelled',
};

export function useUpdateProgramStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      programsApi.updateStatus(id, status),
    onSuccess: (_data, { status }) => {
      queryClient.invalidateQueries({ queryKey: ['programs'] });
      toast.success(`Program ${statusLabels[status] ?? status}`);
    },
    onError: async (error) => {
      let message = 'Failed to update program status';
      if (error instanceof HTTPError) {
        try {
          const body = await error.response.json();
          if (body.message) message = body.message;
        } catch {
          // use default message
        }
      }
      toast.error(message);
    },
  });
}
