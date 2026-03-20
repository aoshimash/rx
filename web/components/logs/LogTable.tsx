import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { Log } from '@/types/api';
import { Link as LinkIcon } from 'lucide-react';
import Link from 'next/link';

interface LogTableProps {
  logs: Log[];
}

export function LogTable({ logs }: LogTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Date</TableHead>
          <TableHead>Exercises</TableHead>
          <TableHead className="hidden sm:table-cell">Plan</TableHead>
          <TableHead className="hidden md:table-cell">Notes</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {logs.map((log) => {
          const performedDate = new Date(log.performed_at).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
          });
          const exerciseCount = log.entries.length;

          return (
            <TableRow key={log.id} className="cursor-pointer">
              <TableCell>
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {performedDate}
                </Link>
              </TableCell>
              <TableCell>
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {exerciseCount} {exerciseCount === 1 ? 'exercise' : 'exercises'}
                </Link>
              </TableCell>
              <TableCell className="hidden sm:table-cell">
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {log.plan_id ? (
                    <span className="flex items-center gap-1 text-muted-foreground">
                      <LinkIcon className="h-3 w-3" />
                      Linked
                    </span>
                  ) : (
                    <span className="text-muted-foreground">—</span>
                  )}
                </Link>
              </TableCell>
              <TableCell className="hidden md:table-cell max-w-[200px]">
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {log.notes ? (
                    <span className="truncate text-muted-foreground">{log.notes}</span>
                  ) : (
                    <span className="text-muted-foreground">—</span>
                  )}
                </Link>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
