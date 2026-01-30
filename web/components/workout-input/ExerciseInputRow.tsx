import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { X } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface ExerciseInputRowProps {
  exerciseName: string;
  sets: number;
  reps: number;
  load: number;
  rpe: number;
  onSetsChange: (value: number) => void;
  onRepsChange: (value: number) => void;
  onLoadChange: (value: number) => void;
  onRpeChange: (value: number) => void;
  onRemove?: () => void;
  planValues?: {
    sets?: number;
    reps?: number;
    load?: number;
    rpe?: number;
  };
}

/**
 * Single exercise input row for workout recording
 * 
 * Shows planned values as placeholders and allows editing actual values
 */
export function ExerciseInputRow({
  exerciseName,
  sets,
  reps,
  load,
  rpe,
  onSetsChange,
  onRepsChange,
  onLoadChange,
  onRpeChange,
  onRemove,
  planValues,
}: ExerciseInputRowProps) {
  return (
    <div className="grid gap-4 border rounded-lg p-4">
      <div className="flex items-center justify-between">
        <h4 className="font-semibold">{exerciseName}</h4>
        {onRemove && (
          <Button variant="ghost" size="sm" onClick={onRemove}>
            <X className="h-4 w-4" />
          </Button>
        )}
      </div>

      {planValues && (
        <div className="text-sm text-muted-foreground">
          Plan: {planValues.sets}×{planValues.reps} @ {planValues.load}kg
          {planValues.rpe && ` (RPE ${planValues.rpe})`}
        </div>
      )}

      <div className="grid grid-cols-4 gap-3">
        <div className="space-y-2">
          <Label htmlFor={`sets-${exerciseName}`}>Sets</Label>
          <Input
            id={`sets-${exerciseName}`}
            type="number"
            value={sets}
            onChange={(e) => onSetsChange(Number(e.target.value))}
            min={1}
            placeholder={planValues?.sets?.toString()}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor={`reps-${exerciseName}`}>Reps</Label>
          <Input
            id={`reps-${exerciseName}`}
            type="number"
            value={reps}
            onChange={(e) => onRepsChange(Number(e.target.value))}
            min={1}
            placeholder={planValues?.reps?.toString()}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor={`load-${exerciseName}`}>Load (kg)</Label>
          <Input
            id={`load-${exerciseName}`}
            type="number"
            step="0.5"
            value={load}
            onChange={(e) => onLoadChange(Number(e.target.value))}
            min={0}
            placeholder={planValues?.load?.toString()}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor={`rpe-${exerciseName}`}>RPE</Label>
          <Input
            id={`rpe-${exerciseName}`}
            type="number"
            value={rpe}
            onChange={(e) => onRpeChange(Number(e.target.value))}
            min={1}
            max={10}
            placeholder={planValues?.rpe?.toString()}
          />
        </div>
      </div>
    </div>
  );
}
