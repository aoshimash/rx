import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { Log } from '@/types/api';
import { Calendar, Dumbbell, Link as LinkIcon } from 'lucide-react';
import Link from 'next/link';

interface LogCardProps {
  log: Log;
}

export function LogCard({ log }: LogCardProps) {
  const performedDate = new Date(log.performed_at).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });

  const exerciseCount = log.entries.length;

  return (
    <Link href={`/logs/${log.id}`}>
      <Card className="hover:border-primary transition-colors cursor-pointer">
        <CardHeader>
          <div className="flex items-start justify-between">
            <CardTitle className="text-lg flex items-center gap-2">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              {performedDate}
            </CardTitle>
            {log.plan_id && (
              <Badge variant="secondary">
                <LinkIcon className="h-3 w-3 mr-1" />
                Linked to Plan
              </Badge>
            )}
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-1 text-sm text-muted-foreground">
            <Dumbbell className="h-4 w-4" />
            <span>
              {exerciseCount} {exerciseCount === 1 ? 'exercise' : 'exercises'}
            </span>
          </div>
          {log.notes && (
            <p className="text-sm text-muted-foreground mt-2 line-clamp-2">{log.notes}</p>
          )}
        </CardContent>
      </Card>
    </Link>
  );
}
