'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Table, TableBody, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import type { FieldDef, LogEntryCreate, ProgramSessionEntry } from '@/types/api';
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

function buildDynamicLogEntry(entry: TableEntry, fieldDefs: FieldDef[]): LogEntryCreate {
  const fields: Record<string, unknown> = {};
  for (const fd of fieldDefs.filter((fd) => fd.type !== 'video')) {
    if (entry.logFieldValues[fd.name] !== undefined) {
      fields[fd.name] = entry.logFieldValues[fd.name];
    }
  }
  const hasVideoField = fieldDefs.some((fd) => fd.type === 'video');
  const sets =
    hasVideoField && entry.videoObjectKey
      ? [{ set_number: 1, fields: {}, video_object_key: entry.videoObjectKey }]
      : undefined;
  return {
    exercise_name: entry.exercise_name,
    fields: Object.keys(fields).length > 0 ? fields : undefined,
    sets,
    notes: entry.notes || undefined,
  };
}

function buildFallbackLogEntry(entry: TableEntry): LogEntryCreate {
  const fields: Record<string, unknown> = { ...entry.fields };
  if (entry.sets !== undefined) fields.sets = entry.sets;
  if (entry.reps !== undefined) fields.reps = entry.reps;
  if (entry.load_kg !== undefined) fields.load_kg = entry.load_kg;
  return {
    exercise_name: entry.exercise_name,
    fields: Object.keys(fields).length > 0 ? fields : undefined,
    notes: entry.notes || undefined,
  };
}

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
  /** FieldDef definitions from the session's FieldGroup.log_fields */
  fieldDefs?: FieldDef[];
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
  fieldDefs,
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

  const buildLogEntry = (entry: TableEntry): LogEntryCreate =>
    fieldDefs && fieldDefs.length > 0
      ? buildDynamicLogEntry(entry, fieldDefs)
      : buildFallbackLogEntry(entry);

  const handleSave = async () => {
    const validEntries = entries.filter((e) => e.exercise_name.trim() !== '');
    if (validEntries.length === 0) return;

    setIsSaving(true);
    try {
      const logEntries: LogEntryCreate[] = validEntries.map(buildLogEntry);

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

  const dynamicHeaders = fieldDefs?.map((fd) => (
    <TableHead key={fd.name} className="min-w-[80px]">
      {fd.name}
      {fd.description && (
        <span className="text-xs font-normal text-muted-foreground ml-1">({fd.description})</span>
      )}
    </TableHead>
  ));

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
              {fieldDefs && fieldDefs.length > 0 ? (
                dynamicHeaders
              ) : (
                <>
                  <TableHead className="w-20 text-right">Sets</TableHead>
                  <TableHead className="w-20 text-right">Reps</TableHead>
                  <TableHead className="w-24 text-right">Load (kg)</TableHead>
                </>
              )}
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
                  fieldDefs={fieldDefs}
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
