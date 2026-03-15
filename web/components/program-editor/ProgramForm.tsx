'use client';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { Exercise, ProgramEntryCreate } from '@/types/api';
import { Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { DeleteConfirmDialog } from './DeleteConfirmDialog';
import { WeekAccordion } from './WeekAccordion';
import { type DayGroup, type WeekGroup, entriesToWeekGroups, weekGroupsToEntries } from './types';

interface ProgramFormProps {
  programName: string;
  programDescription: string;
  initialEntries?: ProgramEntryCreate[];
  availableExercises: Exercise[];
  onNameChange: (name: string) => void;
  onDescriptionChange: (description: string) => void;
  onSave: (entries: ProgramEntryCreate[]) => void;
  onDelete?: () => void;
  isSaving?: boolean;
  isEditing?: boolean;
}

/**
 * Full hierarchical program editor form.
 * Manages week/day/exercise tree internally and converts to flat entries on save.
 */
export function ProgramForm({
  programName,
  programDescription,
  initialEntries,
  availableExercises,
  onNameChange,
  onDescriptionChange,
  onSave,
  onDelete,
  isSaving,
  isEditing,
}: ProgramFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [weekGroups, setWeekGroups] = useState<WeekGroup[]>(() =>
    initialEntries && initialEntries.length > 0
      ? entriesToWeekGroups(initialEntries)
      : [{ name: 'Week 1', days: [] }]
  );

  const handleAddWeek = () => {
    const newWeek: WeekGroup = {
      name: `Week ${weekGroups.length + 1}`,
      days: [],
    };
    setWeekGroups([...weekGroups, newWeek]);
  };

  const handleRemoveWeek = (weekIdx: number) => {
    setWeekGroups(weekGroups.filter((_, idx) => idx !== weekIdx));
  };

  const handleWeekNameChange = (weekIdx: number, name: string) => {
    const updated = [...weekGroups];
    updated[weekIdx] = { ...updated[weekIdx], name };
    setWeekGroups(updated);
  };

  const handleAddDay = (weekIdx: number) => {
    const updated = [...weekGroups];
    const week = updated[weekIdx];
    if (!week) return;

    const newDay: DayGroup = {
      name: `Day ${week.days.length + 1}`,
      exercises: [],
    };
    updated[weekIdx] = { ...week, days: [...week.days, newDay] };
    setWeekGroups(updated);
  };

  const handleRemoveDay = (weekIdx: number, dayIdx: number) => {
    const updated = [...weekGroups];
    const week = updated[weekIdx];
    if (!week) return;

    updated[weekIdx] = { ...week, days: week.days.filter((_, idx) => idx !== dayIdx) };
    setWeekGroups(updated);
  };

  const handleDayNameChange = (weekIdx: number, dayIdx: number, name: string) => {
    const updated = [...weekGroups];
    const week = updated[weekIdx];
    if (!week) return;

    const days = [...week.days];
    const day = days[dayIdx];
    if (!day) return;

    days[dayIdx] = { ...day, name };
    updated[weekIdx] = { ...week, days };
    setWeekGroups(updated);
  };

  const handleAddExercise = (weekIdx: number, dayIdx: number) => {
    const updated = [...weekGroups];
    const week = updated[weekIdx];
    if (!week) return;

    const days = [...week.days];
    const day = days[dayIdx];
    if (!day) return;

    const newExercise: ProgramEntryCreate = {
      name: 'Exercise',
      order: day.exercises.length,
      target_sets: 3,
      target_reps: 10,
      target_rpe: 7,
    };
    days[dayIdx] = { ...day, exercises: [...day.exercises, newExercise] };
    updated[weekIdx] = { ...week, days };
    setWeekGroups(updated);
  };

  const handleRemoveExercise = (weekIdx: number, dayIdx: number, exIdx: number) => {
    const updated = [...weekGroups];
    const week = updated[weekIdx];
    if (!week) return;

    const days = [...week.days];
    const day = days[dayIdx];
    if (!day) return;

    days[dayIdx] = {
      ...day,
      exercises: day.exercises
        .filter((_, idx) => idx !== exIdx)
        .map((e, idx) => ({ ...e, order: idx })),
    };
    updated[weekIdx] = { ...week, days };
    setWeekGroups(updated);
  };

  const handleExerciseChange = (
    weekIdx: number,
    dayIdx: number,
    exIdx: number,
    exerciseId: string
  ) => {
    const updated = [...weekGroups];
    const week = updated[weekIdx];
    if (!week) return;

    const days = [...week.days];
    const day = days[dayIdx];
    if (!day) return;

    const exercises = [...day.exercises];
    const ex = exercises[exIdx];
    if (!ex) return;

    const exercise = availableExercises.find((e) => e.id === exerciseId);
    exercises[exIdx] = { ...ex, exercise_id: exerciseId, name: exercise?.name || 'Exercise' };
    days[dayIdx] = { ...day, exercises };
    updated[weekIdx] = { ...week, days };
    setWeekGroups(updated);
  };

  const handleSetsChange = (weekIdx: number, dayIdx: number, exIdx: number, value: number) => {
    const updated = [...weekGroups];
    const week = updated[weekIdx];
    if (!week) return;

    const days = [...week.days];
    const day = days[dayIdx];
    if (!day) return;

    const exercises = [...day.exercises];
    const ex = exercises[exIdx];
    if (!ex) return;

    exercises[exIdx] = { ...ex, target_sets: value };
    days[dayIdx] = { ...day, exercises };
    updated[weekIdx] = { ...week, days };
    setWeekGroups(updated);
  };

  const handleRepsChange = (weekIdx: number, dayIdx: number, exIdx: number, value: number) => {
    const updated = [...weekGroups];
    const week = updated[weekIdx];
    if (!week) return;

    const days = [...week.days];
    const day = days[dayIdx];
    if (!day) return;

    const exercises = [...day.exercises];
    const ex = exercises[exIdx];
    if (!ex) return;

    exercises[exIdx] = { ...ex, target_reps: value };
    days[dayIdx] = { ...day, exercises };
    updated[weekIdx] = { ...week, days };
    setWeekGroups(updated);
  };

  const handleRpeChange = (weekIdx: number, dayIdx: number, exIdx: number, value: number) => {
    const updated = [...weekGroups];
    const week = updated[weekIdx];
    if (!week) return;

    const days = [...week.days];
    const day = days[dayIdx];
    if (!day) return;

    const exercises = [...day.exercises];
    const ex = exercises[exIdx];
    if (!ex) return;

    exercises[exIdx] = { ...ex, target_rpe: value };
    days[dayIdx] = { ...day, exercises };
    updated[weekIdx] = { ...week, days };
    setWeekGroups(updated);
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Program Details</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="program-name">Program Name</Label>
            <Input
              id="program-name"
              value={programName}
              onChange={(e) => onNameChange(e.target.value)}
              placeholder="e.g., 5/3/1, Starting Strength"
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
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Training Weeks</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <WeekAccordion
            weeks={weekGroups}
            availableExercises={availableExercises}
            onWeekNameChange={handleWeekNameChange}
            onDayNameChange={handleDayNameChange}
            onExerciseChange={handleExerciseChange}
            onSetsChange={handleSetsChange}
            onRepsChange={handleRepsChange}
            onRpeChange={handleRpeChange}
            onRemoveExercise={handleRemoveExercise}
            onAddExercise={handleAddExercise}
            onRemoveDay={handleRemoveDay}
            onAddDay={handleAddDay}
            onRemoveWeek={handleRemoveWeek}
          />
          <Button variant="outline" onClick={handleAddWeek} className="w-full">
            <Plus className="h-4 w-4 mr-2" />
            Add Week
          </Button>
        </CardContent>
      </Card>

      <div className="flex justify-between">
        {isEditing && onDelete && (
          <Button variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
            <Trash2 className="h-4 w-4 mr-2" />
            Delete Program
          </Button>
        )}
        <Button
          onClick={() => onSave(weekGroupsToEntries(weekGroups))}
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
        description="This will permanently delete this training program. This action cannot be undone."
      />
    </div>
  );
}
