'use client';

import { InlineEdit } from '@/components/plan/InlineEdit';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useFieldGroups } from '@/lib/hooks/useFieldGroups';
import type { FieldDef, PlanSession, PlanSessionEntryCreate } from '@/types/api';
import { CalendarDays, ClipboardList, Plus, Trash2, X } from 'lucide-react';
import { useCallback } from 'react';

interface SessionCardProps {
  session: PlanSession;
  programName?: string;
  onRecord: (session: PlanSession) => void;
  onDelete: (sessionId: string) => void;
  onUpdate: (
    sessionId: string,
    data: {
      session_name: string;
      order: number;
      field_group_id?: string;
      date?: string;
      source_program_id?: string;
      source_session_id?: string;
      entries?: PlanSessionEntryCreate[];
    }
  ) => Promise<void>;
}

function parseFieldValue(value: string, type: FieldDef['type']): string | number {
  if (type === 'number') {
    const num = Number(value);
    return Number.isNaN(num) ? value : num;
  }
  return value;
}

/** Collect unique field keys across all entries, preserving insertion order. */
function collectFieldKeys(session: PlanSession): string[] {
  const keys: string[] = [];
  const seen = new Set<string>();
  for (const entry of session.entries) {
    if (!entry.fields) continue;
    for (const key of Object.keys(entry.fields)) {
      if (!seen.has(key)) {
        seen.add(key);
        keys.push(key);
      }
    }
  }
  return keys;
}

function resolveFields(
  rawFields: Record<string, unknown> | undefined,
  programFields: FieldDef[]
): Record<string, unknown> | undefined {
  if (!rawFields) return undefined;
  const result: Record<string, unknown> = {};
  for (const [k, v] of Object.entries(rawFields)) {
    if (v == null || String(v).trim() === '') continue;
    const fieldDef = programFields.find((f) => f.name === k);
    result[k] = fieldDef ? parseFieldValue(String(v), fieldDef.type) : v;
  }
  return Object.keys(result).length > 0 ? result : undefined;
}

function entriesToCreatePayload(
  entries: PlanSession['entries'],
  programFields: FieldDef[]
): PlanSessionEntryCreate[] | undefined {
  const sorted = [...entries].sort((a, b) => a.order - b.order);
  if (sorted.length === 0) return undefined;
  return sorted.map((e, i) => ({
    exercise_name: e.exercise_name,
    order: i,
    fields: resolveFields(e.fields, programFields),
    notes: e.notes || undefined,
  }));
}

function entryToCreate(e: PlanSession['entries'][number], order: number): PlanSessionEntryCreate {
  return {
    exercise_name: e.exercise_name,
    order,
    fields: e.fields ? { ...e.fields } : undefined,
    notes: e.notes || undefined,
  };
}

export function SessionCard({
  session,
  programName,
  onRecord,
  onDelete,
  onUpdate,
}: SessionCardProps) {
  const { data: fieldGroupsData } = useFieldGroups();
  const fieldGroups = fieldGroupsData?.data ?? [];
  const selectedGroup = fieldGroups.find((g) => g.id === session.field_group_id);
  const programFields = selectedGroup?.program_fields ?? [];

  const saveSession = useCallback(
    (
      patch: Partial<{
        session_name: string;
        date: string | undefined;
        field_group_id: string | undefined;
        entries: PlanSessionEntryCreate[];
      }>
    ) => {
      const entries = patch.entries ?? entriesToCreatePayload(session.entries, programFields);
      onUpdate(session.id, {
        session_name: patch.session_name ?? session.session_name,
        order: session.order,
        date: patch.date !== undefined ? patch.date : session.date || undefined,
        field_group_id:
          patch.field_group_id !== undefined
            ? patch.field_group_id
            : session.field_group_id || undefined,
        source_program_id: session.source_program_id || undefined,
        source_session_id: session.source_session_id || undefined,
        entries,
      });
    },
    [session, programFields, onUpdate]
  );

  const handleSessionNameSave = useCallback(
    (value: string) => {
      if (!value) return;
      saveSession({ session_name: value });
    },
    [saveSession]
  );

  const handleDateSave = useCallback(
    (value: string) => {
      saveSession({ date: value || undefined });
    },
    [saveSession]
  );

  const handleFieldGroupChange = useCallback(
    (value: string) => {
      saveSession({ field_group_id: value || undefined });
    },
    [saveSession]
  );

  const handleExerciseNameSave = useCallback(
    (entryId: string, value: string) => {
      if (!value) return;
      const updatedEntries = [...session.entries]
        .sort((a, b) => a.order - b.order)
        .map((e, i) => ({
          ...entryToCreate(e, i),
          exercise_name: e.id === entryId ? value : e.exercise_name,
        }));
      saveSession({ entries: updatedEntries });
    },
    [session.entries, saveSession]
  );

  const handleFieldValueSave = useCallback(
    (entryId: string, fieldName: string, value: string) => {
      const updatedEntries = [...session.entries]
        .sort((a, b) => a.order - b.order)
        .map((e, i) => {
          const base = entryToCreate(e, i);
          if (e.id !== entryId) return base;
          const fields = { ...(e.fields ?? {}) };
          const fieldDef = programFields.find((f) => f.name === fieldName);
          if (value.trim()) {
            fields[fieldName] = fieldDef ? parseFieldValue(value, fieldDef.type) : value;
          } else {
            delete fields[fieldName];
          }
          return { ...base, fields: Object.keys(fields).length > 0 ? fields : undefined };
        });
      saveSession({ entries: updatedEntries });
    },
    [session.entries, programFields, saveSession]
  );

  const handleAddExercise = useCallback(() => {
    const sorted = [...session.entries].sort((a, b) => a.order - b.order);
    const newEntries: PlanSessionEntryCreate[] = [
      ...sorted.map((e, i) => entryToCreate(e, i)),
      { exercise_name: 'New Exercise', order: sorted.length },
    ];
    saveSession({ entries: newEntries });
  }, [session.entries, saveSession]);

  const handleRemoveExercise = useCallback(
    (entryId: string) => {
      const remaining = [...session.entries]
        .sort((a, b) => a.order - b.order)
        .filter((e) => e.id !== entryId)
        .map((e, i) => entryToCreate(e, i));
      saveSession({ entries: remaining.length > 0 ? remaining : [] });
    },
    [session.entries, saveSession]
  );

  const sortedEntries = [...session.entries].sort((a, b) => a.order - b.order);
  const fieldKeys = collectFieldKeys(session);

  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 min-w-0">
            <InlineEdit
              value={session.session_name}
              onSave={handleSessionNameSave}
              placeholder="Session name"
              className="font-semibold"
            />
            {session.source_program_id && programName ? (
              <span className="text-xs text-muted-foreground shrink-0">from {programName}</span>
            ) : !session.source_program_id ? (
              <span className="text-xs text-muted-foreground italic shrink-0">manual</span>
            ) : null}
          </div>
          <div className="flex items-center gap-2 shrink-0">
            <Button
              variant="outline"
              size="sm"
              className="h-7 text-xs gap-1"
              onClick={() => onRecord(session)}
            >
              <ClipboardList className="h-3.5 w-3.5" />
              Record
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6 text-muted-foreground hover:text-destructive"
              onClick={() => {
                if (window.confirm(`Remove "${session.session_name}" from plan?`)) {
                  onDelete(session.id);
                }
              }}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
        </div>

        <div className="flex items-center gap-3 mt-1">
          <span className="flex items-center gap-1 text-xs text-muted-foreground">
            <CalendarDays className="h-3 w-3" />
            <InlineEdit
              value={session.date ?? ''}
              onSave={handleDateSave}
              type="date"
              emptyDisplay="No date"
              className="text-xs"
              inputClassName="text-xs w-[140px]"
            />
          </span>
          <span className="text-xs text-muted-foreground">
            <Select value={session.field_group_id ?? ''} onValueChange={handleFieldGroupChange}>
              <SelectTrigger className="h-6 text-xs border-none shadow-none px-1 hover:bg-muted">
                <SelectValue placeholder="Field group..." />
              </SelectTrigger>
              <SelectContent>
                {fieldGroups.map((fg) => (
                  <SelectItem key={fg.id} value={fg.id}>
                    {fg.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </span>
        </div>
      </CardHeader>

      <CardContent className="pt-0">
        {sortedEntries.length === 0 ? (
          <p className="text-sm text-muted-foreground">No exercises</p>
        ) : (
          <table className="text-sm border-collapse w-full">
            <thead>
              <tr className="text-xs text-muted-foreground">
                <th className="text-left font-medium pb-1 pr-4 w-[180px]">Exercise</th>
                {fieldKeys.map((key) => (
                  <th key={key} className="text-left font-medium pb-1 px-2 w-[80px]">
                    {key}
                  </th>
                ))}
                <th className="w-[32px]" />
              </tr>
            </thead>
            <tbody>
              {sortedEntries.map((entry) => (
                <tr key={entry.id}>
                  <td className="pr-4 py-0.5">
                    <InlineEdit
                      value={entry.exercise_name}
                      onSave={(v) => handleExerciseNameSave(entry.id, v)}
                      placeholder="Exercise name"
                    />
                  </td>
                  {fieldKeys.map((key) => (
                    <td key={key} className="px-2 py-0.5 tabular-nums">
                      <InlineEdit
                        value={entry.fields?.[key] != null ? String(entry.fields[key]) : ''}
                        onSave={(v) => handleFieldValueSave(entry.id, key, v)}
                        type={
                          programFields.find((f) => f.name === key)?.type === 'number'
                            ? 'number'
                            : 'text'
                        }
                        emptyDisplay="—"
                        className="text-muted-foreground"
                      />
                    </td>
                  ))}
                  <td className="py-0.5">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-6 w-6 text-muted-foreground hover:text-destructive"
                      onClick={() => handleRemoveExercise(entry.id)}
                    >
                      <Trash2 className="h-3 w-3" />
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        <Button
          variant="outline"
          size="sm"
          className="h-7 text-xs mt-2"
          onClick={handleAddExercise}
        >
          <Plus className="h-3 w-3 mr-1" />
          Add Exercise
        </Button>
      </CardContent>
    </Card>
  );
}
