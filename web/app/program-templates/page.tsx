'use client';

import { ProgramTemplateCard } from '@/components/program-templates/ProgramTemplateCard';
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
import { useCreateProgramTemplate, useProgramTemplates } from '@/lib/hooks/useProgramTemplates';
import type { ProgramTemplateEntryCreate } from '@/types/api';
import { Archive, Plus } from 'lucide-react';
import { useState } from 'react';

export default function ProgramTemplatesPage() {
  const [showArchived, setShowArchived] = useState(false);
  const { data: templatesData, isLoading } = useProgramTemplates(showArchived);
  const createTemplate = useCreateProgramTemplate();
  const templates = templatesData?.data || [];

  const [open, setOpen] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [notes, setNotes] = useState('');

  const handleSave = async (entries: ProgramTemplateEntryCreate[]) => {
    await createTemplate.mutateAsync({
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
          <h1 className="text-3xl font-bold">Templates</h1>
          <p className="text-muted-foreground mt-1">
            Reusable training templates with RPE/% prescriptions.
          </p>
        </div>
        <Button
          variant={showArchived ? 'secondary' : 'outline'}
          size="sm"
          onClick={() => setShowArchived(!showArchived)}
        >
          <Archive className="h-4 w-4 mr-2" />
          {showArchived ? 'Hide Archived' : 'Show Archived'}
        </Button>
      </div>

      {templates.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">
            No templates yet. Create your first training template.
          </p>
          <Button onClick={() => setOpen(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Create Template
          </Button>
        </div>
      ) : (
        <>
          <div className="grid gap-4 md:grid-cols-2">
            {templates.map((template) => (
              <ProgramTemplateCard key={template.id} template={template} />
            ))}
          </div>
          <div className="mt-6">
            <Button variant="outline" onClick={() => setOpen(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Create Template
            </Button>
          </div>
        </>
      )}

      <Dialog open={open} onOpenChange={handleOpenChange}>
        <DialogContent className="sm:max-w-4xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Create Template</DialogTitle>
            <DialogDescription>
              Define sessions and exercise prescriptions (RPE / %1RM)
            </DialogDescription>
          </DialogHeader>
          <ProgramForm
            programName={name}
            programDescription={description}
            programNotes={notes}
            onNameChange={setName}
            onDescriptionChange={setDescription}
            onNotesChange={setNotes}
            onSave={handleSave}
            isSaving={createTemplate.isPending}
          />
        </DialogContent>
      </Dialog>
    </main>
  );
}
