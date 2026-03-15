'use client';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { PlanEntryCreate } from '@/types/api';
import { Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { DeleteConfirmDialog } from './DeleteConfirmDialog';
import { WeekAccordion } from './WeekAccordion';
import { type DayGroup, type WeekGroup, entriesToWeekGroups, weekGroupsToEntries } from './types';

interface PlanFormProps {
  planName: string;
  planDescription: string;
  initialEntries?: PlanEntryCreate[];
  onNameChange: (name: string) => void;
  onDescriptionChange: (description: string) => void;
  onSave: (entries: PlanEntryCreate[]) => void;
  onDelete?: () => void;
  isSaving?: boolean;
  isEditing?: boolean;
}

export function PlanForm({
  planName,
  planDescription,
  initialEntries,
  onNameChange,
  onDescriptionChange,
  onSave,
  onDelete,
  isSaving,
  isEditing,
}: PlanFormProps) {
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
    const week = updated[weekIdx];
    if (!week) return;
    updated[weekIdx] = { ...week, name };
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

    const newExercise: PlanEntryCreate = {
      exercise_name: 'Exercise',
      order: day.exercises.length,
      sets: 3,
      reps: 10,
      rpe: 7,
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

  const handleExerciseNameChange = (
    weekIdx: number,
    dayIdx: number,
    exIdx: number,
    name: string
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

    exercises[exIdx] = { ...ex, exercise_name: name };
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

    exercises[exIdx] = { ...ex, sets: value };
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

    exercises[exIdx] = { ...ex, reps: value };
    days[dayIdx] = { ...day, exercises };
    updated[weekIdx] = { ...week, days };
    setWeekGroups(updated);
  };

  const handleLoadKgChange = (weekIdx: number, dayIdx: number, exIdx: number, value: number) => {
    const updated = [...weekGroups];
    const week = updated[weekIdx];
    if (!week) return;

    const days = [...week.days];
    const day = days[dayIdx];
    if (!day) return;

    const exercises = [...day.exercises];
    const ex = exercises[exIdx];
    if (!ex) return;

    exercises[exIdx] = { ...ex, load_kg: value };
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

    exercises[exIdx] = { ...ex, rpe: value };
    days[dayIdx] = { ...day, exercises };
    updated[weekIdx] = { ...week, days };
    setWeekGroups(updated);
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Plan Details</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="plan-name">Plan Name</Label>
            <Input
              id="plan-name"
              value={planName}
              onChange={(e) => onNameChange(e.target.value)}
              placeholder="e.g., 5/3/1, Starting Strength"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="plan-description">Description</Label>
            <Input
              id="plan-description"
              value={planDescription}
              onChange={(e) => onDescriptionChange(e.target.value)}
              placeholder="Brief description of the plan"
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
            onWeekNameChange={handleWeekNameChange}
            onDayNameChange={handleDayNameChange}
            onExerciseNameChange={handleExerciseNameChange}
            onSetsChange={handleSetsChange}
            onRepsChange={handleRepsChange}
            onLoadKgChange={handleLoadKgChange}
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
            Delete Plan
          </Button>
        )}
        <Button
          onClick={() => onSave(weekGroupsToEntries(weekGroups))}
          disabled={isSaving || !planName}
          className={!isEditing ? 'ml-auto' : ''}
        >
          {isSaving ? 'Saving...' : isEditing ? 'Update Plan' : 'Create Plan'}
        </Button>
      </div>

      <DeleteConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        onConfirm={() => {
          onDelete?.();
          setDeleteDialogOpen(false);
        }}
        title="Delete Plan?"
        description="This will permanently delete this training plan. This action cannot be undone."
      />
    </div>
  );
}
