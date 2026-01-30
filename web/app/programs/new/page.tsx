'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ProgramForm } from '@/components/program-editor/ProgramForm';
import { useExercises } from '@/lib/hooks/useExercises';
import { useCreateProgram } from '@/lib/hooks/usePrograms';
import { Skeleton } from '@/components/ui/skeleton';
import type { ProgramNodeCreate } from '@/types/api';

export default function NewProgramPage() {
  const router = useRouter();
  const { data: exercisesData, isLoading: exercisesLoading } = useExercises();
  const createProgram = useCreateProgram();

  const [programName, setProgramName] = useState('');
  const [programDescription, setProgramDescription] = useState('');
  const [weeks, setWeeks] = useState<ProgramNodeCreate[]>([
    {
      name: 'Week 1',
      node_type: 'week',
      order: 0,
      children: [],
    },
  ]);

  const handleSave = async () => {
    await createProgram.mutateAsync({
      name: programName,
      description: programDescription,
      root_nodes: weeks,
    });
    router.push('/programs');
  };

  if (exercisesLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <Skeleton className="h-[600px] w-full" />
      </main>
    );
  }

  const exercises = exercisesData?.data || [];

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Create Training Program</h1>
        <p className="text-muted-foreground mt-1">
          Define weeks, days, and exercise prescriptions
        </p>
      </div>

      <ProgramForm
        programName={programName}
        programDescription={programDescription}
        weeks={weeks}
        availableExercises={exercises}
        onNameChange={setProgramName}
        onDescriptionChange={setProgramDescription}
        onWeeksChange={setWeeks}
        onSave={handleSave}
        isSaving={createProgram.isPending}
      />
    </main>
  );
}
