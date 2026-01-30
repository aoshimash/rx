'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Plus, Trash2 } from 'lucide-react';
import { WeekAccordion } from './WeekAccordion';
import { DeleteConfirmDialog } from './DeleteConfirmDialog';
import type { Exercise, ProgramNodeCreate } from '@/types/api';

interface ProgramFormProps {
  programName: string;
  programDescription: string;
  weeks: ProgramNodeCreate[];
  availableExercises: Exercise[];
  onNameChange: (name: string) => void;
  onDescriptionChange: (description: string) => void;
  onWeeksChange: (weeks: ProgramNodeCreate[]) => void;
  onSave: () => void;
  onDelete?: () => void;
  isSaving?: boolean;
  isEditing?: boolean;
}

/**
 * Full hierarchical program editor form
 */
export function ProgramForm({
  programName,
  programDescription,
  weeks,
  availableExercises,
  onNameChange,
  onDescriptionChange,
  onWeeksChange,
  onSave,
  onDelete,
  isSaving,
  isEditing,
}: ProgramFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const handleAddWeek = () => {
    const newWeek: ProgramNodeCreate = {
      name: `Week ${weeks.length + 1}`,
      node_type: 'week',
      order: weeks.length,
      children: [],
    };
    onWeeksChange([...weeks, newWeek]);
  };

  const handleRemoveWeek = (weekIdx: number) => {
    const updated = weeks.filter((_, idx) => idx !== weekIdx);
    onWeeksChange(updated.map((w, idx) => ({ ...w, order: idx })));
  };

  const handleWeekNameChange = (weekIdx: number, name: string) => {
    const updated = [...weeks];
    updated[weekIdx] = { ...updated[weekIdx], name };
    onWeeksChange(updated);
  };

  const handleAddDay = (weekIdx: number) => {
    const updated = [...weeks];
    const currentDays = updated[weekIdx].children || [];
    const newDay: ProgramNodeCreate = {
      name: `Day ${currentDays.length + 1}`,
      node_type: 'day',
      order: currentDays.length,
      children: [],
    };
    updated[weekIdx] = {
      ...updated[weekIdx],
      children: [...currentDays, newDay],
    };
    onWeeksChange(updated);
  };

  const handleRemoveDay = (weekIdx: number, dayIdx: number) => {
    const updated = [...weeks];
    const days = updated[weekIdx].children || [];
    updated[weekIdx] = {
      ...updated[weekIdx],
      children: days.filter((_, idx) => idx !== dayIdx).map((d, idx) => ({ ...d, order: idx })),
    };
    onWeeksChange(updated);
  };

  const handleDayNameChange = (weekIdx: number, dayIdx: number, name: string) => {
    const updated = [...weeks];
    const days = [...(updated[weekIdx].children || [])];
    days[dayIdx] = { ...days[dayIdx], name };
    updated[weekIdx] = { ...updated[weekIdx], children: days };
    onWeeksChange(updated);
  };

  const handleAddExercise = (weekIdx: number, dayIdx: number) => {
    const updated = [...weeks];
    const days = [...(updated[weekIdx].children || [])];
    const exercises = days[dayIdx].children || [];
    const newExercise: ProgramNodeCreate = {
      name: 'Exercise',
      node_type: 'exercise',
      order: exercises.length,
      target_sets: 3,
      target_reps: 10,
      target_rpe: 7,
    };
    days[dayIdx] = {
      ...days[dayIdx],
      children: [...exercises, newExercise],
    };
    updated[weekIdx] = { ...updated[weekIdx], children: days };
    onWeeksChange(updated);
  };

  const handleRemoveExercise = (weekIdx: number, dayIdx: number, exIdx: number) => {
    const updated = [...weeks];
    const days = [...(updated[weekIdx].children || [])];
    const exercises = days[dayIdx].children || [];
    days[dayIdx] = {
      ...days[dayIdx],
      children: exercises.filter((_, idx) => idx !== exIdx).map((e, idx) => ({ ...e, order: idx })),
    };
    updated[weekIdx] = { ...updated[weekIdx], children: days };
    onWeeksChange(updated);
  };

  const handleExerciseChange = (
    weekIdx: number,
    dayIdx: number,
    exIdx: number,
    exerciseId: string
  ) => {
    const updated = [...weeks];
    const days = [...(updated[weekIdx].children || [])];
    const exercises = [...(days[dayIdx].children || [])];
    const exercise = availableExercises.find((e) => e.id === exerciseId);
    exercises[exIdx] = {
      ...exercises[exIdx],
      exercise_id: exerciseId,
      name: exercise?.name || 'Exercise',
    };
    days[dayIdx] = { ...days[dayIdx], children: exercises };
    updated[weekIdx] = { ...updated[weekIdx], children: days };
    onWeeksChange(updated);
  };

  const handleSetsChange = (weekIdx: number, dayIdx: number, exIdx: number, value: number) => {
    const updated = [...weeks];
    const days = [...(updated[weekIdx].children || [])];
    const exercises = [...(days[dayIdx].children || [])];
    exercises[exIdx] = { ...exercises[exIdx], target_sets: value };
    days[dayIdx] = { ...days[dayIdx], children: exercises };
    updated[weekIdx] = { ...updated[weekIdx], children: days };
    onWeeksChange(updated);
  };

  const handleRepsChange = (weekIdx: number, dayIdx: number, exIdx: number, value: number) => {
    const updated = [...weeks];
    const days = [...(updated[weekIdx].children || [])];
    const exercises = [...(days[dayIdx].children || [])];
    exercises[exIdx] = { ...exercises[exIdx], target_reps: value };
    days[dayIdx] = { ...days[dayIdx], children: exercises };
    updated[weekIdx] = { ...updated[weekIdx], children: days };
    onWeeksChange(updated);
  };

  const handleRpeChange = (weekIdx: number, dayIdx: number, exIdx: number, value: number) => {
    const updated = [...weeks];
    const days = [...(updated[weekIdx].children || [])];
    const exercises = [...(days[dayIdx].children || [])];
    exercises[exIdx] = { ...exercises[exIdx], target_rpe: value };
    days[dayIdx] = { ...days[dayIdx], children: exercises };
    updated[weekIdx] = { ...updated[weekIdx], children: days };
    onWeeksChange(updated);
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
            weeks={weeks}
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
          <Button
            variant="destructive"
            onClick={() => setDeleteDialogOpen(true)}
          >
            <Trash2 className="h-4 w-4 mr-2" />
            Delete Program
          </Button>
        )}
        <Button
          onClick={onSave}
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
