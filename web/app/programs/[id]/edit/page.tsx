'use client';

import { ProgramForm } from '@/components/program-editor/ProgramForm';
import { Skeleton } from '@/components/ui/skeleton';
import { useExercises } from '@/lib/hooks/useExercises';
import { useDeleteProgram, useProgram, useUpdateProgram } from '@/lib/hooks/usePrograms';
import type { ProgramEntryCreate } from '@/types/api';
import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

export default function EditProgramPage() {
  const router = useRouter();
  const params = useParams();
  const programId = params.id as string;

  const { data: exercisesData, isLoading: exercisesLoading } = useExercises();
  const { data: program, isLoading: programLoading } = useProgram(programId);
  const updateProgram = useUpdateProgram();
  const deleteProgram = useDeleteProgram();

  const [programName, setProgramName] = useState('');
  const [programDescription, setProgramDescription] = useState('');

  useEffect(() => {
    if (program) {
      setProgramName(program.name);
      setProgramDescription(program.description || '');
    }
  }, [program]);

  const handleSave = async (entries: ProgramEntryCreate[]) => {
    await updateProgram.mutateAsync({
      id: programId,
      data: {
        name: programName,
        description: programDescription,
        entries,
      },
    });
    router.push('/programs');
  };

  const handleDelete = async () => {
    await deleteProgram.mutateAsync(programId);
    router.push('/programs');
  };

  if (exercisesLoading || programLoading) {
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
        <h1 className="text-3xl font-bold">Edit Training Program</h1>
        <p className="text-muted-foreground mt-1">Modify weeks, days, and exercise prescriptions</p>
      </div>

      <ProgramForm
        key={programId}
        programName={programName}
        programDescription={programDescription}
        initialEntries={program?.entries}
        availableExercises={exercises}
        onNameChange={setProgramName}
        onDescriptionChange={setProgramDescription}
        onSave={handleSave}
        onDelete={handleDelete}
        isSaving={updateProgram.isPending}
        isEditing
      />
    </main>
  );
}
