import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useProgramTemplate } from '@/lib/hooks/useProgramTemplates';
import { useDeleteProgram } from '@/lib/hooks/usePrograms';
import type { Program } from '@/types/api';
import { Eye, Trash2 } from 'lucide-react';
import Link from 'next/link';

interface ProgramCardProps {
  program: Program;
}

function TemplateInfo({ program }: { program: Program }) {
  const { data: template } = useProgramTemplate(program.program_template_id ?? null);
  const targetWeights = program.metadata?.target_weights as Record<string, number> | undefined;

  if (!template && !targetWeights) return null;

  return (
    <div className="text-sm text-muted-foreground space-y-0.5">
      {template && <p>Template: {template.name}</p>}
      {targetWeights && (
        <p>
          Targets:{' '}
          {Object.entries(targetWeights)
            .map(([exercise, weight]) => `${exercise}: ${weight}kg`)
            .join(', ')}
        </p>
      )}
    </div>
  );
}

export function ProgramCard({ program }: ProgramCardProps) {
  const deleteProgram = useDeleteProgram();
  const isCompleted = program.status === 'completed';

  return (
    <Card className={isCompleted ? 'opacity-50' : undefined}>
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
        <div className="space-y-1">
          {program.program_template_id && <TemplateInfo program={program} />}
          <p className="text-sm text-muted-foreground">
            Created {new Date(program.created_at).toLocaleDateString()}
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
