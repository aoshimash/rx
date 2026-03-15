'use client';

import { ScheduleModal } from '@/components/schedule/ScheduleModal';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { DiffStatus } from '@/lib/utils/diff';
import { calculateDiff } from '@/lib/utils/diff';
import { generatePlanContext } from '@/lib/utils/planContext';
import type { DaySchedule } from '@/lib/utils/schedule';
import { useScheduleStore } from '@/stores/schedule';
import type { Log, Plan, PlanEntry } from '@/types/api';
import { addWeeks, endOfWeek, format, startOfWeek } from 'date-fns';
import { Calendar as CalendarIcon, ChevronLeft, ChevronRight } from 'lucide-react';
import { useMemo, useState } from 'react';
import { DayAccordion } from './DayAccordion';
import { UnplannedLogs } from './UnplannedLogs';

interface WeekViewProps {
  plan: Plan | null;
  logs: Log[];
}

export function WeekView({ plan, logs }: WeekViewProps) {
  const [weekOffset, setWeekOffset] = useState(0);
  const [scheduleModalOpen, setScheduleModalOpen] = useState(false);
  const { getSchedule, setSchedule } = useScheduleStore();

  const { weekStart, weekEnd } = useMemo(() => {
    const today = new Date();
    const baseDate = addWeeks(today, weekOffset);
    return {
      weekStart: startOfWeek(baseDate, { weekStartsOn: 1 }),
      weekEnd: endOfWeek(baseDate, { weekStartsOn: 1 }),
    };
  }, [weekOffset]);

  const schedule = plan ? getSchedule(plan.id) : undefined;

  // Find logs linked to this plan
  const planLogs = useMemo(() => {
    if (!plan) return [];
    return logs.filter((l) => l.plan_id === plan.id);
  }, [plan, logs]);

  const days = useMemo(() => {
    if (!plan) return [];

    const entries = plan.entries || [];

    const weekGroups = new Map<string, PlanEntry[]>();
    const ungroupedEntries: PlanEntry[] = [];

    for (const entry of entries) {
      const weekKey = entry.metadata?.week !== undefined ? String(entry.metadata.week) : null;
      if (weekKey !== null) {
        if (!weekGroups.has(weekKey)) weekGroups.set(weekKey, []);
        weekGroups.get(weekKey)?.push(entry);
      } else {
        ungroupedEntries.push(entry);
      }
    }

    const firstWeekEntries =
      weekGroups.size > 0 ? weekGroups.values().next().value : ungroupedEntries;

    if (!firstWeekEntries || firstWeekEntries.length === 0) return [];

    const dayGroups = new Map<string, PlanEntry[]>();
    const ungroupedDayEntries: PlanEntry[] = [];

    for (const entry of firstWeekEntries) {
      const dayKey = entry.metadata?.day !== undefined ? String(entry.metadata.day) : null;
      if (dayKey !== null) {
        if (!dayGroups.has(dayKey)) dayGroups.set(dayKey, []);
        dayGroups.get(dayKey)?.push(entry);
      } else {
        ungroupedDayEntries.push(entry);
      }
    }

    const result = [...dayGroups.entries()].map(([dayName, dayEntries], idx) => {
      const scheduledDate = schedule?.find((s) => s.dayIndex === idx)?.date;

      // Find a log that covers this day's exercises (by exercise_name match)
      const dayExerciseNames = new Set(dayEntries.map((e) => e.exercise_name));
      const matchingLog = planLogs.find((l) =>
        l.entries.some((e) => dayExerciseNames.has(e.exercise_name))
      );

      const status = calculateDayStatus(dayEntries, matchingLog);

      const firstEntry = dayEntries[0];
      const planContext = firstEntry ? generatePlanContext(firstEntry.id, plan) : [];

      return {
        dayName,
        date: scheduledDate,
        planEntries: dayEntries,
        log: matchingLog,
        status,
        planContext,
      };
    });

    if (ungroupedDayEntries.length > 0) {
      const dayExerciseNames = new Set(ungroupedDayEntries.map((e) => e.exercise_name));
      const matchingLog = planLogs.find((l) =>
        l.entries.some((e) => dayExerciseNames.has(e.exercise_name))
      );
      result.push({
        dayName: 'Ungrouped',
        date: undefined,
        planEntries: ungroupedDayEntries,
        log: matchingLog,
        status: calculateDayStatus(ungroupedDayEntries, matchingLog),
        planContext: [],
      });
    }

    return result;
  }, [plan, planLogs, schedule]);

  // Get unplanned logs (not linked to any plan)
  const unplannedLogs = useMemo(() => {
    if (!plan) return logs;
    return logs.filter((l) => !l.plan_id);
  }, [plan, logs]);

  const handlePrevWeek = () => setWeekOffset((prev) => prev - 1);
  const handleNextWeek = () => setWeekOffset((prev) => prev + 1);
  const handleScheduleGenerated = (newSchedule: DaySchedule[]) => {
    if (plan) {
      setSchedule(plan.id, newSchedule);
    }
  };

  if (!plan) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center text-muted-foreground">
            No plan selected. Create or select a plan to view training week.
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
            <div className="text-center text-muted-foreground py-4">No days found in plan</div>
          )}
        </CardContent>
      </Card>

      <UnplannedLogs logs={unplannedLogs} />

      <ScheduleModal
        open={scheduleModalOpen}
        onOpenChange={setScheduleModalOpen}
        plan={plan}
        onScheduleGenerated={handleScheduleGenerated}
      />
    </div>
  );
}

function calculateDayStatus(entries: PlanEntry[], log?: Log): DiffStatus {
  if (!log) return 'pending';

  if (entries.length === 0) {
    return log.entries.length > 0 ? 'unplanned' : 'pending';
  }

  let hasMatch = false;
  let hasDiff = false;

  for (const entry of entries) {
    const actual = log.entries.find((e) => e.exercise_name === entry.exercise_name);
    const diff = calculateDiff(entry, actual);

    if (diff.status === 'match') hasMatch = true;
    if (diff.status === 'diff') hasDiff = true;
  }

  if (hasDiff) return 'diff';
  if (hasMatch) return 'match';
  return 'pending';
}
