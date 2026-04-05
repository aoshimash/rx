import { plansApi } from '@/lib/api/plans';
import type { PlanCreate, PlanSessionCreate, PlanSessionUpdate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { HTTPError } from 'ky';
import { toast } from 'sonner';

export function usePlan() {
  return useQuery({
    queryKey: ['plan'],
    queryFn: () => plansApi.get(),
    retry: (failureCount, error) => {
      if (error instanceof HTTPError && error.response.status === 404) return false;
      return failureCount < 1;
    },
  });
}

export function useCreatePlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: PlanCreate) => plansApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan'] });
    },
  });
}

export function useAddPlanSessions() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (sessions: PlanSessionCreate[]) => {
      try {
        return await plansApi.addSessions(sessions);
      } catch (error) {
        if (error instanceof HTTPError && error.response.status === 404) {
          // Plan doesn't exist yet — create it with the sessions
          return await plansApi.create({ sessions });
        }
        throw error;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan'] });
    },
    onError: async (error) => {
      let message = 'Failed to add sessions to plan';
      if (error instanceof HTTPError) {
        try {
          const body = await error.response.json();
          if (body.message) message = body.message;
        } catch {
          // use default
        }
      }
      toast.error(message);
    },
  });
}

export function useUpdatePlanSession() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ sessionId, data }: { sessionId: string; data: PlanSessionUpdate }) =>
      plansApi.updateSession(sessionId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan'] });
    },
    onError: async (error) => {
      let message = 'Failed to update session';
      if (error instanceof HTTPError) {
        try {
          const body = await error.response.json();
          if (body.message) message = body.message;
        } catch {
          // use default
        }
      }
      toast.error(message);
    },
  });
}

export function useDeletePlanSession() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (sessionId: string) => plansApi.deleteSession(sessionId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan'] });
    },
  });
}

export function useExpandProgram() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (programId: string) => plansApi.expandProgram(programId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan'] });
      toast.success('Sessions added to plan');
    },
    onError: async (error) => {
      let message = 'Failed to expand program';
      if (error instanceof HTTPError) {
        try {
          const body = await error.response.json();
          if (body.message) message = body.message;
        } catch {
          // use default
        }
      }
      toast.error(message);
    },
  });
}
