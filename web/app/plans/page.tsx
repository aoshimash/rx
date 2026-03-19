'use client';

import type { LogSaveContext } from '@/components/log-input/LogModal';
import { LogModal } from '@/components/log-input/LogModal';
import { PlanForm } from '@/components/plan-editor/PlanForm';
import { entriesToSessionGroups } from '@/components/plan-editor/types';
import { PlanSelector } from '@/components/plans/PlanSelector';
import { SessionList } from '@/components/plans/SessionList';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Skeleton } from '@/components/ui/skeleton';
import { useCreateLog, useLogs } from '@/lib/hooks/useLogs';
import { useCreatePlan, usePlan, usePlans } from '@/lib/hooks/usePlans';
import { detectNextSession, sortSessionsByNext } from '@/lib/utils/next-session';
import { usePlanStore } from '@/stores/plan';
import type { LogEntryCreate, PlanEntryCreate } from '@/types/api';
import { Edit, Plus } from 'lucide-react';
import Link from 'next/link';
import { useEffect, useMemo, useState } from 'react';

export default function PlansPage() {
  const { data: plansData, isLoading: plansLoading } = usePlans();
  const { data: logsData } = useLogs();
  const createPlan = useCreatePlan();
  const createLog = useCreateLog();
  const { selectedPlanId, setSelectedPlan } = usePlanStore();

  const plans = plansData?.data || [];
  const logs = logsData?.data || [];

  // Auto-select first plan if none selected
  useEffect(() => {
    if (plans.length > 0 && !selectedPlanId && plans[0]) {
      setSelectedPlan(plans[0].id);
    }
  }, [plans, selectedPlanId, setSelectedPlan]);

  const { data: selectedPlan } = usePlan(selectedPlanId);

  const sessions = useMemo(() => {
    if (!selectedPlan?.entries) return [];
    return entriesToSessionGroups(selectedPlan.entries);
  }, [selectedPlan]);

  const sessionStatuses = useMemo(() => {
    if (!selectedPlanId || sessions.length === 0) return [];
    return sortSessionsByNext(detectNextSession(sessions, logs, selectedPlanId));
  }, [sessions, logs, selectedPlanId]);

  // Create Plan dialog state
  const [createOpen, setCreateOpen] = useState(false);
  const [planName, setPlanName] = useState('');
  const [planDescription, setPlanDescription] = useState('');

  // Log modal state
  const [logModalOpen, setLogModalOpen] = useState(false);
  const [logSessionName, setLogSessionName] = useState<string | undefined>();
  const [logDayEntries, setLogDayEntries] = useState<PlanEntryCreate[] | undefined>();

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

  const handleRecordLog = (sessionName: string) => {
    const session = sessions.find((s) => s.name === sessionName);
    setLogSessionName(sessionName);
    setLogDayEntries(session?.exercises);
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
      notes: notes || undefined,
      metadata: context?.sessionName ? { session_name: context.sessionName } : undefined,
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
        <>
          <div className="mb-4 flex items-center gap-3">
            <PlanSelector
              plans={plans}
              selectedPlanId={selectedPlanId}
              onSelect={setSelectedPlan}
            />
            {selectedPlanId && (
              <Link href={`/plans/${selectedPlanId}/edit`}>
                <Button variant="ghost" size="sm">
                  <Edit className="h-4 w-4 mr-1" />
                  Edit
                </Button>
              </Link>
            )}
          </div>

          {sessions.length === 0 ? (
            <div className="text-center py-8">
              <p className="text-muted-foreground">
                This plan has no sessions.{' '}
                {selectedPlanId && (
                  <Link href={`/plans/${selectedPlanId}/edit`} className="underline">
                    Add sessions
                  </Link>
                )}
              </p>
            </div>
          ) : (
            <SessionList statuses={sessionStatuses} onRecordLog={handleRecordLog} />
          )}
        </>
      )}

      <Dialog open={createOpen} onOpenChange={handleCreateOpenChange}>
        <DialogContent className="sm:max-w-3xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Create Plan</DialogTitle>
            <DialogDescription>Define sessions and exercise prescriptions</DialogDescription>
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
        dayEntries={logDayEntries}
        planId={selectedPlanId ?? undefined}
        sessionName={logSessionName}
        planContext={
          selectedPlan && logSessionName ? [selectedPlan.name, logSessionName] : undefined
        }
        onSave={handleSaveLog}
      />
    </main>
  );
}
