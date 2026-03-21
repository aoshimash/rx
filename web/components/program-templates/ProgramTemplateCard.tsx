import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  useArchiveProgramTemplate,
  useDuplicateProgramTemplate,
} from '@/lib/hooks/useProgramTemplates';
import type { ProgramTemplate } from '@/types/api';
import { Archive, Copy, Eye } from 'lucide-react';
import Link from 'next/link';

interface ProgramTemplateCardProps {
  template: ProgramTemplate;
}

export function ProgramTemplateCard({ template }: ProgramTemplateCardProps) {
  const archiveTemplate = useArchiveProgramTemplate();
  const duplicateTemplate = useDuplicateProgramTemplate();
  const isArchived = !!template.archived_at;

  const entries = template.entries || [];
  const sessionNames = new Set(
    entries
      .map((e) => (e.metadata?.session as string) ?? undefined)
      .filter((s): s is string => s !== undefined)
  );

  return (
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
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center gap-4 text-sm text-muted-foreground">
          {sessionNames.size > 0 ? (
            <span>{sessionNames.size} sessions</span>
          ) : entries.length > 0 ? (
            <span>{entries.length} entries</span>
          ) : (
            <span>No entries</span>
          )}
          <span>Created {new Date(template.created_at).toLocaleDateString()}</span>
        </div>
      </CardContent>
    </Card>
  );
}
