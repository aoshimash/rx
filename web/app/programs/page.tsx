'use client';

import { ProgramCard } from '@/components/programs/ProgramCard';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { usePrograms } from '@/lib/hooks/usePrograms';
import { useProgramStore } from '@/stores/program';
import { Plus } from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

export default function ProgramsPage() {
  const router = useRouter();
  const { data: programsData, isLoading } = usePrograms();
  const { selectedProgramId, setSelectedProgram } = useProgramStore();

  const programs = programsData?.data || [];

  const handleSelectProgram = (programId: string) => {
    setSelectedProgram(programId);
    router.push('/');
  };

  if (isLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <Skeleton className="h-[200px]" />
          <Skeleton className="h-[200px]" />
          <Skeleton className="h-[200px]" />
        </div>
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Training Programs</h1>
          <p className="text-muted-foreground mt-1">
            Select a program to view in Week View
          </p>
        </div>
        <Link href="/programs/new">
          <Button>
            <Plus className="h-4 w-4 mr-2" />
            New Program
          </Button>
        </Link>
      </div>

      {programs.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">
            No programs yet. Create your first training program.
          </p>
          <Link href="/programs/new">
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              Create Program
            </Button>
          </Link>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {programs.map((program) => (
            <ProgramCard
              key={program.id}
              program={program}
              onSelect={() => handleSelectProgram(program.id)}
              isSelected={selectedProgramId === program.id}
            />
          ))}
        </div>
      )}
    </main>
  );
}
