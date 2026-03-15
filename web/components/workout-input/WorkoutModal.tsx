'use client';

import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { EntryType, ProgramEntry, WorkoutEntryCreate } from '@/types/api';
import { useEffect, useState } from 'react';
import { AddExerciseButton } from './AddExerciseButton';
import { ExerciseInputRow } from './ExerciseInputRow';

interface WorkoutModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  dayEntries?: ProgramEntry[];
  programContext?: string[];
  onSave: (
    entries: WorkoutEntryCreate[],
    notes: string,
    programContext?: string[]
  ) => Promise<void>;
}

interface EntryInput {
  id: string;
  exercise_id: string;
  exercise_name: string;
  entry_type: EntryType;
  sets: number;
  reps: number;
  load: number;
  rpe: number;
  program_node_id?: string;
  plan?: {
    sets?: number;
    reps?: number;
    load?: number;
    rpe?: number;
  };
}

/**
 * Modal for recording workout results
 *
 * Pre-populates planned exercises from program entries
 * Allows adding unplanned exercises
 * Includes session notes field
 */
export function WorkoutModal({
  open,
  onOpenChange,
  dayEntries,
  programContext,
  onSave,
}: WorkoutModalProps) {
  const [entries, setEntries] = useState<EntryInput[]>([]);
  const [sessionNotes, setSessionNotes] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  // Initialize entries from day entries when modal opens
  useEffect(() => {
    if (open && dayEntries && dayEntries.length > 0) {
      const initialEntries: EntryInput[] = dayEntries.map((entry, idx) => ({
        id: `${entry.id}-${idx}`,
        exercise_id: entry.exercise_id || '',
        exercise_name: entry.name,
        entry_type: 'Main',
        sets: entry.target_sets || 3,
        reps: entry.target_reps || 10,
        load: 0,
        rpe: entry.target_rpe || 7,
        program_node_id: entry.id,
        plan: {
          sets: entry.target_sets,
          reps: entry.target_reps,
          load: undefined,
          rpe: entry.target_rpe,
        },
      }));
      setEntries(initialEntries);
    }
  }, [open, dayEntries]);

  const handleAddExercise = () => {
    const newEntry: EntryInput = {
      id: `new-${Date.now()}`,
      exercise_id: '', // Will need exercise selector in future
      exercise_name: 'New Exercise',
      entry_type: null,
      sets: 3,
      reps: 10,
      load: 0,
      rpe: 7,
    };
    setEntries([...entries, newEntry]);
  };

  const handleRemoveEntry = (id: string) => {
    setEntries(entries.filter((e) => e.id !== id));
  };

  const handleSave = async () => {
    setIsSaving(true);
    try {
      const workoutEntries: WorkoutEntryCreate[] = entries.map((entry) => ({
        exercise_id: entry.exercise_id,
        display_name: entry.exercise_name,
        entry_type: entry.entry_type,
        sets: entry.sets,
        reps: entry.reps,
        load_kg: entry.load,
        rpe: entry.rpe,
        program_node_id: entry.program_node_id,
        plan_snapshot: entry.plan
          ? {
              target_sets: entry.plan.sets,
              target_reps: entry.plan.reps,
              target_load_kg: entry.plan.load,
              target_rpe: entry.plan.rpe,
            }
          : undefined,
      }));

      await onSave(workoutEntries, sessionNotes, programContext);
      onOpenChange(false);
      setEntries([]);
      setSessionNotes('');
    } catch (error) {
      console.error('Failed to save workout:', error);
    } finally {
      setIsSaving(false);
    }
  };

  const dayName = dayEntries?.[0]
    ? (dayEntries[0].metadata?.day as string) || dayEntries[0].name
    : undefined;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            Record Workout
            {dayName && <span className="ml-2 text-muted-foreground">- {dayName}</span>}
          </DialogTitle>
          {programContext && programContext.length > 0 && (
            <p className="text-sm text-muted-foreground">{programContext.join(' > ')}</p>
          )}
        </DialogHeader>

        <div className="space-y-4">
          {entries.map((entry) => (
            <ExerciseInputRow
              key={entry.id}
              exerciseName={entry.exercise_name}
              entryType={entry.entry_type}
              sets={entry.sets}
              reps={entry.reps}
              load={entry.load}
              rpe={entry.rpe}
              onEntryTypeChange={(value) => {
                setEntries(
                  entries.map((e) => (e.id === entry.id ? { ...e, entry_type: value } : e))
                );
              }}
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
              onRemove={!entry.program_node_id ? () => handleRemoveEntry(entry.id) : undefined}
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
              {isSaving ? 'Saving...' : 'Save Workout'}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
