'use client';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import type { PlanSession } from '@/types/api';
import { CalendarDays, X } from 'lucide-react';

interface SessionCardProps {
  session: PlanSession;
  programName?: string;
  onLog: (session: PlanSession) => void;
  onDelete: (sessionId: string) => void;
}

/** Collect unique field keys across all entries, preserving insertion order. */
function collectFieldKeys(session: PlanSession): string[] {
  const keys: string[] = [];
  const seen = new Set<string>();
  for (const entry of session.entries) {
    if (!entry.fields) continue;
    for (const key of Object.keys(entry.fields)) {
      if (!seen.has(key)) {
        seen.add(key);
        keys.push(key);
      }
    }
  }
  return keys;
}

function formatFieldValue(value: unknown): string {
  if (value == null) return '—';
  if (typeof value === 'boolean') return value ? 'Yes' : 'No';
  return String(value);
}

export function SessionCard({ session, programName, onLog, onDelete }: SessionCardProps) {
  const sortedEntries = [...session.entries].sort((a, b) => a.order - b.order);
  const fieldKeys = collectFieldKeys(session);

  return (
    <Card
      className="cursor-pointer transition-colors hover:border-primary"
      onClick={() => onLog(session)}
    >
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 min-w-0">
            <span className="font-semibold truncate">{session.session_name}</span>
            {session.source_program_id && programName ? (
              <span className="text-xs text-muted-foreground shrink-0">from {programName}</span>
            ) : !session.source_program_id ? (
              <span className="text-xs text-muted-foreground italic shrink-0">manual</span>
            ) : null}
          </div>
          <div className="flex items-center gap-2 shrink-0">
            {session.date && (
              <span className="flex items-center gap-1 text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded-full">
                <CalendarDays className="h-3 w-3" />
                {session.date}
              </span>
            )}
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6 text-muted-foreground hover:text-destructive"
              onClick={(e) => {
                e.stopPropagation();
                onDelete(session.id);
              }}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        {sortedEntries.length === 0 ? (
          <p className="text-sm text-muted-foreground">No exercises</p>
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
              {sortedEntries.map((entry) => (
                <tr key={entry.id} className="text-muted-foreground">
                  <td className="pr-4 py-0.5">{entry.exercise_name}</td>
                  {fieldKeys.map((key) => (
                    <td key={key} className="px-2 py-0.5 tabular-nums">
                      {formatFieldValue(entry.fields?.[key])}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </CardContent>
    </Card>
  );
}
