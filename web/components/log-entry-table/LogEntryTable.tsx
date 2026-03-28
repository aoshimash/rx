'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Table, TableBody, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import type { LogEntryCreate, ProgramSessionEntry } from '@/types/api';
import {
  DndContext,
  type DragEndEvent,
  KeyboardSensor,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import { restrictToVerticalAxis } from '@dnd-kit/modifiers';
import {
  SortableContext,
  arrayMove,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { Plus } from 'lucide-react';
import { useState } from 'react';
import { LogEntryRow } from './LogEntryRow';
import { type TableEntry, createEmptyTableEntry, createTableEntryFromPlan } from './types';

interface LogEntryTableProps {
  initialEntries?: ProgramSessionEntry[];
  /** Pre-existing log entries for edit mode */
  existingEntries?: TableEntry[];
  onSave: (data: {
    entries: LogEntryCreate[];
    notes: string;
    startedAt?: string;
    finishedAt?: string;
  }) => Promise<void>;
  onCancel: () => void;
  saveLabel?: string;
  initialNotes?: string;
  initialStartedAt?: string;
  initialFinishedAt?: string;
}

export function LogEntryTable({
  initialEntries,
  existingEntries,
  onSave,
  onCancel,
  saveLabel = 'Save Log',
  initialNotes = '',
  initialStartedAt = '',
  initialFinishedAt = '',
}: LogEntryTableProps) {
  const [entries, setEntries] = useState<TableEntry[]>(() => {
    if (existingEntries) return existingEntries;
    if (initialEntries) return initialEntries.map((e, i) => createTableEntryFromPlan(e, i));
    return [createEmptyTableEntry()];
  });
  const [sessionNotes, setSessionNotes] = useState(initialNotes);
  const [startedAt, setStartedAt] = useState(initialStartedAt);
  const [finishedAt, setFinishedAt] = useState(initialFinishedAt);
  const [isSaving, setIsSaving] = useState(false);

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
  );

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;

    setEntries((prev) => {
      const oldIndex = prev.findIndex((e) => e.id === active.id);
      const newIndex = prev.findIndex((e) => e.id === over.id);
      return arrayMove(prev, oldIndex, newIndex);
    });
  };

  const handleUpdate = (id: string, field: keyof TableEntry, value: unknown) => {
    setEntries((prev) => prev.map((e) => (e.id === id ? { ...e, [field]: value } : e)));
  };

  const handleRemove = (id: string) => {
    setEntries((prev) => prev.filter((e) => e.id !== id));
  };

  const handleAddExercise = () => {
    setEntries((prev) => [...prev, createEmptyTableEntry()]);
  };

  const handleSave = async () => {
    const validEntries = entries.filter((e) => e.exercise_name.trim() !== '');
    if (validEntries.length === 0) return;

    setIsSaving(true);
    try {
      const logEntries: LogEntryCreate[] = validEntries.map((entry) => ({
        exercise_name: entry.exercise_name,
        sets: entry.sets,
        reps: entry.reps,
        load_kg: entry.load_kg,
        rpe: entry.rpe,
        notes: entry.notes || undefined,
      }));

      await onSave({
        entries: logEntries,
        notes: sessionNotes,
        startedAt: startedAt ? new Date(startedAt).toISOString() : undefined,
        finishedAt: finishedAt ? new Date(finishedAt).toISOString() : undefined,
      });
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-2 gap-4 max-w-md">
        <div className="space-y-2">
          <Label htmlFor="session-start">Session Start</Label>
          <Input
            id="session-start"
            type="datetime-local"
            value={startedAt}
            onChange={(e) => setStartedAt(e.target.value)}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="session-end">Session End</Label>
          <Input
            id="session-end"
            type="datetime-local"
            value={finishedAt}
            onChange={(e) => setFinishedAt(e.target.value)}
          />
        </div>
      </div>

      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        modifiers={[restrictToVerticalAxis]}
        onDragEnd={handleDragEnd}
      >
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-8" />
              <TableHead>Exercise</TableHead>
              <TableHead className="w-20 text-right">Sets</TableHead>
              <TableHead className="w-20 text-right">Reps</TableHead>
              <TableHead className="w-24 text-right">Load (kg)</TableHead>
              <TableHead className="w-20 text-right">RPE</TableHead>
              <TableHead>Notes</TableHead>
              <TableHead className="w-10" />
            </TableRow>
          </TableHeader>
          <SortableContext items={entries.map((e) => e.id)} strategy={verticalListSortingStrategy}>
            <TableBody>
              {entries.map((entry) => (
                <LogEntryRow
                  key={entry.id}
                  entry={entry}
                  onUpdate={handleUpdate}
                  onRemove={handleRemove}
                />
              ))}
            </TableBody>
          </SortableContext>
        </Table>
      </DndContext>

      <Button variant="outline" onClick={handleAddExercise} className="w-full">
        <Plus className="h-4 w-4 mr-2" />
        Add Exercise
      </Button>

      <div className="space-y-2">
        <Label htmlFor="session-notes">Session Notes</Label>
        <Input
          id="session-notes"
          placeholder="How did the session go?"
          value={sessionNotes}
          onChange={(e) => setSessionNotes(e.target.value)}
        />
      </div>

      <div className="flex justify-end gap-2 pt-4 border-t">
        <Button variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button onClick={handleSave} disabled={isSaving || entries.length === 0}>
          {isSaving ? 'Saving...' : saveLabel}
        </Button>
      </div>
    </div>
  );
}
