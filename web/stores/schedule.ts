import type { DaySchedule } from '@/lib/utils/schedule';
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface ScheduleState {
  schedules: Record<string, DaySchedule[]>; // planId -> schedule
  setSchedule: (planId: string, schedule: DaySchedule[]) => void;
  getSchedule: (planId: string) => DaySchedule[] | undefined;
  clearSchedule: (planId: string) => void;
}

export const useScheduleStore = create<ScheduleState>()(
  persist(
    (set, get) => ({
      schedules: {},
      setSchedule: (planId, schedule) =>
        set((state) => ({
          schedules: { ...state.schedules, [planId]: schedule },
        })),
      getSchedule: (planId) => get().schedules[planId],
      clearSchedule: (planId) =>
        set((state) => {
          const { [planId]: _, ...rest } = state.schedules;
          return { schedules: rest };
        }),
    }),
    {
      name: 'rx-schedule-storage',
    }
  )
);
