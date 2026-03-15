'use client';

import { ProgramForm } from '@/components/program-editor/ProgramForm';
import { Skeleton } from '@/components/ui/skeleton';
import { useExercises } from '@/lib/hooks/useExercises';
import { useCreateProgram } from '@/lib/hooks/usePrograms';
import type { ProgramEntryCreate } from '@/types/api';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function NewProgramPage() {
  const router = useRouter();
  const { data: exercisesData, isLoading: exercisesLoading } = useExercises();
  const createProgram = useCreateProgram();

  const [programName, setProgramName] = useState('');
  const [programDescription, setProgramDescription] = useState('');

  const handleSave = async (entries: ProgramEntryCreate[]) => {
    await createProgram.mutateAsync({
      name: programName,
      description: programDescription,
      entries,
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
        <p className="text-muted-foreground mt-1">Define weeks, days, and exercise prescriptions</p>
      </div>

      <ProgramForm
        programName={programName}
        programDescription={programDescription}
        availableExercises={exercises}
        onNameChange={setProgramName}
        onDescriptionChange={setProgramDescription}
        onSave={handleSave}
        isSaving={createProgram.isPending}
      />
    </main>
  );
}
