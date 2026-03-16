'use client';

import { ProgramForm } from '@/components/programs/ProgramForm';
import { Skeleton } from '@/components/ui/skeleton';
import { useDeleteProgram, useProgram, useUpdateProgram } from '@/lib/hooks/usePrograms';
import type { ProgramEntryCreate } from '@/types/api';
import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

export default function EditProgramPage() {
  const router = useRouter();
  const params = useParams();
  const programId = params.id as string;

  const { data: program, isLoading } = useProgram(programId);
  const updateProgram = useUpdateProgram();
  const deleteProgram = useDeleteProgram();

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [notes, setNotes] = useState('');

  useEffect(() => {
    if (program) {
      setName(program.name);
      setDescription(program.description || '');
      setNotes(program.notes || '');
    }
  }, [program]);

  const handleSave = async (entries: ProgramEntryCreate[]) => {
    await updateProgram.mutateAsync({
      id: programId,
      data: {
        name,
        description: description || undefined,
        notes: notes || undefined,
        entries,
      },
    });
    router.push('/programs');
  };

  const handleDelete = async () => {
    await deleteProgram.mutateAsync(programId);
    router.push('/programs');
  };

  if (isLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <Skeleton className="h-[600px] w-full" />
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Edit Program</h1>
        <p className="text-muted-foreground mt-1">Modify sessions and exercise prescriptions</p>
      </div>

      <ProgramForm
        key={programId}
        programName={name}
        programDescription={description}
        programNotes={notes}
        initialEntries={program?.entries}
        onNameChange={setName}
        onDescriptionChange={setDescription}
        onNotesChange={setNotes}
        onSave={handleSave}
        onDelete={handleDelete}
        isSaving={updateProgram.isPending}
        isEditing
      />
    </main>
  );
}
