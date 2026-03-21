'use client';

import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { LogEntryCreate, ProgramSessionEntryCreate } from '@/types/api';
import { useEffect, useState } from 'react';
import { AddExerciseButton } from './AddExerciseButton';
import { ExerciseInputRow } from './ExerciseInputRow';

export interface LogSaveContext {
  programId?: string;
  sessionName?: string;
  startedAt?: string;
  finishedAt?: string;
}

interface LogModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  dayEntries?: ProgramSessionEntryCreate[];
  sessionContext?: string[];
  programId?: string;
  sessionName?: string;
  onSave: (entries: LogEntryCreate[], notes: string, context?: LogSaveContext) => Promise<void>;
}

interface EntryInput {
  id: string;
  exercise_name: string;
  sets: number;
  reps: number;
  load: number;
  rpe: number;
  startedAt: string;
  finishedAt: string;
  fromPlan: boolean;
  plan?: {
    sets?: number;
    reps?: number;
    load_kg?: number;
    rpe?: number;
  };
}

export function LogModal({
  open,
  onOpenChange,
  dayEntries,
  sessionContext,
  programId,
  sessionName,
  onSave,
}: LogModalProps) {
  const [entries, setEntries] = useState<EntryInput[]>([]);
  const [sessionNotes, setSessionNotes] = useState('');
  const [sessionStartedAt, setSessionStartedAt] = useState('');
  const [sessionFinishedAt, setSessionFinishedAt] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (open && dayEntries && dayEntries.length > 0) {
      const initialEntries: EntryInput[] = dayEntries.map((entry, idx) => ({
        id: `plan-${idx}`,
        exercise_name: entry.exercise_name,
        sets: entry.sets || 3,
        reps: entry.reps || 10,
        load: entry.load_kg || 0,
        rpe: entry.rpe || 7,
        startedAt: '',
        finishedAt: '',
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
      startedAt: '',
      finishedAt: '',
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
        started_at: entry.startedAt ? new Date(entry.startedAt).toISOString() : undefined,
        finished_at: entry.finishedAt ? new Date(entry.finishedAt).toISOString() : undefined,
      }));

      const context: LogSaveContext = {
        programId,
        sessionName,
        startedAt: sessionStartedAt ? new Date(sessionStartedAt).toISOString() : undefined,
        finishedAt: sessionFinishedAt ? new Date(sessionFinishedAt).toISOString() : undefined,
      };
      await onSave(logEntries, sessionNotes, context);
      onOpenChange(false);
      setEntries([]);
      setSessionNotes('');
      setSessionStartedAt('');
      setSessionFinishedAt('');
    } catch (error) {
      console.error('Failed to save log:', error);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            Record Log
            {sessionName && <span className="ml-2 text-muted-foreground">- {sessionName}</span>}
          </DialogTitle>
          {sessionContext && sessionContext.length > 0 && (
            <p className="text-sm text-muted-foreground">{sessionContext.join(' > ')}</p>
          )}
        </DialogHeader>

        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-3 p-4 bg-muted/50 rounded-lg">
            <div className="space-y-2">
              <Label htmlFor="session-started-at">Session Start</Label>
              <Input
                id="session-started-at"
                type="datetime-local"
                value={sessionStartedAt}
                onChange={(e) => setSessionStartedAt(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="session-finished-at">Session End</Label>
              <Input
                id="session-finished-at"
                type="datetime-local"
                value={sessionFinishedAt}
                onChange={(e) => setSessionFinishedAt(e.target.value)}
              />
            </div>
          </div>

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
              startedAt={entry.startedAt}
              finishedAt={entry.finishedAt}
              onStartedAtChange={(value) => {
                setEntries(
                  entries.map((e) => (e.id === entry.id ? { ...e, startedAt: value } : e))
                );
              }}
              onFinishedAtChange={(value) => {
                setEntries(
                  entries.map((e) => (e.id === entry.id ? { ...e, finishedAt: value } : e))
                );
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
