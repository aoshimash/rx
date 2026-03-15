'use client';

import { ScheduleModal } from '@/components/schedule/ScheduleModal';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { DiffStatus } from '@/lib/utils/diff';
import { calculateDiff } from '@/lib/utils/diff';
import { generateProgramContext } from '@/lib/utils/programContext';
import type { DaySchedule } from '@/lib/utils/schedule';
import { useScheduleStore } from '@/stores/schedule';
import type { Program, ProgramEntry, Workout } from '@/types/api';
import { addWeeks, endOfWeek, format, startOfWeek } from 'date-fns';
import { Calendar as CalendarIcon, ChevronLeft, ChevronRight } from 'lucide-react';
import { useMemo, useState } from 'react';
import { DayAccordion } from './DayAccordion';
import { UnplannedWorkouts } from './UnplannedWorkouts';

interface WeekViewProps {
  program: Program | null;
  workouts: Workout[];
}

/**
 * Week View component with navigation and Plan/Actual/Diff display
 *
 * Shows:
 * - Week navigation (Prev/Next)
 * - Entries grouped by metadata.week, then by metadata.day
 * - Ungrouped entries (no metadata.week) shown in a separate section
 * - Schedule configuration
 */
export function WeekView({ program, workouts }: WeekViewProps) {
  const [weekOffset, setWeekOffset] = useState(0);
  const [scheduleModalOpen, setScheduleModalOpen] = useState(false);
  const { getSchedule, setSchedule } = useScheduleStore();

  // Calculate current week range
  const { weekStart, weekEnd } = useMemo(() => {
    const today = new Date();
    const baseDate = addWeeks(today, weekOffset);
    return {
      weekStart: startOfWeek(baseDate, { weekStartsOn: 1 }), // Monday
      weekEnd: endOfWeek(baseDate, { weekStartsOn: 1 }),
    };
  }, [weekOffset]);

  // Get schedule for current program
  const schedule = program ? getSchedule(program.id) : undefined;

  // Build day data from program entries and workouts
  const days = useMemo(() => {
    if (!program) return [];

    const entries = program.entries || [];

    // Group entries by metadata.week, then by metadata.day
    // For the week view, we show the first week group (or all ungrouped)
    const weekGroups = new Map<string, ProgramEntry[]>();
    const ungroupedEntries: ProgramEntry[] = [];

    for (const entry of entries) {
      const weekKey = entry.metadata?.week !== undefined ? String(entry.metadata.week) : null;
      if (weekKey !== null) {
        if (!weekGroups.has(weekKey)) weekGroups.set(weekKey, []);
        weekGroups.get(weekKey)?.push(entry);
      } else {
        ungroupedEntries.push(entry);
      }
    }

    // Show the first week group's days
    const firstWeekEntries =
      weekGroups.size > 0 ? weekGroups.values().next().value : ungroupedEntries;

    if (!firstWeekEntries || firstWeekEntries.length === 0) return [];

    // Group this week's entries by metadata.day
    const dayGroups = new Map<string, ProgramEntry[]>();
    const ungroupedDayEntries: ProgramEntry[] = [];

    for (const entry of firstWeekEntries) {
      const dayKey = entry.metadata?.day !== undefined ? String(entry.metadata.day) : null;
      if (dayKey !== null) {
        if (!dayGroups.has(dayKey)) dayGroups.set(dayKey, []);
        dayGroups.get(dayKey)?.push(entry);
      } else {
        ungroupedDayEntries.push(entry);
      }
    }

    // Build day data from day groups
    const result = [...dayGroups.entries()].map(([dayName, dayEntries], idx) => {
      // Find workout for this day (matching first entry's program_node_id)
      const firstEntry = dayEntries[0];
      const workout = workouts.find((w) => w.program_node_id === firstEntry?.id);

      // Get scheduled date if available
      const scheduledDate = schedule?.find((s) => s.dayIndex === idx)?.date;

      // Calculate overall day status
      const status = calculateDayStatus(dayEntries, workout);

      // Get program context from workout or generate from first entry
      const programContext =
        workout?.program_context ||
        (firstEntry ? generateProgramContext(firstEntry.id, program) : []);

      return {
        dayName,
        date: scheduledDate,
        programEntries: dayEntries,
        workout,
        status,
        programContext,
      };
    });

    // Add ungrouped entries as a single "day" if any exist
    if (ungroupedDayEntries.length > 0) {
      const workout = workouts.find((w) =>
        ungroupedDayEntries.some((e) => w.program_node_id === e.id)
      );
      result.push({
        dayName: 'Ungrouped',
        date: undefined,
        programEntries: ungroupedDayEntries,
        workout,
        status: calculateDayStatus(ungroupedDayEntries, workout),
        programContext: [],
      });
    }

    return result;
  }, [program, workouts, schedule]);

  // Get unplanned workouts (not linked to any program entry)
  const unplannedWorkouts = useMemo(() => {
    if (!program) return workouts;

    const entryIds = new Set((program.entries || []).map((e) => e.id));
    return workouts.filter((w) => !w.program_node_id || !entryIds.has(w.program_node_id));
  }, [program, workouts]);

  const handlePrevWeek = () => setWeekOffset((prev) => prev - 1);
  const handleNextWeek = () => setWeekOffset((prev) => prev + 1);
  const handleScheduleGenerated = (newSchedule: DaySchedule[]) => {
    if (program) {
      setSchedule(program.id, newSchedule);
    }
  };

  if (!program) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center text-muted-foreground">
            No program selected. Create or select a program to view training week.
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Week View</CardTitle>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" onClick={() => setScheduleModalOpen(true)}>
                <CalendarIcon className="h-4 w-4 mr-2" />
                Schedule
              </Button>
              <Button variant="outline" size="sm" onClick={handlePrevWeek}>
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <span className="text-sm font-medium min-w-[200px] text-center">
                {format(weekStart, 'MMM d')} - {format(weekEnd, 'MMM d, yyyy')}
              </span>
              <Button variant="outline" size="sm" onClick={handleNextWeek}>
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {days.length > 0 ? (
            <DayAccordion days={days} />
          ) : (
            <div className="text-center text-muted-foreground py-4">No days found in program</div>
          )}
        </CardContent>
      </Card>

      <UnplannedWorkouts workouts={unplannedWorkouts} />

      <ScheduleModal
        open={scheduleModalOpen}
        onOpenChange={setScheduleModalOpen}
        program={program}
        onScheduleGenerated={handleScheduleGenerated}
      />
    </div>
  );
}

/**
 * Calculate overall status for a day based on its entries
 */
function calculateDayStatus(entries: ProgramEntry[], workout?: Workout): DiffStatus {
  if (!workout) return 'pending';

  if (entries.length === 0) {
    return workout.entries.length > 0 ? 'unplanned' : 'pending';
  }

  let hasMatch = false;
  let hasDiff = false;

  for (const entry of entries) {
    const actual = workout.entries.find((e) => e.program_node_id === entry.id);
    const planSnapshot = {
      target_sets: entry.target_sets,
      target_reps: entry.target_reps,
      target_load_kg: undefined,
      target_rpe: entry.target_rpe,
    };

    const diff = calculateDiff(planSnapshot, actual);

    if (diff.status === 'match') hasMatch = true;
    if (diff.status === 'diff') hasDiff = true;
  }

  if (hasDiff) return 'diff';
  if (hasMatch) return 'match';
  return 'pending';
}
