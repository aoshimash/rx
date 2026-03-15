import type { LogEntry, PlanEntry } from '@/types/api';

/**
 * Diff status indicator
 * - match: Plan and actual values match
 * - diff: Plan and actual values differ
 * - pending: Planned but no actual log yet
 * - unplanned: No plan exists (log was unplanned)
 */
export type DiffStatus = 'match' | 'diff' | 'pending' | 'unplanned';

export interface DiffResult {
  status: DiffStatus;
  differences: string[];
}

/**
 * Calculate diff between a plan entry and an actual log entry.
 *
 * Only Sets, Reps, and Load are considered for diff.
 * RPE differences do NOT count as diff.
 */
export function calculateDiff(
  plan: PlanEntry | null | undefined,
  actual: LogEntry | null | undefined
): DiffResult {
  if (!plan) {
    return { status: 'unplanned', differences: [] };
  }

  if (!actual) {
    return { status: 'pending', differences: [] };
  }

  const diffs: string[] = [];

  if (plan.sets !== undefined && actual.sets !== undefined && plan.sets !== actual.sets) {
    const delta = actual.sets - plan.sets;
    diffs.push(`Sets ${delta > 0 ? '+' : ''}${delta}`);
  }

  if (plan.reps !== undefined && actual.reps !== undefined && plan.reps !== actual.reps) {
    const delta = actual.reps - plan.reps;
    diffs.push(`Reps ${delta > 0 ? '+' : ''}${delta}`);
  }

  if (
    plan.load_kg !== undefined &&
    actual.load_kg !== undefined &&
    plan.load_kg !== actual.load_kg
  ) {
    const delta = actual.load_kg - plan.load_kg;
    diffs.push(`Load ${delta > 0 ? '+' : ''}${delta}kg`);
  }

  return {
    status: diffs.length > 0 ? 'diff' : 'match',
    differences: diffs,
  };
}

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

export function getStatusVariant(
  status: DiffStatus
): 'default' | 'secondary' | 'destructive' | 'outline' {
  switch (status) {
    case 'match':
      return 'default';
    case 'diff':
      return 'destructive';
    case 'pending':
      return 'secondary';
    case 'unplanned':
      return 'outline';
  }
}
