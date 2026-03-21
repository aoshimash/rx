import type { Log, Plan } from '@/types/api';

export interface PlanStatus {
  plan: Plan;
  isNext: boolean;
  completedCount: number;
  lastPerformedAt?: string;
}

/**
 * Group plans by cycle_id.
 * Plans without a cycle_id are each treated as their own independent group.
 */
export function groupPlansByCycle(plans: Plan[]): Map<string, Plan[]> {
  const groups = new Map<string, Plan[]>();
  for (const plan of plans) {
    const key = plan.cycle_id ?? `standalone:${plan.id}`;
    const group = groups.get(key) ?? [];
    group.push(plan);
    groups.set(key, group);
  }
  // Sort each group by created_at to preserve creation order (e.g., Day A, Day B, Day C)
  for (const group of groups.values()) {
    group.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
  }
  return groups;
}

/**
 * Detect the next plan within a group of plans that share the same cycle_id.
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

export interface ProgramGroup {
  groupKey: string;
  programName: string;
  isStandalone: boolean;
  statuses: PlanStatus[];
}

/**
 * Find the active program group key based on most recent log activity.
 * The active program is the one whose plans have the most recent log.
 * Falls back to the first group if no logs exist.
 */
export function findActiveProgram(groups: Map<string, Plan[]>, logs: Log[]): string | null {
  if (groups.size === 0) return null;

  const planToGroup = new Map<string, string>();
  for (const [groupKey, groupPlans] of groups) {
    for (const plan of groupPlans) {
      planToGroup.set(plan.id, groupKey);
    }
  }

  const mostRecentLog = logs
    .filter((log) => log.plan_id && planToGroup.has(log.plan_id))
    .sort((a, b) => new Date(b.performed_at).getTime() - new Date(a.performed_at).getTime())[0];

  if (mostRecentLog?.plan_id) {
    return planToGroup.get(mostRecentLog.plan_id) ?? null;
  }

  return groups.keys().next().value ?? null;
}

/**
 * Build PlanStatus arrays for all program groups with a single global NEXT plan.
 * Only the active program's plans participate in NEXT detection.
 * Groups are ordered: active program first, then remaining groups.
 */
export function buildGlobalPlanStatuses(plans: Plan[], logs: Log[]): ProgramGroup[] {
  const grouped = groupPlansByCycle(plans);
  const activeKey = findActiveProgram(grouped, logs);

  const result: ProgramGroup[] = [];

  for (const [groupKey, groupPlans] of grouped) {
    const isStandalone = groupKey.startsWith('standalone:');
    const isActive = groupKey === activeKey;

    let statuses: PlanStatus[];
    if (isActive) {
      statuses = sortPlansByNext(detectNextPlan(groupPlans, logs));
    } else {
      statuses = buildInactiveStatuses(groupPlans, logs);
    }

    const programName = isStandalone
      ? groupPlans[0]?.name || 'Unnamed'
      : groupPlans[0]?.name.split(' - ')[0] || 'Unnamed';

    result.push({ groupKey, programName, isStandalone, statuses });
  }

  // Sort: active program first
  result.sort((a, b) => {
    if (a.groupKey === activeKey) return -1;
    if (b.groupKey === activeKey) return 1;
    return 0;
  });

  return result;
}

/**
 * Build PlanStatus for plans in a non-active program group.
 * All plans get isNext: false, but still compute completedCount and lastPerformedAt.
 */
function buildInactiveStatuses(plans: Plan[], logs: Log[]): PlanStatus[] {
  const planIds = new Set(plans.map((p) => p.id));
  const relevantLogs = logs.filter((log) => log.plan_id && planIds.has(log.plan_id));

  const planStats = new Map<string, { count: number; lastAt?: string }>();
  for (const log of relevantLogs) {
    if (!log.plan_id) continue;
    const existing = planStats.get(log.plan_id) || { count: 0 };
    existing.count++;
    if (!existing.lastAt || log.performed_at > (existing.lastAt ?? '')) {
      existing.lastAt = log.performed_at;
    }
    planStats.set(log.plan_id, existing);
  }

  return plans.map((plan) => {
    const stats = planStats.get(plan.id);
    return {
      plan,
      isNext: false,
      completedCount: stats?.count ?? 0,
      lastPerformedAt: stats?.lastAt,
    };
  });
}
