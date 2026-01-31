'use client';

import { ScheduleModal } from '@/components/schedule/ScheduleModal';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { DiffStatus } from '@/lib/utils/diff';
import { calculateDiff } from '@/lib/utils/diff';
import { generateProgramContext } from '@/lib/utils/programContext';
import type { DaySchedule } from '@/lib/utils/schedule';
import { useScheduleStore } from '@/stores/schedule';
import type { Program, ProgramNode, Workout } from '@/types/api';
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
 * - Collapsible days with exercise tables
 * - Unplanned workouts section
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

  // Build day data from program and workouts
  const days = useMemo(() => {
    if (!program) return [];

    // Get week nodes from program
    const weekNodes = program.root_nodes?.filter((node) => node.node_type === 'week') || [];

    // For now, show first week's days (multi-week support is future work)
    const firstWeek = weekNodes[0];
    if (!firstWeek) return [];

    const dayNodes = firstWeek.children?.filter((node) => node.node_type === 'day') || [];

    return dayNodes.map((dayNode, idx) => {
      // Find workout for this day (matching by program_node_id)
      const workout = workouts.find((w) => w.program_node_id === dayNode.id);

      // Get scheduled date if available
      const scheduledDate = schedule?.find((s) => s.dayIndex === idx)?.date;

      // Calculate overall day status
      const status = calculateDayStatus(dayNode, workout);

      // Get program context from workout or generate from program structure
      const programContext =
        workout?.program_context || generateProgramContext(dayNode.id, program);

      return {
        dayName: dayNode.name,
        date: scheduledDate,
        programNode: dayNode,
        workout,
        status,
        programContext,
      };
    });
  }, [program, workouts, schedule]);

  // Get unplanned workouts (no program_node_id or node not in current program)
  const unplannedWorkouts = useMemo(() => {
    if (!program) return workouts;

    const programNodeIds = new Set(
      program.root_nodes?.flatMap(
        (week) => week.children?.flatMap((day) => day.children?.map((ex) => ex.id)) || []
      ) || []
    );

    return workouts.filter((w) => !w.program_node_id || !programNodeIds.has(w.program_node_id));
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
 * Calculate overall status for a day based on its exercises
 */
function calculateDayStatus(dayNode: ProgramNode, workout?: Workout): DiffStatus {
  if (!workout) return 'pending';

  const exerciseNodes = dayNode.children?.filter((child) => child.node_type === 'exercise') || [];

  if (exerciseNodes.length === 0) {
    return workout.entries.length > 0 ? 'unplanned' : 'pending';
  }

  let hasMatch = false;
  let hasDiff = false;

  for (const node of exerciseNodes) {
    const entry = workout.entries.find((e) => e.program_node_id === node.id);
    const planSnapshot = {
      target_sets: node.target_sets,
      target_reps: node.target_reps,
      target_load_kg: undefined, // Load calculation deferred
      target_rpe: node.target_rpe,
    };

    const diff = calculateDiff(planSnapshot, entry);

    if (diff.status === 'match') hasMatch = true;
    if (diff.status === 'diff') hasDiff = true;
  }

  if (hasDiff) return 'diff';
  if (hasMatch) return 'match';
  return 'pending';
}
