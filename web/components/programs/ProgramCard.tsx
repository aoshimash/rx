import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { Program } from '@/types/api';
import { Calendar, Edit } from 'lucide-react';
import Link from 'next/link';

interface ProgramCardProps {
  program: Program;
  onSelect?: () => void;
  isSelected?: boolean;
}

/**
 * Program list item card
 */
export function ProgramCard({ program, onSelect, isSelected }: ProgramCardProps) {
  const weekCount = program.root_nodes?.filter((node) => node.node_type === 'week').length || 0;
  const dayCount =
    program.root_nodes?.reduce((count, week) => {
      const days = week.children?.filter((child) => child.node_type === 'day') || [];
      return count + days.length;
    }, 0) || 0;

  return (
    <Card className={isSelected ? 'border-primary' : ''}>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div>
            <CardTitle className="text-xl">{program.name}</CardTitle>
            {program.description && (
              <p className="text-sm text-muted-foreground mt-1">{program.description}</p>
            )}
          </div>
          <Link href={`/programs/${program.id}/edit`}>
            <Button variant="ghost" size="sm">
              <Edit className="h-4 w-4" />
            </Button>
          </Link>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center gap-4 text-sm text-muted-foreground">
          <div className="flex items-center gap-1">
            <Calendar className="h-4 w-4" />
            <span>{weekCount} weeks</span>
          </div>
          <div>
            <span>{dayCount} training days</span>
          </div>
        </div>
        {onSelect && (
          <Button
            onClick={onSelect}
            variant={isSelected ? 'default' : 'outline'}
            className="w-full mt-4"
          >
            {isSelected ? 'Selected' : 'Select Program'}
          </Button>
        )}
      </CardContent>
    </Card>
  );
}
