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
import { useFieldGroups } from '@/lib/hooks/useFieldGroups';
import type { FieldDef } from '@/types/api';
import { ExternalLink, Lock, Plus, Trash2, X } from 'lucide-react';
import Link from 'next/link';
import { useState } from 'react';

// ============================================================================
// Types
// ============================================================================

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
  field_group_id?: string;
  exercises: SessionExercise[];
}

// ============================================================================
// Conversion helpers
// ============================================================================

function entryToExercise(entry: ProgramFormEntry): SessionExercise {
  return {
    exercise_name: entry.exercise_name,
    fields: entry.fields,
  };
}

function entriesToSessionGroups(entries: ProgramFormEntry[]): SessionGroup[] {
  const sessionMap = new Map<string, SessionExercise[]>();
  const sessionDates = new Map<string, string>();
  const sessionFieldGroupIds = new Map<string, string>();
  const sessionOrder: string[] = [];

  for (const entry of entries) {
    const sessionName = (entry.metadata?.session as string) ?? 'Session 1';
    if (!sessionMap.has(sessionName)) {
      sessionMap.set(sessionName, []);
      sessionOrder.push(sessionName);
      const date = entry.metadata?.date as string | undefined;
      if (date) sessionDates.set(sessionName, date);
      const fgId = entry.metadata?.field_group_id as string | undefined;
      if (fgId) sessionFieldGroupIds.set(sessionName, fgId);
    }
    sessionMap.get(sessionName)?.push(entryToExercise(entry));
  }

  return sessionOrder.map((name) => ({
    name,
    date: sessionDates.get(name),
    field_group_id: sessionFieldGroupIds.get(name),
    exercises: sessionMap.get(name) || [],
  }));
}

function exerciseToEntry(
  ex: SessionExercise,
  sessionName: string,
  sessionDate: string | undefined,
  fieldGroupId: string | undefined,
  order: number
): ProgramFormEntry {
  const metadata: Record<string, unknown> = { session: sessionName };
  if (sessionDate) metadata.date = sessionDate;
  if (fieldGroupId) metadata.field_group_id = fieldGroupId;
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
      entries.push(exerciseToEntry(ex, session.name, session.date, session.field_group_id, order));
      order++;
    }
  }

  return entries;
}

// ============================================================================
// ProgramExerciseRow
// ============================================================================

interface ProgramExerciseRowProps {
  exercise: SessionExercise;
  onChange: (updated: SessionExercise) => void;
  onRemove: () => void;
  fields: FieldDef[];
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
    const updated = { ...exercise.fields };
    if (raw) {
      updated[key] = isNumber ? Number(raw) : raw;
    } else {
      delete updated[key];
    }
    onChange({ ...exercise, fields: updated });
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

      {fields.length > 0 && (
        <div
          className="grid gap-3"
          style={{ gridTemplateColumns: `repeat(${fields.length}, minmax(0, 1fr))` }}
        >
          {fields.map((field) => {
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
                  value={(exercise.fields?.[field.name] as string | number | undefined) ?? ''}
                  onChange={(e) => updateField(field.name, e.target.value, field.type === 'number')}
                  placeholder={field.description}
                  disabled={disabled}
                />
              </div>
            );
          })}
        </div>
      )}
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
  onNameChange: (name: string) => void;
  onNotesChange: (notes: string) => void;
  onSave: (entries: ProgramFormEntry[]) => void;
  onDelete?: () => void;
  isSaving?: boolean;
  isEditing?: boolean;
  lockedSessionNames?: Set<string>;
}

export function ProgramForm({
  programName,
  programNotes,
  initialEntries,
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

  const { data: fieldGroupsData } = useFieldGroups();
  const fieldGroups = fieldGroupsData?.data ?? [];

  const getFieldsForSession = (session: SessionGroup): FieldDef[] => {
    if (!session.field_group_id) return [];
    const group = fieldGroups.find((g) => g.id === session.field_group_id);
    return group?.program_fields ?? [];
  };

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

  const handleSessionFieldGroupChange = (idx: number, fieldGroupId: string | undefined) => {
    const updated = [...sessions];
    const session = updated[idx];
    if (!session) return;
    updated[idx] = { ...session, field_group_id: fieldGroupId };
    setSessions(updated);
  };

  const handleAddExercise = (sessionIdx: number) => {
    const updated = [...sessions];
    const session = updated[sessionIdx];
    if (!session) return;
    updated[sessionIdx] = {
      ...session,
      exercises: [...session.exercises, { exercise_name: '' }],
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
      </div>

      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <p className="text-sm font-semibold">Sessions</p>
          <Link
            href="/settings#field-groups"
            className="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground"
          >
            Manage Field Groups
            <ExternalLink className="h-3 w-3" />
          </Link>
        </div>
        <Accordion type="multiple" className="w-full space-y-3">
          {sessions.map((session, sessionIdx) => {
            const isLocked = lockedSessionNames?.has(session.name) ?? false;
            const sessionFields = getFieldsForSession(session);
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
                      {!isLocked && (
                        <div className="space-y-2">
                          <Label className="text-xs">Field Group</Label>
                          <Select
                            value={session.field_group_id ?? ''}
                            onValueChange={(v) =>
                              handleSessionFieldGroupChange(sessionIdx, v || undefined)
                            }
                          >
                            <SelectTrigger className="h-8 w-64">
                              <SelectValue placeholder="Select field group..." />
                            </SelectTrigger>
                            <SelectContent>
                              {fieldGroups.map((fg) => (
                                <SelectItem key={fg.id} value={fg.id}>
                                  {fg.name}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </div>
                      )}

                      <div className="space-y-3">
                        {session.exercises.map((exercise, exIdx) => (
                          <ProgramExerciseRow
                            key={exIdx}
                            exercise={exercise}
                            onChange={(updated) => handleExerciseChange(sessionIdx, exIdx, updated)}
                            onRemove={() => handleRemoveExercise(sessionIdx, exIdx)}
                            fields={sessionFields}
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
          onClick={() => onSave(sessionGroupsToEntries(sessions))}
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
