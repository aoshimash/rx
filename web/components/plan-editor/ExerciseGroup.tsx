'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import type { PlanEntryCreate } from '@/types/api';
import { Plus, Trash2 } from 'lucide-react';
import { ExerciseRow, type SetType } from './ExerciseRow';

interface ExerciseGroupProps {
  exerciseName: string;
  entries: { entry: PlanEntryCreate; flatIndex: number }[];
  onRename: (newName: string) => void;
  onRemoveExercise: () => void;
  onAddSet: () => void;
  onSetTypeChange: (flatIndex: number, value: SetType | undefined) => void;
  onSetsChange: (flatIndex: number, value: number) => void;
  onRepsChange: (flatIndex: number, value: number) => void;
  onLoadKgChange: (flatIndex: number, value: number) => void;
  onRpeChange: (flatIndex: number, value: number) => void;
  onRemoveSet: (flatIndex: number) => void;
}

export function ExerciseGroup({
  exerciseName,
  entries,
  onRename,
  onRemoveExercise,
  onAddSet,
  onSetTypeChange,
  onSetsChange,
  onRepsChange,
  onLoadKgChange,
  onRpeChange,
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
          <span className="w-[110px] text-xs text-muted-foreground">タイプ</span>
          <span className="w-14 text-xs text-muted-foreground text-center">Sets</span>
          <span className="w-4" />
          <span className="w-14 text-xs text-muted-foreground text-center">Reps</span>
          <span className="w-20 text-xs text-muted-foreground text-center">Load (kg)</span>
          <span className="w-16 text-xs text-muted-foreground text-center">RPE</span>
        </div>
      )}

      {/* Set rows */}
      <div className="space-y-1.5">
        {entries.map(({ entry, flatIndex }) => (
          <ExerciseRow
            key={flatIndex}
            setType={(entry.metadata?.set_type as SetType) ?? undefined}
            sets={entry.sets}
            reps={entry.reps}
            loadKg={entry.load_kg}
            rpe={entry.rpe}
            onSetTypeChange={(val) => onSetTypeChange(flatIndex, val)}
            onSetsChange={(val) => onSetsChange(flatIndex, val)}
            onRepsChange={(val) => onRepsChange(flatIndex, val)}
            onLoadKgChange={(val) => onLoadKgChange(flatIndex, val)}
            onRpeChange={(val) => onRpeChange(flatIndex, val)}
            onRemove={() => onRemoveSet(flatIndex)}
          />
        ))}
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
