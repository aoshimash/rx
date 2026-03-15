import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import type { DiffStatus } from '@/lib/utils/diff';
import { formatPlanContext } from '@/lib/utils/planContext';
import type { Log, LogEntry, PlanEntry } from '@/types/api';
import { ExerciseTable } from './ExerciseTable';
import { StatusBadge } from './StatusBadge';

interface DayData {
  dayName: string;
  date?: string;
  planEntries?: PlanEntry[];
  log?: Log;
  status: DiffStatus;
  planContext?: string[];
}

interface DayAccordionProps {
  days: DayData[];
}

export function DayAccordion({ days }: DayAccordionProps) {
  return (
    <Accordion type="multiple" className="w-full">
      {days.map((day, idx) => {
        const exercises = buildExerciseRows(day);

        return (
          <AccordionItem key={idx} value={`day-${idx}`}>
            <AccordionTrigger className="hover:no-underline">
              <div className="flex items-center justify-between w-full pr-4">
                <div className="flex flex-col items-start gap-1">
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
                  {day.planContext && day.planContext.length > 0 && (
                    <span className="text-xs text-muted-foreground">
                      {formatPlanContext(day.planContext)}
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

function buildExerciseRows(day: DayData) {
  const rows: Array<{
    exerciseName: string;
    plan: { sets?: number; reps?: number; load_kg?: number; rpe?: number } | null;
    actual: { sets?: number; reps?: number; load_kg?: number; rpe?: number } | null;
    planEntry?: PlanEntry;
    logEntry?: LogEntry;
  }> = [];

  const plannedEntries = day.planEntries || [];
  const actualEntries = day.log?.entries || [];

  // Map actual entries by exercise_name
  const actualMap = new Map(actualEntries.map((entry) => [entry.exercise_name, entry]));

  for (const planEntry of plannedEntries) {
    const actual = actualMap.get(planEntry.exercise_name);

    rows.push({
      exerciseName: planEntry.exercise_name,
      plan: {
        sets: planEntry.sets,
        reps: planEntry.reps,
        load_kg: planEntry.load_kg,
        rpe: planEntry.rpe,
      },
      actual: actual
        ? {
            sets: actual.sets,
            reps: actual.reps,
            load_kg: actual.load_kg,
            rpe: actual.rpe,
          }
        : null,
      planEntry,
      logEntry: actual,
    });

    actualMap.delete(planEntry.exercise_name);
  }

  // Add unplanned exercises
  for (const entry of actualMap.values()) {
    rows.push({
      exerciseName: entry.exercise_name,
      plan: null,
      actual: {
        sets: entry.sets,
        reps: entry.reps,
        load_kg: entry.load_kg,
        rpe: entry.rpe,
      },
      logEntry: entry,
    });
  }

  return rows;
}
