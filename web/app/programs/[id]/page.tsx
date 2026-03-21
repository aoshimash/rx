'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { useDeleteProgram, useProgram } from '@/lib/hooks/usePrograms';
import { Trash2 } from 'lucide-react';
import { useParams, useRouter } from 'next/navigation';

export default function ProgramDetailPage() {
  const params = useParams();
  const router = useRouter();
  const programId = params.id as string;
  const { data: program, isLoading } = useProgram(programId);
  const deleteProgram = useDeleteProgram();

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

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-3xl font-bold">{program.name}</h1>
            <Badge variant={program.status === 'active' ? 'default' : 'secondary'}>
              {program.status}
            </Badge>
          </div>
          {program.notes && <p className="text-muted-foreground mt-1">{program.notes}</p>}
        </div>
        <Button variant="outline" onClick={handleDelete} disabled={deleteProgram.isPending}>
          <Trash2 className="h-4 w-4 mr-2" />
          Delete
        </Button>
      </div>

      {program.sessions.length === 0 ? (
        <p className="text-muted-foreground">No sessions in this program.</p>
      ) : (
        <div className="space-y-4">
          {program.sessions
            .slice()
            .sort((a, b) => a.order - b.order)
            .map((session) => (
              <Card key={session.id}>
                <CardHeader className="pb-2">
                  <div className="flex items-center gap-2">
                    <CardTitle className="text-base">{session.session_name}</CardTitle>
                    {session.date && (
                      <span className="text-sm text-muted-foreground">{session.date}</span>
                    )}
                  </div>
                </CardHeader>
                <CardContent>
                  {session.entries.length === 0 ? (
                    <p className="text-sm text-muted-foreground">No exercises</p>
                  ) : (
                    <div className="divide-y">
                      {session.entries
                        .slice()
                        .sort((a, b) => a.order - b.order)
                        .map((entry) => (
                          <div
                            key={entry.id}
                            className="flex items-baseline gap-3 py-1.5 first:pt-0 last:pb-0"
                          >
                            <span className="w-40 shrink-0 font-medium text-sm truncate">
                              {entry.exercise_name}
                            </span>
                            <span className="text-sm text-muted-foreground">
                              {[
                                entry.sets != null && `${entry.sets}sets`,
                                entry.reps != null && `${entry.reps}reps`,
                                entry.load_kg != null && `${entry.load_kg}kg`,
                                entry.rpe != null && `RPE${entry.rpe}`,
                              ]
                                .filter(Boolean)
                                .join(' ')}
                            </span>
                          </div>
                        ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
        </div>
      )}
    </main>
  );
}
