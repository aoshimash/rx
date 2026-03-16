'use client';

import { PlanForm } from '@/components/plan-editor/PlanForm';
import { Skeleton } from '@/components/ui/skeleton';
import { useDeletePlan, usePlan, useUpdatePlan } from '@/lib/hooks/usePlans';
import type { PlanEntryCreate } from '@/types/api';
import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

export default function EditPlanPage() {
  const router = useRouter();
  const params = useParams();
  const planId = params.id as string;

  const { data: plan, isLoading: planLoading } = usePlan(planId);
  const updatePlan = useUpdatePlan();
  const deletePlan = useDeletePlan();

  const [planName, setPlanName] = useState('');
  const [planDescription, setPlanDescription] = useState('');

  useEffect(() => {
    if (plan) {
      setPlanName(plan.name);
      setPlanDescription(plan.description || '');
    }
  }, [plan]);

  const handleSave = async (entries: PlanEntryCreate[]) => {
    await updatePlan.mutateAsync({
      id: planId,
      data: {
        name: planName,
        description: planDescription,
        entries,
      },
    });
    router.push('/plans');
  };

  const handleDelete = async () => {
    await deletePlan.mutateAsync(planId);
    router.push('/plans');
  };

  if (planLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <Skeleton className="h-[600px] w-full" />
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Edit Training Plan</h1>
        <p className="text-muted-foreground mt-1">Modify sessions and exercise prescriptions</p>
      </div>

      <PlanForm
        key={planId}
        planName={planName}
        planDescription={planDescription}
        initialEntries={plan?.entries}
        onNameChange={setPlanName}
        onDescriptionChange={setPlanDescription}
        onSave={handleSave}
        onDelete={handleDelete}
        isSaving={updatePlan.isPending}
        isEditing
      />
    </main>
  );
}
