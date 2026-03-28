'use client';

import { ProgramForm } from '@/components/programs/ProgramForm';
import { useLoggedSessions, useProgram, useUpdateProgram } from '@/lib/hooks/usePrograms';
import type {
  Program,
  ProgramSessionCreate,
  ProgramSessionEntryCreate,
  ProgramTemplateEntryCreate,
  ProgramUpdate,
} from '@/types/api';
import { useParams, useRouter } from 'next/navigation';
import { useMemo, useState } from 'react';

function programToEntries(program: Program): ProgramTemplateEntryCreate[] {
  const entries: ProgramTemplateEntryCreate[] = [];
  let order = 0;

  for (const session of program.sessions) {
    for (const entry of session.entries) {
      const metadata: Record<string, unknown> = { session: session.session_name };
      if (session.date) metadata.date = session.date;
      if (entry.metadata?.label) metadata.label = entry.metadata.label;
      if (entry.load_kg != null) metadata.weight_kg = entry.load_kg;

      entries.push({
        exercise_name: entry.exercise_name,
        order,
        sets: entry.sets,
        reps: entry.reps,
        rpe: entry.rpe,
        metadata,
      });
      order++;
    }
  }

  return entries;
}

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

function entriesToProgramUpdate(
  name: string,
  notes: string,
  entries: ProgramTemplateEntryCreate[]
): ProgramUpdate {
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

  const sessions: ProgramSessionCreate[] = sessionOrder.map((sessionName, idx) => {
    const data = sessionMap.get(sessionName);
    const session: ProgramSessionCreate = {
      session_name: sessionName,
      order: idx,
      entries: data?.entries ?? [],
    };
    if (data?.date) session.date = data.date;
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
  const { data: loggedSessionsData, isLoading: sessionsLoading } = useLoggedSessions(programId);
  const updateProgram = useUpdateProgram();

  const [name, setName] = useState<string | null>(null);
  const [notes, setNotes] = useState<string | null>(null);

  const lockedSessionNames = useMemo(() => {
    if (!loggedSessionsData?.sessions) return new Set<string>();
    return new Set(loggedSessionsData.sessions.map((s) => s.session_name));
  }, [loggedSessionsData]);

  const initialEntries = useMemo(() => {
    if (!program) return undefined;
    return programToEntries(program);
  }, [program]);

  if (programLoading || sessionsLoading) {
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

  const handleSave = (entries: ProgramTemplateEntryCreate[]) => {
    const programName = name ?? program.name;
    const programNotes = notes ?? program.notes ?? '';
    const data = entriesToProgramUpdate(programName, programNotes, entries);

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
      {lockedSessionNames.size > 0 && (
        <p className="text-sm text-muted-foreground mb-4">
          Sessions with logged data are locked and cannot be modified.
        </p>
      )}
      <ProgramForm
        programName={name ?? program.name}
        programDescription=""
        programNotes={notes ?? program.notes ?? ''}
        initialEntries={initialEntries}
        onNameChange={setName}
        onDescriptionChange={() => {}}
        onNotesChange={setNotes}
        onSave={handleSave}
        isSaving={updateProgram.isPending}
        isEditing
        lockedSessionNames={lockedSessionNames}
      />
    </main>
  );
}
