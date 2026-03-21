'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { X } from 'lucide-react';

interface ExerciseRowProps {
  label?: string;
  sets?: number;
  reps?: number;
  loadKg?: number;
  rpe?: number;
  onLabelChange: (value: string | undefined) => void;
  onSetsChange: (value: number) => void;
  onRepsChange: (value: number) => void;
  onLoadKgChange: (value: number) => void;
  onRpeChange: (value: number) => void;
  onRemove: () => void;
}

export function ExerciseRow({
  label,
  sets,
  reps,
  loadKg,
  rpe,
  onLabelChange,
  onSetsChange,
  onRepsChange,
  onLoadKgChange,
  onRpeChange,
  onRemove,
}: ExerciseRowProps) {
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
        value={sets ?? ''}
        onChange={(e) => onSetsChange(Number(e.target.value))}
        min={1}
        placeholder="-"
        className="w-14 h-8 text-center text-sm"
      />
      <span className="text-muted-foreground text-sm">×</span>
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
        step="0.5"
        value={loadKg ?? ''}
        onChange={(e) => onLoadKgChange(Number(e.target.value))}
        min={0}
        placeholder="kg"
        className="w-20 h-8 text-center text-sm"
      />
      <Input
        type="number"
        value={rpe ?? ''}
        onChange={(e) => onRpeChange(Number(e.target.value))}
        min={1}
        max={10}
        placeholder="RPE"
        className="w-16 h-8 text-center text-sm"
      />

      <Button variant="ghost" size="sm" className="h-8 w-8 p-0 shrink-0" onClick={onRemove}>
        <X className="h-3 w-3" />
      </Button>
    </div>
  );
}
