'use client';

import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
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
import type { FieldDef, PlanSessionCreate, PlanSessionEntryCreate } from '@/types/api';
import { Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';

interface AddSessionDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onAdd: (session: PlanSessionCreate) => void;
  nextOrder: number;
}

interface EntryDraft {
  id: string;
  exercise_name: string;
  fields: Record<string, string>;
}

function createEmptyEntry(): EntryDraft {
  return { id: crypto.randomUUID(), exercise_name: '', fields: {} };
}

function parseFieldValue(value: string, type: FieldDef['type']): string | number {
  if (type === 'number') {
    const num = Number(value);
    return Number.isNaN(num) ? value : num;
  }
  return value;
}

export function AddSessionDialog({ open, onOpenChange, onAdd, nextOrder }: AddSessionDialogProps) {
  const [sessionName, setSessionName] = useState('');
  const [date, setDate] = useState('');
  const [fieldGroupId, setFieldGroupId] = useState<string>('');
  const [entries, setEntries] = useState<EntryDraft[]>([createEmptyEntry()]);

  const { data: fieldGroupsData } = useFieldGroups();
  const fieldGroups = fieldGroupsData?.data ?? [];
  const selectedGroup = fieldGroups.find((g) => g.id === fieldGroupId);
  const programFields = selectedGroup?.program_fields ?? [];

  const handleAdd = () => {
    if (!sessionName.trim()) return;

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
      };
    });

    onAdd({
      session_name: sessionName.trim(),
      order: nextOrder,
      date: date || undefined,
      field_group_id: fieldGroupId || undefined,
      entries: planEntries.length > 0 ? planEntries : undefined,
    });

    setSessionName('');
    setDate('');
    setFieldGroupId('');
    setEntries([createEmptyEntry()]);
    onOpenChange(false);
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

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Add Session</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="session-name">Session Name</Label>
            <Input
              id="session-name"
              placeholder="e.g., Upper Body A"
              value={sessionName}
              onChange={(e) => setSessionName(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="session-date">Date (optional)</Label>
            <Input
              id="session-date"
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label>Field Group</Label>
            <Select value={fieldGroupId} onValueChange={setFieldGroupId}>
              <SelectTrigger>
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

          <div className="space-y-2">
            <Label>Exercises</Label>
            {entries.map((entry) => (
              <div key={entry.id} className="space-y-2 rounded border p-2">
                <div className="flex items-center gap-2">
                  <Input
                    placeholder="Exercise name"
                    value={entry.exercise_name}
                    onChange={(e) => updateEntry(entry.id, 'exercise_name', e.target.value)}
                    className="flex-1"
                  />
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 shrink-0"
                    onClick={() => removeEntry(entry.id)}
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </Button>
                </div>
                {programFields.length > 0 && (
                  <div className="flex items-center gap-2">
                    {programFields.map((field) => (
                      <Input
                        key={field.name}
                        placeholder={field.name}
                        value={entry.fields[field.name] ?? ''}
                        onChange={(e) => updateEntryField(entry.id, field.name, e.target.value)}
                        type={field.type === 'number' ? 'number' : 'text'}
                        className="w-20"
                      />
                    ))}
                  </div>
                )}
              </div>
            ))}
            <Button
              variant="outline"
              size="sm"
              onClick={() => setEntries((prev) => [...prev, createEmptyEntry()])}
            >
              <Plus className="h-3.5 w-3.5 mr-1" />
              Add Exercise
            </Button>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleAdd} disabled={!sessionName.trim()}>
            Add Session
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
