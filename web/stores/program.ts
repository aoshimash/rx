import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface ProgramSelectionState {
  selectedProgramId: string | null;
  setSelectedProgram: (programId: string) => void;
  clearSelectedProgram: () => void;
}

/**
 * Program selection store for current program
 * 
 * Persists the currently selected program across sessions
 */
export const useProgramStore = create<ProgramSelectionState>()(
  persist(
    (set) => ({
      selectedProgramId: null,
      setSelectedProgram: (programId) => set({ selectedProgramId: programId }),
      clearSelectedProgram: () => set({ selectedProgramId: null }),
    }),
    {
      name: 'optel-program-selection',
    }
  )
);
