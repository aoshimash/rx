import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface PlanSelectionState {
  selectedPlanId: string | null;
  setSelectedPlan: (planId: string) => void;
  clearSelectedPlan: () => void;
}

export const usePlanStore = create<PlanSelectionState>()(
  persist(
    (set) => ({
      selectedPlanId: null,
      setSelectedPlan: (planId) => set({ selectedPlanId: planId }),
      clearSelectedPlan: () => set({ selectedPlanId: null }),
    }),
    {
      name: 'rx-plan-selection',
    }
  )
);
