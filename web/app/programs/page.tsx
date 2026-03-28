'use client';

import { CreateProgramDialog } from '@/components/programs/CreateProgramDialog';
import {
  CreatedProgramCard,
  FinishedProgramCard,
  OngoingProgramCard,
} from '@/components/programs/ProgramCard';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { usePrograms } from '@/lib/hooks/usePrograms';
import { Plus } from 'lucide-react';
import { useState } from 'react';

const SHOW_FINISHED_KEY = 'programs:showFinished';
const OLD_SHOW_COMPLETED_KEY = 'programs:showCompleted';

export default function ProgramsPage() {
  const { data: programsData, isLoading } = usePrograms();
  const [showFinished, setShowFinished] = useState(() => {
    if (typeof window === 'undefined') return false;
    const oldValue = localStorage.getItem(OLD_SHOW_COMPLETED_KEY);
    if (oldValue !== null) {
      localStorage.setItem(SHOW_FINISHED_KEY, oldValue);
      localStorage.removeItem(OLD_SHOW_COMPLETED_KEY);
      return oldValue === 'true';
    }
    return localStorage.getItem(SHOW_FINISHED_KEY) === 'true';
  });
  const [open, setOpen] = useState(false);

  const allPrograms = [...(programsData?.data || [])].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  );
  const ongoingPrograms = allPrograms.filter((p) => p.status === 'ongoing');
  const createdPrograms = allPrograms.filter((p) => p.status === 'created');
  const finishedPrograms = allPrograms.filter(
    (p) => p.status === 'completed' || p.status === 'cancelled'
  );

  const toggleShowFinished = () => {
    const next = !showFinished;
    setShowFinished(next);
    localStorage.setItem(SHOW_FINISHED_KEY, String(next));
  };

  if (isLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <div className="space-y-4">
          <Skeleton className="h-[120px]" />
          <Skeleton className="h-[120px]" />
        </div>
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-start justify-between">
        <div>
          <h1 className="text-3xl font-bold">Programs</h1>
          <p className="text-muted-foreground mt-1">
            Concrete training programs with scheduled sessions.
          </p>
        </div>
      </div>

      {allPrograms.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">
            No programs yet. Create your first training program.
          </p>
          <Button onClick={() => setOpen(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Create Program
          </Button>
        </div>
      ) : (
        <>
          {ongoingPrograms.length > 0 && (
            <div className="mb-8">
              <h2 className="text-xs uppercase tracking-wider text-muted-foreground mb-3">
                Ongoing
              </h2>
              <div className="space-y-4">
                {ongoingPrograms.map((program) => (
                  <OngoingProgramCard key={program.id} program={program} />
                ))}
              </div>
            </div>
          )}

          {createdPrograms.length > 0 && (
            <div className="mb-8">
              <h2 className="text-xs uppercase tracking-wider text-muted-foreground mb-3">
                Programs
              </h2>
              <div className="grid grid-cols-2 gap-4">
                {createdPrograms.map((program) => (
                  <CreatedProgramCard key={program.id} program={program} />
                ))}
              </div>
            </div>
          )}

          {finishedPrograms.length > 0 && (
            <div className="mb-8">
              <Button variant="ghost" size="sm" onClick={toggleShowFinished}>
                {showFinished
                  ? `Hide finished (${finishedPrograms.length})`
                  : `Show finished (${finishedPrograms.length})`}
              </Button>
              {showFinished && (
                <div className="grid grid-cols-2 gap-4 mt-3">
                  {finishedPrograms.map((program) => (
                    <FinishedProgramCard key={program.id} program={program} />
                  ))}
                </div>
              )}
            </div>
          )}

          <div className="mt-6">
            <Button variant="outline" onClick={() => setOpen(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Create Program
            </Button>
          </div>
        </>
      )}

      <CreateProgramDialog open={open} onOpenChange={setOpen} />
    </main>
  );
}
