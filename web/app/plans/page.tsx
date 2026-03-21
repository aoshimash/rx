'use client';

import type { LogSaveContext } from '@/components/log-input/LogModal';
import { LogModal } from '@/components/log-input/LogModal';
import { PlanForm } from '@/components/plan-editor/PlanForm';
import { PlanList } from '@/components/plans/PlanList';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Skeleton } from '@/components/ui/skeleton';
import { plansApi } from '@/lib/api/plans';
import { useCreateLog, useLogs } from '@/lib/hooks/useLogs';
import { useCreatePlan, usePlans } from '@/lib/hooks/usePlans';
import { buildGlobalPlanStatuses } from '@/lib/utils/next-session';
import type { LogEntryCreate, PlanEntryCreate } from '@/types/api';
import { useQueryClient } from '@tanstack/react-query';
import { Plus } from 'lucide-react';
import { useMemo, useState } from 'react';

export default function PlansPage() {
  const { data: plansData, isLoading: plansLoading } = usePlans();
  const { data: logsData } = useLogs();
  const createPlan = useCreatePlan();
  const createLog = useCreateLog();
  const queryClient = useQueryClient();

  const plans = plansData?.data || [];
  const logs = logsData?.data || [];

  // Build plan statuses with a single global NEXT plan
  const programGroups = useMemo(() => buildGlobalPlanStatuses(plans, logs), [plans, logs]);

  // Create Plan dialog state
  const [createOpen, setCreateOpen] = useState(false);
  const [planName, setPlanName] = useState('');
  const [planDescription, setPlanDescription] = useState('');

  // Log modal state
  const [logModalOpen, setLogModalOpen] = useState(false);
  const [logPlan, setLogPlan] = useState<
    { id: string; name: string; entries: PlanEntryCreate[] } | undefined
  >();

  const handleCreatePlan = async (entries: PlanEntryCreate[]) => {
    await createPlan.mutateAsync({
      name: planName,
      description: planDescription,
      entries,
    });
    setCreateOpen(false);
    setPlanName('');
    setPlanDescription('');
  };

  const handleCreateOpenChange = (value: boolean) => {
    setCreateOpen(value);
    if (!value) {
      setPlanName('');
      setPlanDescription('');
    }
  };

  const handleRecordLog = async (planId: string) => {
    const fullPlan = await queryClient.ensureQueryData({
      queryKey: ['plans', planId],
      queryFn: () => plansApi.get(planId),
    });
    setLogPlan({
      id: fullPlan.id,
      name: fullPlan.name,
      entries: (fullPlan.entries || []).map((e) => ({
        exercise_name: e.exercise_name,
        order: e.order,
        sets: e.sets,
        reps: e.reps,
        load_kg: e.load_kg,
        rpe: e.rpe,
        notes: e.notes,
        metadata: e.metadata,
      })),
    });
    setLogModalOpen(true);
  };

  const handleSaveLog = async (
    entries: LogEntryCreate[],
    notes: string,
    context?: LogSaveContext
  ) => {
    await createLog.mutateAsync({
      plan_id: context?.planId,
      performed_at: new Date().toISOString(),
      started_at: context?.startedAt,
      finished_at: context?.finishedAt,
      notes: notes || undefined,
      entries,
    });
  };

  if (plansLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <Skeleton className="h-[200px]" />
        <Skeleton className="h-[200px]" />
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Plans</h1>
          <p className="text-muted-foreground mt-1">Manage your training plans</p>
        </div>
        <Button variant="outline" size="sm" onClick={() => setCreateOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Create Plan
        </Button>
      </div>

      {plans.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">
            No plans yet. Create your first training plan.
          </p>
          <Button onClick={() => setCreateOpen(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Create Plan
          </Button>
        </div>
      ) : (
        <div className="space-y-8">
          {programGroups.map(({ groupKey, programName, isStandalone, statuses }) => (
            <section key={groupKey}>
              {programGroups.length > 1 && !isStandalone && (
                <h2 className="text-lg font-semibold mb-3">
                  {programName}
                  <span className="text-sm font-normal text-muted-foreground ml-2">
                    {statuses.length} session{statuses.length !== 1 ? 's' : ''}
                  </span>
                </h2>
              )}
              <PlanList statuses={statuses} onRecordLog={handleRecordLog} />
            </section>
          ))}
        </div>
      )}

      <Dialog open={createOpen} onOpenChange={handleCreateOpenChange}>
        <DialogContent className="sm:max-w-3xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Create Plan</DialogTitle>
            <DialogDescription>
              Define exercise prescriptions for a single session
            </DialogDescription>
          </DialogHeader>
          <PlanForm
            planName={planName}
            planDescription={planDescription}
            onNameChange={setPlanName}
            onDescriptionChange={setPlanDescription}
            onSave={handleCreatePlan}
            isSaving={createPlan.isPending}
          />
        </DialogContent>
      </Dialog>

      <LogModal
        open={logModalOpen}
        onOpenChange={setLogModalOpen}
        dayEntries={logPlan?.entries}
        planId={logPlan?.id}
        sessionName={logPlan?.name}
        planContext={logPlan ? [logPlan.name] : undefined}
        onSave={handleSaveLog}
      />
    </main>
  );
}
