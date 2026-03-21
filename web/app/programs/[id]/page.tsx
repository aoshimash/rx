'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { useProgram } from '@/lib/hooks/usePrograms';
import type { ProgramEntry } from '@/types/api';
import { ArrowRightLeft, Edit } from 'lucide-react';
import Link from 'next/link';
import { useParams } from 'next/navigation';

type ExerciseGroup = { name: string; entries: ProgramEntry[] };
type SessionGroup = { name: string; exerciseGroups: ExerciseGroup[] };

function formatEntryText(entry: ProgramEntry): string {
  const parts: string[] = [];
  if (entry.rpe != null) parts.push(`RPE${entry.rpe}`);
  if (entry.reps != null) parts.push(`${entry.reps}reps`);
  if (entry.sets != null) parts.push(`${entry.sets}sets`);
  if (entry.percent_1rm != null) parts.push(`${Math.round(entry.percent_1rm * 100)}%`);
  const weightKg = entry.metadata?.weight_kg as number | undefined;
  if (weightKg != null) parts.push(`${weightKg}kg`);
  return parts.join(' ');
}

function groupBySession(entries: ProgramEntry[]): SessionGroup[] {
  const sorted = [...entries].sort((a, b) => a.order - b.order);
  const sessionOrder: string[] = [];
  const sessionMap = new Map<string, ProgramEntry[]>();

  for (const entry of sorted) {
    const session = (entry.metadata?.session as string) || '';
    if (!sessionMap.has(session)) {
      sessionOrder.push(session);
      sessionMap.set(session, []);
    }
    sessionMap.get(session)?.push(entry);
  }

  return sessionOrder.map((session) => ({
    name: session,
    exerciseGroups: groupByExercise(sessionMap.get(session) ?? []),
  }));
}

function groupByExercise(entries: ProgramEntry[]): ExerciseGroup[] {
  const groups: ExerciseGroup[] = [];
  const map = new Map<string, ExerciseGroup>();
  for (const entry of entries) {
    if (!map.has(entry.exercise_name)) {
      const g: ExerciseGroup = { name: entry.exercise_name, entries: [] };
      groups.push(g);
      map.set(entry.exercise_name, g);
    }
    map.get(entry.exercise_name)?.entries.push(entry);
  }
  return groups;
}

export default function ProgramDetailPage() {
  const params = useParams();
  const programId = params.id as string;
  const { data: program, isLoading } = useProgram(programId);

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

  const sessions = groupBySession(program.entries || []);

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">{program.name}</h1>
          {program.description && (
            <p className="text-muted-foreground mt-1">{program.description}</p>
          )}
          {program.notes && <p className="text-sm text-muted-foreground mt-1">{program.notes}</p>}
        </div>
        <div className="flex gap-2">
          <Link href={`/programs/${program.id}/edit`}>
            <Button variant="outline">
              <Edit className="h-4 w-4 mr-2" />
              Edit
            </Button>
          </Link>
          <Link href={`/programs/${program.id}/convert`}>
            <Button>
              <ArrowRightLeft className="h-4 w-4 mr-2" />
              Convert to Plan
            </Button>
          </Link>
        </div>
      </div>

      {sessions.length === 0 ? (
        <p className="text-muted-foreground">No entries in this program.</p>
      ) : (
        <div className="space-y-4">
          {sessions.map((session) => (
            <Card key={session.name}>
              {session.name && (
                <CardHeader className="pb-2">
                  <CardTitle className="text-base">{session.name}</CardTitle>
                </CardHeader>
              )}
              <CardContent className={session.name ? '' : 'pt-4'}>
                <div className="divide-y">
                  {session.exerciseGroups.map((group) => (
                    <div
                      key={`${session.name}-${group.name}`}
                      className="flex items-baseline gap-3 py-1.5 first:pt-0 last:pb-0"
                    >
                      <span className="w-40 shrink-0 font-medium text-sm truncate">
                        {group.name}
                      </span>
                      <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-muted-foreground">
                        {group.entries.map((entry) => {
                          const label = entry.metadata?.label as string | undefined;
                          return (
                            <span key={entry.id} className="inline-flex items-center gap-1.5">
                              {label && <Badge variant="outline">{label}</Badge>}
                              {formatEntryText(entry)}
                            </span>
                          );
                        })}
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </main>
  );
}
