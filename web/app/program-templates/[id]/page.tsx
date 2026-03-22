'use client';

import { DeleteConfirmDialog } from '@/components/plan-editor/DeleteConfirmDialog';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { Skeleton } from '@/components/ui/skeleton';
import {
  useArchiveProgramTemplate,
  useDeleteProgramTemplate,
  useDuplicateProgramTemplate,
  useProgramTemplate,
  useUnarchiveProgramTemplate,
} from '@/lib/hooks/useProgramTemplates';
import { useProgramsByTemplateId } from '@/lib/hooks/usePrograms';
import type { ProgramTemplateEntry } from '@/types/api';
import { Archive, ArchiveRestore, ArrowLeft, Copy, Info, Trash2, User } from 'lucide-react';
import { useParams, useRouter } from 'next/navigation';
import { useState } from 'react';

type ExerciseGroup = { name: string; entries: ProgramTemplateEntry[] };
type SessionGroup = { sessionName: string; exerciseGroups: ExerciseGroup[] };

function getSessionName(entry: ProgramTemplateEntry): string {
  const session = entry.metadata?.session;
  if (typeof session === 'string') return session;
  return '';
}

function groupBySession(entries: ProgramTemplateEntry[]): SessionGroup[] {
  const sorted = [...entries].sort((a, b) => a.order - b.order);
  const sessionOrder: string[] = [];
  const sessionMap = new Map<string, ProgramTemplateEntry[]>();

  for (const entry of sorted) {
    const session = getSessionName(entry);
    if (!sessionMap.has(session)) {
      sessionOrder.push(session);
      sessionMap.set(session, []);
    }
    sessionMap.get(session)?.push(entry);
  }

  return sessionOrder.map((sessionName) => ({
    sessionName,
    exerciseGroups: groupByExercise(sessionMap.get(sessionName) ?? []),
  }));
}

function groupByExercise(entries: ProgramTemplateEntry[]): ExerciseGroup[] {
  const groups: ExerciseGroup[] = [];
  const map = new Map<string, ExerciseGroup>();
  for (const entry of entries) {
    if (!map.has(entry.exercise_name)) {
      const g: ExerciseGroup = { name: entry.exercise_name, entries: [] };
      groups.push(g);
      map.set(entry.exercise_name, g);
    }
    map.get(entry.exercise_name)?.entries.push(entry);
  }
  return groups;
}

export default function ProgramTemplateDetailPage() {
  const params = useParams();
  const router = useRouter();
  const templateId = params.id as string;
  const { data: template, isLoading } = useProgramTemplate(templateId);
  const archiveTemplate = useArchiveProgramTemplate();
  const unarchiveTemplate = useUnarchiveProgramTemplate();
  const duplicateTemplate = useDuplicateProgramTemplate();
  const deleteTemplate = useDeleteProgramTemplate();
  const { data: programsData, isLoading: programsLoading } = useProgramsByTemplateId(templateId);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  if (isLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <Skeleton className="h-[400px] w-full" />
      </main>
    );
  }

  if (!template) {
    return (
      <main className="container mx-auto p-6">
        <p className="text-muted-foreground">Template not found.</p>
      </main>
    );
  }

  const sessionGroups = groupBySession(template.entries || []);
  const isArchived = !!template.archived_at;
  const hasPrograms = (programsData?.data ?? []).length > 0;

  const handleDuplicate = async () => {
    await duplicateTemplate.mutateAsync(templateId);
    router.push('/program-templates');
  };

  const handleArchiveToggle = async () => {
    if (isArchived) {
      await unarchiveTemplate.mutateAsync(templateId);
    } else {
      await archiveTemplate.mutateAsync(templateId);
      router.push('/program-templates');
    }
  };

  const handleDelete = async () => {
    await deleteTemplate.mutateAsync(templateId);
    router.push('/program-templates');
  };

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6">
        <button
          type="button"
          onClick={() => router.back()}
          className="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-4"
        >
          <ArrowLeft className="h-4 w-4" />
          Back
        </button>
        <div className="flex items-start justify-between">
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-3xl font-bold">{template.name}</h1>
              {isArchived && <Badge variant="secondary">Archived</Badge>}
            </div>
            {template.description && (
              <p className="text-muted-foreground mt-1">{template.description}</p>
            )}
            {template.notes && (
              <p className="text-sm text-muted-foreground mt-1">{template.notes}</p>
            )}
            <div className="flex items-center gap-4 mt-2">
              {template.created_by && (
                <span className="flex items-center gap-1 text-sm text-muted-foreground">
                  <User className="h-3 w-3" />
                  {template.created_by}
                </span>
              )}
              <span className="text-sm text-muted-foreground">
                Created {new Date(template.created_at).toLocaleDateString()}
              </span>
            </div>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={handleDuplicate}
              disabled={duplicateTemplate.isPending}
            >
              <Copy className="h-4 w-4 mr-2" />
              Duplicate
            </Button>
            <Button
              variant="outline"
              onClick={handleArchiveToggle}
              disabled={archiveTemplate.isPending || unarchiveTemplate.isPending}
            >
              {isArchived ? (
                <>
                  <ArchiveRestore className="h-4 w-4 mr-2" />
                  Unarchive
                </>
              ) : (
                <>
                  <Archive className="h-4 w-4 mr-2" />
                  Archive
                </>
              )}
            </Button>
            <div className="flex items-center gap-1">
              <Button
                variant="outline"
                onClick={() => setDeleteDialogOpen(true)}
                disabled={deleteTemplate.isPending || hasPrograms || programsLoading}
                className="text-destructive hover:text-destructive"
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Delete
              </Button>
              {hasPrograms && (
                <Popover>
                  <PopoverTrigger asChild>
                    <button type="button" className="text-muted-foreground hover:text-foreground">
                      <Info className="h-4 w-4" />
                    </button>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto px-3 py-2 text-xs">
                    This template is used by programs and cannot be deleted.
                  </PopoverContent>
                </Popover>
              )}
            </div>
          </div>
        </div>
      </div>

      {sessionGroups.length === 0 ? (
        <p className="text-muted-foreground">No entries in this template.</p>
      ) : (
        <div className="space-y-4">
          {sessionGroups.map((session) => (
            <Card key={session.sessionName || '__default__'}>
              {session.sessionName && (
                <CardHeader className="pb-2">
                  <CardTitle className="text-base">{session.sessionName}</CardTitle>
                </CardHeader>
              )}
              <CardContent className={session.sessionName ? '' : 'pt-4'}>
                <div className="divide-y">
                  {session.exerciseGroups.map((group) => {
                    const hasLabel = group.entries.some((e) => e.metadata?.label);
                    const hasPercent1rm = group.entries.some((e) => e.percent_1rm != null);
                    return (
                      <div
                        key={`${session.sessionName}-${group.name}`}
                        className="py-2 first:pt-0 last:pb-0"
                      >
                        <p className="font-medium text-sm mb-1.5">{group.name}</p>
                        <table className="w-full text-sm">
                          <thead>
                            <tr className="text-xs text-muted-foreground">
                              {hasLabel && <th className="text-left font-normal pb-1 w-16" />}
                              <th className="text-right font-normal pb-1 pr-4">RPE</th>
                              <th className="text-right font-normal pb-1 pr-4">Reps</th>
                              <th className="text-right font-normal pb-1 pr-4">Sets</th>
                              {hasPercent1rm && (
                                <th className="text-right font-normal pb-1">%1RM</th>
                              )}
                            </tr>
                          </thead>
                          <tbody>
                            {group.entries.map((entry) => {
                              const label = entry.metadata?.label as string | undefined;
                              return (
                                <tr key={entry.id} className="text-muted-foreground">
                                  {hasLabel && (
                                    <td className="text-xs pr-3 py-0.5">{label ?? ''}</td>
                                  )}
                                  <td className="text-right tabular-nums pr-4 py-0.5">
                                    {entry.rpe ?? '—'}
                                  </td>
                                  <td className="text-right tabular-nums pr-4 py-0.5">
                                    {entry.reps ?? '—'}
                                  </td>
                                  <td className="text-right tabular-nums pr-4 py-0.5">
                                    {entry.sets ?? '—'}
                                  </td>
                                  {hasPercent1rm && (
                                    <td className="text-right tabular-nums py-0.5">
                                      {entry.percent_1rm != null
                                        ? `${Math.round(entry.percent_1rm * 100)}%`
                                        : '—'}
                                    </td>
                                  )}
                                </tr>
                              );
                            })}
                          </tbody>
                        </table>
                      </div>
                    );
                  })}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <DeleteConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        onConfirm={handleDelete}
        title="Delete template?"
        description="This will permanently delete this template. This action cannot be undone."
      />
    </main>
  );
}
