'use client';

import { ExportButton } from '@/components/export/ExportButton';
import { LogModal } from '@/components/log-input/LogModal';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import { WeekView } from '@/components/week-view/WeekView';
import { useCreateLog, useLogs } from '@/lib/hooks/useLogs';
import { usePlan, usePlans } from '@/lib/hooks/usePlans';
import { generatePlanContext } from '@/lib/utils/planContext';
import { usePlanStore } from '@/stores/plan';
import type { LogEntryCreate, PlanEntry } from '@/types/api';
import { Plus } from 'lucide-react';
import Link from 'next/link';
import { useMemo, useState } from 'react';

export default function Home() {
  const [logModalOpen, setLogModalOpen] = useState(false);
  const [selectedDayEntries, setSelectedDayEntries] = useState<PlanEntry[] | undefined>();
  const { selectedPlanId, setSelectedPlan } = usePlanStore();

  const { data: plansData, isLoading: plansLoading } = usePlans();
  const { data: selectedPlan, isLoading: selectedPlanLoading } = usePlan(selectedPlanId);
  const { data: logsData, isLoading: logsLoading } = useLogs();
  const createLog = useCreateLog();

  const plans = plansData?.data || [];
  const plan = selectedPlan || (plans.length > 0 ? plans[0] : null);
  const logs = logsData?.data || [];

  const firstDayEntries = useMemo(() => {
    const entries = plan?.entries;
    if (!entries || entries.length === 0) return undefined;
    const firstEntry = entries[0];
    if (!firstEntry) return undefined;
    const dayKey = firstEntry.metadata?.day;
    const weekKey = firstEntry.metadata?.week;
    if (!dayKey) return [firstEntry];
    return entries.filter((e) => e.metadata?.day === dayKey && e.metadata?.week === weekKey);
  }, [plan]);

  const selectedPlanContext = useMemo(() => {
    const refEntry = selectedDayEntries?.[0];
    if (!refEntry || !plan) return undefined;
    return generatePlanContext(refEntry.id, plan);
  }, [selectedDayEntries, plan]);

  const handleSaveLog = async (entries: LogEntryCreate[], notes: string) => {
    await createLog.mutateAsync({
      performed_at: new Date().toISOString(),
      plan_id: plan?.id,
      notes,
      entries,
    });
  };

  const handlePlanChange = (planId: string) => {
    setSelectedPlan(planId);
  };

  const handleOpenLogModal = () => {
    setSelectedDayEntries(firstDayEntries);
    setLogModalOpen(true);
  };

  if (plansLoading || selectedPlanLoading || logsLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <Skeleton className="h-[400px] w-full" />
      </main>
    );
  }

  if (plans.length === 0) {
    return (
      <main className="container mx-auto p-6">
        <div className="text-center py-12">
          <h1 className="text-3xl font-bold mb-4">Welcome to Rx</h1>
          <p className="text-muted-foreground mb-6">
            Create your first training plan to get started.
          </p>
          <Link href="/plans/new">
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              Create Plan
            </Button>
          </Link>
        </div>
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <h1 className="text-3xl font-bold">Training Week</h1>
          {plans.length > 1 && (
            <Select value={plan?.id || ''} onValueChange={handlePlanChange}>
              <SelectTrigger className="w-[250px]">
                <SelectValue placeholder="Select plan" />
              </SelectTrigger>
              <SelectContent>
                {plans.map((p) => (
                  <SelectItem key={p.id} value={p.id}>
                    {p.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
          {plans.length === 1 && plan && <p className="text-muted-foreground">Plan: {plan.name}</p>}
        </div>
        <div className="flex items-center gap-2">
          <ExportButton logs={logs} plan={plan || null} />
          <Button onClick={handleOpenLogModal}>
            <Plus className="h-4 w-4 mr-2" />
            Record Log
          </Button>
        </div>
      </div>

      <WeekView plan={plan || null} logs={logs} />

      <LogModal
        open={logModalOpen}
        onOpenChange={setLogModalOpen}
        dayEntries={selectedDayEntries}
        planContext={selectedPlanContext}
        onSave={handleSaveLog}
      />
    </main>
  );
}
