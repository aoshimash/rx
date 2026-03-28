'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { TableCell, TableRow } from '@/components/ui/table';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { GripVertical, X } from 'lucide-react';
import { EditableCell } from './EditableCell';
import type { TableEntry } from './types';

interface LogEntryRowProps {
  entry: TableEntry;
  onUpdate: (id: string, field: keyof TableEntry, value: unknown) => void;
  onRemove: (id: string) => void;
}

export function LogEntryRow({ entry, onUpdate, onRemove }: LogEntryRowProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: entry.id,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const hasPlan = !!entry.plan;

  return (
    <TableRow ref={setNodeRef} style={style}>
      <TableCell className="w-8 px-1">
        <button
          type="button"
          className="cursor-grab active:cursor-grabbing text-muted-foreground hover:text-foreground"
          {...attributes}
          {...listeners}
        >
          <GripVertical className="h-4 w-4" />
        </button>
      </TableCell>

      <TableCell className="font-medium min-w-[160px]">
        {hasPlan ? (
          <span>{entry.exercise_name}</span>
        ) : (
          <Input
            value={entry.exercise_name}
            onChange={(e) => onUpdate(entry.id, 'exercise_name', e.target.value)}
            placeholder="Exercise name"
            className="h-8"
          />
        )}
      </TableCell>

      <TableCell className="w-20">
        <EditableCell
          value={entry.sets}
          onChange={(v) => {
            onUpdate(entry.id, 'sets', v);
            onUpdate(entry.id, 'setsEdited', true);
          }}
          defaultReadOnly={hasPlan && !entry.setsEdited}
          displayText={entry.plan?.sets?.toString()}
          isEdited={entry.setsEdited}
          min={1}
          step={1}
        />
      </TableCell>

      <TableCell className="w-20">
        <EditableCell
          value={entry.reps}
          onChange={(v) => {
            onUpdate(entry.id, 'reps', v);
            onUpdate(entry.id, 'repsEdited', true);
          }}
          defaultReadOnly={hasPlan && !entry.repsEdited}
          displayText={entry.plan?.reps?.toString()}
          isEdited={entry.repsEdited}
          min={1}
          step={1}
        />
      </TableCell>

      <TableCell className="w-24">
        <EditableCell
          value={entry.load_kg}
          onChange={(v) => onUpdate(entry.id, 'load_kg', v)}
          placeholder={entry.plan?.load_kg?.toString()}
          min={0}
          step={0.5}
        />
      </TableCell>

      <TableCell className="w-20">
        <EditableCell
          value={entry.rpe}
          onChange={(v) => onUpdate(entry.id, 'rpe', v)}
          placeholder={entry.plan?.rpe?.toString()}
          min={1}
          max={10}
          step={0.5}
        />
      </TableCell>

      <TableCell className="min-w-[120px]">
        <Input
          value={entry.notes}
          onChange={(e) => onUpdate(entry.id, 'notes', e.target.value)}
          placeholder="Notes"
          className="h-8"
        />
      </TableCell>

      <TableCell className="w-10 px-1">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onRemove(entry.id)}
          className="h-8 w-8 p-0"
        >
          <X className="h-4 w-4" />
        </Button>
      </TableCell>
    </TableRow>
  );
}
