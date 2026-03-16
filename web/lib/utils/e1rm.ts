/**
 * Estimated 1RM calculation utilities.
 *
 * Brzycki formula: e1RM = load × (36 / (37 - reps))
 * Valid for reps < 37.
 *
 * RPE correction: RPE 10 = 0 reps in reserve, RPE 9 = 1 RIR, etc.
 * effectiveReps = reps + (10 - rpe)
 * e1RM = load × (36 / (37 - effectiveReps))
 */

export function calculateE1rm(load: number, reps: number, rpe?: number): number | null {
  if (load <= 0 || reps <= 0) return null;

  const effectiveReps = rpe !== undefined ? reps + (10 - rpe) : reps;

  if (effectiveReps >= 37) return null;

  return load * (36 / (37 - effectiveReps));
}
