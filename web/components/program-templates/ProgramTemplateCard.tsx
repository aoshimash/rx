import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { ProgramTemplate } from '@/types/api';
import { Calendar, User } from 'lucide-react';
import Link from 'next/link';

interface ProgramTemplateCardProps {
  template: ProgramTemplate;
}

function formatDuration(template: ProgramTemplate): string | null {
  const parts: string[] = [];
  if (template.weeks) parts.push(`${template.weeks}w`);
  if (template.days_per_week) parts.push(`${template.days_per_week} days/w`);
  return parts.length > 0 ? parts.join(' · ') : null;
}

export function ProgramTemplateCard({ template }: ProgramTemplateCardProps) {
  const isArchived = !!template.archived_at;
  const duration = formatDuration(template);

  return (
    <Link href={`/program-templates/${template.id}`} className="block">
      <Card
        className={`transition-colors hover:bg-accent/50 cursor-pointer ${isArchived ? 'opacity-60' : ''}`}
      >
        <CardHeader>
          <div className="flex items-center gap-2">
            <CardTitle className="text-xl">{template.name}</CardTitle>
            {isArchived && <Badge variant="secondary">Archived</Badge>}
          </div>
          {template.description && (
            <p className="text-sm text-muted-foreground">{template.description}</p>
          )}
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-3 text-sm text-muted-foreground">
            {duration && (
              <span className="flex items-center gap-1">
                <Calendar className="h-3 w-3" />
                {duration}
              </span>
            )}
            {template.created_by && (
              <span className="flex items-center gap-1">
                <User className="h-3 w-3" />
                {template.created_by}
              </span>
            )}
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}
