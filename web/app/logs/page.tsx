'use client';

import { ExportButton } from '@/components/export/ExportButton';
import type { LogSaveContext } from '@/components/log-input/LogModal';
import { LogModal } from '@/components/log-input/LogModal';
import { LogTable } from '@/components/logs/LogTable';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { useCreateLog, useLogs } from '@/lib/hooks/useLogs';
import { usePlans } from '@/lib/hooks/usePlans';
import { usePrograms } from '@/lib/hooks/usePrograms';
import type { LogEntryCreate, Plan, Program } from '@/types/api';
import { Plus } from 'lucide-react';
import { useMemo, useState } from 'react';

export default function LogsPage() {
  const { data: logsData, isLoading: logsLoading, error: logsError } = useLogs();
  const { data: plansData, isLoading: plansLoading } = usePlans();
  const { data: programsData, isLoading: programsLoading } = usePrograms();
  const createLog = useCreateLog();
  const [modalOpen, setModalOpen] = useState(false);

  const isLoading = logsLoading || plansLoading || programsLoading;

  const logs = logsData?.data || [];
  const sortedLogs = [...logs].sort(
    (a, b) => new Date(b.performed_at).getTime() - new Date(a.performed_at).getTime()
  );

  const planMap = useMemo(() => {
    const map = new Map<string, Plan>();
    for (const plan of plansData?.data || []) {
      map.set(plan.id, plan);
    }
    return map;
  }, [plansData]);

  const programMap = useMemo(() => {
    const map = new Map<string, Program>();
    for (const program of programsData?.data || []) {
      map.set(program.id, program);
    }
    return map;
  }, [programsData]);

  const handleSaveLog = async (
    entries: LogEntryCreate[],
    notes: string,
    context?: LogSaveContext
  ) => {
    await createLog.mutateAsync({
      performed_at: new Date().toISOString(),
      started_at: context?.startedAt,
      finished_at: context?.finishedAt,
      notes: notes || undefined,
      entries,
    });
  };

  if (isLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <div className="space-y-2">
          <Skeleton className="h-10 w-full" />
          <Skeleton className="h-12 w-full" />
          <Skeleton className="h-12 w-full" />
          <Skeleton className="h-12 w-full" />
        </div>
      </main>
    );
  }

  if (logsError) {
    return (
      <main className="container mx-auto p-6">
        <p className="text-destructive">Failed to load logs. Please try again later.</p>
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Logs</h1>
          <p className="text-muted-foreground mt-1">View and record your sessions</p>
        </div>
        <ExportButton logs={logs} plan={null} />
      </div>

      {sortedLogs.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">
            No logs yet. Record your first training session.
          </p>
          <Button onClick={() => setModalOpen(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Record Log
          </Button>
        </div>
      ) : (
        <LogTable logs={sortedLogs} planMap={planMap} programMap={programMap} />
      )}

      <LogModal open={modalOpen} onOpenChange={setModalOpen} onSave={handleSaveLog} />
    </main>
  );
}
