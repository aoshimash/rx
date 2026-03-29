import { Card, CardContent, CardHeader } from '@/components/ui/card';
import type { Log, LogEntry } from '@/types/api';
import Link from 'next/link';

interface LogCardProps {
  log: Log;
  programName?: string;
}

interface SnapshotEntry {
  exercise_name: string;
  fields?: Record<string, unknown>;
}

function parseSnapshot(snapshot: Record<string, unknown> | undefined): SnapshotEntry[] {
  if (!snapshot?.entries || !Array.isArray(snapshot.entries)) return [];
  return snapshot.entries as SnapshotEntry[];
}

function collectFieldKeys(entries: LogEntry[], snapshotEntries: SnapshotEntry[]): string[] {
  const keys: string[] = [];
  const seen = new Set<string>();
  const addKeys = (fields?: Record<string, unknown>) => {
    if (!fields) return;
    for (const key of Object.keys(fields)) {
      if (!seen.has(key)) {
        seen.add(key);
        keys.push(key);
      }
    }
  };
  for (const e of snapshotEntries) addKeys(e.fields);
  for (const e of entries) addKeys(e.fields);
  return keys;
}

function formatValue(value: unknown): string {
  if (value == null) return '—';
  if (typeof value === 'boolean') return value ? 'Yes' : 'No';
  return String(value);
}

function DiffCell({ planValue, actualValue }: { planValue: unknown; actualValue: unknown }) {
  const pStr = formatValue(planValue);
  const aStr = formatValue(actualValue);

  if (pStr === aStr) return <span className="tabular-nums">{aStr}</span>;

  const pNum = typeof planValue === 'number' ? planValue : null;
  const aNum = typeof actualValue === 'number' ? actualValue : null;
  const color =
    pNum != null && aNum != null
      ? aNum > pNum
        ? 'text-green-600'
        : 'text-red-500'
      : '';

  return (
    <span className="tabular-nums">
      <span className="line-through text-muted-foreground mr-1">{pStr}</span>
      <span className={color}>{aStr}</span>
    </span>
  );
}

export function LogCard({ log, programName }: LogCardProps) {
  const performedDate = new Date(log.performed_at).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });

  const snapshotEntries = parseSnapshot(log.plan_snapshot);
  const hasSnapshot = snapshotEntries.length > 0;
  const sortedEntries = [...log.entries].sort((a, b) => a.order - b.order);
  const fieldKeys = collectFieldKeys(sortedEntries, snapshotEntries);
  const snapshotByName = new Map(snapshotEntries.map((e) => [e.exercise_name, e]));

  return (
    <Link href={`/logs/${log.id}`} className="block">
      <Card className="cursor-pointer transition-colors hover:border-primary">
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 min-w-0">
              <span className="font-semibold truncate">
                {log.session_name ?? 'Untitled Session'}
              </span>
              {log.program_id && programName ? (
                <span className="text-xs text-muted-foreground shrink-0">from {programName}</span>
              ) : !log.program_id ? (
                <span className="text-xs text-muted-foreground italic shrink-0">ad-hoc</span>
              ) : null}
            </div>
            <span className="text-sm text-muted-foreground shrink-0">{performedDate}</span>
          </div>
        </CardHeader>
        <CardContent className="pt-0">
          {sortedEntries.length === 0 ? (
            <p className="text-sm text-muted-foreground">No entries</p>
          ) : (
            <table className="text-sm border-collapse">
              <thead>
                <tr className="text-xs text-muted-foreground">
                  <th className="text-left font-medium pb-1 pr-4 w-[150px]">Exercise</th>
                  {fieldKeys.map((key) => (
                    <th key={key} className="text-left font-medium pb-1 px-2 w-[60px]">
                      {key}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {sortedEntries.map((entry) => {
                  const snapshot = snapshotByName.get(entry.exercise_name);
                  return (
                    <tr key={entry.id} className="text-muted-foreground">
                      <td className="pr-4 py-0.5">{entry.exercise_name}</td>
                      {fieldKeys.map((key) => (
                        <td key={key} className="px-2 py-0.5">
                          {hasSnapshot && snapshot ? (
                            <DiffCell
                              planValue={snapshot.fields?.[key]}
                              actualValue={entry.fields?.[key]}
                            />
                          ) : (
                            <span className="tabular-nums">
                              {formatValue(entry.fields?.[key])}
                            </span>
                          )}
                        </td>
                      ))}
                    </tr>
                  );
                })}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>
    </Link>
  );
}
