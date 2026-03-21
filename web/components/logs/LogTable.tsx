import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { Log, Plan, Program } from '@/types/api';
import Link from 'next/link';

interface LogTableProps {
  logs: Log[];
  planMap: Map<string, Plan>;
  programMap: Map<string, Program>;
}

export function LogTable({ logs, planMap, programMap }: LogTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Date</TableHead>
          <TableHead>Program</TableHead>
          <TableHead>Plan</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {logs.map((log) => {
          const performedDate = new Date(log.performed_at).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
          });
          const plan = log.plan_id ? planMap.get(log.plan_id) : undefined;
          const program = plan?.program_id ? programMap.get(plan.program_id) : undefined;

          return (
            <TableRow key={log.id} className="cursor-pointer">
              <TableCell>
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {performedDate}
                </Link>
              </TableCell>
              <TableCell>
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {program ? program.name : <span className="text-muted-foreground">—</span>}
                </Link>
              </TableCell>
              <TableCell>
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {plan ? plan.name : <span className="text-muted-foreground">—</span>}
                </Link>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
