import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { Log } from '@/types/api';
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
          <TableHead>Session</TableHead>
          <TableHead>Entries</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {logs.map((log) => {
          const performedDate = new Date(log.performed_at).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
          });

          return (
            <TableRow key={log.id} className="cursor-pointer">
              <TableCell>
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {performedDate}
                </Link>
              </TableCell>
              <TableCell>
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {log.session_name ?? <span className="text-muted-foreground">—</span>}
                </Link>
              </TableCell>
              <TableCell>
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {log.entries.length} exercise{log.entries.length !== 1 ? 's' : ''}
                </Link>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
