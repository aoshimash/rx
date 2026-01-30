import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { ExerciseTable } from './ExerciseTable';
import { StatusBadge } from './StatusBadge';
import type { ProgramNode, Workout } from '@/types/api';
import type { DiffStatus } from '@/lib/utils/diff';

interface DayData {
  dayName: string;
  date?: string;
  programNode?: ProgramNode;
  workout?: Workout;
  status: DiffStatus;
}

interface DayAccordionProps {
  days: DayData[];
}

/**
 * Collapsible day accordion for Week View
 * 
 * Each day shows:
 * - Day name and date (if scheduled)
 * - Overall status badge
 * - Exercise table with Plan/Actual/Diff
 */
export function DayAccordion({ days }: DayAccordionProps) {
  return (
    <Accordion type="multiple" className="w-full">
      {days.map((day, idx) => {
        const exercises = buildExerciseRows(day);

        return (
          <AccordionItem key={idx} value={`day-${idx}`}>
            <AccordionTrigger className="hover:no-underline">
              <div className="flex items-center justify-between w-full pr-4">
                <div className="flex items-center gap-3">
                  <span className="font-semibold">{day.dayName}</span>
                  {day.date && (
                    <span className="text-sm text-muted-foreground">
                      {new Date(day.date).toLocaleDateString('en-US', {
                        month: 'short',
                        day: 'numeric',
                      })}
                    </span>
                  )}
                </div>
                <StatusBadge status={day.status} />
              </div>
            </AccordionTrigger>
            <AccordionContent>
              <div className="pt-4">
                <ExerciseTable exercises={exercises} />
              </div>
            </AccordionContent>
          </AccordionItem>
        );
      })}
    </Accordion>
  );
}

/**
 * Build exercise rows from program node and workout data
 */
function buildExerciseRows(day: DayData) {
  const rows: Array<{
    exerciseName: string;
    plan: { sets?: number; reps?: number; load?: number; rpe?: number } | null;
    actual: { sets: number; reps: number; load: number; rpe: number } | null;
    entry?: WorkoutEntry;
    node?: ProgramNode;
  }> = [];

  // Get planned exercises from program node
  const plannedExercises =
    day.programNode?.children?.filter((child) => child.node_type === 'exercise') || [];

  // Get actual exercises from workout
  const actualEntries = day.workout?.entries || [];

  // Create a map of exercise_id to actual entries
  const actualMap = new Map(
    actualEntries.map((entry) => [entry.exercise_id, entry])
  );

  // Add all planned exercises
  for (const node of plannedExercises) {
    const actual = node.exercise_id ? actualMap.get(node.exercise_id) : undefined;
    
    rows.push({
      exerciseName: node.name,
      plan: {
        sets: node.target_sets,
        reps: node.target_reps,
        load: node.percent_1rm ? undefined : undefined, // Load calculation deferred
        rpe: node.target_rpe,
      },
      actual: actual
        ? {
            sets: actual.sets,
            reps: actual.reps,
            load: actual.load_kg,
            rpe: actual.rpe,
          }
        : null,
      entry: actual,
      node,
    });

    // Remove from map to track unplanned
    if (node.exercise_id) {
      actualMap.delete(node.exercise_id);
    }
  }

  // Add unplanned exercises (remaining in actualMap)
  for (const entry of actualMap.values()) {
    rows.push({
      exerciseName: entry.display_name || 'Unknown Exercise',
      plan: null,
      actual: {
        sets: entry.sets,
        reps: entry.reps,
        load: entry.load_kg,
        rpe: entry.rpe,
      },
      entry,
    });
  }

  return rows;
}
