'use client';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
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
import type { FieldDef, PlanSession, PlanSessionEntryCreate } from '@/types/api';
import { CalendarDays, Check, Pencil, Plus, Trash2, X } from 'lucide-react';
import { useEffect, useState } from 'react';

interface SessionCardProps {
  session: PlanSession;
  programName?: string;
  isNew?: boolean;
  onLog: (session: PlanSession) => void;
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
  onNewHandled?: () => void;
}

interface EntryDraft {
  id: string;
  exercise_name: string;
  fields: Record<string, string>;
  notes: string;
}

function createEmptyEntry(): EntryDraft {
  return { id: crypto.randomUUID(), exercise_name: '', fields: {}, notes: '' };
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

function formatFieldValue(value: unknown): string {
  if (value == null) return '—';
  if (typeof value === 'boolean') return value ? 'Yes' : 'No';
  return String(value);
}

function entriesToDrafts(session: PlanSession): EntryDraft[] {
  if (session.entries.length === 0) return [];
  return [...session.entries]
    .sort((a, b) => a.order - b.order)
    .map((e) => ({
      id: e.id,
      exercise_name: e.exercise_name,
      fields: Object.fromEntries(
        Object.entries(e.fields ?? {}).map(([k, v]) => [k, String(v ?? '')])
      ),
      notes: e.notes ?? '',
    }));
}

export function SessionCard({
  session,
  programName,
  isNew,
  onLog,
  onDelete,
  onUpdate,
  onNewHandled,
}: SessionCardProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [sessionName, setSessionName] = useState(session.session_name);
  const [date, setDate] = useState(session.date ?? '');
  const [fieldGroupId, setFieldGroupId] = useState(session.field_group_id ?? '');
  const [entries, setEntries] = useState<EntryDraft[]>(entriesToDrafts(session));
  const [isSaving, setIsSaving] = useState(false);

  const { data: fieldGroupsData } = useFieldGroups();
  const fieldGroups = fieldGroupsData?.data ?? [];
  const selectedGroup = fieldGroups.find((g) => g.id === fieldGroupId);
  const programFields = selectedGroup?.program_fields ?? [];

  // Auto-enter edit mode for newly created sessions
  useEffect(() => {
    if (isNew) {
      setIsEditing(true);
      onNewHandled?.();
    }
  }, [isNew, onNewHandled]);

  const startEditing = () => {
    setSessionName(session.session_name);
    setDate(session.date ?? '');
    setFieldGroupId(session.field_group_id ?? '');
    setEntries(entriesToDrafts(session));
    setIsEditing(true);
  };

  const cancelEditing = () => {
    setIsEditing(false);
  };

  const handleSave = async () => {
    if (!sessionName.trim()) return;
    setIsSaving(true);
    try {
      const validEntries = entries.filter((e) => e.exercise_name.trim() !== '');
      const planEntries: PlanSessionEntryCreate[] = validEntries.map((e, i) => {
        const fields: Record<string, unknown> = {};
        for (const [k, v] of Object.entries(e.fields)) {
          if (v.trim() === '') continue;
          const fieldDef = programFields.find((f) => f.name === k);
          fields[k] = fieldDef ? parseFieldValue(v, fieldDef.type) : v;
        }
        return {
          exercise_name: e.exercise_name.trim(),
          order: i,
          fields: Object.keys(fields).length > 0 ? fields : undefined,
          notes: e.notes.trim() || undefined,
        };
      });

      await onUpdate(session.id, {
        session_name: sessionName.trim(),
        order: session.order,
        date: date || undefined,
        field_group_id: fieldGroupId || undefined,
        source_program_id: session.source_program_id || undefined,
        source_session_id: session.source_session_id || undefined,
        entries: planEntries.length > 0 ? planEntries : undefined,
      });
      setIsEditing(false);
    } finally {
      setIsSaving(false);
    }
  };

  const updateEntry = (id: string, field: string, value: string) => {
    setEntries((prev) => prev.map((e) => (e.id === id ? { ...e, [field]: value } : e)));
  };

  const updateEntryField = (id: string, fieldName: string, value: string) => {
    setEntries((prev) =>
      prev.map((e) => (e.id === id ? { ...e, fields: { ...e.fields, [fieldName]: value } } : e))
    );
  };

  const removeEntry = (id: string) => {
    setEntries((prev) => prev.filter((e) => e.id !== id));
  };

  // ─── Edit Mode ───
  if (isEditing) {
    return (
      <Card>
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <Input
              value={sessionName}
              onChange={(e) => setSessionName(e.target.value)}
              placeholder="Session name"
              className="font-semibold text-base h-8 max-w-[300px]"
              autoFocus
            />
            <div className="flex items-center gap-1 shrink-0">
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7 text-primary"
                onClick={handleSave}
                disabled={!sessionName.trim() || isSaving}
              >
                <Check className="h-4 w-4" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7 text-muted-foreground"
                onClick={cancelEditing}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent className="pt-0 space-y-3">
          <div className="flex gap-3">
            <div className="space-y-1 flex-1">
              <Label className="text-xs">Date</Label>
              <Input
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
                className="h-8"
              />
            </div>
            <div className="space-y-1 flex-1">
              <Label className="text-xs">Field Group</Label>
              <Select value={fieldGroupId} onValueChange={setFieldGroupId}>
                <SelectTrigger className="h-8">
                  <SelectValue placeholder="Select..." />
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
          </div>

          <div className="space-y-2">
            <Label className="text-xs">Exercises</Label>
            {entries.map((entry) => (
              <div key={entry.id} className="flex items-center gap-2">
                <Input
                  placeholder="Exercise name"
                  value={entry.exercise_name}
                  onChange={(e) => updateEntry(entry.id, 'exercise_name', e.target.value)}
                  className="flex-1 h-8"
                />
                {programFields.map((field) => (
                  <Input
                    key={field.name}
                    placeholder={field.name}
                    value={entry.fields[field.name] ?? ''}
                    onChange={(e) => updateEntryField(entry.id, field.name, e.target.value)}
                    type={field.type === 'number' ? 'number' : 'text'}
                    className="w-20 h-8"
                  />
                ))}
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7 shrink-0 text-muted-foreground hover:text-destructive"
                  onClick={() => removeEntry(entry.id)}
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </Button>
              </div>
            ))}
            <Button
              variant="outline"
              size="sm"
              className="h-7 text-xs"
              onClick={() => setEntries((prev) => [...prev, createEmptyEntry()])}
            >
              <Plus className="h-3 w-3 mr-1" />
              Add Exercise
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }

  // ─── View Mode ───
  const sortedEntries = [...session.entries].sort((a, b) => a.order - b.order);
  const fieldKeys = collectFieldKeys(session);

  return (
    <Card
      className="cursor-pointer transition-colors hover:border-primary"
      onClick={() => onLog(session)}
    >
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 min-w-0">
            <span className="font-semibold truncate">{session.session_name}</span>
            {session.source_program_id && programName ? (
              <span className="text-xs text-muted-foreground shrink-0">from {programName}</span>
            ) : !session.source_program_id ? (
              <span className="text-xs text-muted-foreground italic shrink-0">manual</span>
            ) : null}
          </div>
          <div className="flex items-center gap-2 shrink-0">
            {session.date && (
              <span className="flex items-center gap-1 text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded-full">
                <CalendarDays className="h-3 w-3" />
                {session.date}
              </span>
            )}
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6 text-muted-foreground hover:text-primary"
              onClick={(e) => {
                e.stopPropagation();
                startEditing();
              }}
            >
              <Pencil className="h-3.5 w-3.5" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6 text-muted-foreground hover:text-destructive"
              onClick={(e) => {
                e.stopPropagation();
                if (window.confirm(`Remove "${session.session_name}" from plan?`)) {
                  onDelete(session.id);
                }
              }}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        {sortedEntries.length === 0 ? (
          <p className="text-sm text-muted-foreground">No exercises</p>
        ) : (
          <table className="text-sm border-collapse">
            <thead>
              <tr className="text-xs text-muted-foreground">
                <th className="text-left font-medium pb-1 pr-4 w-[150px]">Exercise</th>
                {fieldKeys.map((key) => (
                  <th key={key} className="text-left font-medium pb-1 px-2 w-[60px]">
                    {key}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {sortedEntries.map((entry) => (
                <tr key={entry.id} className="text-muted-foreground">
                  <td className="pr-4 py-0.5">{entry.exercise_name}</td>
                  {fieldKeys.map((key) => (
                    <td key={key} className="px-2 py-0.5 tabular-nums">
                      {formatFieldValue(entry.fields?.[key])}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </CardContent>
    </Card>
  );
}
