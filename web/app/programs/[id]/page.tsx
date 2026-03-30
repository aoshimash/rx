'use client';

import { SessionSelector } from '@/components/programs/SessionSelector';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Skeleton } from '@/components/ui/skeleton';
import { useAddPlanSessions } from '@/lib/hooks/usePlans';
import { useDeleteProgram, useProgram } from '@/lib/hooks/usePrograms';
import type { PlanSessionCreate, Program, ProgramSession, ProgramSessionEntry } from '@/types/api';
import { CalendarDays, ClipboardPen, Copy, Download, Pencil, Share2, Trash2 } from 'lucide-react';
import Link from 'next/link';
import { useParams, useRouter } from 'next/navigation';
import { useState } from 'react';

function programToExportJson(program: Program): string {
  const payload = {
    rx_version: '1',
    name: program.name,
    ...(program.notes ? { notes: program.notes } : {}),
    sessions: program.sessions.map((s) => ({
      session_name: s.session_name,
      order: s.order,
      ...(s.date ? { date: s.date } : {}),
      entries: s.entries.map((e) => ({
        exercise_name: e.exercise_name,
        order: e.order,
        ...(e.fields ? { fields: e.fields } : {}),
        ...(e.notes ? { notes: e.notes } : {}),
      })),
    })),
  };
  return JSON.stringify(payload, null, 2);
}

type ExerciseGroup = { name: string; entries: ProgramSessionEntry[] };

function groupByExercise(entries: ProgramSessionEntry[]): ExerciseGroup[] {
  const groups: ExerciseGroup[] = [];
  const map = new Map<string, ExerciseGroup>();
  for (const entry of [...entries].sort((a, b) => a.order - b.order)) {
    if (!map.has(entry.exercise_name)) {
      const g: ExerciseGroup = { name: entry.exercise_name, entries: [] };
      groups.push(g);
      map.set(entry.exercise_name, g);
    }
    map.get(entry.exercise_name)?.entries.push(entry);
  }
  return groups;
}

function ProgramSessionCard({
  session,
  programId,
}: {
  session: ProgramSession;
  programId: string;
}) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 flex-wrap">
            <CardTitle className="text-base">{session.session_name}</CardTitle>
            {session.date && (
              <span className="flex items-center gap-1 text-sm font-medium text-foreground">
                <CalendarDays className="h-3.5 w-3.5" />
                {session.date}
              </span>
            )}
          </div>
          <Button variant="outline" size="sm" asChild>
            <Link
              href={`/logs/new?programId=${programId}&session=${encodeURIComponent(session.session_name)}`}
            >
              <ClipboardPen className="h-4 w-4 mr-1" />
              Record
            </Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {session.entries.length === 0 ? (
          <p className="text-sm text-muted-foreground">No exercises</p>
        ) : (
          <div className="divide-y">
            {groupByExercise(session.entries).map((group) => (
              <div key={group.name} className="py-2 first:pt-0 last:pb-0">
                <p className="font-medium text-sm mb-1">{group.name}</p>
                <table className="w-full text-sm">
                  <thead>
                    <tr className="text-xs text-muted-foreground">
                      {group.entries.some((e) => e.fields?.label) && (
                        <th className="text-left font-normal pb-1 w-16" />
                      )}
                      {group.entries.some((e) => e.fields?.load_kg != null) && (
                        <th className="text-right font-normal pb-1 pr-4 w-20">Load</th>
                      )}
                      <th className="text-right font-normal pb-1 pr-4 w-16">Reps</th>
                      <th className="text-right font-normal pb-1 pr-4 w-16">Sets</th>
                    </tr>
                  </thead>
                  <tbody>
                    {group.entries.map((entry) => {
                      const label = entry.fields?.label as string | undefined;
                      const hasLabel = group.entries.some((e) => e.fields?.label);
                      const hasLoad = group.entries.some((e) => e.fields?.load_kg != null);
                      const loadKg = entry.fields?.load_kg as number | null | undefined;
                      const reps = entry.fields?.reps as number | null | undefined;
                      const sets = entry.fields?.sets as number | null | undefined;
                      return (
                        <tr key={entry.id} className="text-muted-foreground">
                          {hasLabel && <td className="text-xs pr-3 py-0.5">{label ?? ''}</td>}
                          {hasLoad && (
                            <td className="text-right tabular-nums pr-4 py-0.5">
                              {loadKg != null ? `${loadKg}kg` : '—'}
                            </td>
                          )}
                          <td className="text-right tabular-nums pr-4 py-0.5">{reps ?? '—'}</td>
                          <td className="text-right tabular-nums pr-4 py-0.5">{sets ?? '—'}</td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
                {group.entries.some((e) => e.notes) && (
                  <div className="mt-1 space-y-0.5">
                    {group.entries
                      .filter((e) => e.notes)
                      .map((e) => (
                        <p key={e.id} className="text-xs text-muted-foreground">
                          {e.notes}
                        </p>
                      ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

export default function ProgramDetailPage() {
  const params = useParams();
  const router = useRouter();
  const programId = params.id as string;
  const { data: program, isLoading } = useProgram(programId);
  const deleteProgram = useDeleteProgram();
  const addPlanSessions = useAddPlanSessions();
  const [copied, setCopied] = useState(false);
  const [selectorOpen, setSelectorOpen] = useState(false);

  if (isLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <Skeleton className="h-[400px] w-full" />
      </main>
    );
  }

  if (!program) {
    return (
      <main className="container mx-auto p-6">
        <p className="text-muted-foreground">Program not found.</p>
      </main>
    );
  }

  const handleDelete = async () => {
    await deleteProgram.mutateAsync(programId);
    router.push('/');
  };

  const handleAddToPlan = async (sessions: PlanSessionCreate[]) => {
    await addPlanSessions.mutateAsync(sessions);
    setSelectorOpen(false);
    router.push('/');
  };

  const handleCopyToClipboard = async () => {
    await navigator.clipboard.writeText(programToExportJson(program));
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleDownload = () => {
    const json = programToExportJson(program);
    const blob = new Blob([json], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${program.name}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const sortedSessions = program.sessions.slice().sort((a, b) => a.order - b.order);

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">{program.name}</h1>
          {program.notes && <p className="text-muted-foreground mt-1">{program.notes}</p>}
        </div>
        <div className="flex items-center gap-2">
          <Button onClick={() => setSelectorOpen(true)}>Add to Plan</Button>
          <Button variant="outline" size="sm" asChild>
            <Link href={`/programs/${programId}/edit`}>
              <Pencil className="h-4 w-4 mr-2" />
              Edit
            </Link>
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm">
                <Share2 className="h-4 w-4 mr-2" />
                Export
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={handleCopyToClipboard}>
                <Copy className="h-4 w-4 mr-2" />
                {copied ? 'Copied!' : 'Copy to clipboard'}
              </DropdownMenuItem>
              <DropdownMenuItem onClick={handleDownload}>
                <Download className="h-4 w-4 mr-2" />
                Download .json
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          <Button variant="outline" onClick={handleDelete} disabled={deleteProgram.isPending}>
            <Trash2 className="h-4 w-4 mr-2" />
            Delete
          </Button>
        </div>
      </div>

      {sortedSessions.length === 0 ? (
        <p className="text-muted-foreground">No sessions in this program.</p>
      ) : (
        <div className="space-y-3">
          {sortedSessions.map((session) => (
            <ProgramSessionCard key={session.id} session={session} programId={programId} />
          ))}
        </div>
      )}

      <SessionSelector
        open={selectorOpen}
        onOpenChange={setSelectorOpen}
        sessions={program.sessions}
        programId={programId}
        onConfirm={handleAddToPlan}
        isPending={addPlanSessions.isPending}
      />
    </main>
  );
}
