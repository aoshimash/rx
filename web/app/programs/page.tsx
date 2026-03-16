'use client';

import { ProgramCard } from '@/components/programs/ProgramCard';
import { ProgramForm } from '@/components/programs/ProgramForm';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Skeleton } from '@/components/ui/skeleton';
import { useCreateProgram, usePrograms } from '@/lib/hooks/usePrograms';
import type { ProgramEntryCreate } from '@/types/api';
import { Plus } from 'lucide-react';
import { useState } from 'react';

export default function ProgramsPage() {
  const { data: programsData, isLoading } = usePrograms();
  const createProgram = useCreateProgram();
  const programs = programsData?.data || [];

  const [open, setOpen] = useState(false);
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
    setOpen(false);
    setName('');
    setDescription('');
    setNotes('');
  };

  const handleOpenChange = (value: boolean) => {
    setOpen(value);
    if (!value) {
      setName('');
      setDescription('');
      setNotes('');
    }
  };

  if (isLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <div className="grid gap-4 md:grid-cols-2">
          <Skeleton className="h-[200px]" />
          <Skeleton className="h-[200px]" />
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
            Reusable training templates. Convert to a Plan with target weights.
          </p>
        </div>
        <Button onClick={() => setOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Create Program
        </Button>
      </div>

      {programs.length === 0 ? (
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
        <div className="grid gap-4 md:grid-cols-2">
          {programs.map((program) => (
            <ProgramCard key={program.id} program={program} />
          ))}
        </div>
      )}

      <Dialog open={open} onOpenChange={handleOpenChange}>
        <DialogContent className="sm:max-w-4xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Create Program</DialogTitle>
            <DialogDescription>Define sessions and exercise prescriptions</DialogDescription>
          </DialogHeader>
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
        </DialogContent>
      </Dialog>
    </main>
  );
}
