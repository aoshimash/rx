'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useLogs } from '@/lib/hooks/useLogs';
import { format, parseISO } from 'date-fns';
import Link from 'next/link';

export function RecentLogs() {
  const { data, isLoading } = useLogs();

  const recentLogs = data?.data
    .slice()
    .sort((a, b) => new Date(b.performed_at).getTime() - new Date(a.performed_at).getTime())
    .slice(0, 5);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Recent Logs</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading && <p className="text-sm text-muted-foreground">Loading...</p>}
        {!isLoading && (!recentLogs || recentLogs.length === 0) && (
          <p className="text-sm text-muted-foreground">No logs yet.</p>
        )}
        <ul className="space-y-2">
          {recentLogs?.map((log) => (
            <li key={log.id}>
              <Link
                href={`/logs/${log.id}`}
                className="flex items-center justify-between rounded-md px-3 py-2 text-sm hover:bg-accent transition-colors"
              >
                <span className="font-medium">
                  {format(parseISO(log.performed_at), 'MMM d, yyyy')}
                </span>
                <span className="text-muted-foreground">{log.entries.length} exercises</span>
              </Link>
            </li>
          ))}
        </ul>
      </CardContent>
    </Card>
  );
}
