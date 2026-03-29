import type { LogEntry, ProgramSessionEntry } from '@/types/api';

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
  plan: ProgramSessionEntry | null | undefined,
  actual: LogEntry | null | undefined
): DiffResult {
  if (!plan) {
    return { status: 'unplanned', differences: [] };
  }

  if (!actual) {
    return { status: 'pending', differences: [] };
  }

  const diffs: string[] = [];

  const planSets = plan.fields?.sets as number | undefined;
  const actualSets = actual.fields?.sets as number | undefined;
  if (planSets !== undefined && actualSets !== undefined && planSets !== actualSets) {
    const delta = actualSets - planSets;
    diffs.push(`Sets ${delta > 0 ? '+' : ''}${delta}`);
  }

  const planReps = plan.fields?.reps as number | undefined;
  const actualReps = actual.fields?.reps as number | undefined;
  if (planReps !== undefined && actualReps !== undefined && planReps !== actualReps) {
    const delta = actualReps - planReps;
    diffs.push(`Reps ${delta > 0 ? '+' : ''}${delta}`);
  }

  const planLoadKg = plan.fields?.load_kg as number | undefined;
  const actualLoadKg = actual.fields?.load_kg as number | undefined;
  if (planLoadKg !== undefined && actualLoadKg !== undefined && planLoadKg !== actualLoadKg) {
    const delta = actualLoadKg - planLoadKg;
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
