'use client';

import { Button } from '@/components/ui/button';
import { calculateE1rm } from '@/lib/utils/e1rm';
import type { PlanEntryCreate } from '@/types/api';
import { Plus } from 'lucide-react';
import { ExerciseGroup } from './ExerciseGroup';

interface ExerciseTableProps {
  exercises: PlanEntryCreate[];
  onChange: (exercises: PlanEntryCreate[]) => void;
}

type GroupData = {
  name: string;
  indices: number[];
};

function groupByExercise(exercises: PlanEntryCreate[]): GroupData[] {
  const groups: GroupData[] = [];
  const nameToGroup = new Map<string, GroupData>();

  exercises.forEach((ex, i) => {
    const name = ex.exercise_name;
    if (!nameToGroup.has(name)) {
      const group: GroupData = { name, indices: [] };
      groups.push(group);
      nameToGroup.set(name, group);
    }
    nameToGroup.get(name)?.indices.push(i);
  });

  return groups;
}

/** Given e1RM, reps, and RPE, solve for load. Returns null if invalid. */
function loadFromE1rm(e1rm: number, reps: number, rpe: number): number | null {
  const effectiveReps = reps + (10 - rpe);
  if (effectiveReps >= 37) return null;
  return (e1rm * (37 - effectiveReps)) / 36;
}

/** Given e1RM, reps, and load, solve for RPE. Returns null if invalid. */
function rpeFromE1rm(e1rm: number, reps: number, load: number): number | null {
  // e1rm = load × 36 / (37 - reps - (10 - rpe))
  // 27 - reps + rpe = load × 36 / e1rm
  // rpe = load × 36 / e1rm - 27 + reps
  const rpe = (load * 36) / e1rm - 27 + reps;
  if (rpe < 1 || rpe > 10) return null;
  return rpe;
}

function roundToHalf(n: number): number {
  return Math.round(n * 2) / 2;
}

export function ExerciseTable({ exercises, onChange }: ExerciseTableProps) {
  const groups = groupByExercise(exercises);

  const reindex = (arr: PlanEntryCreate[]): PlanEntryCreate[] =>
    arr.map((e, i) => ({ ...e, order: i }));

  const updateEntry = (index: number, updater: (e: PlanEntryCreate) => PlanEntryCreate) => {
    const next = [...exercises];
    if (next[index]) next[index] = updater(next[index]);
    onChange(next);
  };

  const handleRename = (oldName: string, newName: string) => {
    onChange(
      exercises.map((e) => (e.exercise_name === oldName ? { ...e, exercise_name: newName } : e))
    );
  };

  const handleRemoveExercise = (name: string) => {
    onChange(reindex(exercises.filter((e) => e.exercise_name !== name)));
  };

  const handleRemoveSet = (index: number) => {
    onChange(reindex(exercises.filter((_, i) => i !== index)));
  };

  const handleAddSet = (exerciseName: string) => {
    const lastIdx = exercises.map((e) => e.exercise_name).lastIndexOf(exerciseName);
    const newEntry: PlanEntryCreate = {
      exercise_name: exerciseName,
      order: exercises.length,
      sets: 3,
      reps: 5,
    };
    const next = [...exercises];
    next.splice(lastIdx + 1, 0, newEntry);
    onChange(reindex(next));
  };

  const handleAddExercise = () => {
    onChange([
      ...exercises,
      { exercise_name: 'Exercise', order: exercises.length, sets: 3, reps: 5 },
    ]);
  };

  const handleLabelChange = (index: number, value: string | undefined) => {
    updateEntry(index, (e) => ({
      ...e,
      metadata: { ...e.metadata, label: value },
    }));
  };

  const handleRpeChange = (index: number, newRpe: number | undefined) => {
    const entry = exercises[index];
    if (!entry) return;

    // If all three old values exist, recalculate load to keep e1RM constant
    if (newRpe != null && entry.load_kg != null && entry.reps != null && entry.rpe != null) {
      const e1rm = calculateE1rm(entry.load_kg, entry.reps, entry.rpe);
      if (e1rm != null) {
        const newLoad = loadFromE1rm(e1rm, entry.reps, newRpe);
        if (newLoad != null) {
          updateEntry(index, (e) => ({ ...e, rpe: newRpe, load_kg: roundToHalf(newLoad) }));
          return;
        }
      }
    }
    updateEntry(index, (e) => ({ ...e, rpe: newRpe }));
  };

  const handleLoadKgChange = (index: number, newLoad: number | undefined) => {
    const entry = exercises[index];
    if (!entry) return;

    // If all three old values exist, recalculate RPE to keep e1RM constant
    if (newLoad != null && entry.rpe != null && entry.reps != null && entry.load_kg != null) {
      const e1rm = calculateE1rm(entry.load_kg, entry.reps, entry.rpe);
      if (e1rm != null) {
        const newRpe = rpeFromE1rm(e1rm, entry.reps, newLoad);
        if (newRpe != null) {
          updateEntry(index, (e) => ({ ...e, load_kg: newLoad, rpe: roundToHalf(newRpe) }));
          return;
        }
      }
    }
    updateEntry(index, (e) => ({ ...e, load_kg: newLoad }));
  };

  return (
    <div className="space-y-3">
      {groups.map((group) => (
        <ExerciseGroup
          key={group.name}
          exerciseName={group.name}
          entries={group.indices.flatMap((i) => {
            const entry = exercises[i];
            return entry ? [{ entry, flatIndex: i }] : [];
          })}
          onRename={(newName) => handleRename(group.name, newName)}
          onRemoveExercise={() => handleRemoveExercise(group.name)}
          onAddSet={() => handleAddSet(group.name)}
          onLabelChange={handleLabelChange}
          onRpeChange={handleRpeChange}
          onLoadKgChange={handleLoadKgChange}
          onRepsChange={(i, val) => updateEntry(i, (e) => ({ ...e, reps: val }))}
          onSetsChange={(i, val) => updateEntry(i, (e) => ({ ...e, sets: val }))}
          onRemoveSet={handleRemoveSet}
        />
      ))}
      <Button variant="outline" onClick={handleAddExercise} className="w-full">
        <Plus className="h-4 w-4 mr-2" />
        Add Exercise
      </Button>
    </div>
  );
}
