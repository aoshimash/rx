import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { Log, LogEntry } from '@/types/api';
import Link from 'next/link';

interface LogTableProps {
  logs: Log[];
  programMap?: Map<string, string>; // program_id → program name
}

function totalVolumeKg(entries: LogEntry[]): number | null {
  const volumes = entries
    .filter((e) => e.fields?.load_kg != null && e.fields?.sets != null && e.fields?.reps != null)
    .map(
      (e) => (e.fields?.sets as number) * (e.fields?.reps as number) * (e.fields?.load_kg as number)
    );
  if (volumes.length === 0) return null;
  return volumes.reduce((a, b) => a + b, 0);
}

function durationMinutes(log: Log): number | null {
  if (!log.started_at || !log.finished_at) return null;
  const diff = new Date(log.finished_at).getTime() - new Date(log.started_at).getTime();
  return Math.round(diff / 60000);
}

const dash = <span className="text-muted-foreground">—</span>;

export function LogTable({ logs, programMap }: LogTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Date</TableHead>
          <TableHead>Program</TableHead>
          <TableHead>Session</TableHead>
          <TableHead className="text-right">Volume</TableHead>
          <TableHead className="text-right">Duration</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {logs.map((log) => {
          const performedDate = new Date(log.performed_at).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
          });
          const programName = log.program_id ? (programMap?.get(log.program_id) ?? null) : null;
          const volume = totalVolumeKg(log.entries);
          const duration = durationMinutes(log);

          return (
            <TableRow key={log.id} className="cursor-pointer">
              <TableCell>
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {performedDate}
                </Link>
              </TableCell>
              <TableCell>
                {log.program_id && programName ? (
                  <Link href={`/programs/${log.program_id}`} className="hover:underline text-sm">
                    {programName}
                  </Link>
                ) : (
                  dash
                )}
              </TableCell>
              <TableCell>
                <Link href={`/logs/${log.id}`} className="block w-full text-sm">
                  {log.session_name ?? dash}
                </Link>
              </TableCell>
              <TableCell className="text-right">
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {volume != null ? (
                    <span className="tabular-nums">
                      {volume.toLocaleString('en-US', { maximumFractionDigits: 0 })} kg
                    </span>
                  ) : (
                    dash
                  )}
                </Link>
              </TableCell>
              <TableCell className="text-right">
                <Link href={`/logs/${log.id}`} className="block w-full">
                  {duration != null ? <span className="tabular-nums">{duration} min</span> : dash}
                </Link>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
