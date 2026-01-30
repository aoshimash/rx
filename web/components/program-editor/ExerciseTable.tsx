import { Button } from '@/components/ui/button';
import { Plus } from 'lucide-react';
import { ExerciseRow } from './ExerciseRow';
import type { Exercise, ProgramNodeCreate } from '@/types/api';

interface ExerciseTableProps {
  exercises: ProgramNodeCreate[];
  availableExercises: Exercise[];
  onExerciseChange: (index: number, exerciseId: string) => void;
  onSetsChange: (index: number, value: number) => void;
  onRepsChange: (index: number, value: number) => void;
  onRpeChange: (index: number, value: number) => void;
  onRemove: (index: number) => void;
  onAdd: () => void;
}

/**
 * Exercise table for day's exercise list in program editor
 */
export function ExerciseTable({
  exercises,
  availableExercises,
  onExerciseChange,
  onSetsChange,
  onRepsChange,
  onRpeChange,
  onRemove,
  onAdd,
}: ExerciseTableProps) {
  return (
    <div className="space-y-3">
      {exercises.map((exercise, idx) => (
        <ExerciseRow
          key={idx}
          exerciseId={exercise.exercise_id}
          exerciseName={exercise.name}
          targetSets={exercise.target_sets}
          targetReps={exercise.target_reps}
          targetRpe={exercise.target_rpe}
          exercises={availableExercises}
          onExerciseChange={(id) => onExerciseChange(idx, id)}
          onSetsChange={(value) => onSetsChange(idx, value)}
          onRepsChange={(value) => onRepsChange(idx, value)}
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
