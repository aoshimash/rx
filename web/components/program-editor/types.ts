import type { ProgramEntryCreate } from '@/types/api';

export interface DayGroup {
  name: string;
  exercises: ProgramEntryCreate[];
}

export interface WeekGroup {
  name: string;
  days: DayGroup[];
}

export function weekGroupsToEntries(weekGroups: WeekGroup[]): ProgramEntryCreate[] {
  const entries: ProgramEntryCreate[] = [];
  let order = 0;
  for (const week of weekGroups) {
    for (const day of week.days) {
      for (const ex of day.exercises) {
        entries.push({
          ...ex,
          order: order++,
          metadata: { ...ex.metadata, week: week.name, day: day.name },
        });
      }
    }
  }
  return entries;
}

export function entriesToWeekGroups(entries: ProgramEntryCreate[]): WeekGroup[] {
  if (entries.length === 0) {
    return [{ name: 'Week 1', days: [] }];
  }

  const weekMap = new Map<string, Map<string, ProgramEntryCreate[]>>();
  const weekOrder: string[] = [];
  const dayOrderMap = new Map<string, string[]>();

  for (const entry of entries) {
    const weekName = (entry.metadata?.week as string) || 'Week 1';
    const dayName = (entry.metadata?.day as string) || 'Day 1';

    if (!weekMap.has(weekName)) {
      weekMap.set(weekName, new Map());
      weekOrder.push(weekName);
      dayOrderMap.set(weekName, []);
    }
    const weekDays = weekMap.get(weekName);
    if (!weekDays) continue;
    if (!weekDays.has(dayName)) {
      weekDays.set(dayName, []);
      dayOrderMap.get(weekName)?.push(dayName);
    }
    weekDays.get(dayName)?.push({ ...entry });
  }

  return weekOrder.map((weekName) => ({
    name: weekName,
    days: (dayOrderMap.get(weekName) || []).map((dayName) => ({
      name: dayName,
      exercises: weekMap.get(weekName)?.get(dayName) || [],
    })),
  }));
}
