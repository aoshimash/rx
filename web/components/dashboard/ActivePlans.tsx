'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { usePrograms } from '@/lib/hooks/usePrograms';
import Link from 'next/link';

export function ActivePlans() {
  const { data, isLoading } = usePrograms('active');

  const programs = data?.data || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Active Programs</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading && <p className="text-sm text-muted-foreground">Loading...</p>}
        {!isLoading && programs.length === 0 && (
          <p className="text-sm text-muted-foreground">
            No active programs.{' '}
            <Link href="/programs" className="underline">
              Create one
            </Link>
          </p>
        )}
        <ul className="space-y-2">
          {programs.slice(0, 5).map((program) => (
            <li key={program.id}>
              <Link
                href={`/programs/${program.id}`}
                className="flex items-center justify-between rounded-md px-3 py-2 text-sm hover:bg-accent transition-colors"
              >
                <span className="font-medium">{program.name}</span>
                <span className="text-muted-foreground">
                  {program.sessions.length} session{program.sessions.length !== 1 ? 's' : ''}
                </span>
              </Link>
            </li>
          ))}
        </ul>
      </CardContent>
    </Card>
  );
}
