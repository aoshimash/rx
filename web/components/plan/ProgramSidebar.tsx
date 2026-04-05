'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { usePrograms } from '@/lib/hooks/usePrograms';
import { Plus } from 'lucide-react';
import Link from 'next/link';

export function ProgramSidebar() {
  const { data, isLoading } = usePrograms();
  const programs = data?.data ?? [];

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <h3 className="font-semibold text-sm">Programs</h3>
        <Button variant="default" size="sm" className="h-7 text-xs" asChild>
          <Link href="/programs/new">
            <Plus className="h-3 w-3 mr-1" />
            New
          </Link>
        </Button>
      </div>

      {isLoading ? (
        <div className="space-y-2">
          <Skeleton className="h-14 w-full" />
          <Skeleton className="h-14 w-full" />
        </div>
      ) : programs.length === 0 ? (
        <p className="text-sm text-muted-foreground">
          No programs yet. Create one to quickly add sessions to your plan.
        </p>
      ) : (
        <div className="space-y-1.5">
          {programs.map((program) => (
            <Link
              key={program.id}
              href={`/programs/${program.id}`}
              className="block rounded-md border p-2.5 text-sm hover:border-primary transition-colors"
            >
              <div className="flex items-center gap-1.5">
                <span className="font-medium truncate">{program.name}</span>
                {program.status === 'draft' && (
                  <Badge variant="secondary" className="text-[10px] px-1.5 py-0">
                    Draft
                  </Badge>
                )}
              </div>
              <div className="text-xs text-muted-foreground mt-0.5">
                {program.sessions.length} session{program.sessions.length !== 1 ? 's' : ''}
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
