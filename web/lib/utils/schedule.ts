import { addDays, isWeekend, parseISO } from 'date-fns';

export interface ScheduleOptions {
  startDate: string;
  skipWeekends: boolean;
  avoidConsecutive: boolean;
}

export interface DaySchedule {
  dayIndex: number;
  date: string;
}

/**
 * Generate schedule dates for program days
 *
 * @param programDayCount - Number of training days in the program
 * @param options - Schedule configuration
 * @returns Array of day schedules with assigned dates
 */
export function generateSchedule(programDayCount: number, options: ScheduleOptions): DaySchedule[] {
  const { startDate, skipWeekends, avoidConsecutive } = options;
  const schedule: DaySchedule[] = [];
  let currentDate = parseISO(startDate);

  for (let dayIndex = 0; dayIndex < programDayCount; dayIndex++) {
    // Skip weekends if enabled
    while (skipWeekends && isWeekend(currentDate)) {
      currentDate = addDays(currentDate, 1);
    }

    schedule.push({
      dayIndex,
      date: currentDate.toISOString(),
    });

    // Move to next day
    currentDate = addDays(currentDate, 1);

    // Add rest day if avoiding consecutive days
    if (avoidConsecutive) {
      currentDate = addDays(currentDate, 1);
    }

    // Skip weekends after adding rest day
    if (skipWeekends) {
      while (isWeekend(currentDate)) {
        currentDate = addDays(currentDate, 1);
      }
    }
  }

  return schedule;
}
