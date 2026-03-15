import { Button } from '@/components/ui/button';
import type { PlanEntryCreate } from '@/types/api';
import { Plus } from 'lucide-react';
import { ExerciseRow } from './ExerciseRow';

interface ExerciseTableProps {
  exercises: PlanEntryCreate[];
  onExerciseNameChange: (index: number, name: string) => void;
  onSetsChange: (index: number, value: number) => void;
  onRepsChange: (index: number, value: number) => void;
  onLoadKgChange: (index: number, value: number) => void;
  onRpeChange: (index: number, value: number) => void;
  onRemove: (index: number) => void;
  onAdd: () => void;
}

export function ExerciseTable({
  exercises,
  onExerciseNameChange,
  onSetsChange,
  onRepsChange,
  onLoadKgChange,
  onRpeChange,
  onRemove,
  onAdd,
}: ExerciseTableProps) {
  return (
    <div className="space-y-3">
      {exercises.map((exercise, idx) => (
        <ExerciseRow
          key={idx}
          exerciseName={exercise.exercise_name}
          sets={exercise.sets}
          reps={exercise.reps}
          loadKg={exercise.load_kg}
          rpe={exercise.rpe}
          onExerciseNameChange={(name) => onExerciseNameChange(idx, name)}
          onSetsChange={(value) => onSetsChange(idx, value)}
          onRepsChange={(value) => onRepsChange(idx, value)}
          onLoadKgChange={(value) => onLoadKgChange(idx, value)}
          onRpeChange={(value) => onRpeChange(idx, value)}
          onRemove={() => onRemove(idx)}
        />
      ))}
      <Button variant="outline" onClick={onAdd} className="w-full">
        <Plus className="h-4 w-4 mr-2" />
        Add Exercise
      </Button>
    </div>
  );
}
