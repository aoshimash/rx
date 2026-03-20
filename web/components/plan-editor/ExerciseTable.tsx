'use client';

import { Button } from '@/components/ui/button';
import type { PlanEntryCreate } from '@/types/api';
import { Plus } from 'lucide-react';
import { ExerciseGroup } from './ExerciseGroup';
import type { SetType } from './ExerciseRow';

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

  const handleSetTypeChange = (index: number, value: SetType | undefined) => {
    updateEntry(index, (e) => ({
      ...e,
      metadata: { ...e.metadata, set_type: value },
    }));
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
          onSetTypeChange={handleSetTypeChange}
          onSetsChange={(i, val) => updateEntry(i, (e) => ({ ...e, sets: val }))}
          onRepsChange={(i, val) => updateEntry(i, (e) => ({ ...e, reps: val }))}
          onLoadKgChange={(i, val) => updateEntry(i, (e) => ({ ...e, load_kg: val }))}
          onRpeChange={(i, val) => updateEntry(i, (e) => ({ ...e, rpe: val }))}
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
