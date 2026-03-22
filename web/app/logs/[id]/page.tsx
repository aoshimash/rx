'use client';

import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { useLog } from '@/lib/hooks/useLogs';
import type { LogEntry } from '@/types/api';
import { ArrowLeft } from 'lucide-react';
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
      </div>

      <div className="space-y-4">
        {groups.map((group) => (
          <Card key={group.name}>
            <CardHeader className="pb-2">
              <CardTitle className="text-base">{group.name}</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {group.entries.map((entry) => {
                  const label = entry.metadata?.label as string | undefined;

                  return (
                    <div key={entry.id} className="rounded-md border px-3 py-2">
                      <div className="flex items-center justify-between mb-1.5">
                        <div>
                          {label ? (
                            <Badge variant="outline" className="text-xs">
                              {label}
                            </Badge>
                          ) : (
                            <span className="text-xs text-muted-foreground">-</span>
                          )}
                        </div>
                      </div>
                      <div className="grid grid-cols-4 gap-3 text-sm">
                        <div>
                          <span className="text-muted-foreground block text-xs">Sets</span>
                          <span className="font-medium">{entry.sets ?? '-'}</span>
                        </div>
                        <div>
                          <span className="text-muted-foreground block text-xs">Reps</span>
                          <span className="font-medium">{entry.reps ?? '-'}</span>
                        </div>
                        <div>
                          <span className="text-muted-foreground block text-xs">Load</span>
                          <span className="font-medium">
                            {entry.load_kg != null ? `${entry.load_kg}kg` : '-'}
                          </span>
                        </div>
                        <div>
                          <span className="text-muted-foreground block text-xs">RPE</span>
                          <span className="font-medium">{entry.rpe ?? '-'}</span>
                        </div>
                      </div>
                      {entry.notes && (
                        <p className="text-xs text-muted-foreground mt-2">{entry.notes}</p>
                      )}
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </main>
  );
}
