import { DeleteConfirmDialog } from '@/components/plan-editor/DeleteConfirmDialog';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  useArchiveProgramTemplate,
  useDeleteProgramTemplate,
  useDuplicateProgramTemplate,
} from '@/lib/hooks/useProgramTemplates';
import { useProgramsByTemplateId } from '@/lib/hooks/usePrograms';
import type { ProgramTemplate } from '@/types/api';
import { Archive, Copy, Eye, Trash2 } from 'lucide-react';
import Link from 'next/link';
import { useState } from 'react';

interface ProgramTemplateCardProps {
  template: ProgramTemplate;
}

export function ProgramTemplateCard({ template }: ProgramTemplateCardProps) {
  const archiveTemplate = useArchiveProgramTemplate();
  const duplicateTemplate = useDuplicateProgramTemplate();
  const deleteTemplate = useDeleteProgramTemplate();
  const isArchived = !!template.archived_at;

  const { data: programsData, isLoading: programsLoading } = useProgramsByTemplateId(template.id);
  const programs = programsData?.items ?? [];
  const hasPrograms = programs.length > 0;

  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  return (
    <>
      <Card className={isArchived ? 'opacity-60' : ''}>
        <CardHeader>
          <div className="flex items-start justify-between">
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <CardTitle className="text-xl">{template.name}</CardTitle>
                {isArchived && <Badge variant="secondary">Archived</Badge>}
              </div>
              {template.description && (
                <p className="text-sm text-muted-foreground">{template.description}</p>
              )}
            </div>
            <div className="flex gap-1">
              <Link href={`/program-templates/${template.id}`}>
                <Button variant="ghost" size="sm">
                  <Eye className="h-4 w-4" />
                </Button>
              </Link>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => duplicateTemplate.mutate(template.id)}
                disabled={duplicateTemplate.isPending}
              >
                <Copy className="h-4 w-4" />
              </Button>
              {!isArchived && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => archiveTemplate.mutate(template.id)}
                  disabled={archiveTemplate.isPending}
                >
                  <Archive className="h-4 w-4" />
                </Button>
              )}
              {!isArchived && !hasPrograms && !programsLoading && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setDeleteDialogOpen(true)}
                  disabled={deleteTemplate.isPending}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              )}
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {hasPrograms && (
              <div className="text-sm">
                <span className="text-muted-foreground">Programs: </span>
                {programs.map((p) => p.name).join(', ')}
              </div>
            )}
            <div className="text-sm text-muted-foreground">
              Created {new Date(template.created_at).toLocaleDateString()}
            </div>
          </div>
        </CardContent>
      </Card>
      <DeleteConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        onConfirm={() => {
          deleteTemplate.mutate(template.id);
          setDeleteDialogOpen(false);
        }}
        title="Delete template?"
        description="This will permanently delete this template. This action cannot be undone."
      />
    </>
  );
}
