'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { PlanEntryCreate } from '@/types/api';
import { Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { DeleteConfirmDialog } from './DeleteConfirmDialog';
import { SessionAccordion } from './SessionAccordion';
import { type SessionGroup, entriesToSessionGroups, sessionGroupsToEntries } from './types';

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
  const [sessions, setSessions] = useState<SessionGroup[]>(() =>
    initialEntries && initialEntries.length > 0
      ? entriesToSessionGroups(initialEntries)
      : [{ name: 'Day 1', exercises: [] }]
  );

  const handleAddSession = () => {
    const newSession: SessionGroup = {
      name: `Day ${sessions.length + 1}`,
      exercises: [],
    };
    setSessions([...sessions, newSession]);
  };

  const handleSessionChange = (index: number, updated: SessionGroup) => {
    const newSessions = [...sessions];
    newSessions[index] = updated;
    setSessions(newSessions);
  };

  const handleSessionDelete = (index: number) => {
    setSessions(sessions.filter((_, idx) => idx !== index));
  };

  return (
    <div className="space-y-6">
      <div className="space-y-4">
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
      </div>

      <div className="space-y-4">
        <p className="text-sm font-semibold">Training Sessions</p>
        <SessionAccordion
          sessions={sessions}
          onChange={handleSessionChange}
          onDelete={handleSessionDelete}
        />
        <Button variant="outline" onClick={handleAddSession} className="w-full">
          <Plus className="h-4 w-4 mr-2" />
          Add Session
        </Button>
      </div>

      <div className="flex justify-between">
        {isEditing && onDelete && (
          <Button variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
            <Trash2 className="h-4 w-4 mr-2" />
            Delete Plan
          </Button>
        )}
        <Button
          onClick={() => onSave(sessionGroupsToEntries(sessions))}
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
