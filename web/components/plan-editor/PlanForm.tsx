'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { PlanEntryCreate } from '@/types/api';
import { Trash2 } from 'lucide-react';
import { useState } from 'react';
import { DeleteConfirmDialog } from './DeleteConfirmDialog';
import { ExerciseTable } from './ExerciseTable';

interface PlanFormProps {
  planName: string;
  planDescription: string;
  planDate?: string;
  initialEntries?: PlanEntryCreate[];
  onNameChange: (name: string) => void;
  onDescriptionChange: (description: string) => void;
  onDateChange?: (date: string | undefined) => void;
  onSave: (entries: PlanEntryCreate[]) => void;
  onDelete?: () => void;
  isSaving?: boolean;
  isEditing?: boolean;
}

export function PlanForm({
  planName,
  planDescription,
  planDate,
  initialEntries,
  onNameChange,
  onDescriptionChange,
  onDateChange,
  onSave,
  onDelete,
  isSaving,
  isEditing,
}: PlanFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [exercises, setExercises] = useState<PlanEntryCreate[]>(
    initialEntries && initialEntries.length > 0 ? initialEntries : []
  );

  return (
    <div className="space-y-6">
      <div className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="plan-name">Plan Name</Label>
          <Input
            id="plan-name"
            value={planName}
            onChange={(e) => onNameChange(e.target.value)}
            placeholder="e.g., 5/3/1 - Day A, Starting Strength"
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
        {onDateChange && (
          <div className="space-y-2">
            <Label htmlFor="plan-date">Date (optional)</Label>
            <Input
              id="plan-date"
              type="date"
              value={planDate || ''}
              onChange={(e) => onDateChange(e.target.value || undefined)}
            />
          </div>
        )}
      </div>

      <div className="space-y-2">
        <Label>Exercises</Label>
        <ExerciseTable exercises={exercises} onChange={setExercises} />
      </div>

      <div className="flex justify-between">
        {isEditing && onDelete && (
          <Button variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
            <Trash2 className="h-4 w-4 mr-2" />
            Delete Plan
          </Button>
        )}
        <Button
          onClick={() => onSave(exercises)}
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
