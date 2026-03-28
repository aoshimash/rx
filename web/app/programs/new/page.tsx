'use client';

import { ProgramForm } from '@/components/programs/ProgramForm';
import { useCreateProgram } from '@/lib/hooks/usePrograms';
import type {
  ProgramSessionCreate,
  ProgramSessionEntryCreate,
  ProgramTemplateEntryCreate,
} from '@/types/api';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { toast } from 'sonner';

function convertEntryToProgramEntry(
  entry: ProgramTemplateEntryCreate,
  order: number
): ProgramSessionEntryCreate {
  const programEntry: ProgramSessionEntryCreate = {
    exercise_name: entry.exercise_name,
    order,
    sets: entry.sets,
    reps: entry.reps,
    rpe: entry.rpe,
  };

  const weightKg = entry.metadata?.weight_kg as number | undefined;
  if (weightKg != null) {
    programEntry.load_kg = weightKg;
  }

  if (entry.metadata) {
    const { session: _s, date: _d, weight_kg: _w, label, ...rest } = entry.metadata;
    if (label || Object.keys(rest).length > 0) {
      programEntry.metadata = { ...(label ? { label } : {}), ...rest };
    }
  }

  return programEntry;
}

function entriesToSessions(entries: ProgramTemplateEntryCreate[]): ProgramSessionCreate[] {
  const sessionMap = new Map<string, { date?: string; entries: ProgramSessionEntryCreate[] }>();
  const sessionOrder: string[] = [];

  for (const entry of entries) {
    const sessionName = (entry.metadata?.session as string) ?? 'Session 1';
    if (!sessionMap.has(sessionName)) {
      sessionMap.set(sessionName, {
        date: entry.metadata?.date as string | undefined,
        entries: [],
      });
      sessionOrder.push(sessionName);
    }

    const sessionData = sessionMap.get(sessionName);
    if (sessionData) {
      sessionData.entries.push(convertEntryToProgramEntry(entry, sessionData.entries.length));
    }
  }

  return sessionOrder.map((sessionName, idx) => {
    const data = sessionMap.get(sessionName);
    const session: ProgramSessionCreate = {
      session_name: sessionName,
      order: idx,
      entries: data?.entries ?? [],
    };
    if (data?.date) session.date = data.date;
    return session;
  });
}

export default function NewProgramPage() {
  const router = useRouter();
  const createProgram = useCreateProgram();
  const [name, setName] = useState('');
  const [notes, setNotes] = useState('');

  const handleSave = (entries: ProgramTemplateEntryCreate[]) => {
    const sessions = entriesToSessions(entries);
    createProgram.mutate(
      { name, notes: notes || undefined, sessions },
      {
        onSuccess: () => {
          router.push('/programs');
        },
        onError: (err) => {
          const isConflict =
            err &&
            typeof err === 'object' &&
            'response' in err &&
            (err as { response?: { status?: number } }).response?.status === 409;
          toast.error(
            isConflict ? 'A program with this name already exists' : 'Failed to create program'
          );
        },
      }
    );
  };

  return (
    <main className="container max-w-4xl mx-auto py-8 px-4">
      <h1 className="text-2xl font-bold mb-6">Create Program</h1>
      <ProgramForm
        programName={name}
        programDescription=""
        programNotes={notes}
        onNameChange={setName}
        onDescriptionChange={() => {}}
        onNotesChange={setNotes}
        onSave={handleSave}
        isSaving={createProgram.isPending}
      />
    </main>
  );
}
