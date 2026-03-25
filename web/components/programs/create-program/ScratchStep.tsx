'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { ProgramSessionCreate } from '@/types/api';
import { Plus, X } from 'lucide-react';
import { useState } from 'react';

interface SessionInput {
  session_name: string;
  order: number;
}

interface ScratchStepProps {
  onBack: () => void;
  onSubmit: (data: {
    name: string;
    notes?: string;
    sessions: ProgramSessionCreate[];
  }) => void;
  isPending: boolean;
  nameError?: string;
}

export function ScratchStep({ onBack, onSubmit, isPending, nameError }: ScratchStepProps) {
  const [name, setName] = useState('');
  const [notes, setNotes] = useState('');
  const [sessions, setSessions] = useState<SessionInput[]>([{ session_name: '', order: 0 }]);

  const handleAddSession = () => {
    setSessions([...sessions, { session_name: '', order: sessions.length }]);
  };

  const handleRemoveSession = (idx: number) => {
    setSessions(sessions.filter((_, i) => i !== idx).map((s, i) => ({ ...s, order: i })));
  };

  const handleSessionNameChange = (idx: number, value: string) => {
    const updated = [...sessions];
    const s = updated[idx];
    if (s) updated[idx] = { ...s, session_name: value };
    setSessions(updated);
  };

  const handleSubmit = () => {
    const programSessions: ProgramSessionCreate[] = sessions
      .filter((s) => s.session_name.trim())
      .map((s) => ({ session_name: s.session_name.trim(), order: s.order }));
    onSubmit({
      name,
      notes: notes || undefined,
      sessions: programSessions,
    });
  };

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="scratch-name">Program Name</Label>
        <Input
          id="scratch-name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="e.g., SBD Block 1"
        />
        {nameError && <p className="text-sm text-destructive">{nameError}</p>}
      </div>
      <div className="space-y-2">
        <Label htmlFor="scratch-notes">Notes</Label>
        <Input
          id="scratch-notes"
          value={notes}
          onChange={(e) => setNotes(e.target.value)}
          placeholder="Optional notes"
        />
      </div>
      <div className="space-y-2">
        <Label>Sessions</Label>
        <div className="space-y-2">
          {sessions.map((session, idx) => (
            <div key={idx} className="flex gap-2">
              <Input
                value={session.session_name}
                onChange={(e) => handleSessionNameChange(idx, e.target.value)}
                placeholder={`e.g., Week 1 Day ${idx + 1}`}
              />
              {sessions.length > 1 && (
                <Button variant="ghost" size="sm" onClick={() => handleRemoveSession(idx)}>
                  <X className="h-4 w-4" />
                </Button>
              )}
            </div>
          ))}
        </div>
        <Button variant="outline" size="sm" onClick={handleAddSession} className="w-full">
          <Plus className="h-4 w-4 mr-2" />
          Add Session
        </Button>
      </div>
      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onBack}>
          Back
        </Button>
        <Button onClick={handleSubmit} disabled={isPending || !name.trim()}>
          {isPending ? 'Creating...' : 'Create Program'}
        </Button>
      </div>
    </div>
  );
}
