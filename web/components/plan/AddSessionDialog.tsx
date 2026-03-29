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
import type { PlanSessionCreate, PlanSessionEntryCreate } from '@/types/api';
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

export function AddSessionDialog({ open, onOpenChange, onAdd, nextOrder }: AddSessionDialogProps) {
  const [sessionName, setSessionName] = useState('');
  const [date, setDate] = useState('');
  const [entries, setEntries] = useState<EntryDraft[]>([createEmptyEntry()]);

  const handleAdd = () => {
    if (!sessionName.trim()) return;

    const validEntries = entries.filter((e) => e.exercise_name.trim() !== '');
    const planEntries: PlanSessionEntryCreate[] = validEntries.map((e, i) => ({
      exercise_name: e.exercise_name.trim(),
      order: i,
      fields:
        Object.keys(e.fields).length > 0
          ? Object.fromEntries(
              Object.entries(e.fields)
                .filter(([, v]) => v.trim() !== '')
                .map(([k, v]) => {
                  const num = Number(v);
                  return [k, Number.isNaN(num) ? v : num];
                })
            )
          : undefined,
    }));

    onAdd({
      session_name: sessionName.trim(),
      order: nextOrder,
      date: date || undefined,
      entries: planEntries.length > 0 ? planEntries : undefined,
    });

    setSessionName('');
    setDate('');
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
            <Label>Exercises</Label>
            {entries.map((entry) => (
              <div key={entry.id} className="flex items-center gap-2">
                <Input
                  placeholder="Exercise name"
                  value={entry.exercise_name}
                  onChange={(e) => updateEntry(entry.id, 'exercise_name', e.target.value)}
                  className="flex-1"
                />
                <Input
                  placeholder="load_kg"
                  value={entry.fields.load_kg ?? ''}
                  onChange={(e) => updateEntryField(entry.id, 'load_kg', e.target.value)}
                  className="w-20"
                />
                <Input
                  placeholder="sets"
                  value={entry.fields.sets ?? ''}
                  onChange={(e) => updateEntryField(entry.id, 'sets', e.target.value)}
                  className="w-16"
                />
                <Input
                  placeholder="reps"
                  value={entry.fields.reps ?? ''}
                  onChange={(e) => updateEntryField(entry.id, 'reps', e.target.value)}
                  className="w-16"
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
            Add to Plan
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
