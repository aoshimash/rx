import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { StatusBadge } from './StatusBadge';
import { calculateDiff } from '@/lib/utils/diff';
import type { ProgramNode, WorkoutEntry } from '@/types/api';

interface ExerciseRow {
  exerciseName: string;
  plan: {
    sets?: number;
    reps?: number;
    load?: number;
    rpe?: number;
  } | null;
  actual: {
    sets: number;
    reps: number;
    load: number;
    rpe: number;
  } | null;
  entry?: WorkoutEntry;
  node?: ProgramNode;
}

interface ExerciseTableProps {
  exercises: ExerciseRow[];
}

/**
 * Exercise table showing Plan/Actual/Diff columns
 * 
 * Displays planned exercises from program nodes and actual workout entries
 * with diff status indicators.
 */
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
          const planSnapshot = row.plan
            ? {
                target_sets: row.plan.sets,
                target_reps: row.plan.reps,
                target_load_kg: row.plan.load,
                target_rpe: row.plan.rpe,
              }
            : null;

          const diff = calculateDiff(planSnapshot, row.entry);

          return (
            <TableRow key={idx}>
              <TableCell className="font-medium">{row.exerciseName}</TableCell>
              <TableCell className="text-center">
                {row.plan ? (
                  <div className="space-y-1 text-sm">
                    <div>
                      {row.plan.sets}×{row.plan.reps} @ {row.plan.load}kg
                    </div>
                    {row.plan.rpe && (
                      <div className="text-xs text-muted-foreground">
                        RPE {row.plan.rpe}
                      </div>
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
                      {row.actual.sets}×{row.actual.reps} @ {row.actual.load}kg
                    </div>
                    {row.actual.rpe && (
                      <div className="text-xs text-muted-foreground">
                        RPE {row.actual.rpe}
                      </div>
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
