'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { usePlans } from '@/lib/hooks/usePlans';
import { groupPlansByCycle } from '@/lib/utils/next-session';
import Link from 'next/link';
import { useMemo } from 'react';

export function ActivePlans() {
  const { data, isLoading } = usePlans();

  const plans = data?.data || [];

  const groups = useMemo(() => {
    const grouped = groupPlansByCycle(plans);
    const result: { groupKey: string; name: string; count: number }[] = [];
    for (const [groupKey, groupPlans] of grouped) {
      const isStandalone = groupKey.startsWith('standalone:');
      const name = isStandalone
        ? groupPlans[0]?.name || 'Unnamed'
        : groupPlans[0]?.name.split(' - ')[0] || 'Unnamed';
      result.push({ groupKey, name, count: groupPlans.length });
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
            <li key={group.groupKey}>
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
