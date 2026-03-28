'use client';

import { LogEntryTable } from '@/components/log-entry-table/LogEntryTable';
import { Skeleton } from '@/components/ui/skeleton';
import { useCreateLog } from '@/lib/hooks/useLogs';
import { useProgram } from '@/lib/hooks/usePrograms';
import type { LogEntryCreate } from '@/types/api';
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { Suspense } from 'react';

function NewLogPageContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const programId = searchParams.get('programId');
  const sessionName = searchParams.get('session');
  const { data: program, isLoading: programLoading } = useProgram(programId);
  const createLog = useCreateLog();

  const session = program?.sessions.find((s) => s.session_name === sessionName);
  const planEntries = session?.entries?.slice().sort((a, b) => a.order - b.order);

  const backHref = programId ? `/programs/${programId}` : '/logs';

  const handleSave = async (data: {
    entries: LogEntryCreate[];
    notes: string;
    startedAt?: string;
    finishedAt?: string;
  }) => {
    await createLog.mutateAsync({
      program_id: programId ?? undefined,
      session_name: sessionName ?? undefined,
      performed_at: new Date().toISOString(),
      started_at: data.startedAt,
      finished_at: data.finishedAt,
      notes: data.notes || undefined,
      entries: data.entries,
    });
    router.push(backHref);
  };

  if (programId && programLoading) {
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
        {program && (
          <p className="text-muted-foreground mt-1">
            {program.name}
            {sessionName && <span> — {sessionName}</span>}
          </p>
        )}
      </div>

      <LogEntryTable
        initialEntries={planEntries}
        onSave={handleSave}
        onCancel={() => router.push(backHref)}
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
