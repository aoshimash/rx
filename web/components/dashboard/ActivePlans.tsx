'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { usePlans } from '@/lib/hooks/usePlans';
import { groupPlansByProgram } from '@/lib/utils/next-session';
import Link from 'next/link';
import { useMemo } from 'react';

export function ActivePlans() {
  const { data, isLoading } = usePlans();

  const plans = data?.data || [];

  const groups = useMemo(() => {
    const grouped = groupPlansByProgram(plans);
    const result: { programId: string | null; name: string; count: number }[] = [];
    for (const [programId, groupPlans] of grouped) {
      const name = programId
        ? groupPlans[0]?.name.split(' - ')[0] || 'Unnamed'
        : 'Standalone Plans';
      result.push({ programId, name, count: groupPlans.length });
    }
    return result.slice(0, 5);
  }, [plans]);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Plans</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading && <p className="text-sm text-muted-foreground">Loading...</p>}
        {!isLoading && plans.length === 0 && (
          <p className="text-sm text-muted-foreground">
            No plans yet.{' '}
            <Link href="/plans" className="underline">
              Create one
            </Link>
          </p>
        )}
        <ul className="space-y-2">
          {groups.map((group) => (
            <li key={group.programId ?? 'standalone'}>
              <Link
                href="/plans"
                className="flex items-center justify-between rounded-md px-3 py-2 text-sm hover:bg-accent transition-colors"
              >
                <span className="font-medium">{group.name}</span>
                <span className="text-muted-foreground">
                  {group.count} session{group.count !== 1 ? 's' : ''}
                </span>
              </Link>
            </li>
          ))}
        </ul>
      </CardContent>
    </Card>
  );
}
