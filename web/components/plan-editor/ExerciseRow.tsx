import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { X } from 'lucide-react';

interface ExerciseRowProps {
  exerciseName: string;
  sets?: number;
  reps?: number;
  loadKg?: number;
  rpe?: number;
  onExerciseNameChange: (name: string) => void;
  onSetsChange: (value: number) => void;
  onRepsChange: (value: number) => void;
  onLoadKgChange: (value: number) => void;
  onRpeChange: (value: number) => void;
  onRemove: () => void;
}

export function ExerciseRow({
  exerciseName,
  sets,
  reps,
  loadKg,
  rpe,
  onExerciseNameChange,
  onSetsChange,
  onRepsChange,
  onLoadKgChange,
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

      <Input
        value={exerciseName}
        onChange={(e) => onExerciseNameChange(e.target.value)}
        placeholder="e.g., Squat, Bench Press"
      />

      <div className="grid grid-cols-4 gap-3">
        <div className="space-y-2">
          <Label>Sets</Label>
          <Input
            type="number"
            value={sets ?? ''}
            onChange={(e) => onSetsChange(Number(e.target.value))}
            min={1}
            placeholder="3"
          />
        </div>
        <div className="space-y-2">
          <Label>Reps</Label>
          <Input
            type="number"
            value={reps ?? ''}
            onChange={(e) => onRepsChange(Number(e.target.value))}
            min={1}
            placeholder="10"
          />
        </div>
        <div className="space-y-2">
          <Label>Load (kg)</Label>
          <Input
            type="number"
            step="0.5"
            value={loadKg ?? ''}
            onChange={(e) => onLoadKgChange(Number(e.target.value))}
            min={0}
            placeholder="0"
          />
        </div>
        <div className="space-y-2">
          <Label>RPE</Label>
          <Input
            type="number"
            value={rpe ?? ''}
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
