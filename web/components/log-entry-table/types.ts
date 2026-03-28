import type { ProgramSessionEntry } from '@/types/api';

export interface TableEntry {
  id: string;
  exercise_name: string;
  sets: number | undefined;
  reps: number | undefined;
  load_kg: number | undefined;
  rpe: number | undefined;
  notes: string;
  plan?: {
    sets?: number;
    reps?: number;
    load_kg?: number;
    rpe?: number;
  };
  /** true when sets or reps were manually edited away from plan values */
  setsEdited: boolean;
  repsEdited: boolean;
}

export function createTableEntryFromPlan(entry: ProgramSessionEntry, index: number): TableEntry {
  return {
    id: `plan-${index}-${entry.id}`,
    exercise_name: entry.exercise_name,
    sets: entry.sets,
    reps: entry.reps,
    load_kg: entry.load_kg,
    rpe: entry.rpe,
    notes: '',
    plan: {
      sets: entry.sets,
      reps: entry.reps,
      load_kg: entry.load_kg,
      rpe: entry.rpe,
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
    rpe: undefined,
    notes: '',
    setsEdited: false,
    repsEdited: false,
  };
}
