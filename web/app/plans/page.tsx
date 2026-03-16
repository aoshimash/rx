'use client';

import { PlanForm } from '@/components/plan-editor/PlanForm';
import { PlanCard } from '@/components/plans/PlanCard';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Skeleton } from '@/components/ui/skeleton';
import { useCreatePlan, usePlans } from '@/lib/hooks/usePlans';
import { usePlanStore } from '@/stores/plan';
import type { PlanEntryCreate } from '@/types/api';
import { Plus } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function PlansPage() {
  const router = useRouter();
  const { data: plansData, isLoading } = usePlans();
  const createPlan = useCreatePlan();
  const { selectedPlanId, setSelectedPlan } = usePlanStore();

  const plans = plansData?.data || [];

  const [open, setOpen] = useState(false);
  const [planName, setPlanName] = useState('');
  const [planDescription, setPlanDescription] = useState('');

  const handleSelectPlan = (planId: string) => {
    setSelectedPlan(planId);
    router.push('/');
  };

  const handleSave = async (entries: PlanEntryCreate[]) => {
    await createPlan.mutateAsync({
      name: planName,
      description: planDescription,
      entries,
    });
    setOpen(false);
    setPlanName('');
    setPlanDescription('');
  };

  const handleOpenChange = (value: boolean) => {
    setOpen(value);
    if (!value) {
      setPlanName('');
      setPlanDescription('');
    }
  };

  if (isLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <Skeleton className="h-[200px]" />
          <Skeleton className="h-[200px]" />
          <Skeleton className="h-[200px]" />
        </div>
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-start justify-between">
        <div>
          <h1 className="text-3xl font-bold">Training Plans</h1>
          <p className="text-muted-foreground mt-1">Manage your training plans</p>
        </div>
        <Button onClick={() => setOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Create Plan
        </Button>
      </div>

      {plans.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">
            No plans yet. Create your first training plan.
          </p>
          <Button onClick={() => setOpen(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Create Plan
          </Button>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {plans.map((plan) => (
            <PlanCard
              key={plan.id}
              plan={plan}
              onSelect={() => handleSelectPlan(plan.id)}
              isSelected={selectedPlanId === plan.id}
            />
          ))}
        </div>
      )}

      <Dialog open={open} onOpenChange={handleOpenChange}>
        <DialogContent className="sm:max-w-3xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Create Training Plan</DialogTitle>
            <DialogDescription>Define sessions and exercise prescriptions</DialogDescription>
          </DialogHeader>
          <PlanForm
            planName={planName}
            planDescription={planDescription}
            onNameChange={setPlanName}
            onDescriptionChange={setPlanDescription}
            onSave={handleSave}
            isSaving={createPlan.isPending}
          />
        </DialogContent>
      </Dialog>
    </main>
  );
}
