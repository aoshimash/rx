'use client';

import { ConvertToPlanForm } from '@/components/programs/ConvertToPlanForm';
import { Skeleton } from '@/components/ui/skeleton';
import { useProgram } from '@/lib/hooks/usePrograms';
import { useParams } from 'next/navigation';

export default function ConvertProgramPage() {
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

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Convert to Plan</h1>
        <p className="text-muted-foreground mt-1">
          Set target weights for each exercise to generate a concrete training plan from &quot;
          {program.name}&quot;
        </p>
      </div>

      <ConvertToPlanForm program={program} />
    </main>
  );
}
