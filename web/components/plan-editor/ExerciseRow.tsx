'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { X } from 'lucide-react';

interface ExerciseRowProps {
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

export function ExerciseRow({
  label,
  loadKg,
  reps,
  sets,
  onLabelChange,
  onLoadKgChange,
  onRepsChange,
  onSetsChange,
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
