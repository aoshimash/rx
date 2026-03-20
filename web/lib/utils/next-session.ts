import type { Log, Plan } from '@/types/api';

export interface PlanStatus {
  plan: Plan;
  isNext: boolean;
  completedCount: number;
  lastPerformedAt?: string;
}

/**
 * Group plans by program_id.
 * Plans without a program_id are each treated as their own group.
 */
export function groupPlansByProgram(plans: Plan[]): Map<string | null, Plan[]> {
  const groups = new Map<string | null, Plan[]>();
  for (const plan of plans) {
    const key = plan.program_id ?? null;
    const group = groups.get(key) ?? [];
    group.push(plan);
    groups.set(key, group);
  }
  return groups;
}

/**
 * Detect the next plan within a group of plans that share the same program_id.
 * Cycles through plans based on which was most recently logged.
 */
export function detectNextPlan(plans: Plan[], logs: Log[]): PlanStatus[] {
  if (plans.length === 0) return [];

  const planIds = new Set(plans.map((p) => p.id));
  const relevantLogs = logs
    .filter((log) => log.plan_id && planIds.has(log.plan_id))
    .sort((a, b) => new Date(b.performed_at).getTime() - new Date(a.performed_at).getTime());

  // Count completions and find last performed date per plan
  const planStats = new Map<string, { count: number; lastAt?: string }>();
  for (const log of relevantLogs) {
    if (!log.plan_id) continue;
    const existing = planStats.get(log.plan_id) || { count: 0 };
    existing.count++;
    if (!existing.lastAt) {
      existing.lastAt = log.performed_at;
    }
    planStats.set(log.plan_id, existing);
  }

  // Determine next plan index
  let nextIndex = 0;
  const lastLog = relevantLogs[0];
  if (lastLog?.plan_id) {
    const lastIndex = plans.findIndex((p) => p.id === lastLog.plan_id);
    if (lastIndex >= 0) {
      nextIndex = (lastIndex + 1) % plans.length;
    }
  }

  return plans.map((plan, index) => {
    const stats = planStats.get(plan.id);
    return {
      plan,
      isNext: index === nextIndex,
      completedCount: stats?.count ?? 0,
      lastPerformedAt: stats?.lastAt,
    };
  });
}

/**
 * Rotate plan statuses so that the "next" plan appears first.
 */
export function sortPlansByNext(statuses: PlanStatus[]): PlanStatus[] {
  if (statuses.length === 0) return [];
  const nextIndex = statuses.findIndex((s) => s.isNext);
  if (nextIndex <= 0) return statuses;
  return [...statuses.slice(nextIndex), ...statuses.slice(0, nextIndex)];
}
