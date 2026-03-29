import type { ProgramSessionEntry } from '@/types/api';

export interface TableEntry {
  id: string;
  exercise_name: string;
  sets: number | undefined;
  reps: number | undefined;
  load_kg: number | undefined;
  notes: string;
  fields: Record<string, unknown>;
  plan?: {
    fields?: Record<string, unknown>;
  };
  /** true when sets or reps were manually edited away from plan values */
  setsEdited: boolean;
  repsEdited: boolean;
}

export function createTableEntryFromPlan(entry: ProgramSessionEntry, index: number): TableEntry {
  const sets = entry.fields?.sets as number | undefined;
  const reps = entry.fields?.reps as number | undefined;
  const load_kg = entry.fields?.load_kg as number | undefined;
  return {
    id: `plan-${index}-${entry.id}`,
    exercise_name: entry.exercise_name,
    sets,
    reps,
    load_kg,
    notes: '',
    fields: entry.fields ?? {},
    plan: {
      fields: entry.fields,
    },
    setsEdited: false,
    repsEdited: false,
  };
}

export function createEmptyTableEntry(): TableEntry {
  return {
    id: `new-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
    exercise_name: '',
    sets: undefined,
    reps: undefined,
    load_kg: undefined,
    notes: '',
    fields: {},
    setsEdited: false,
    repsEdited: false,
  };
}
