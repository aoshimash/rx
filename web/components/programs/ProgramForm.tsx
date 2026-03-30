'use client';

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import { DeleteConfirmDialog } from '@/components/ui/delete-confirm-dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { usePrograms } from '@/lib/hooks/usePrograms';
import type { Program } from '@/types/api';
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
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Copy, GripVertical, Lock, Plus, Trash2, X } from 'lucide-react';
import { useState } from 'react';

// ============================================================================
// Types
// ============================================================================

export type CustomFieldType = 'text' | 'number' | 'select';

export type BuiltinFieldKey = 'load_kg' | 'reps' | 'sets';

export interface CustomFieldDef {
  name: string;
  type: CustomFieldType;
  options?: string[];
  description?: string;
  builtin?: BuiltinFieldKey;
}

const DEFAULT_BUILTIN_FIELDS: CustomFieldDef[] = [
  { name: 'Load (kg)', type: 'number', builtin: 'load_kg' },
  { name: 'Reps', type: 'number', builtin: 'reps' },
  { name: 'Sets', type: 'number', builtin: 'sets' },
];

function ensureBuiltinFields(fields: CustomFieldDef[]): CustomFieldDef[] {
  const result = [...fields];
  for (const builtin of DEFAULT_BUILTIN_FIELDS) {
    if (!result.some((f) => f.builtin === builtin.builtin)) {
      result.push(builtin);
    }
  }
  return result;
}

/**
 * Intermediate entry format used by the program form for session-based editing.
 * Converted to/from ProgramSession structures when saving/loading.
 */
export interface ProgramFormEntry {
  exercise_name: string;
  order: number;
  fields?: Record<string, unknown>;
  notes?: string;
  metadata?: Record<string, unknown>;
}

interface SessionExercise {
  exercise_name: string;
  fields?: Record<string, unknown>;
}

interface SessionGroup {
  name: string;
  date?: string;
  exercises: SessionExercise[];
}

// ============================================================================
// Conversion helpers
// ============================================================================

function entryToExercise(entry: ProgramFormEntry): SessionExercise {
  const metadata: Record<string, string> = {};
  if (entry.metadata) {
    for (const [k, v] of Object.entries(entry.metadata)) {
      if (k !== 'session' && k !== 'date' && typeof v === 'string') {
        metadata[k] = v;
      }
    }
  }
  return {
    exercise_name: entry.exercise_name,
    fields: entry.fields,
  };
}

function entriesToSessionGroups(entries: ProgramFormEntry[]): SessionGroup[] {
  const sessionMap = new Map<string, SessionExercise[]>();
  const sessionDates = new Map<string, string>();
  const sessionOrder: string[] = [];

  for (const entry of entries) {
    const sessionName = (entry.metadata?.session as string) ?? 'Session 1';
    if (!sessionMap.has(sessionName)) {
      sessionMap.set(sessionName, []);
      sessionOrder.push(sessionName);
      const date = entry.metadata?.date as string | undefined;
      if (date) sessionDates.set(sessionName, date);
    }
    sessionMap.get(sessionName)?.push(entryToExercise(entry));
  }

  return sessionOrder.map((name) => ({
    name,
    date: sessionDates.get(name),
    exercises: sessionMap.get(name) || [],
  }));
}

function exerciseToEntry(
  ex: SessionExercise,
  sessionName: string,
  sessionDate: string | undefined,
  order: number
): ProgramFormEntry {
  const metadata: Record<string, unknown> = { session: sessionName };
  if (sessionDate) metadata.date = sessionDate;
  return {
    exercise_name: ex.exercise_name,
    order,
    fields: ex.fields,
    metadata,
  };
}

function sessionGroupsToEntries(sessions: SessionGroup[]): ProgramFormEntry[] {
  const entries: ProgramFormEntry[] = [];
  let order = 0;

  for (const session of sessions) {
    for (const ex of session.exercises) {
      entries.push(exerciseToEntry(ex, session.name, session.date, order));
      order++;
    }
  }

  return entries;
}

// ============================================================================
// CustomFieldsEditor
// ============================================================================

function SelectOptionsEditor({
  options,
  onChange,
}: {
  options: string[];
  onChange: (options: string[]) => void;
}) {
  const [newOption, setNewOption] = useState('');

  const addOption = () => {
    const trimmed = newOption.trim();
    if (trimmed && !options.includes(trimmed)) {
      onChange([...options, trimmed]);
      setNewOption('');
    }
  };

  return (
    <div className="flex flex-wrap items-center gap-1.5">
      {options.map((opt) => (
        <span
          key={opt}
          className="inline-flex items-center gap-0.5 rounded bg-muted px-1.5 py-0.5 text-xs"
        >
          {opt}
          <button
            type="button"
            onClick={() => onChange(options.filter((o) => o !== opt))}
            className="text-muted-foreground hover:text-foreground cursor-pointer"
          >
            <X className="h-2.5 w-2.5" />
          </button>
        </span>
      ))}
      <Input
        value={newOption}
        onChange={(e) => setNewOption(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === 'Enter') {
            e.preventDefault();
            addOption();
          }
        }}
        placeholder="Add option..."
        className="h-6 w-24 text-xs"
      />
    </div>
  );
}

function SortableFieldRow({
  field,
  onUpdate,
  onRemove,
}: {
  field: CustomFieldDef;
  onUpdate: (updated: CustomFieldDef) => void;
  onRemove: () => void;
}) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: field.name,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <div ref={setNodeRef} style={style} className="rounded-md border p-2 space-y-1">
      <div className="flex items-start gap-2">
        <button
          type="button"
          className="text-muted-foreground hover:text-foreground cursor-grab active:cursor-grabbing mt-0.5 shrink-0"
          {...attributes}
          {...listeners}
        >
          <GripVertical className="h-3.5 w-3.5" />
        </button>
        <div className="flex items-center gap-2 flex-1 min-w-0">
          <span className="text-sm font-medium whitespace-nowrap">{field.name}</span>
          <span className="text-xs text-muted-foreground">({field.type})</span>
          {field.type === 'select' && (
            <SelectOptionsEditor
              options={field.options ?? []}
              onChange={(options) => onUpdate({ ...field, options })}
            />
          )}
        </div>
        {!field.builtin && (
          <button
            type="button"
            onClick={onRemove}
            className="text-muted-foreground hover:text-foreground cursor-pointer mt-0.5 shrink-0"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        )}
      </div>
      <Input
        value={field.description ?? ''}
        onChange={(e) => onUpdate({ ...field, description: e.target.value || undefined })}
        placeholder="Description (e.g., unit, scale, format)"
        className="h-7 text-xs ml-6"
      />
    </div>
  );
}

function extractCustomFields(program: Program): CustomFieldDef[] {
  const raw = program.metadata?.custom_fields;
  if (!Array.isArray(raw)) return [];
  return raw as CustomFieldDef[];
}

function CustomFieldsEditor({
  fields,
  onChange,
}: {
  fields: CustomFieldDef[];
  onChange: (fields: CustomFieldDef[]) => void;
}) {
  const [newFieldName, setNewFieldName] = useState('');
  const [newFieldType, setNewFieldType] = useState<CustomFieldType>('text');
  const { data: programsData } = usePrograms();

  const programsWithFields = programsData?.data?.filter((p) => extractCustomFields(p).length > 0);

  const handleCopyFrom = (programId: string) => {
    const program = programsData?.data?.find((p) => p.id === programId);
    if (!program) return;
    const copiedFields = extractCustomFields(program);
    onChange(ensureBuiltinFields(copiedFields));
  };

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
  );

  const addField = () => {
    const trimmed = newFieldName.trim();
    if (trimmed && !fields.some((f) => f.name === trimmed)) {
      const field: CustomFieldDef = { name: trimmed, type: newFieldType };
      if (newFieldType === 'select') field.options = [];
      onChange([...fields, field]);
      setNewFieldName('');
      setNewFieldType('text');
    }
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;
    const oldIndex = fields.findIndex((f) => f.name === active.id);
    const newIndex = fields.findIndex((f) => f.name === over.id);
    onChange(arrayMove(fields, oldIndex, newIndex));
  };

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <Label>Custom Fields</Label>
        {programsWithFields && programsWithFields.length > 0 && (
          <Select onValueChange={handleCopyFrom}>
            <SelectTrigger className="h-8 w-auto gap-2 text-xs">
              <Copy className="h-3.5 w-3.5" />
              <SelectValue placeholder="Copy from..." />
            </SelectTrigger>
            <SelectContent>
              {programsWithFields.map((p) => (
                <SelectItem key={p.id} value={p.id}>
                  {p.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        )}
      </div>
      {fields.length > 0 && (
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          modifiers={[restrictToVerticalAxis]}
          onDragEnd={handleDragEnd}
        >
          <SortableContext items={fields.map((f) => f.name)} strategy={verticalListSortingStrategy}>
            <div className="space-y-2">
              {fields.map((field, idx) => (
                <SortableFieldRow
                  key={field.name}
                  field={field}
                  onUpdate={(updated) => {
                    const next = [...fields];
                    next[idx] = updated;
                    onChange(next);
                  }}
                  onRemove={() => onChange(fields.filter((_, i) => i !== idx))}
                />
              ))}
            </div>
          </SortableContext>
        </DndContext>
      )}
      <div className="flex items-center gap-2">
        <Input
          value={newFieldName}
          onChange={(e) => setNewFieldName(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              e.preventDefault();
              addField();
            }
          }}
          placeholder="Field name"
          className="h-8 w-40"
        />
        <Select value={newFieldType} onValueChange={(v) => setNewFieldType(v as CustomFieldType)}>
          <SelectTrigger className="h-8 w-28">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="text">Text</SelectItem>
            <SelectItem value="number">Number</SelectItem>
            <SelectItem value="select">Select</SelectItem>
          </SelectContent>
        </Select>
        <Button variant="ghost" size="sm" onClick={addField} disabled={!newFieldName.trim()}>
          <Plus className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}

// ============================================================================
// ProgramExerciseRow
// ============================================================================

const BUILTIN_CONFIG: Record<
  BuiltinFieldKey,
  { min?: number; step?: string; placeholder: string }
> = {
  load_kg: { min: 0, step: '0.5', placeholder: '100' },
  reps: { min: 1, placeholder: '10' },
  sets: { min: 1, placeholder: '3' },
};

interface ProgramExerciseRowProps {
  exercise: SessionExercise;
  onChange: (updated: SessionExercise) => void;
  onRemove: () => void;
  fields: CustomFieldDef[];
  disabled?: boolean;
}

function ProgramExerciseRow({
  exercise,
  onChange,
  onRemove,
  fields,
  disabled,
}: ProgramExerciseRowProps) {
  const updateField = (key: string, raw: string, isNumber: boolean) => {
    const fields = { ...exercise.fields };
    if (raw) {
      fields[key] = isNumber ? Number(raw) : raw;
    } else {
      delete fields[key];
    }
    onChange({ ...exercise, fields });
  };

  return (
    <div className="grid gap-4 border rounded-lg p-4">
      <div className="flex items-center justify-between">
        <Label>Exercise</Label>
        {!disabled && (
          <Button variant="ghost" size="sm" onClick={onRemove}>
            <X className="h-4 w-4" />
          </Button>
        )}
      </div>

      <Input
        value={exercise.exercise_name}
        onChange={(e) => onChange({ ...exercise, exercise_name: e.target.value })}
        placeholder="e.g., Squat, Bench Press"
        disabled={disabled}
      />

      <div
        className="grid gap-3"
        style={{ gridTemplateColumns: `repeat(${fields.length}, minmax(0, 1fr))` }}
      >
        {fields.map((field) => {
          if (field.builtin) {
            const config = BUILTIN_CONFIG[field.builtin];
            return (
              <div key={field.builtin} className="space-y-2">
                <Label>{field.name}</Label>
                <Input
                  type="number"
                  value={(exercise.fields?.[field.builtin] as number | undefined) ?? ''}
                  onChange={(e) => updateField(field.builtin as string, e.target.value, true)}
                  min={config.min}
                  step={config.step}
                  placeholder={config.placeholder}
                  disabled={disabled}
                />
              </div>
            );
          }

          if (field.type === 'select') {
            return (
              <div key={field.name} className="space-y-2">
                <Label>{field.name}</Label>
                <div className="flex items-center gap-1">
                  <Select
                    value={(exercise.fields?.[field.name] as string | undefined) ?? ''}
                    onValueChange={(v) => updateField(field.name, v, false)}
                    disabled={disabled}
                  >
                    <SelectTrigger className="flex-1">
                      <SelectValue placeholder="Select..." />
                    </SelectTrigger>
                    <SelectContent>
                      {(field.options ?? []).map((opt) => (
                        <SelectItem key={opt} value={opt}>
                          {opt}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {!!exercise.fields?.[field.name] && !disabled && (
                    <button
                      type="button"
                      onClick={() => updateField(field.name, '', false)}
                      className="text-muted-foreground hover:text-foreground cursor-pointer shrink-0"
                    >
                      <X className="h-3.5 w-3.5" />
                    </button>
                  )}
                </div>
              </div>
            );
          }

          return (
            <div key={field.name} className="space-y-2">
              <Label>{field.name}</Label>
              <Input
                type={field.type === 'number' ? 'number' : 'text'}
                value={(exercise.fields?.[field.name] as string | undefined) ?? ''}
                onChange={(e) => updateField(field.name, e.target.value, field.type === 'number')}
                disabled={disabled}
              />
            </div>
          );
        })}
      </div>
    </div>
  );
}

// ============================================================================
// ProgramForm
// ============================================================================

interface ProgramFormProps {
  programName: string;
  programNotes: string;
  initialEntries?: ProgramFormEntry[];
  initialCustomFields?: CustomFieldDef[];
  onNameChange: (name: string) => void;
  onNotesChange: (notes: string) => void;
  onSave: (entries: ProgramFormEntry[], customFields: CustomFieldDef[]) => void;
  onDelete?: () => void;
  isSaving?: boolean;
  isEditing?: boolean;
  lockedSessionNames?: Set<string>;
}

export function ProgramForm({
  programName,
  programNotes,
  initialEntries,
  initialCustomFields,
  onNameChange,
  onNotesChange,
  onSave,
  onDelete,
  isSaving,
  isEditing,
  lockedSessionNames,
}: ProgramFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [sessions, setSessions] = useState<SessionGroup[]>(() =>
    initialEntries && initialEntries.length > 0
      ? entriesToSessionGroups(initialEntries)
      : [{ name: '', exercises: [] }]
  );

  const [customFields, setCustomFields] = useState<CustomFieldDef[]>(() =>
    ensureBuiltinFields(initialCustomFields ?? [])
  );

  const handleAddSession = () => {
    setSessions([...sessions, { name: '', exercises: [] }]);
  };

  const handleRemoveSession = (idx: number) => {
    setSessions(sessions.filter((_, i) => i !== idx));
  };

  const handleSessionNameChange = (idx: number, name: string) => {
    const updated = [...sessions];
    const session = updated[idx];
    if (!session) return;
    updated[idx] = { ...session, name };
    setSessions(updated);
  };

  const handleAddExercise = (sessionIdx: number) => {
    const updated = [...sessions];
    const session = updated[sessionIdx];
    if (!session) return;
    updated[sessionIdx] = {
      ...session,
      exercises: [...session.exercises, { exercise_name: '', fields: { sets: 3, reps: 10 } }],
    };
    setSessions(updated);
  };

  const handleExerciseChange = (sessionIdx: number, exIdx: number, exercise: SessionExercise) => {
    const updated = [...sessions];
    const session = updated[sessionIdx];
    if (!session) return;
    const exercises = [...session.exercises];
    exercises[exIdx] = exercise;
    updated[sessionIdx] = { ...session, exercises };
    setSessions(updated);
  };

  const handleRemoveExercise = (sessionIdx: number, exIdx: number) => {
    const updated = [...sessions];
    const session = updated[sessionIdx];
    if (!session) return;
    updated[sessionIdx] = {
      ...session,
      exercises: session.exercises.filter((_, i) => i !== exIdx),
    };
    setSessions(updated);
  };

  return (
    <div className="space-y-6">
      <div className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="program-name">Program Name</Label>
          <Input
            id="program-name"
            value={programName}
            onChange={(e) => onNameChange(e.target.value)}
            placeholder="e.g., 5/3/1 BBB, GZCL"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="program-notes">Notes</Label>
          <Input
            id="program-notes"
            value={programNotes}
            onChange={(e) => onNotesChange(e.target.value)}
            placeholder="Additional notes"
          />
        </div>
        <CustomFieldsEditor fields={customFields} onChange={setCustomFields} />
      </div>

      <div className="space-y-4">
        <p className="text-sm font-semibold">Sessions</p>
        <Accordion type="multiple" className="w-full space-y-3">
          {sessions.map((session, sessionIdx) => {
            const isLocked = lockedSessionNames?.has(session.name) ?? false;
            return (
              <div
                key={sessionIdx}
                className={`border rounded-lg px-4 ${isLocked ? 'opacity-60' : ''}`}
              >
                <AccordionItem value={`session-${sessionIdx}`} className="border-0">
                  <div className="flex items-center gap-2 py-4">
                    <AccordionTrigger className="hover:no-underline p-0 shrink-0 [&>svg]:ml-0">
                      <span className="sr-only">Toggle {session.name || 'session'}</span>
                    </AccordionTrigger>
                    <div className="flex items-center justify-between flex-1 min-w-0">
                      {isLocked ? (
                        <div className="flex items-center gap-2">
                          <Lock className="h-4 w-4 text-muted-foreground" />
                          <span className="font-semibold">{session.name}</span>
                          <span className="text-xs text-muted-foreground">(completed)</span>
                        </div>
                      ) : (
                        <Input
                          value={session.name}
                          onChange={(e) => handleSessionNameChange(sessionIdx, e.target.value)}
                          onClick={(e) => e.stopPropagation()}
                          placeholder="e.g., Block1 Week2 Day3, Week1 Day2"
                          className="font-semibold border-none shadow-none p-0 h-auto focus-visible:ring-0 bg-transparent"
                        />
                      )}
                      {!isLocked && (
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleRemoveSession(sessionIdx);
                          }}
                        >
                          <X className="h-4 w-4" />
                        </Button>
                      )}
                    </div>
                  </div>
                  <AccordionContent>
                    <div className="space-y-4 pt-4">
                      <div className="space-y-3">
                        {session.exercises.map((exercise, exIdx) => (
                          <ProgramExerciseRow
                            key={exIdx}
                            exercise={exercise}
                            onChange={(updated) => handleExerciseChange(sessionIdx, exIdx, updated)}
                            onRemove={() => handleRemoveExercise(sessionIdx, exIdx)}
                            fields={customFields}
                            disabled={isLocked}
                          />
                        ))}
                      </div>

                      {!isLocked && (
                        <Button
                          variant="outline"
                          onClick={() => handleAddExercise(sessionIdx)}
                          className="w-full"
                        >
                          <Plus className="h-4 w-4 mr-2" />
                          Add Exercise
                        </Button>
                      )}
                    </div>
                  </AccordionContent>
                </AccordionItem>
              </div>
            );
          })}
        </Accordion>

        <Button variant="outline" onClick={handleAddSession} className="w-full">
          <Plus className="h-4 w-4 mr-2" />
          Add Session
        </Button>
      </div>

      <div className="flex justify-between">
        {isEditing && onDelete && (
          <Button variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
            <Trash2 className="h-4 w-4 mr-2" />
            Delete Program
          </Button>
        )}
        <Button
          onClick={() => onSave(sessionGroupsToEntries(sessions), customFields)}
          disabled={isSaving || !programName}
          className={!isEditing ? 'ml-auto' : ''}
        >
          {isSaving ? 'Saving...' : isEditing ? 'Update Program' : 'Create Program'}
        </Button>
      </div>

      <DeleteConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        onConfirm={() => {
          onDelete?.();
          setDeleteDialogOpen(false);
        }}
        title="Delete Program?"
        description="This will permanently delete this program. This action cannot be undone."
      />
    </div>
  );
}
