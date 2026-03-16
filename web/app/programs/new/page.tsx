'use client';

import { ProgramForm } from '@/components/programs/ProgramForm';
import { useCreateProgram } from '@/lib/hooks/usePrograms';
import type { ProgramEntryCreate } from '@/types/api';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function NewProgramPage() {
  const router = useRouter();
  const createProgram = useCreateProgram();

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [notes, setNotes] = useState('');

  const handleSave = async (entries: ProgramEntryCreate[]) => {
    await createProgram.mutateAsync({
      name,
      description: description || undefined,
      notes: notes || undefined,
      entries,
    });
    router.push('/programs');
  };

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Create Program</h1>
        <p className="text-muted-foreground mt-1">Define sessions and exercise prescriptions</p>
      </div>

      <ProgramForm
        programName={name}
        programDescription={description}
        programNotes={notes}
        onNameChange={setName}
        onDescriptionChange={setDescription}
        onNotesChange={setNotes}
        onSave={handleSave}
        isSaving={createProgram.isPending}
      />
    </main>
  );
}
