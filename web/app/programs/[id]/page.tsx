'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Skeleton } from '@/components/ui/skeleton';
import {
  useDeleteProgram,
  useLoggedSessions,
  useProgram,
  useUpdateProgramStatus,
} from '@/lib/hooks/usePrograms';
import type { LoggedSession, Program, ProgramSession, ProgramSessionEntry } from '@/types/api';
import {
  CalendarDays,
  ClipboardList,
  ClipboardPen,
  Clock,
  Copy,
  Download,
  Share2,
  Trash2,
} from 'lucide-react';
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
        ...(e.sets != null ? { sets: e.sets } : {}),
        ...(e.reps != null ? { reps: e.reps } : {}),
        ...(e.load_kg != null ? { load_kg: e.load_kg } : {}),
        ...(e.rpe != null ? { rpe: e.rpe } : {}),
        ...(e.notes ? { notes: e.notes } : {}),
      })),
    })),
  };
  return JSON.stringify(payload, null, 2);
}

function buildCompletedSessionMap(sessions: LoggedSession[]): Map<string, LoggedSession> {
  return new Map(sessions.map((s) => [s.session_name, s]));
}

function durationMinutes(startedAt: string, finishedAt: string): number {
  return Math.round((new Date(finishedAt).getTime() - new Date(startedAt).getTime()) / 60000);
}

function formatDate(isoString: string): string {
  return new Date(isoString).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
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

function sessionCardClassName(isCompleted: boolean, isNext: boolean): string | undefined {
  if (isCompleted) return 'opacity-50';
  if (isNext) return 'border-2 border-primary';
  return undefined;
}

function SessionCard({
  session,
  programId,
  isCompleted,
  isNext,
  loggedSession,
}: {
  session: ProgramSession;
  programId: string;
  isCompleted: boolean;
  isNext: boolean;
  loggedSession?: LoggedSession;
}) {
  const duration =
    loggedSession?.started_at && loggedSession?.finished_at
      ? durationMinutes(loggedSession.started_at, loggedSession.finished_at)
      : null;

  return (
    <Card className={sessionCardClassName(isCompleted, isNext)}>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 flex-wrap">
            <CardTitle className="text-base">{session.session_name}</CardTitle>
            {isCompleted && loggedSession ? (
              <>
                <span className="flex items-center gap-1 text-sm text-muted-foreground">
                  <CalendarDays className="h-3.5 w-3.5" />
                  {formatDate(loggedSession.performed_at)}
                </span>
                {duration !== null && (
                  <span className="flex items-center gap-1 text-sm text-muted-foreground">
                    <Clock className="h-3.5 w-3.5" />
                    {duration}min
                  </span>
                )}
              </>
            ) : (
              session.date && (
                <span className="flex items-center gap-1 text-sm font-medium text-foreground">
                  <CalendarDays className="h-3.5 w-3.5" />
                  {session.date}
                </span>
              )
            )}
          </div>
          {isCompleted && loggedSession ? (
            <Button variant="outline" size="sm" asChild>
              <Link href={`/logs/${loggedSession.log_id}`}>
                <ClipboardList className="h-4 w-4 mr-1" />
                View Log
              </Link>
            </Button>
          ) : (
            !isCompleted && (
              <Button variant="outline" size="sm" asChild>
                <Link
                  href={`/logs/new?programId=${programId}&session=${encodeURIComponent(session.session_name)}`}
                >
                  <ClipboardPen className="h-4 w-4 mr-1" />
                  Record
                </Link>
              </Button>
            )
          )}
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
                          <td className="text-right tabular-nums pr-4 py-0.5">
                            {entry.rpe ?? '—'}
                          </td>
                          {hasLoad && (
                            <td className="text-right tabular-nums pr-4 py-0.5">
                              {entry.load_kg != null ? `${entry.load_kg}kg` : '—'}
                            </td>
                          )}
                          <td className="text-right tabular-nums pr-4 py-0.5">
                            {entry.reps ?? '—'}
                          </td>
                          <td className="text-right tabular-nums pr-4 py-0.5">
                            {entry.sets ?? '—'}
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
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
  const { data: loggedSessions } = useLoggedSessions(programId);
  const deleteProgram = useDeleteProgram();
  const updateStatus = useUpdateProgramStatus();
  const [copied, setCopied] = useState(false);

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
    router.push('/programs');
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
  const completedSessionMap = buildCompletedSessionMap(loggedSessions?.sessions ?? []);
  const allSessionsLogged =
    program.sessions.length > 0 &&
    program.sessions.every((s) => completedSessionMap.has(s.session_name));

  let foundNextSession = false;
  const sessionsWithStatus = sortedSessions.map((session) => {
    const loggedSession = completedSessionMap.get(session.session_name);
    const isCompleted = loggedSession !== undefined;
    const isNext = !isCompleted && !foundNextSession;
    if (isNext) foundNextSession = true;
    return { session, isCompleted, isNext, loggedSession };
  });

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-3xl font-bold">{program.name}</h1>
            <Badge variant={program.status === 'ongoing' ? 'default' : 'secondary'}>
              {program.status}
            </Badge>
          </div>
          {program.notes && <p className="text-muted-foreground mt-1">{program.notes}</p>}
        </div>
        <div className="flex items-center gap-2">
          {program.status === 'created' && (
            <Button
              onClick={() => updateStatus.mutate({ id: programId, status: 'ongoing' })}
              disabled={updateStatus.isPending}
            >
              {updateStatus.isPending ? 'Starting...' : 'Start Program'}
            </Button>
          )}
          {program.status === 'ongoing' && (
            <>
              <Button
                onClick={() => updateStatus.mutate({ id: programId, status: 'completed' })}
                disabled={updateStatus.isPending || !allSessionsLogged}
              >
                {updateStatus.isPending ? 'Completing...' : 'Complete Program'}
              </Button>
              <Button
                variant="outline"
                onClick={() => updateStatus.mutate({ id: programId, status: 'cancelled' })}
                disabled={updateStatus.isPending}
              >
                Cancel Program
              </Button>
            </>
          )}
          {program.status === 'cancelled' && (
            <Button
              variant="outline"
              onClick={() => updateStatus.mutate({ id: programId, status: 'ongoing' })}
              disabled={updateStatus.isPending}
            >
              {updateStatus.isPending ? 'Resuming...' : 'Resume Program'}
            </Button>
          )}
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

      {sessionsWithStatus.length === 0 ? (
        <p className="text-muted-foreground">No sessions in this program.</p>
      ) : (
        <div className="space-y-3">
          {sessionsWithStatus.map(({ session, isCompleted, isNext, loggedSession }) => (
            <SessionCard
              key={session.id}
              session={session}
              programId={programId}
              isCompleted={isCompleted}
              isNext={isNext}
              loggedSession={loggedSession}
            />
          ))}
        </div>
      )}
    </main>
  );
}
