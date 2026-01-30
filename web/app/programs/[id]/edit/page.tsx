'use client';

import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { ProgramForm } from '@/components/program-editor/ProgramForm';
import { useExercises } from '@/lib/hooks/useExercises';
import { useProgram, useUpdateProgram, useDeleteProgram } from '@/lib/hooks/usePrograms';
import { Skeleton } from '@/components/ui/skeleton';
import type { ProgramNodeCreate } from '@/types/api';

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
  const [weeks, setWeeks] = useState<ProgramNodeCreate[]>([]);

  // Initialize form from loaded program
  useEffect(() => {
    if (program) {
      setProgramName(program.name);
      setProgramDescription(program.description || '');
      setWeeks(program.root_nodes || []);
    }
  }, [program]);

  const handleSave = async () => {
    await updateProgram.mutateAsync({
      id: programId,
      data: {
        name: programName,
        description: programDescription,
        root_nodes: weeks,
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
        <p className="text-muted-foreground mt-1">
          Modify weeks, days, and exercise prescriptions
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
        onDelete={handleDelete}
        isSaving={updateProgram.isPending}
        isEditing
      />
    </main>
  );
}
