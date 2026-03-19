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

const SET_TYPE_LABELS: Record<string, string> = {
  top: 'Top',
  main: 'メイン',
  backoff: 'バックオフ',
};

type ExerciseGroup = { name: string; entries: ProgramEntry[] };

function groupByExercise(entries: ProgramEntry[]): ExerciseGroup[] {
  const groups: ExerciseGroup[] = [];
  const map = new Map<string, ExerciseGroup>();
  for (const entry of [...entries].sort((a, b) => a.order - b.order)) {
    if (!map.has(entry.exercise_name)) {
      const g: ExerciseGroup = { name: entry.exercise_name, entries: [] };
      groups.push(g);
      map.set(entry.exercise_name, g);
    }
    map.get(entry.exercise_name)!.entries.push(entry);
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

  const groups = groupByExercise(program.entries || []);

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

      {groups.length === 0 ? (
        <p className="text-muted-foreground">No entries in this program.</p>
      ) : (
        <div className="space-y-4">
          {groups.map((group) => (
            <Card key={group.name}>
              <CardHeader className="pb-2">
                <CardTitle className="text-base">{group.name}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {group.entries.map((entry) => {
                    const setType = entry.metadata?.set_type as string | undefined;
                    return (
                      <div
                        key={entry.id}
                        className="flex items-center justify-between rounded-md px-2 py-1.5 hover:bg-muted/50"
                      >
                        <div className="w-24">
                          {setType ? (
                            <Badge variant="secondary" className="text-xs">
                              {SET_TYPE_LABELS[setType] ?? setType}
                            </Badge>
                          ) : (
                            <span className="text-xs text-muted-foreground">-</span>
                          )}
                        </div>
                        <div className="flex gap-2">
                          {entry.sets != null && (
                            <Badge variant="outline">{entry.sets} sets</Badge>
                          )}
                          {entry.reps != null && (
                            <Badge variant="outline">{entry.reps} reps</Badge>
                          )}
                          {entry.rpe != null && (
                            <Badge variant="outline">RPE {entry.rpe}</Badge>
                          )}
                          {entry.percent_1rm != null && (
                            <Badge variant="outline">
                              {Math.round(entry.percent_1rm * 100)}% 1RM
                            </Badge>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </main>
  );
}
