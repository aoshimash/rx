import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import type { DiffStatus } from '@/lib/utils/diff';
import { formatProgramContext } from '@/lib/utils/programContext';
import type { EntryType, ProgramEntry, Workout, WorkoutEntry } from '@/types/api';
import { ExerciseTable } from './ExerciseTable';
import { StatusBadge } from './StatusBadge';

interface DayData {
  dayName: string;
  date?: string;
  programEntries?: ProgramEntry[];
  workout?: Workout;
  status: DiffStatus;
  programContext?: string[];
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
                  {day.programContext && day.programContext.length > 0 && (
                    <span className="text-xs text-muted-foreground">
                      {formatProgramContext(day.programContext)}
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
 * Build exercise rows from program entries and workout data
 */
function buildExerciseRows(day: DayData) {
  const rows: Array<{
    exerciseName: string;
    entryType?: EntryType;
    plan: { sets?: number; reps?: number; load?: number; rpe?: number } | null;
    actual: { sets: number; reps: number; load: number; rpe: number } | null;
    entry?: WorkoutEntry;
    programEntry?: ProgramEntry;
  }> = [];

  const plannedEntries = day.programEntries || [];
  const actualEntries = day.workout?.entries || [];

  // Create a map of exercise_id to actual entries
  const actualMap = new Map(actualEntries.map((entry) => [entry.exercise_id, entry]));

  // Add all planned entries
  for (const programEntry of plannedEntries) {
    const actual = programEntry.exercise_id
      ? actualMap.get(programEntry.exercise_id)
      : undefined;

    rows.push({
      exerciseName: programEntry.name,
      entryType: actual?.entry_type,
      plan: {
        sets: programEntry.target_sets,
        reps: programEntry.target_reps,
        load: undefined,
        rpe: programEntry.target_rpe,
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
      programEntry,
    });

    // Remove from map to track unplanned
    if (programEntry.exercise_id) {
      actualMap.delete(programEntry.exercise_id);
    }
  }

  // Add unplanned exercises (remaining in actualMap)
  for (const entry of actualMap.values()) {
    rows.push({
      exerciseName: entry.display_name || 'Unknown Exercise',
      entryType: entry.entry_type,
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
