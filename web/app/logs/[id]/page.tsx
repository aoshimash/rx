'use client';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { useLog } from '@/lib/hooks/useLogs';
import type { LogEntry } from '@/types/api';
import { ArrowLeft, Pencil } from 'lucide-react';
import Link from 'next/link';
import { useParams } from 'next/navigation';

type LogExerciseGroup = { name: string; entries: LogEntry[] };

function groupByExercise(entries: LogEntry[]): LogExerciseGroup[] {
  const groups: LogExerciseGroup[] = [];
  const map = new Map<string, LogExerciseGroup>();
  for (const entry of [...entries].sort((a, b) => a.order - b.order)) {
    if (!map.has(entry.exercise_name)) {
      const g: LogExerciseGroup = { name: entry.exercise_name, entries: [] };
      groups.push(g);
      map.set(entry.exercise_name, g);
    }
    map.get(entry.exercise_name)!.entries.push(entry);
  }
  return groups;
}

export default function LogDetailPage() {
  const params = useParams();
  const logId = params.id as string;

  const { data: log, isLoading: logLoading, error: logError } = useLog(logId);

  const performedDate = log
    ? new Date(log.performed_at).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
    : '';

  const formatTime = (iso: string) =>
    new Date(iso).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    });

  if (logLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-8 w-[200px]" />
        <Skeleton className="h-12 w-[400px]" />
        <Skeleton className="h-[400px] w-full" />
      </main>
    );
  }

  if (logError || !log) {
    return (
      <main className="container mx-auto p-6">
        <p className="text-destructive">Failed to load log. Please try again later.</p>
        <Link href="/logs" className="text-primary hover:underline mt-4 inline-block">
          Back to Logs
        </Link>
      </main>
    );
  }

  const groups = groupByExercise(log.entries);

  return (
    <main className="container mx-auto p-6">
      <Link
        href="/logs"
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-4"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Logs
      </Link>

      <div className="mb-6">
        <h1 className="text-3xl font-bold">{performedDate}</h1>
        {log.session_name && (
          <p className="text-muted-foreground mt-1">Session: {log.session_name}</p>
        )}
        {(log.started_at || log.finished_at) && (
          <p className="text-muted-foreground mt-1 text-sm">
            {log.started_at && <>Start: {formatTime(log.started_at)}</>}
            {log.started_at && log.finished_at && <span className="mx-2">·</span>}
            {log.finished_at && <>End: {formatTime(log.finished_at)}</>}
            {log.started_at && log.finished_at && (
              <>
                <span className="mx-2">·</span>
                {Math.round(
                  (new Date(log.finished_at).getTime() - new Date(log.started_at).getTime()) / 60000
                )}{' '}
                min
              </>
            )}
          </p>
        )}
        {log.notes && <p className="text-muted-foreground mt-2">{log.notes}</p>}
        <div className="flex items-center gap-2 mt-2">
          <Button variant="outline" size="sm" asChild>
            <Link href={`/logs/${logId}/edit`}>
              <Pencil className="h-4 w-4 mr-1" />
              Edit
            </Link>
          </Button>
        </div>
      </div>

      <div className="space-y-4">
        {groups.map((group) => (
          <Card key={group.name}>
            <CardHeader className="pb-2">
              <CardTitle className="text-base">{group.name}</CardTitle>
            </CardHeader>
            <CardContent>
              <table className="w-full text-sm">
                <thead>
                  <tr className="text-xs text-muted-foreground">
                    {group.entries.some((e) => e.metadata?.label) && (
                      <th className="text-left font-normal pb-1 w-16" />
                    )}
                    <th className="text-right font-normal pb-1 pr-4">RPE</th>
                    {group.entries.some((e) => e.load_kg != null) && (
                      <th className="text-right font-normal pb-1 pr-4">Load</th>
                    )}
                    <th className="text-right font-normal pb-1 pr-4">Reps</th>
                    <th className="text-right font-normal pb-1 pr-4">Sets</th>
                  </tr>
                </thead>
                <tbody>
                  {group.entries.map((entry) => {
                    const label = entry.metadata?.label as string | undefined;
                    const hasLabel = group.entries.some((e) => e.metadata?.label);
                    const hasLoad = group.entries.some((e) => e.load_kg != null);
                    return (
                      <tr key={entry.id} className="text-muted-foreground">
                        {hasLabel && <td className="text-xs pr-3 py-0.5">{label ?? ''}</td>}
                        <td className="text-right tabular-nums pr-4 py-0.5">{entry.rpe ?? '—'}</td>
                        {hasLoad && (
                          <td className="text-right tabular-nums pr-4 py-0.5">
                            {entry.load_kg != null ? `${entry.load_kg}kg` : '—'}
                          </td>
                        )}
                        <td className="text-right tabular-nums pr-4 py-0.5">{entry.reps ?? '—'}</td>
                        <td className="text-right tabular-nums pr-4 py-0.5">{entry.sets ?? '—'}</td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </CardContent>
          </Card>
        ))}
      </div>
    </main>
  );
}
