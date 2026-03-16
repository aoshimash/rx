import type { PlanEntryCreate } from '@/types/api';

export interface SessionGroup {
  name: string; // e.g., "Day 1", "Day 2"
  date?: string; // ISO date string, optional
  exercises: PlanEntryCreate[];
}

export function sessionGroupsToEntries(sessions: SessionGroup[]): PlanEntryCreate[] {
  const entries: PlanEntryCreate[] = [];
  let order = 0;
  for (const session of sessions) {
    for (const ex of session.exercises) {
      entries.push({
        ...ex,
        order: order++,
        metadata: { ...ex.metadata, session: session.name, date: session.date },
      });
    }
  }
  return entries;
}

export function entriesToSessionGroups(entries: PlanEntryCreate[]): SessionGroup[] {
  if (entries.length === 0) {
    return [{ name: 'Day 1', exercises: [] }];
  }

  const sessionMap = new Map<string, { date?: string; exercises: PlanEntryCreate[] }>();
  const sessionOrder: string[] = [];

  for (const entry of entries) {
    const sessionName = (entry.metadata?.session as string) || 'Day 1';
    const date = entry.metadata?.date as string | undefined;
    if (!sessionMap.has(sessionName)) {
      sessionMap.set(sessionName, { date, exercises: [] });
      sessionOrder.push(sessionName);
    }
    sessionMap.get(sessionName)?.exercises.push({ ...entry });
  }

  return sessionOrder.map((name) => ({
    name,
    date: sessionMap.get(name)?.date,
    exercises: sessionMap.get(name)?.exercises || [],
  }));
}
