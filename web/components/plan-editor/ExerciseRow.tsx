'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Link2, X } from 'lucide-react';

interface ExerciseRowProps {
  label?: string;
  rpe?: number;
  loadKg?: number;
  reps?: number;
  sets?: number;
  onLabelChange: (value: string | undefined) => void;
  onRpeChange: (value: number | undefined) => void;
  onLoadKgChange: (value: number | undefined) => void;
  onRepsChange: (value: number) => void;
  onSetsChange: (value: number) => void;
  onRemove: () => void;
  linked?: boolean;
}

export function ExerciseRow({
  label,
  rpe,
  loadKg,
  reps,
  sets,
  onLabelChange,
  onRpeChange,
  onLoadKgChange,
  onRepsChange,
  onSetsChange,
  onRemove,
  linked,
}: ExerciseRowProps) {
  return (
    <div className="flex items-center gap-2">
      <Input
        value={label ?? ''}
        onChange={(e) => onLabelChange(e.target.value || undefined)}
        placeholder="-"
        className="w-[110px] h-8 text-xs"
      />

      <div
        className={`flex items-center gap-1 rounded-md px-1 ${linked ? 'bg-muted/50 border border-dashed' : ''}`}
      >
        <Input
          type="number"
          value={rpe ?? ''}
          onChange={(e) => onRpeChange(e.target.value ? Number(e.target.value) : undefined)}
          min={1}
          max={10}
          step={0.5}
          placeholder="RPE"
          className="w-16 h-8 text-center text-sm"
        />
        {linked && <Link2 className="h-3 w-3 text-muted-foreground shrink-0" />}
        <Input
          type="number"
          step="0.5"
          value={loadKg ?? ''}
          onChange={(e) => onLoadKgChange(e.target.value ? Number(e.target.value) : undefined)}
          min={0}
          placeholder="kg"
          className="w-20 h-8 text-center text-sm"
        />
      </div>

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
