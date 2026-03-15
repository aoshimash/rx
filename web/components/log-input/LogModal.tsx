'use client';

import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { LogEntryCreate, PlanEntry } from '@/types/api';
import { useEffect, useState } from 'react';
import { AddExerciseButton } from './AddExerciseButton';
import { ExerciseInputRow } from './ExerciseInputRow';

interface LogModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  dayEntries?: PlanEntry[];
  planContext?: string[];
  onSave: (entries: LogEntryCreate[], notes: string) => Promise<void>;
}

interface EntryInput {
  id: string;
  exercise_name: string;
  sets: number;
  reps: number;
  load: number;
  rpe: number;
  fromPlan: boolean;
  plan?: {
    sets?: number;
    reps?: number;
    load_kg?: number;
    rpe?: number;
  };
}

export function LogModal({ open, onOpenChange, dayEntries, planContext, onSave }: LogModalProps) {
  const [entries, setEntries] = useState<EntryInput[]>([]);
  const [sessionNotes, setSessionNotes] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (open && dayEntries && dayEntries.length > 0) {
      const initialEntries: EntryInput[] = dayEntries.map((entry, idx) => ({
        id: `${entry.id}-${idx}`,
        exercise_name: entry.exercise_name,
        sets: entry.sets || 3,
        reps: entry.reps || 10,
        load: entry.load_kg || 0,
        rpe: entry.rpe || 7,
        fromPlan: true,
        plan: {
          sets: entry.sets,
          reps: entry.reps,
          load_kg: entry.load_kg,
          rpe: entry.rpe,
        },
      }));
      setEntries(initialEntries);
    }
  }, [open, dayEntries]);

  const handleAddExercise = () => {
    const newEntry: EntryInput = {
      id: `new-${Date.now()}`,
      exercise_name: 'New Exercise',
      sets: 3,
      reps: 10,
      load: 0,
      rpe: 7,
      fromPlan: false,
    };
    setEntries([...entries, newEntry]);
  };

  const handleRemoveEntry = (id: string) => {
    setEntries(entries.filter((e) => e.id !== id));
  };

  const handleSave = async () => {
    setIsSaving(true);
    try {
      const logEntries: LogEntryCreate[] = entries.map((entry) => ({
        exercise_name: entry.exercise_name,
        sets: entry.sets,
        reps: entry.reps,
        load_kg: entry.load,
        rpe: entry.rpe,
      }));

      await onSave(logEntries, sessionNotes);
      onOpenChange(false);
      setEntries([]);
      setSessionNotes('');
    } catch (error) {
      console.error('Failed to save log:', error);
    } finally {
      setIsSaving(false);
    }
  };

  const dayName = dayEntries?.[0]
    ? (dayEntries[0].metadata?.day as string) || dayEntries[0].exercise_name
    : undefined;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            Record Log
            {dayName && <span className="ml-2 text-muted-foreground">- {dayName}</span>}
          </DialogTitle>
          {planContext && planContext.length > 0 && (
            <p className="text-sm text-muted-foreground">{planContext.join(' > ')}</p>
          )}
        </DialogHeader>

        <div className="space-y-4">
          {entries.map((entry) => (
            <ExerciseInputRow
              key={entry.id}
              exerciseName={entry.exercise_name}
              sets={entry.sets}
              reps={entry.reps}
              load={entry.load}
              rpe={entry.rpe}
              onSetsChange={(value) => {
                setEntries(entries.map((e) => (e.id === entry.id ? { ...e, sets: value } : e)));
              }}
              onRepsChange={(value) => {
                setEntries(entries.map((e) => (e.id === entry.id ? { ...e, reps: value } : e)));
              }}
              onLoadChange={(value) => {
                setEntries(entries.map((e) => (e.id === entry.id ? { ...e, load: value } : e)));
              }}
              onRpeChange={(value) => {
                setEntries(entries.map((e) => (e.id === entry.id ? { ...e, rpe: value } : e)));
              }}
              onRemove={!entry.fromPlan ? () => handleRemoveEntry(entry.id) : undefined}
              planValues={entry.plan}
            />
          ))}

          <AddExerciseButton onClick={handleAddExercise} />

          <div className="space-y-2 pt-4 border-t">
            <Label htmlFor="session-notes">Session Notes</Label>
            <Input
              id="session-notes"
              placeholder="How did the session go?"
              value={sessionNotes}
              onChange={(e) => setSessionNotes(e.target.value)}
            />
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button onClick={handleSave} disabled={isSaving || entries.length === 0}>
              {isSaving ? 'Saving...' : 'Save Log'}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
