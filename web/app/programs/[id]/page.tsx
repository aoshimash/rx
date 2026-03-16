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

function groupBySession(entries: ProgramEntry[]): Map<string, ProgramEntry[]> {
  const groups = new Map<string, ProgramEntry[]>();
  for (const entry of entries) {
    const session = (entry.metadata?.session as string) ?? 'Ungrouped';
    if (!groups.has(session)) {
      groups.set(session, []);
    }
    groups.get(session)!.push(entry);
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

  const entries = program.entries || [];
  const sessionGroups = groupBySession(entries);

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

      {entries.length === 0 ? (
        <p className="text-muted-foreground">No entries in this program.</p>
      ) : (
        <div className="space-y-4">
          {Array.from(sessionGroups.entries()).map(([sessionName, sessionEntries]) => (
            <Card key={sessionName}>
              <CardHeader>
                <CardTitle>{sessionName}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {sessionEntries
                    .sort((a, b) => a.order - b.order)
                    .map((entry) => (
                      <div
                        key={entry.id}
                        className="flex items-center justify-between border rounded-lg p-3"
                      >
                        <div>
                          <span className="font-medium">{entry.exercise_name}</span>
                        </div>
                        <div className="flex gap-2">
                          {entry.sets && <Badge variant="secondary">{entry.sets} sets</Badge>}
                          {entry.reps && <Badge variant="secondary">{entry.reps} reps</Badge>}
                          {entry.rpe && <Badge variant="outline">RPE {entry.rpe}</Badge>}
                          {entry.percent_1rm !== undefined && (
                            <Badge variant="outline">
                              {Math.round(entry.percent_1rm * 100)}% 1RM
                            </Badge>
                          )}
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
