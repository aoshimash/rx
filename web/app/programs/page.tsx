'use client';

import { ProgramCard } from '@/components/programs/ProgramCard';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Skeleton } from '@/components/ui/skeleton';
import { useCreateProgram, usePrograms } from '@/lib/hooks/usePrograms';
import type { ProgramSessionCreate } from '@/types/api';
import { Plus, X } from 'lucide-react';
import { useState } from 'react';

const SHOW_COMPLETED_KEY = 'programs:showCompleted';

interface SessionInput {
  session_name: string;
  order: number;
}

export default function ProgramsPage() {
  const { data: programsData, isLoading } = usePrograms();
  const createProgram = useCreateProgram();
  const [showCompleted, setShowCompleted] = useState(
    () => typeof window !== 'undefined' && localStorage.getItem(SHOW_COMPLETED_KEY) === 'true'
  );
  const [open, setOpen] = useState(false);
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

  const handleSave = async () => {
    const programSessions: ProgramSessionCreate[] = sessions
      .filter((s) => s.session_name.trim())
      .map((s) => ({ session_name: s.session_name.trim(), order: s.order }));

    await createProgram.mutateAsync({
      name,
      notes: notes || undefined,
      sessions: programSessions,
    });
    setOpen(false);
    setName('');
    setNotes('');
    setSessions([{ session_name: '', order: 0 }]);
  };

  const handleOpenChange = (value: boolean) => {
    setOpen(value);
    if (!value) {
      setName('');
      setNotes('');
      setSessions([{ session_name: '', order: 0 }]);
    }
  };

  const allPrograms = [...(programsData?.data || [])].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  );
  const programs = showCompleted
    ? allPrograms
    : allPrograms.filter((p) => p.status !== 'completed');
  const completedCount = allPrograms.filter((p) => p.status === 'completed').length;

  const toggleShowCompleted = () => {
    const next = !showCompleted;
    setShowCompleted(next);
    localStorage.setItem(SHOW_COMPLETED_KEY, String(next));
  };

  if (isLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <div className="space-y-4">
          <Skeleton className="h-[120px]" />
          <Skeleton className="h-[120px]" />
        </div>
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-start justify-between">
        <div>
          <h1 className="text-3xl font-bold">Programs</h1>
          <p className="text-muted-foreground mt-1">
            Concrete training programs with scheduled sessions.
          </p>
        </div>
        {completedCount > 0 && (
          <Button variant="ghost" size="sm" onClick={toggleShowCompleted}>
            {showCompleted
              ? `Hide completed (${completedCount})`
              : `Show completed (${completedCount})`}
          </Button>
        )}
      </div>

      {programs.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">
            No programs yet. Create your first training program.
          </p>
          <Button onClick={() => setOpen(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Create Program
          </Button>
        </div>
      ) : (
        <>
          <div className="space-y-4">
            {programs.map((program) => (
              <ProgramCard key={program.id} program={program} />
            ))}
          </div>
          <div className="mt-6">
            <Button variant="outline" onClick={() => setOpen(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Create Program
            </Button>
          </div>
        </>
      )}

      <Dialog open={open} onOpenChange={handleOpenChange}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Create Program</DialogTitle>
            <DialogDescription>Define the program name and session names.</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="program-name">Program Name</Label>
              <Input
                id="program-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g., SBD Block 1"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="program-notes">Notes</Label>
              <Input
                id="program-notes"
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
              <Button variant="outline" onClick={() => setOpen(false)}>
                Cancel
              </Button>
              <Button onClick={handleSave} disabled={createProgram.isPending || !name.trim()}>
                {createProgram.isPending ? 'Creating...' : 'Create Program'}
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </main>
  );
}
