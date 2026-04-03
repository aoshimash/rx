'use client';

import { EmptyState } from '@/components/plan/EmptyState';
import { ProgramSidebar } from '@/components/plan/ProgramSidebar';
import { SessionCard } from '@/components/plan/SessionCard';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Skeleton } from '@/components/ui/skeleton';
import {
  useAddPlanSessions,
  useDeletePlanSession,
  usePlan,
  useUpdatePlanSession,
} from '@/lib/hooks/usePlans';
import { usePrograms } from '@/lib/hooks/usePrograms';
import type { PlanSession, PlanSessionUpdate } from '@/types/api';
import { HTTPError } from 'ky';
import { Plus } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useCallback, useRef, useState } from 'react';

export default function PlanPage() {
  const router = useRouter();
  const { data: plan, isLoading, error } = usePlan();
  const { data: programsData } = usePrograms();
  const addSessions = useAddPlanSessions();
  const deleteSession = useDeletePlanSession();
  const updateSession = useUpdatePlanSession();
  const [isAdding, setIsAdding] = useState(false);
  const [newSessionName, setNewSessionName] = useState('');
  const [newSessionId, setNewSessionId] = useState<string | null>(null);
  const nameInputRef = useRef<HTMLInputElement>(null);

  const noPlan = error instanceof HTTPError && error.response.status === 404;
  const sessions = plan?.sessions?.slice().sort((a, b) => a.order - b.order) ?? [];
  const programMap = new Map((programsData?.data ?? []).map((p) => [p.id, p.name]));

  const handleLog = (session: PlanSession) => {
    const params = new URLSearchParams();
    params.set('planSessionId', session.id);
    if (session.source_program_id) params.set('programId', session.source_program_id);
    params.set('session', session.session_name);
    router.push(`/logs/new?${params}`);
  };

  const handleDelete = async (sessionId: string) => {
    await deleteSession.mutateAsync(sessionId);
  };

  const handleUpdate = async (sessionId: string, data: PlanSessionUpdate) => {
    await updateSession.mutateAsync({ sessionId, data });
  };

  const nextOrder = sessions.length > 0 ? Math.max(...sessions.map((s) => s.order)) + 1 : 0;

  const startAdding = () => {
    setIsAdding(true);
    setNewSessionName('');
    setTimeout(() => nameInputRef.current?.focus(), 0);
  };

  const cancelAdding = () => {
    setIsAdding(false);
    setNewSessionName('');
  };

  const handleCreateSession = async () => {
    const name = newSessionName.trim();
    if (!name) return;
    setIsAdding(false);
    setNewSessionName('');
    const result = await addSessions.mutateAsync([{ session_name: name, order: nextOrder }]);
    // Find the newly created session to auto-open edit mode
    const created = result.sessions
      ?.slice()
      .sort((a, b) => b.order - a.order)
      .find((s) => s.session_name === name);
    if (created) {
      setNewSessionId(created.id);
    }
  };

  const handleNewHandled = useCallback(() => {
    setNewSessionId(null);
  }, []);

  if (isLoading) {
    return (
      <main className="flex flex-1">
        <div className="flex-[7] p-6 space-y-3">
          <Skeleton className="h-8 w-[100px]" />
          <Skeleton className="h-[120px] w-full" />
          <Skeleton className="h-[120px] w-full" />
        </div>
        <div className="flex-[3] p-4 border-l space-y-2">
          <Skeleton className="h-6 w-[80px]" />
          <Skeleton className="h-14 w-full" />
          <Skeleton className="h-14 w-full" />
        </div>
      </main>
    );
  }

  return (
    <main className="flex flex-1">
      {/* Sessions (70%) */}
      <div className="flex-[7] p-6">
        <div className="mb-4">
          <h1 className="text-2xl font-bold">Plan</h1>
          {sessions.length > 0 && (
            <p className="text-sm text-muted-foreground">
              {sessions.length} session{sessions.length !== 1 ? 's' : ''} queued
            </p>
          )}
        </div>

        {noPlan || sessions.length === 0 ? (
          isAdding ? null : (
            <EmptyState onAddSession={startAdding} />
          )
        ) : (
          <div className="space-y-3">
            {sessions.map((session) => (
              <SessionCard
                key={session.id}
                session={session}
                programName={
                  session.source_program_id ? programMap.get(session.source_program_id) : undefined
                }
                isNew={session.id === newSessionId}
                onLog={handleLog}
                onDelete={handleDelete}
                onUpdate={handleUpdate}
                onNewHandled={handleNewHandled}
              />
            ))}
          </div>
        )}

        {isAdding ? (
          <div className="mt-3 flex items-center gap-2">
            <Input
              ref={nameInputRef}
              placeholder="Session name (e.g., Upper Body A)"
              value={newSessionName}
              onChange={(e) => setNewSessionName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') handleCreateSession();
                if (e.key === 'Escape') cancelAdding();
              }}
              className="flex-1"
            />
            <Button
              size="sm"
              onClick={handleCreateSession}
              disabled={!newSessionName.trim() || addSessions.isPending}
            >
              Create
            </Button>
            <Button size="sm" variant="ghost" onClick={cancelAdding}>
              Cancel
            </Button>
          </div>
        ) : (
          <Button variant="outline" className="w-full mt-3 border-dashed" onClick={startAdding}>
            <Plus className="h-4 w-4 mr-2" />
            Add Session
          </Button>
        )}
      </div>

      {/* Programs Sidebar (30%) */}
      <div className="flex-[3] p-4 border-l">
        <ProgramSidebar />
      </div>
    </main>
  );
}
