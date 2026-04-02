'use client';

import { LogEntryTable } from '@/components/log-entry-table/LogEntryTable';
import { Skeleton } from '@/components/ui/skeleton';
import { useFieldGroup } from '@/lib/hooks/useFieldGroups';
import { useCreateLog } from '@/lib/hooks/useLogs';
import { useDeletePlanSession, usePlan } from '@/lib/hooks/usePlans';
import { useProgram } from '@/lib/hooks/usePrograms';
import type { LogEntryCreate } from '@/types/api';
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { Suspense } from 'react';

function NewLogPageContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const planSessionId = searchParams.get('planSessionId');
  const programId = searchParams.get('programId');
  const sessionName = searchParams.get('session');

  const { data: plan, isLoading: planLoading } = usePlan();
  const { data: program, isLoading: programLoading } = useProgram(programId);
  const createLog = useCreateLog();
  const deletePlanSession = useDeletePlanSession();

  const planSession = planSessionId
    ? plan?.sessions.find((s) => s.id === planSessionId)
    : undefined;

  const programSession = program?.sessions.find((s) => s.session_name === sessionName);

  const fieldGroupId = planSession?.field_group_id ?? programSession?.field_group_id ?? null;
  const { data: fieldGroup, isLoading: fieldGroupLoading } = useFieldGroup(fieldGroupId);

  const isDataLoading =
    (planSessionId && planLoading) ||
    (programId && programLoading) ||
    (fieldGroupId && fieldGroupLoading);

  const templateEntries = planSession
    ? planSession.entries.slice().sort((a, b) => a.order - b.order)
    : programSession?.entries?.slice().sort((a, b) => a.order - b.order);

  const displayName = planSession?.session_name ?? sessionName;
  const displayProgramName = program?.name;
  const backHref = planSessionId ? '/' : programId ? `/programs/${programId}` : '/logs';

  const handleSave = async (data: {
    entries: LogEntryCreate[];
    notes: string;
    startedAt?: string;
    finishedAt?: string;
  }) => {
    const planSnapshot = planSession
      ? {
          session_name: planSession.session_name,
          entries: planSession.entries.map((e) => ({
            exercise_name: e.exercise_name,
            order: e.order,
            fields: e.fields,
            notes: e.notes,
          })),
        }
      : undefined;

    await createLog.mutateAsync({
      program_id: planSession?.source_program_id ?? programId ?? undefined,
      session_name: displayName ?? undefined,
      performed_at: new Date().toISOString(),
      started_at: data.startedAt,
      finished_at: data.finishedAt,
      notes: data.notes || undefined,
      plan_snapshot: planSnapshot,
      entries: data.entries,
    });

    if (planSessionId) {
      await deletePlanSession.mutateAsync(planSessionId);
    }

    router.push(backHref);
  };

  if (isDataLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-8 w-[200px]" />
        <Skeleton className="h-[400px] w-full" />
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <Link
        href={backHref}
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-4"
      >
        <ArrowLeft className="h-4 w-4" />
        Back
      </Link>

      <div className="mb-6">
        <h1 className="text-3xl font-bold">Record Log</h1>
        {(displayProgramName || displayName) && (
          <p className="text-muted-foreground mt-1">
            {displayProgramName}
            {displayName && <span> — {displayName}</span>}
          </p>
        )}
        {fieldGroup && (
          <p className="text-xs text-muted-foreground mt-1">Fields: {fieldGroup.name}</p>
        )}
      </div>

      <LogEntryTable
        initialEntries={templateEntries}
        onSave={handleSave}
        onCancel={() => router.push(backHref)}
        fieldDefs={fieldGroup?.log_fields}
      />
    </main>
  );
}

function NewLogPageFallback() {
  return (
    <main className="container mx-auto p-6 space-y-4">
      <Skeleton className="h-8 w-[200px]" />
      <Skeleton className="h-[400px] w-full" />
    </main>
  );
}

export default function NewLogPage() {
  return (
    <Suspense fallback={<NewLogPageFallback />}>
      <NewLogPageContent />
    </Suspense>
  );
}
