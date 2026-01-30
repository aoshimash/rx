import type { PlanSnapshot, WorkoutEntry } from '@/types/api';

/**
 * Diff status indicator
 * - match: Plan and actual values match (✓)
 * - diff: Plan and actual values differ (≠)
 * - pending: Planned but no actual workout yet (○)
 * - unplanned: No plan exists (workout was unplanned)
 */
export type DiffStatus = 'match' | 'diff' | 'pending' | 'unplanned';

/**
 * Diff calculation result
 */
export interface DiffResult {
  status: DiffStatus;
  differences: string[];
}

/**
 * Calculate diff between plan and actual workout
 * 
 * Per spec: Only Sets, Reps, and Load are considered for diff.
 * RPE differences do NOT count as diff.
 */
export function calculateDiff(
  plan: PlanSnapshot | null | undefined,
  actual: WorkoutEntry | null | undefined
): DiffResult {
  // No plan = unplanned workout
  if (!plan) {
    return { status: 'unplanned', differences: [] };
  }

  // Plan exists but no actual workout = pending
  if (!actual) {
    return { status: 'pending', differences: [] };
  }

  const diffs: string[] = [];

  // Compare sets
  if (plan.target_sets !== undefined && plan.target_sets !== actual.sets) {
    const delta = actual.sets - plan.target_sets;
    diffs.push(`Sets ${delta > 0 ? '+' : ''}${delta}`);
  }

  // Compare reps
  if (plan.target_reps !== undefined && plan.target_reps !== actual.reps) {
    const delta = actual.reps - plan.target_reps;
    diffs.push(`Reps ${delta > 0 ? '+' : ''}${delta}`);
  }

  // Compare load
  if (plan.target_load_kg !== undefined && plan.target_load_kg !== actual.load_kg) {
    const delta = actual.load_kg - plan.target_load_kg;
    diffs.push(`Load ${delta > 0 ? '+' : ''}${delta}kg`);
  }

  // Note: RPE is NOT included in diff calculation per spec

  return {
    status: diffs.length > 0 ? 'diff' : 'match',
    differences: diffs,
  };
}

/**
 * Get status icon for diff result
 */
export function getStatusIcon(status: DiffStatus): string {
  switch (status) {
    case 'match':
      return '✓';
    case 'diff':
      return '≠';
    case 'pending':
      return '○';
    case 'unplanned':
      return '📝';
  }
}

/**
 * Get status color variant for Badge component
 */
export function getStatusVariant(
  status: DiffStatus
): 'default' | 'secondary' | 'destructive' | 'outline' {
  switch (status) {
    case 'match':
      return 'default'; // green
    case 'diff':
      return 'destructive'; // red/yellow
    case 'pending':
      return 'secondary'; // gray
    case 'unplanned':
      return 'outline'; // outlined
  }
}
