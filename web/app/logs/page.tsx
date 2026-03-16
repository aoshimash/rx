'use client';

import { ExportButton } from '@/components/export/ExportButton';
import { LogModal } from '@/components/log-input/LogModal';
import { LogCard } from '@/components/logs/LogCard';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { useCreateLog, useLogs } from '@/lib/hooks/useLogs';
import type { LogEntryCreate } from '@/types/api';
import { Plus } from 'lucide-react';
import { useState } from 'react';

export default function LogsPage() {
  const { data: logsData, isLoading, error } = useLogs();
  const createLog = useCreateLog();
  const [modalOpen, setModalOpen] = useState(false);

  const logs = logsData?.data || [];
  const sortedLogs = [...logs].sort(
    (a, b) => new Date(b.performed_at).getTime() - new Date(a.performed_at).getTime()
  );

  const handleSaveLog = async (entries: LogEntryCreate[], notes: string) => {
    await createLog.mutateAsync({
      performed_at: new Date().toISOString(),
      notes: notes || undefined,
      entries,
    });
  };

  if (isLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <div className="grid gap-4 md:grid-cols-2">
          <Skeleton className="h-[150px]" />
          <Skeleton className="h-[150px]" />
          <Skeleton className="h-[150px]" />
        </div>
      </main>
    );
  }

  if (error) {
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
        <div className="grid gap-4 md:grid-cols-2">
          {sortedLogs.map((log) => (
            <LogCard key={log.id} log={log} />
          ))}
        </div>
      )}

      <LogModal open={modalOpen} onOpenChange={setModalOpen} onSave={handleSaveLog} />
    </main>
  );
}
