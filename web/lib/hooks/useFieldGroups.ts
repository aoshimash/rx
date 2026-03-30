import { fieldGroupsApi } from '@/lib/api/field-groups';
import type { FieldGroupCreate, FieldGroupUpdate } from '@/types/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export function useFieldGroups() {
  return useQuery({
    queryKey: ['field-groups'],
    queryFn: () => fieldGroupsApi.list(),
  });
}

export function useFieldGroup(id: string | null) {
  return useQuery({
    queryKey: ['field-groups', id],
    queryFn: () => fieldGroupsApi.get(id!),
    enabled: !!id,
  });
}

export function useCreateFieldGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: FieldGroupCreate) => fieldGroupsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['field-groups'] });
    },
  });
}

export function useUpdateFieldGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: FieldGroupUpdate }) =>
      fieldGroupsApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['field-groups'] });
    },
  });
}

export function useDeleteFieldGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => fieldGroupsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['field-groups'] });
    },
  });
}
