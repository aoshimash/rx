import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { calculateDiff } from '@/lib/utils/diff';
import type { LogEntry, PlanEntry } from '@/types/api';
import { StatusBadge } from './StatusBadge';

interface ExerciseRow {
  exerciseName: string;
  plan: {
    sets?: number;
    reps?: number;
    load_kg?: number;
    rpe?: number;
  } | null;
  actual: {
    sets?: number;
    reps?: number;
    load_kg?: number;
    rpe?: number;
  } | null;
  planEntry?: PlanEntry;
  logEntry?: LogEntry;
}

interface ExerciseTableProps {
  exercises: ExerciseRow[];
}

export function ExerciseTable({ exercises }: ExerciseTableProps) {
  if (exercises.length === 0) {
    return (
      <div className="text-sm text-muted-foreground text-center py-4">
        No exercises for this day
      </div>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[200px]">Exercise</TableHead>
          <TableHead className="text-center">Plan</TableHead>
          <TableHead className="text-center">Actual</TableHead>
          <TableHead className="text-center">Status</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {exercises.map((row, idx) => {
          const diff = calculateDiff(row.planEntry, row.logEntry);

          return (
            <TableRow key={idx}>
              <TableCell className="font-medium">{row.exerciseName}</TableCell>
              <TableCell className="text-center">
                {row.plan ? (
                  <div className="space-y-1 text-sm">
                    <div>
                      {row.plan.sets}x{row.plan.reps}
                      {row.plan.load_kg !== undefined && ` @ ${row.plan.load_kg}kg`}
                    </div>
                    {row.plan.rpe !== undefined && (
                      <div className="text-xs text-muted-foreground">RPE {row.plan.rpe}</div>
                    )}
                  </div>
                ) : (
                  <span className="text-muted-foreground text-sm">-</span>
                )}
              </TableCell>
              <TableCell className="text-center">
                {row.actual ? (
                  <div className="space-y-1 text-sm">
                    <div>
                      {row.actual.sets}x{row.actual.reps}
                      {row.actual.load_kg !== undefined && ` @ ${row.actual.load_kg}kg`}
                    </div>
                    {row.actual.rpe !== undefined && (
                      <div className="text-xs text-muted-foreground">RPE {row.actual.rpe}</div>
                    )}
                  </div>
                ) : (
                  <span className="text-muted-foreground text-sm">-</span>
                )}
              </TableCell>
              <TableCell className="text-center">
                <StatusBadge status={diff.status} differences={diff.differences} />
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
