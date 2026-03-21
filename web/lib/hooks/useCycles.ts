import { cyclesApi } from '@/lib/api/cycles';
import { useQuery } from '@tanstack/react-query';

export function useCycles(programId?: string) {
  return useQuery({
    queryKey: ['cycles', { programId }],
    queryFn: () => cyclesApi.list({ programId, limit: 100 }),
  });
}

export function useCycle(id: string | null) {
  return useQuery({
    queryKey: ['cycles', id],
    queryFn: () => cyclesApi.get(id!),
    enabled: !!id,
  });
}
