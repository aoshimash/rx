'use client';

import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import type { PlanSessionCreate, PlanSessionEntryCreate, ProgramSession } from '@/types/api';
import { useState } from 'react';

interface SessionSelectorProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  sessions: ProgramSession[];
  programId: string;
  onConfirm: (sessions: PlanSessionCreate[]) => void;
  isPending: boolean;
}

export function SessionSelector({
  open,
  onOpenChange,
  sessions,
  programId,
  onConfirm,
  isPending,
}: SessionSelectorProps) {
  const sortedSessions = [...sessions].sort((a, b) => a.order - b.order);
  const [selected, setSelected] = useState<Set<string>>(
    () => new Set(sortedSessions.map((s) => s.id))
  );

  const toggleSession = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handleConfirm = () => {
    const planSessions: PlanSessionCreate[] = sortedSessions
      .filter((s) => selected.has(s.id))
      .map((s, i) => ({
        session_name: s.session_name,
        order: i,
        date: s.date || undefined,
        source_program_id: programId,
        source_session_id: s.id,
        entries: s.entries.map(
          (e, j): PlanSessionEntryCreate => ({
            exercise_name: e.exercise_name,
            order: j,
            fields: e.fields || undefined,
            notes: e.notes || undefined,
          })
        ),
      }));

    onConfirm(planSessions);
  };

  const handleOpenChange = (isOpen: boolean) => {
    if (isOpen) {
      setSelected(new Set(sortedSessions.map((s) => s.id)));
    }
    onOpenChange(isOpen);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Add to Plan</DialogTitle>
        </DialogHeader>

        <div className="border rounded-md divide-y">
          {sortedSessions.map((session) => (
            <label
              key={session.id}
              className="flex items-center gap-3 px-3 py-2.5 cursor-pointer hover:bg-muted/50"
            >
              <Checkbox
                checked={selected.has(session.id)}
                onCheckedChange={() => toggleSession(session.id)}
              />
              <span className="font-medium text-sm flex-1">{session.session_name}</span>
              <span className="text-xs text-muted-foreground">
                {session.entries.length} exercise{session.entries.length !== 1 ? 's' : ''}
              </span>
            </label>
          ))}
        </div>

        <p className="text-xs text-muted-foreground">
          {selected.size} of {sortedSessions.length} sessions selected
        </p>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleConfirm} disabled={selected.size === 0 || isPending}>
            {isPending ? 'Adding...' : 'Add to Plan'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
