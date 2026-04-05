'use client';

import { ProgramForm, type ProgramFormEntry } from '@/components/programs/ProgramForm';
import { useProgram, useUpdateProgram } from '@/lib/hooks/usePrograms';
import type {
  Program,
  ProgramSessionCreate,
  ProgramSessionEntryCreate,
  ProgramUpdate,
} from '@/types/api';
import { useParams, useRouter } from 'next/navigation';
import { useMemo, useState } from 'react';

function programToEntries(program: Program): ProgramFormEntry[] {
  const entries: ProgramFormEntry[] = [];
  let order = 0;

  for (const session of program.sessions) {
    for (const entry of session.entries) {
      const metadata: Record<string, unknown> = { session: session.session_name };
      if (session.date) metadata.date = session.date;
      if (session.field_group_id) metadata.field_group_id = session.field_group_id;
      if (entry.fields?.label) metadata.label = entry.fields.label;

      entries.push({
        exercise_name: entry.exercise_name,
        order,
        fields: entry.fields,
        notes: entry.notes,
        metadata,
      });
      order++;
    }
  }

  return entries;
}

function convertEntryToProgramEntry(
  entry: ProgramFormEntry,
  order: number
): ProgramSessionEntryCreate {
  const fields = entry.metadata?.label
    ? { ...entry.fields, label: entry.metadata.label }
    : entry.fields;
  return {
    exercise_name: entry.exercise_name,
    order,
    fields,
  };
}

function entriesToProgramUpdate(
  name: string,
  notes: string,
  entries: ProgramFormEntry[]
): ProgramUpdate {
  const sessionMap = new Map<
    string,
    { date?: string; field_group_id?: string; entries: ProgramSessionEntryCreate[] }
  >();
  const sessionOrder: string[] = [];

  for (const entry of entries) {
    const sessionName = (entry.metadata?.session as string) ?? 'Session 1';
    if (!sessionMap.has(sessionName)) {
      sessionMap.set(sessionName, {
        date: entry.metadata?.date as string | undefined,
        field_group_id: entry.metadata?.field_group_id as string | undefined,
        entries: [],
      });
      sessionOrder.push(sessionName);
    }

    const sessionData = sessionMap.get(sessionName);
    if (sessionData) {
      sessionData.entries.push(convertEntryToProgramEntry(entry, sessionData.entries.length));
    }
  }

  const sessions: ProgramSessionCreate[] = sessionOrder.map((sessionName, idx) => {
    const data = sessionMap.get(sessionName);
    const session: ProgramSessionCreate = {
      session_name: sessionName,
      order: idx,
      entries: data?.entries ?? [],
    };
    if (data?.date) session.date = data.date;
    if (data?.field_group_id) session.field_group_id = data.field_group_id;
    return session;
  });

  return {
    name,
    notes: notes || undefined,
    sessions,
  };
}

export default function EditProgramPage() {
  const params = useParams();
  const router = useRouter();
  const programId = params.id as string;

  const { data: program, isLoading: programLoading } = useProgram(programId);
  const updateProgram = useUpdateProgram();

  const [name, setName] = useState<string | null>(null);
  const [notes, setNotes] = useState<string | null>(null);

  const initialEntries = useMemo(() => {
    if (!program) return undefined;
    return programToEntries(program);
  }, [program]);

  if (programLoading) {
    return (
      <main className="container max-w-4xl mx-auto py-8 px-4">
        <p className="text-muted-foreground">Loading...</p>
      </main>
    );
  }

  if (!program) {
    return (
      <main className="container max-w-4xl mx-auto py-8 px-4">
        <p className="text-muted-foreground">Program not found.</p>
      </main>
    );
  }

  const handleSave = (entries: ProgramFormEntry[]) => {
    const programName = name ?? program.name;
    const programNotes = notes ?? program.notes ?? '';
    const data = entriesToProgramUpdate(programName, programNotes, entries);
    data.status = program.status;

    updateProgram.mutate(
      { id: programId, data },
      {
        onSuccess: () => {
          router.push(`/programs/${programId}`);
        },
      }
    );
  };

  return (
    <main className="container max-w-4xl mx-auto py-8 px-4">
      <h1 className="text-2xl font-bold mb-6">Edit Program</h1>
      <ProgramForm
        programName={name ?? program.name}
        programNotes={notes ?? program.notes ?? ''}
        initialEntries={initialEntries}
        onNameChange={setName}
        onNotesChange={setNotes}
        onSave={handleSave}
        isSaving={updateProgram.isPending}
        isEditing
      />
    </main>
  );
}
