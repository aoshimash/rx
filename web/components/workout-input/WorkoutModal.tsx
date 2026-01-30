'use client';

import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { ExerciseInputRow } from './ExerciseInputRow';
import { AddExerciseButton } from './AddExerciseButton';
import { useState } from 'react';
import type { ProgramNode, WorkoutEntryCreate } from '@/types/api';

interface WorkoutModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  dayNode?: ProgramNode;
  onSave: (entries: WorkoutEntryCreate[], notes: string) => Promise<void>;
}

interface EntryInput {
  id: string;
  exercise_id: string;
  exercise_name: string;
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
 * Pre-populates planned exercises from program node
 * Allows adding unplanned exercises
 * Includes session notes field
 */
export function WorkoutModal({
  open,
  onOpenChange,
  dayNode,
  onSave,
}: WorkoutModalProps) {
  const [entries, setEntries] = useState<EntryInput[]>([]);
  const [sessionNotes, setSessionNotes] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  // Initialize entries from day node when modal opens
  useState(() => {
    if (open && dayNode) {
      const exerciseNodes =
        dayNode.children?.filter((child) => child.node_type === 'exercise') || [];

      const initialEntries: EntryInput[] = exerciseNodes.map((node, idx) => ({
        id: `${node.id}-${idx}`,
        exercise_id: node.exercise_id || '',
        exercise_name: node.name,
        sets: node.target_sets || 3,
        reps: node.target_reps || 10,
        load: 0,
        rpe: node.target_rpe || 7,
        program_node_id: node.id,
        plan: {
          sets: node.target_sets,
          reps: node.target_reps,
          load: undefined,
          rpe: node.target_rpe,
        },
      }));

      setEntries(initialEntries);
    }
  });

  const handleAddExercise = () => {
    const newEntry: EntryInput = {
      id: `new-${Date.now()}`,
      exercise_id: '', // Will need exercise selector in future
      exercise_name: 'New Exercise',
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
      const workoutEntries: WorkoutEntryCreate[] = entries.map((entry, idx) => ({
        exercise_id: entry.exercise_id,
        display_name: entry.exercise_name,
        entry_type: 'main',
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

      await onSave(workoutEntries, sessionNotes);
      onOpenChange(false);
      setEntries([]);
      setSessionNotes('');
    } catch (error) {
      console.error('Failed to save workout:', error);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            Record Workout
            {dayNode && <span className="ml-2 text-muted-foreground">- {dayNode.name}</span>}
          </DialogTitle>
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
                setEntries(
                  entries.map((e) =>
                    e.id === entry.id ? { ...e, sets: value } : e
                  )
                );
              }}
              onRepsChange={(value) => {
                setEntries(
                  entries.map((e) =>
                    e.id === entry.id ? { ...e, reps: value } : e
                  )
                );
              }}
              onLoadChange={(value) => {
                setEntries(
                  entries.map((e) =>
                    e.id === entry.id ? { ...e, load: value } : e
                  )
                );
              }}
              onRpeChange={(value) => {
                setEntries(
                  entries.map((e) =>
                    e.id === entry.id ? { ...e, rpe: value } : e
                  )
                );
              }}
              onRemove={
                !entry.program_node_id
                  ? () => handleRemoveEntry(entry.id)
                  : undefined
              }
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
