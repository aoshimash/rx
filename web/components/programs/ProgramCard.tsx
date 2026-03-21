import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useDeleteProgram } from '@/lib/hooks/usePrograms';
import type { Program } from '@/types/api';
import { Eye, Trash2 } from 'lucide-react';
import Link from 'next/link';

interface ProgramCardProps {
  program: Program;
}

export function ProgramCard({ program }: ProgramCardProps) {
  const deleteProgram = useDeleteProgram();

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <CardTitle className="text-xl">{program.name}</CardTitle>
              <Badge variant={program.status === 'active' ? 'default' : 'secondary'}>
                {program.status}
              </Badge>
            </div>
            {program.notes && <p className="text-sm text-muted-foreground">{program.notes}</p>}
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
              onClick={() => deleteProgram.mutate(program.id)}
              disabled={deleteProgram.isPending}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center gap-4 text-sm text-muted-foreground">
          <span>
            {program.sessions.length} session{program.sessions.length !== 1 ? 's' : ''}
          </span>
          <span>Created {new Date(program.created_at).toLocaleDateString()}</span>
        </div>
      </CardContent>
    </Card>
  );
}
