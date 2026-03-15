import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { Plan } from '@/types/api';
import { Calendar, Edit } from 'lucide-react';
import Link from 'next/link';

interface PlanCardProps {
  plan: Plan;
  onSelect?: () => void;
  isSelected?: boolean;
}

export function PlanCard({ plan, onSelect, isSelected }: PlanCardProps) {
  const entries = plan.entries || [];

  const weekCount = new Set(entries.map((e) => e.metadata?.week).filter((w) => w !== undefined))
    .size;

  const dayCount = new Set(
    entries
      .map((e) => {
        const w = e.metadata?.week;
        const d = e.metadata?.day;
        return w !== undefined && d !== undefined ? `${w}::${d}` : String(d);
      })
      .filter((k) => k !== 'undefined')
  ).size;

  return (
    <Card className={isSelected ? 'border-primary' : ''}>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div>
            <CardTitle className="text-xl">{plan.name}</CardTitle>
            {plan.description && (
              <p className="text-sm text-muted-foreground mt-1">{plan.description}</p>
            )}
          </div>
          <Link href={`/plans/${plan.id}/edit`}>
            <Button variant="ghost" size="sm">
              <Edit className="h-4 w-4" />
            </Button>
          </Link>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center gap-4 text-sm text-muted-foreground">
          {weekCount > 0 && (
            <div className="flex items-center gap-1">
              <Calendar className="h-4 w-4" />
              <span>{weekCount} weeks</span>
            </div>
          )}
          {dayCount > 0 && (
            <div>
              <span>{dayCount} training days</span>
            </div>
          )}
          {weekCount === 0 && dayCount === 0 && entries.length > 0 && (
            <div>
              <span>{entries.length} entries</span>
            </div>
          )}
        </div>
        {onSelect && (
          <Button
            onClick={onSelect}
            variant={isSelected ? 'default' : 'outline'}
            className="w-full mt-4"
          >
            {isSelected ? 'Selected' : 'Select Plan'}
          </Button>
        )}
      </CardContent>
    </Card>
  );
}
