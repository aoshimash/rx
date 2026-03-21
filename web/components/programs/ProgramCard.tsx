import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useArchiveProgram, useDuplicateProgram } from '@/lib/hooks/usePrograms';
import type { Program } from '@/types/api';
import { Archive, ArrowRightLeft, Copy, Eye } from 'lucide-react';
import Link from 'next/link';

interface ProgramCardProps {
  program: Program;
}

export function ProgramCard({ program }: ProgramCardProps) {
  const entries = program.entries || [];
  const archiveProgram = useArchiveProgram();
  const duplicateProgram = useDuplicateProgram();
  const isArchived = !!program.archived_at;

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
              <CardTitle className="text-xl">{program.name}</CardTitle>
              {isArchived && <Badge variant="secondary">Archived</Badge>}
            </div>
            {program.description && (
              <p className="text-sm text-muted-foreground">{program.description}</p>
            )}
          </div>
          <div className="flex gap-1">
            <Link href={`/programs/${program.id}`}>
              <Button variant="ghost" size="sm">
                <Eye className="h-4 w-4" />
              </Button>
            </Link>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => duplicateProgram.mutate(program.id)}
              disabled={duplicateProgram.isPending}
            >
              <Copy className="h-4 w-4" />
            </Button>
            {!isArchived && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => archiveProgram.mutate(program.id)}
                disabled={archiveProgram.isPending}
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
          <span>Created {new Date(program.created_at).toLocaleDateString()}</span>
        </div>
        {!isArchived && (
          <div className="flex gap-2 mt-4">
            <Link href={`/programs/${program.id}/convert`}>
              <Button variant="outline" size="sm">
                <ArrowRightLeft className="h-4 w-4 mr-2" />
                Convert to Plan
              </Button>
            </Link>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
