'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import type { PlanEntryCreate } from '@/types/api';
import { Plus, Trash2 } from 'lucide-react';
import { ExerciseRow } from './ExerciseRow';

interface ExerciseGroupProps {
  exerciseName: string;
  entries: { entry: PlanEntryCreate; flatIndex: number }[];
  onRename: (newName: string) => void;
  onRemoveExercise: () => void;
  onAddSet: () => void;
  onLabelChange: (flatIndex: number, value: string | undefined) => void;
  onRpeChange: (flatIndex: number, value: number | undefined) => void;
  onLoadKgChange: (flatIndex: number, value: number | undefined) => void;
  onRepsChange: (flatIndex: number, value: number) => void;
  onSetsChange: (flatIndex: number, value: number) => void;
  onRemoveSet: (flatIndex: number) => void;
}

export function ExerciseGroup({
  exerciseName,
  entries,
  onRename,
  onRemoveExercise,
  onAddSet,
  onLabelChange,
  onRpeChange,
  onLoadKgChange,
  onRepsChange,
  onSetsChange,
  onRemoveSet,
}: ExerciseGroupProps) {
  return (
    <div className="border rounded-lg p-3 space-y-2">
      {/* Header: exercise name + remove button */}
      <div className="flex items-center justify-between gap-2">
        <Input
          value={exerciseName}
          onChange={(e) => onRename(e.target.value)}
          placeholder="e.g., Squat, Bench Press"
          className="font-semibold border-none shadow-none p-0 h-auto focus-visible:ring-0 bg-transparent"
        />
        <Button
          variant="ghost"
          size="sm"
          className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive shrink-0"
          onClick={onRemoveExercise}
        >
          <Trash2 className="h-3.5 w-3.5" />
        </Button>
      </div>

      {/* Column headers */}
      {entries.length > 0 && (
        <div className="flex items-center gap-2 px-0.5">
          <span className="w-[110px] text-xs text-muted-foreground">Label</span>
          <span
            className="text-xs text-muted-foreground text-center"
            style={{ width: 'calc(4rem + 5rem + 0.75rem + 0.5rem + 0.5rem)' }}
          >
            RPE / Load
          </span>
          <span className="w-14 text-xs text-muted-foreground text-center">Reps</span>
          <span className="w-14 text-xs text-muted-foreground text-center">Sets</span>
        </div>
      )}

      {/* Set rows */}
      <div className="space-y-1.5">
        {entries.map(({ entry, flatIndex }) => {
          const hasRpe = entry.rpe != null;
          const hasLoad = entry.load_kg != null;
          const hasReps = entry.reps != null;
          const linked = hasRpe && hasLoad && hasReps;
          return (
            <ExerciseRow
              key={flatIndex}
              label={(entry.metadata?.label as string) ?? undefined}
              rpe={entry.rpe}
              loadKg={entry.load_kg}
              reps={entry.reps}
              sets={entry.sets}
              onLabelChange={(val) => onLabelChange(flatIndex, val)}
              onRpeChange={(val) => onRpeChange(flatIndex, val)}
              onLoadKgChange={(val) => onLoadKgChange(flatIndex, val)}
              onRepsChange={(val) => onRepsChange(flatIndex, val)}
              onSetsChange={(val) => onSetsChange(flatIndex, val)}
              onRemove={() => onRemoveSet(flatIndex)}
              linked={linked}
            />
          );
        })}
      </div>

      {/* Add set button */}
      <Button
        variant="ghost"
        size="sm"
        className="h-7 text-xs text-muted-foreground w-full"
        onClick={onAddSet}
      >
        <Plus className="h-3 w-3 mr-1" />
        セット追加
      </Button>
    </div>
  );
}
