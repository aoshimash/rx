import type { DaySchedule } from '@/lib/utils/schedule';
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface ScheduleState {
  schedules: Record<string, DaySchedule[]>; // programId -> schedule
  setSchedule: (programId: string, schedule: DaySchedule[]) => void;
  getSchedule: (programId: string) => DaySchedule[] | undefined;
  clearSchedule: (programId: string) => void;
}

/**
 * Schedule store for client-side schedule persistence
 *
 * Stores schedule assignments per program
 */
export const useScheduleStore = create<ScheduleState>()(
  persist(
    (set, get) => ({
      schedules: {},
      setSchedule: (programId, schedule) =>
        set((state) => ({
          schedules: { ...state.schedules, [programId]: schedule },
        })),
      getSchedule: (programId) => get().schedules[programId],
      clearSchedule: (programId) =>
        set((state) => {
          const { [programId]: _, ...rest } = state.schedules;
          return { schedules: rest };
        }),
    }),
    {
      name: 'optel-schedule-storage',
    }
  )
);
