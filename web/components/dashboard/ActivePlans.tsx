'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { usePlans } from '@/lib/hooks/usePlans';
import Link from 'next/link';

export function ActivePlans() {
  const { data, isLoading } = usePlans();

  const plans = data?.data.slice(0, 5);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Plans</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading && <p className="text-sm text-muted-foreground">Loading...</p>}
        {!isLoading && (!plans || plans.length === 0) && (
          <p className="text-sm text-muted-foreground">
            No plans yet.{' '}
            <Link href="/plans" className="underline">
              Create one
            </Link>
          </p>
        )}
        <ul className="space-y-2">
          {plans?.map((plan) => {
            const sessionCount = new Set(
              (plan.entries || []).map((e) => (e.metadata?.session as string) || 'Default')
            ).size;
            return (
              <li key={plan.id}>
                <Link
                  href={`/plans/${plan.id}/edit`}
                  className="flex items-center justify-between rounded-md px-3 py-2 text-sm hover:bg-accent transition-colors"
                >
                  <span className="font-medium">{plan.name}</span>
                  <span className="text-muted-foreground">
                    {sessionCount} session{sessionCount !== 1 ? 's' : ''}
                  </span>
                </Link>
              </li>
            );
          })}
        </ul>
      </CardContent>
    </Card>
  );
}
