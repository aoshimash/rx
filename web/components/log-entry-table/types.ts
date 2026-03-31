import type { ProgramSessionEntry } from '@/types/api';

export interface TableEntry {
  id: string;
  exercise_name: string;
  // Fallback fields (used when no FieldGroup is provided)
  sets: number | undefined;
  reps: number | undefined;
  load_kg: number | undefined;
  /** true when sets or reps were manually edited away from plan values */
  setsEdited: boolean;
  repsEdited: boolean;
  // Dynamic field values (keyed by FieldDef.name, used in FieldDef-driven mode)
  logFieldValues: Record<string, unknown>;
  /** Object key for video uploaded via pre-signed URL (video FieldDef only) */
  videoObjectKey: string | undefined;
  notes: string;
  fields: Record<string, unknown>;
  plan?: {
    fields?: Record<string, unknown>;
  };
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
    logFieldValues: {},
    videoObjectKey: undefined,
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
    logFieldValues: {},
    videoObjectKey: undefined,
    setsEdited: false,
    repsEdited: false,
  };
}
