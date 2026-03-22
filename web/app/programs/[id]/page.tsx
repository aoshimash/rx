'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { useDeleteProgram, useLoggedSessions, useProgram, useUpdateProgramStatus } from '@/lib/hooks/usePrograms';
import type { ProgramSession, ProgramSessionEntry } from '@/types/api';
import { Trash2 } from 'lucide-react';
import { useParams, useRouter } from 'next/navigation';

function buildCompletedSessionSet(sessions: string[]): Set<string> {
  return new Set(sessions);
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
  isCompleted,
  isNext,
}: { session: ProgramSession; isCompleted: boolean; isNext: boolean }) {
  return (
    <Card className={sessionCardClassName(isCompleted, isNext)}>
      <CardHeader className="pb-2">
        <div className="flex items-center gap-2">
          <CardTitle className="text-base">{session.session_name}</CardTitle>
          {session.date && <span className="text-sm text-muted-foreground">{session.date}</span>}
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

  const sortedSessions = program.sessions.slice().sort((a, b) => a.order - b.order);
  const completedSessionSet = buildCompletedSessionSet(loggedSessions?.sessions ?? []);
  const allSessionsLogged = program.sessions.length > 0 &&
    program.sessions.every((s) => completedSessionSet.has(s.session_name));

  let foundNextSession = false;
  const sessionsWithStatus = sortedSessions.map((session) => {
    const isCompleted = completedSessionSet.has(session.session_name);
    const isNext = !isCompleted && !foundNextSession;
    if (isNext) foundNextSession = true;
    return { session, isCompleted, isNext };
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
          {sessionsWithStatus.map(({ session, isCompleted, isNext }) => (
            <SessionCard
              key={session.id}
              session={session}
              isCompleted={isCompleted}
              isNext={isNext}
            />
          ))}
        </div>
      )}
    </main>
  );
}
