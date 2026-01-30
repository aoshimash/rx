import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { X } from 'lucide-react';
import { ExerciseCombobox } from './ExerciseCombobox';
import type { Exercise } from '@/types/api';

interface ExerciseRowProps {
  exerciseId?: string;
  targetSets?: number;
  targetReps?: number;
  targetRpe?: number;
  exercises: Exercise[];
  onExerciseChange: (exerciseId: string) => void;
  onSetsChange: (value: number) => void;
  onRepsChange: (value: number) => void;
  onRpeChange: (value: number) => void;
  onRemove: () => void;
}

/**
 * Single exercise prescription row in program editor
 */
export function ExerciseRow({
  exerciseId,
  targetSets,
  targetReps,
  targetRpe,
  exercises,
  onExerciseChange,
  onSetsChange,
  onRepsChange,
  onRpeChange,
  onRemove,
}: ExerciseRowProps) {
  return (
    <div className="grid gap-4 border rounded-lg p-4">
      <div className="flex items-center justify-between">
        <Label>Exercise</Label>
        <Button variant="ghost" size="sm" onClick={onRemove}>
          <X className="h-4 w-4" />
        </Button>
      </div>

      <ExerciseCombobox
        exercises={exercises}
        value={exerciseId}
        onSelect={onExerciseChange}
      />

      <div className="grid grid-cols-3 gap-3">
        <div className="space-y-2">
          <Label>Sets</Label>
          <Input
            type="number"
            value={targetSets || ''}
            onChange={(e) => onSetsChange(Number(e.target.value))}
            min={1}
            placeholder="3"
          />
        </div>
        <div className="space-y-2">
          <Label>Reps</Label>
          <Input
            type="number"
            value={targetReps || ''}
            onChange={(e) => onRepsChange(Number(e.target.value))}
            min={1}
            placeholder="10"
          />
        </div>
        <div className="space-y-2">
          <Label>Target RPE</Label>
          <Input
            type="number"
            value={targetRpe || ''}
            onChange={(e) => onRpeChange(Number(e.target.value))}
            min={1}
            max={10}
            placeholder="7"
          />
        </div>
      </div>
    </div>
  );
}
