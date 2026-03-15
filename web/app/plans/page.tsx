'use client';

import { PlanCard } from '@/components/plans/PlanCard';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { usePlans } from '@/lib/hooks/usePlans';
import { usePlanStore } from '@/stores/plan';
import { Plus } from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

export default function PlansPage() {
  const router = useRouter();
  const { data: plansData, isLoading } = usePlans();
  const { selectedPlanId, setSelectedPlan } = usePlanStore();

  const plans = plansData?.data || [];

  const handleSelectPlan = (planId: string) => {
    setSelectedPlan(planId);
    router.push('/');
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
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Training Plans</h1>
          <p className="text-muted-foreground mt-1">Select a plan to view in Week View</p>
        </div>
        <Link href="/plans/new">
          <Button>
            <Plus className="h-4 w-4 mr-2" />
            New Plan
          </Button>
        </Link>
      </div>

      {plans.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">
            No plans yet. Create your first training plan.
          </p>
          <Link href="/plans/new">
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              Create Plan
            </Button>
          </Link>
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
    </main>
  );
}
