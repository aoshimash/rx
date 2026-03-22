import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { ProgramTemplate } from '@/types/api';
import { User } from 'lucide-react';
import Link from 'next/link';

interface ProgramTemplateCardProps {
  template: ProgramTemplate;
}

export function ProgramTemplateCard({ template }: ProgramTemplateCardProps) {
  const isArchived = !!template.archived_at;

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
          <div className="space-y-1">
            {template.created_by && (
              <div className="flex items-center gap-1 text-sm text-muted-foreground">
                <User className="h-3 w-3" />
                {template.created_by}
              </div>
            )}
            <div className="text-sm text-muted-foreground">
              Created {new Date(template.created_at).toLocaleDateString()}
            </div>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}
