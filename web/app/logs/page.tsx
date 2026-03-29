'use client';

import { ExportButton } from '@/components/export/ExportButton';
import { LogCard } from '@/components/logs/LogCard';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { useLogs } from '@/lib/hooks/useLogs';
import { usePrograms } from '@/lib/hooks/usePrograms';
import { Plus } from 'lucide-react';
import Link from 'next/link';

export default function LogsPage() {
  const { data: logsData, isLoading: logsLoading, error: logsError } = useLogs();
  const { data: programsData } = usePrograms();

  const logs = logsData?.data || [];
  const sortedLogs = [...logs].sort(
    (a, b) => new Date(b.performed_at).getTime() - new Date(a.performed_at).getTime()
  );

  const programMap = new Map<string, string>((programsData?.data ?? []).map((p) => [p.id, p.name]));

  if (logsLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <div className="space-y-3">
          <Skeleton className="h-[120px] w-full" />
          <Skeleton className="h-[120px] w-full" />
          <Skeleton className="h-[120px] w-full" />
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
        <div className="flex items-center gap-2">
          <ExportButton logs={logs} plan={null} />
          <Button asChild>
            <Link href="/logs/new">
              <Plus className="h-4 w-4 mr-2" />
              Record Log
            </Link>
          </Button>
        </div>
      </div>

      {sortedLogs.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">
            No logs yet. Record your first training session.
          </p>
          <Button asChild>
            <Link href="/logs/new">
              <Plus className="h-4 w-4 mr-2" />
              Record Log
            </Link>
          </Button>
        </div>
      ) : (
        <div className="space-y-3">
          {sortedLogs.map((log) => (
            <LogCard
              key={log.id}
              log={log}
              programName={log.program_id ? programMap.get(log.program_id) : undefined}
            />
          ))}
        </div>
      )}
    </main>
  );
}
