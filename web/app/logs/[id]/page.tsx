'use client';

import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { useLog } from '@/lib/hooks/useLogs';
import { usePlan } from '@/lib/hooks/usePlans';
import { calculateDiff, getStatusIcon, getStatusVariant } from '@/lib/utils/diff';
import type { LogEntry, PlanEntry } from '@/types/api';
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';
import { useParams } from 'next/navigation';

function findMatchingPlanEntry(
  logEntry: LogEntry,
  planEntries: PlanEntry[]
): PlanEntry | undefined {
  return planEntries.find((pe) => pe.exercise_name === logEntry.exercise_name);
}

export default function LogDetailPage() {
  const params = useParams();
  const logId = params.id as string;

  const { data: log, isLoading: logLoading, error: logError } = useLog(logId);
  const { data: plan, isLoading: planLoading } = usePlan(log?.plan_id ?? null);

  const performedDate = log
    ? new Date(log.performed_at).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
      })
    : '';

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

  const planEntries = plan?.entries || [];
  const hasPlan = !!log.plan_id;

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
        {log.notes && <p className="text-muted-foreground mt-2">{log.notes}</p>}
        {hasPlan && (
          <div className="mt-2">
            {planLoading ? (
              <Skeleton className="h-6 w-[150px]" />
            ) : plan ? (
              <Link href={`/plans/${plan.id}/edit`}>
                <Badge variant="secondary" className="cursor-pointer">
                  Plan: {plan.name}
                </Badge>
              </Link>
            ) : (
              <Badge variant="outline">Linked Plan</Badge>
            )}
          </div>
        )}
      </div>

      <div className="space-y-3">
        {log.entries.map((entry) => {
          const matchedPlan = hasPlan ? findMatchingPlanEntry(entry, planEntries) : undefined;
          const diff = hasPlan ? calculateDiff(matchedPlan, entry) : null;

          return (
            <Card key={entry.id}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle className="text-base">{entry.exercise_name}</CardTitle>
                  {diff && (
                    <Badge variant={getStatusVariant(diff.status)}>
                      {getStatusIcon(diff.status)}{' '}
                      {diff.status.charAt(0).toUpperCase() + diff.status.slice(1)}
                    </Badge>
                  )}
                </div>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-4 gap-4 text-sm">
                  <div>
                    <span className="text-muted-foreground block">Sets</span>
                    <span className="font-medium">{entry.sets ?? '-'}</span>
                  </div>
                  <div>
                    <span className="text-muted-foreground block">Reps</span>
                    <span className="font-medium">{entry.reps ?? '-'}</span>
                  </div>
                  <div>
                    <span className="text-muted-foreground block">Load</span>
                    <span className="font-medium">
                      {entry.load_kg != null ? `${entry.load_kg}kg` : '-'}
                    </span>
                  </div>
                  <div>
                    <span className="text-muted-foreground block">RPE</span>
                    <span className="font-medium">{entry.rpe ?? '-'}</span>
                  </div>
                </div>
                {diff && diff.differences.length > 0 && (
                  <div className="mt-3 pt-3 border-t">
                    <p className="text-xs text-muted-foreground mb-1">Differences from plan:</p>
                    <div className="flex gap-2 flex-wrap">
                      {diff.differences.map((d) => (
                        <Badge key={d} variant="outline" className="text-xs">
                          {d}
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}
                {entry.notes && <p className="text-sm text-muted-foreground mt-2">{entry.notes}</p>}
              </CardContent>
            </Card>
          );
        })}
      </div>
    </main>
  );
}
