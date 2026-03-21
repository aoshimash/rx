import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { ChevronDown, ChevronRight, X } from 'lucide-react';
import { useState } from 'react';

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
  startedAt?: string;
  finishedAt?: string;
  onStartedAtChange?: (value: string) => void;
  onFinishedAtChange?: (value: string) => void;
  onRemove?: () => void;
  planValues?: {
    sets?: number;
    reps?: number;
    load_kg?: number;
    rpe?: number;
  };
}

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
  startedAt,
  finishedAt,
  onStartedAtChange,
  onFinishedAtChange,
  onRemove,
  planValues,
}: ExerciseInputRowProps) {
  const [showTiming, setShowTiming] = useState(false);

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
          Plan: {planValues.sets}x{planValues.reps}
          {planValues.load_kg !== undefined && ` @ ${planValues.load_kg}kg`}
          {planValues.rpe !== undefined && ` (RPE ${planValues.rpe})`}
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
            placeholder={planValues?.load_kg?.toString()}
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

      {onStartedAtChange && onFinishedAtChange && (
        <div>
          <button
            type="button"
            className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground"
            onClick={() => setShowTiming(!showTiming)}
          >
            {showTiming ? (
              <ChevronDown className="h-3 w-3" />
            ) : (
              <ChevronRight className="h-3 w-3" />
            )}
            Timing
          </button>
          {showTiming && (
            <div className="grid grid-cols-2 gap-3 mt-2">
              <div className="space-y-1">
                <Label className="text-xs" htmlFor={`start-${exerciseName}`}>
                  Start
                </Label>
                <Input
                  id={`start-${exerciseName}`}
                  type="datetime-local"
                  className="text-xs h-8"
                  value={startedAt || ''}
                  onChange={(e) => onStartedAtChange(e.target.value)}
                />
              </div>
              <div className="space-y-1">
                <Label className="text-xs" htmlFor={`end-${exerciseName}`}>
                  End
                </Label>
                <Input
                  id={`end-${exerciseName}`}
                  type="datetime-local"
                  className="text-xs h-8"
                  value={finishedAt || ''}
                  onChange={(e) => onFinishedAtChange(e.target.value)}
                />
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
