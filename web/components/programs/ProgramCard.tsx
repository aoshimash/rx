import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { Program } from '@/types/api';
import { ArrowRightLeft, Edit, Eye } from 'lucide-react';
import Link from 'next/link';

interface ProgramCardProps {
  program: Program;
}

export function ProgramCard({ program }: ProgramCardProps) {
  const entries = program.entries || [];

  const sessionNames = new Set(
    entries
      .map((e) => (e.metadata?.session as string) ?? undefined)
      .filter((s): s is string => s !== undefined)
  );

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="space-y-1">
            <CardTitle className="text-xl">{program.name}</CardTitle>
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
            <Link href={`/programs/${program.id}/edit`}>
              <Button variant="ghost" size="sm">
                <Edit className="h-4 w-4" />
              </Button>
            </Link>
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
        <div className="flex gap-2 mt-4">
          <Link href={`/programs/${program.id}/convert`}>
            <Button variant="outline" size="sm">
              <ArrowRightLeft className="h-4 w-4 mr-2" />
              Convert to Plan
            </Button>
          </Link>
        </div>
      </CardContent>
    </Card>
  );
}
