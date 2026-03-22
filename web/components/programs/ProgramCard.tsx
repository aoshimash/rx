import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useProgramTemplate } from '@/lib/hooks/useProgramTemplates';
import { useLoggedSessions } from '@/lib/hooks/usePrograms';
import type { Program } from '@/types/api';
import { useRouter } from 'next/navigation';

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

function statusBadgeVariant(
  status: string,
): 'default' | 'secondary' | 'outline' | 'destructive' {
  if (status === 'ongoing') return 'default';
  if (status === 'completed') return 'secondary';
  if (status === 'cancelled') return 'destructive';
  return 'outline';
}

export function OngoingProgramCard({ program }: ProgramCardProps) {
  const router = useRouter();
  const { data: loggedSessions } = useLoggedSessions(program.id);

  const totalSessions = program.sessions.length;
  const loggedSet = new Set(loggedSessions?.sessions ?? []);
  const completedCount =
    totalSessions > 0 ? program.sessions.filter((s) => loggedSet.has(s.session_name)).length : 0;
  const progressPct = totalSessions > 0 ? (completedCount / totalSessions) * 100 : 0;

  return (
    <Card
      className="cursor-pointer transition-colors hover:bg-accent/50"
      onClick={() => router.push(`/programs/${program.id}`)}
    >
      <CardHeader>
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <CardTitle className="text-xl">{program.name}</CardTitle>
            <Badge variant="default">ongoing</Badge>
          </div>
          {program.notes && <p className="text-sm text-muted-foreground">{program.notes}</p>}
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {program.program_template_id && <TemplateInfo program={program} />}
          <div className="space-y-1">
            <div className="flex justify-between text-xs text-muted-foreground">
              <span>
                {completedCount} / {totalSessions} sessions
              </span>
            </div>
            <div className="h-2 w-full rounded-full bg-secondary overflow-hidden">
              <div
                className="h-full rounded-full bg-primary transition-all"
                style={{ width: `${progressPct}%` }}
              />
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export function CreatedProgramCard({ program }: ProgramCardProps) {
  const router = useRouter();

  return (
    <Card
      className="cursor-pointer transition-colors hover:bg-accent/50"
      onClick={() => router.push(`/programs/${program.id}`)}
    >
      <CardHeader className="pb-2">
        <CardTitle className="text-base">{program.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex justify-between text-sm text-muted-foreground">
          <span>{program.sessions.length} sessions</span>
          <span>{new Date(program.created_at).toLocaleDateString()}</span>
        </div>
      </CardContent>
    </Card>
  );
}

export function FinishedProgramCard({ program }: ProgramCardProps) {
  const router = useRouter();

  return (
    <Card
      className="cursor-pointer transition-colors hover:bg-accent/50 opacity-50"
      onClick={() => router.push(`/programs/${program.id}`)}
    >
      <CardHeader className="pb-2">
        <div className="flex items-center gap-2">
          <CardTitle className="text-base">{program.name}</CardTitle>
          <Badge variant={statusBadgeVariant(program.status)}>{program.status}</Badge>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex justify-between text-sm text-muted-foreground">
          <span>{program.sessions.length} sessions</span>
          <span>{new Date(program.created_at).toLocaleDateString()}</span>
        </div>
      </CardContent>
    </Card>
  );
}
