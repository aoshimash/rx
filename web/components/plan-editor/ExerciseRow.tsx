'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { X } from 'lucide-react';

interface ExerciseRowEditableProps {
  readOnly?: false;
  label?: string;
  loadKg?: number;
  reps?: number;
  sets?: number;
  onLabelChange: (value: string | undefined) => void;
  onLoadKgChange: (value: number | undefined) => void;
  onRepsChange: (value: number) => void;
  onSetsChange: (value: number) => void;
  onRemove: () => void;
}

interface ExerciseRowReadOnlyProps {
  readOnly: true;
  label?: string;
  rpe?: number;
  loadKg?: number;
  reps?: number;
  sets?: number;
  onRemove: () => void;
}

type ExerciseRowProps = ExerciseRowEditableProps | ExerciseRowReadOnlyProps;

export function ExerciseRow(props: ExerciseRowProps) {
  if (props.readOnly) {
    return (
      <div className="flex items-center gap-2 bg-muted/30 rounded-md px-1 py-0.5">
        {props.label && (
          <Badge variant="outline" className="text-xs shrink-0">
            {props.label}
          </Badge>
        )}
        <span className="text-sm text-muted-foreground min-w-0">
          {[
            props.rpe != null && `RPE${props.rpe}`,
            props.loadKg != null && `${props.loadKg}kg`,
            props.reps != null && `${props.reps}reps`,
            props.sets != null && `${props.sets}sets`,
          ]
            .filter(Boolean)
            .join(' ')}
        </span>
        <div className="ml-auto shrink-0">
          <Button variant="ghost" size="sm" className="h-8 w-8 p-0" onClick={props.onRemove}>
            <X className="h-3 w-3" />
          </Button>
        </div>
      </div>
    );
  }

  const {
    label,
    loadKg,
    reps,
    sets,
    onLabelChange,
    onLoadKgChange,
    onRepsChange,
    onSetsChange,
    onRemove,
  } = props;
  return (
    <div className="flex items-center gap-2">
      <Input
        value={label ?? ''}
        onChange={(e) => onLabelChange(e.target.value || undefined)}
        placeholder="-"
        className="w-[110px] h-8 text-xs"
      />
      <Input
        type="number"
        step="0.5"
        value={loadKg ?? ''}
        onChange={(e) => onLoadKgChange(e.target.value ? Number(e.target.value) : undefined)}
        min={0}
        placeholder="kg"
        className="w-20 h-8 text-center text-sm"
      />
      <Input
        type="number"
        value={reps ?? ''}
        onChange={(e) => onRepsChange(Number(e.target.value))}
        min={1}
        placeholder="-"
        className="w-14 h-8 text-center text-sm"
      />
      <Input
        type="number"
        value={sets ?? ''}
        onChange={(e) => onSetsChange(Number(e.target.value))}
        min={1}
        placeholder="-"
        className="w-14 h-8 text-center text-sm"
      />
      <Button variant="ghost" size="sm" className="h-8 w-8 p-0 shrink-0" onClick={onRemove}>
        <X className="h-3 w-3" />
      </Button>
    </div>
  );
}
