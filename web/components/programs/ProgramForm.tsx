'use client';

import { DeleteConfirmDialog } from '@/components/plan-editor/DeleteConfirmDialog';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { ProgramEntryCreate } from '@/types/api';
import { Plus, Trash2, X } from 'lucide-react';
import { useState } from 'react';

// ============================================================================
// Types
// ============================================================================

interface SessionExercise {
  exercise_name: string;
  sets?: number;
  reps?: number;
  rpe?: number;
  percent_1rm_display?: number; // displayed as 0-100%, stored as 0-1 in API
}

interface SessionGroup {
  name: string;
  exercises: SessionExercise[];
}

// ============================================================================
// Conversion helpers
// ============================================================================

function entriesToSessionGroups(entries: ProgramEntryCreate[]): SessionGroup[] {
  const sessionMap = new Map<string, SessionExercise[]>();
  const sessionOrder: string[] = [];

  for (const entry of entries) {
    const sessionName = (entry.metadata?.session as string) ?? 'Session 1';
    if (!sessionMap.has(sessionName)) {
      sessionMap.set(sessionName, []);
      sessionOrder.push(sessionName);
    }
    sessionMap.get(sessionName)!.push({
      exercise_name: entry.exercise_name,
      sets: entry.sets,
      reps: entry.reps,
      rpe: entry.rpe,
      percent_1rm_display:
        entry.percent_1rm !== undefined ? Math.round(entry.percent_1rm * 100) : undefined,
    });
  }

  return sessionOrder.map((name) => ({
    name,
    exercises: sessionMap.get(name) || [],
  }));
}

function sessionGroupsToEntries(sessions: SessionGroup[]): ProgramEntryCreate[] {
  const entries: ProgramEntryCreate[] = [];
  let order = 0;

  for (const session of sessions) {
    for (const ex of session.exercises) {
      entries.push({
        exercise_name: ex.exercise_name,
        order,
        sets: ex.sets,
        reps: ex.reps,
        rpe: ex.rpe,
        percent_1rm:
          ex.percent_1rm_display !== undefined ? ex.percent_1rm_display / 100 : undefined,
        metadata: { session: session.name },
      });
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
}

function ProgramExerciseRow({ exercise, onChange, onRemove }: ProgramExerciseRowProps) {
  return (
    <div className="grid gap-4 border rounded-lg p-4">
      <div className="flex items-center justify-between">
        <Label>Exercise</Label>
        <Button variant="ghost" size="sm" onClick={onRemove}>
          <X className="h-4 w-4" />
        </Button>
      </div>

      <Input
        value={exercise.exercise_name}
        onChange={(e) => onChange({ ...exercise, exercise_name: e.target.value })}
        placeholder="e.g., Squat, Bench Press"
      />

      <div className="grid grid-cols-4 gap-3">
        <div className="space-y-2">
          <Label>Sets</Label>
          <Input
            type="number"
            value={exercise.sets ?? ''}
            onChange={(e) =>
              onChange({
                ...exercise,
                sets: e.target.value ? Number(e.target.value) : undefined,
              })
            }
            min={1}
            placeholder="3"
          />
        </div>
        <div className="space-y-2">
          <Label>Reps</Label>
          <Input
            type="number"
            value={exercise.reps ?? ''}
            onChange={(e) =>
              onChange({
                ...exercise,
                reps: e.target.value ? Number(e.target.value) : undefined,
              })
            }
            min={1}
            placeholder="10"
          />
        </div>
        <div className="space-y-2">
          <Label>RPE</Label>
          <Input
            type="number"
            value={exercise.rpe ?? ''}
            onChange={(e) =>
              onChange({
                ...exercise,
                rpe: e.target.value ? Number(e.target.value) : undefined,
              })
            }
            min={1}
            max={10}
            placeholder="7"
          />
        </div>
        <div className="space-y-2">
          <Label>%1RM</Label>
          <Input
            type="number"
            value={exercise.percent_1rm_display ?? ''}
            onChange={(e) =>
              onChange({
                ...exercise,
                percent_1rm_display: e.target.value ? Number(e.target.value) : undefined,
              })
            }
            min={0}
            max={100}
            placeholder="75"
          />
        </div>
      </div>
    </div>
  );
}

// ============================================================================
// ProgramForm
// ============================================================================

interface ProgramFormProps {
  programName: string;
  programDescription: string;
  programNotes: string;
  initialEntries?: ProgramEntryCreate[];
  onNameChange: (name: string) => void;
  onDescriptionChange: (description: string) => void;
  onNotesChange: (notes: string) => void;
  onSave: (entries: ProgramEntryCreate[]) => void;
  onDelete?: () => void;
  isSaving?: boolean;
  isEditing?: boolean;
}

export function ProgramForm({
  programName,
  programDescription,
  programNotes,
  initialEntries,
  onNameChange,
  onDescriptionChange,
  onNotesChange,
  onSave,
  onDelete,
  isSaving,
  isEditing,
}: ProgramFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [sessions, setSessions] = useState<SessionGroup[]>(() =>
    initialEntries && initialEntries.length > 0
      ? entriesToSessionGroups(initialEntries)
      : [{ name: '', exercises: [] }]
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
      exercises: [...session.exercises, { exercise_name: '', sets: 3, reps: 10, rpe: 7 }],
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
          <Label htmlFor="program-description">Description</Label>
          <Input
            id="program-description"
            value={programDescription}
            onChange={(e) => onDescriptionChange(e.target.value)}
            placeholder="Brief description of the program"
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
        <p className="text-sm font-semibold">Sessions</p>
        <Accordion type="multiple" className="w-full">
          {sessions.map((session, sessionIdx) => (
            <AccordionItem key={sessionIdx} value={`session-${sessionIdx}`}>
              <AccordionTrigger className="hover:no-underline">
                <div className="flex items-center justify-between w-full pr-4">
                  <span
                    className={
                      session.name ? 'font-semibold' : 'font-semibold text-muted-foreground'
                    }
                  >
                    {session.name || 'Untitled'}
                  </span>
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
                </div>
              </AccordionTrigger>
              <AccordionContent>
                <div className="space-y-4 pt-4">
                  <div className="space-y-2">
                    <Label>Session Name</Label>
                    <Input
                      value={session.name}
                      onChange={(e) => handleSessionNameChange(sessionIdx, e.target.value)}
                      placeholder="e.g., Block1 Week2 Day3, Week1 Day2"
                    />
                  </div>

                  <div className="space-y-3">
                    {session.exercises.map((exercise, exIdx) => (
                      <ProgramExerciseRow
                        key={exIdx}
                        exercise={exercise}
                        onChange={(updated) => handleExerciseChange(sessionIdx, exIdx, updated)}
                        onRemove={() => handleRemoveExercise(sessionIdx, exIdx)}
                      />
                    ))}
                  </div>

                  <Button
                    variant="outline"
                    onClick={() => handleAddExercise(sessionIdx)}
                    className="w-full"
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    Add Exercise
                  </Button>
                </div>
              </AccordionContent>
            </AccordionItem>
          ))}
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
