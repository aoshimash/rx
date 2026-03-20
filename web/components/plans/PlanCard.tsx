import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { PlanStatus } from '@/lib/utils/next-session';
import { Edit } from 'lucide-react';
import Link from 'next/link';
import { ExerciseSummary } from './ExerciseSummary';

interface PlanCardProps {
  status: PlanStatus;
  onRecordLog: () => void;
}

export function PlanCard({ status, onRecordLog }: PlanCardProps) {
  const { plan, isNext, completedCount, lastPerformedAt } = status;

  const lastDate = lastPerformedAt
    ? new Date(lastPerformedAt).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
      })
    : null;

  const displayName = plan.session_name || plan.name;

  return (
    <Card className={isNext ? 'border-primary' : ''}>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <CardTitle className="text-base">{displayName}</CardTitle>
            {plan.date && <span className="text-xs text-muted-foreground">{plan.date}</span>}
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
        <ExerciseSummary exercises={plan.entries || []} />
        <div className="flex gap-2">
          <Button variant={isNext ? 'default' : 'outline'} size="sm" onClick={onRecordLog}>
            Record Log
          </Button>
          <Link href={`/plans/${plan.id}/edit`}>
            <Button variant="ghost" size="sm">
              <Edit className="h-4 w-4 mr-1" />
              Edit
            </Button>
          </Link>
        </div>
      </CardContent>
    </Card>
  );
}
