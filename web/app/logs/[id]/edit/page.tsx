'use client';

import { LogEntryTable } from '@/components/log-entry-table/LogEntryTable';
import type { TableEntry } from '@/components/log-entry-table/types';
import { Skeleton } from '@/components/ui/skeleton';
import { useLog, useUpdateLog } from '@/lib/hooks/useLogs';
import type { LogEntryCreate } from '@/types/api';
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';
import { useParams, useRouter } from 'next/navigation';

/** Convert ISO string to datetime-local input value (YYYY-MM-DDTHH:MM) */
function toDatetimeLocalValue(iso: string): string {
  const date = new Date(iso);
  const pad = (n: number) => n.toString().padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`;
}

function logToTableEntries(log: {
  entries: Array<{
    id: string;
    exercise_name: string;
    order: number;
    sets?: number;
    reps?: number;
    load_kg?: number;
    rpe?: number;
    notes?: string;
  }>;
}): TableEntry[] {
  return log.entries
    .slice()
    .sort((a, b) => a.order - b.order)
    .map((entry) => ({
      id: entry.id,
      exercise_name: entry.exercise_name,
      sets: entry.sets,
      reps: entry.reps,
      load_kg: entry.load_kg,
      rpe: entry.rpe,
      notes: entry.notes ?? '',
      setsEdited: true,
      repsEdited: true,
    }));
}

export default function EditLogPage() {
  const params = useParams();
  const router = useRouter();
  const logId = params.id as string;
  const { data: log, isLoading } = useLog(logId);
  const updateLog = useUpdateLog();

  if (isLoading || !log) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-8 w-[200px]" />
        <Skeleton className="h-[400px] w-full" />
      </main>
    );
  }

  const handleSave = async (data: {
    entries: LogEntryCreate[];
    notes: string;
    startedAt?: string;
    finishedAt?: string;
  }) => {
    await updateLog.mutateAsync({
      id: logId,
      data: {
        program_id: log.program_id,
        session_name: log.session_name,
        performed_at: log.performed_at,
        started_at: data.startedAt ?? log.started_at,
        finished_at: data.finishedAt ?? log.finished_at,
        notes: data.notes || log.notes || undefined,
        entries: data.entries,
      },
    });
    router.push(`/logs/${logId}`);
  };

  return (
    <main className="container mx-auto p-6">
      <Link
        href={`/logs/${logId}`}
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-4"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Log
      </Link>

      <div className="mb-6">
        <h1 className="text-3xl font-bold">Edit Log</h1>
        {log.session_name && <p className="text-muted-foreground mt-1">{log.session_name}</p>}
      </div>

      <LogEntryTable
        existingEntries={logToTableEntries(log)}
        onSave={handleSave}
        onCancel={() => router.push(`/logs/${logId}`)}
        saveLabel="Update Log"
        initialNotes={log.notes ?? ''}
        initialStartedAt={log.started_at ? toDatetimeLocalValue(log.started_at) : ''}
        initialFinishedAt={log.finished_at ? toDatetimeLocalValue(log.finished_at) : ''}
      />
    </main>
  );
}
