import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { SessionStatus } from '@/lib/utils/next-session';
import { ExerciseSummary } from './ExerciseSummary';

interface SessionCardProps {
  status: SessionStatus;
  onRecordLog: () => void;
}

export function SessionCard({ status, onRecordLog }: SessionCardProps) {
  const { session, isNext, completedCount, lastPerformedAt } = status;

  const lastDate = lastPerformedAt
    ? new Date(lastPerformedAt).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
      })
    : null;

  return (
    <Card className={isNext ? 'border-primary' : ''}>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <CardTitle className="text-base">{session.name}</CardTitle>
            {isNext && (
              <Badge variant="default" className="text-xs">
                NEXT
              </Badge>
            )}
          </div>
          <div className="flex items-center gap-3 text-xs text-muted-foreground">
            {completedCount > 0 && (
              <span>
                {completedCount}x done{lastDate && ` · ${lastDate}`}
              </span>
            )}
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        <ExerciseSummary exercises={session.exercises} />
        <Button variant={isNext ? 'default' : 'outline'} size="sm" onClick={onRecordLog}>
          Record Log
        </Button>
      </CardContent>
    </Card>
  );
}
