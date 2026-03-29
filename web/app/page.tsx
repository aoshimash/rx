'use client';

import { AddSessionDialog } from '@/components/plan/AddSessionDialog';
import { EmptyState } from '@/components/plan/EmptyState';
import { ProgramSidebar } from '@/components/plan/ProgramSidebar';
import { SessionCard } from '@/components/plan/SessionCard';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { useAddPlanSessions, useDeletePlanSession, usePlan } from '@/lib/hooks/usePlans';
import { usePrograms } from '@/lib/hooks/usePrograms';
import type { PlanSession, PlanSessionCreate } from '@/types/api';
import { HTTPError } from 'ky';
import { Plus } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function PlanPage() {
  const router = useRouter();
  const { data: plan, isLoading, error } = usePlan();
  const { data: programsData } = usePrograms();
  const addSessions = useAddPlanSessions();
  const deleteSession = useDeletePlanSession();
  const [addDialogOpen, setAddDialogOpen] = useState(false);

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

  const handleAddSession = async (session: PlanSessionCreate) => {
    await addSessions.mutateAsync([session]);
  };

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

  const nextOrder = sessions.length > 0 ? Math.max(...sessions.map((s) => s.order)) + 1 : 0;

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
          <EmptyState onAddSession={() => setAddDialogOpen(true)} />
        ) : (
          <div className="space-y-3">
            {sessions.map((session) => (
              <SessionCard
                key={session.id}
                session={session}
                programName={
                  session.source_program_id ? programMap.get(session.source_program_id) : undefined
                }
                onLog={handleLog}
                onDelete={handleDelete}
              />
            ))}
          </div>
        )}

        {sessions.length > 0 && (
          <Button
            variant="outline"
            className="w-full mt-3 border-dashed"
            onClick={() => setAddDialogOpen(true)}
          >
            <Plus className="h-4 w-4 mr-2" />
            Add Session
          </Button>
        )}

        <AddSessionDialog
          open={addDialogOpen}
          onOpenChange={setAddDialogOpen}
          onAdd={handleAddSession}
          nextOrder={nextOrder}
        />
      </div>

      {/* Programs Sidebar (30%) */}
      <div className="flex-[3] p-4 border-l">
        <ProgramSidebar />
      </div>
    </main>
  );
}
