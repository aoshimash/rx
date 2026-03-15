import type { Plan } from '@/types/api';

/**
 * Generate plan context from a plan entry's metadata.
 *
 * @param entryId - The ID of the plan entry
 * @param plan - The plan containing the entry
 * @returns Array of context strings derived from metadata (e.g., ["Week 1", "Day 1"])
 */
export function generatePlanContext(entryId: string, plan: Plan): string[] {
  const entry = plan.entries?.find((e) => e.id === entryId);
  if (!entry?.metadata) return [];

  const context: string[] = [];
  if (entry.metadata.week !== undefined) context.push(String(entry.metadata.week));
  if (entry.metadata.day !== undefined) context.push(String(entry.metadata.day));
  return context;
}

/**
 * Format plan context as a breadcrumb string
 */
export function formatPlanContext(context: string[], separator = ' > '): string {
  return context.join(separator);
}
