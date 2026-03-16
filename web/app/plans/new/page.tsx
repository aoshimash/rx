'use client';

import { PlanForm } from '@/components/plan-editor/PlanForm';
import { useCreatePlan } from '@/lib/hooks/usePlans';
import type { PlanEntryCreate } from '@/types/api';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function NewPlanPage() {
  const router = useRouter();
  const createPlan = useCreatePlan();

  const [planName, setPlanName] = useState('');
  const [planDescription, setPlanDescription] = useState('');

  const handleSave = async (entries: PlanEntryCreate[]) => {
    await createPlan.mutateAsync({
      name: planName,
      description: planDescription,
      entries,
    });
    router.push('/plans');
  };

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Create Training Plan</h1>
        <p className="text-muted-foreground mt-1">Define sessions and exercise prescriptions</p>
      </div>

      <PlanForm
        planName={planName}
        planDescription={planDescription}
        onNameChange={setPlanName}
        onDescriptionChange={setPlanDescription}
        onSave={handleSave}
        isSaving={createPlan.isPending}
      />
    </main>
  );
}
