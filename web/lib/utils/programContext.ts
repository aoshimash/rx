import type { Program } from '@/types/api';

/**
 * Generate program context from a program entry's metadata.
 *
 * Extracts week/day context from the entry's metadata field.
 *
 * @param entryId - The ID of the program entry
 * @param program - The program containing the entry
 * @returns Array of context strings derived from metadata (e.g., ["Week 1", "Day 1"])
 */
export function generateProgramContext(entryId: string, program: Program): string[] {
  const entry = program.entries?.find((e) => e.id === entryId);
  if (!entry?.metadata) return [];

  const context: string[] = [];
  if (entry.metadata.week !== undefined) context.push(String(entry.metadata.week));
  if (entry.metadata.day !== undefined) context.push(String(entry.metadata.day));
  return context;
}

/**
 * Format program context as a breadcrumb string
 *
 * @param context - Array of context strings
 * @param separator - Separator between levels (default: " > ")
 * @returns Formatted breadcrumb string (e.g., "Week 1 > Day 1")
 */
export function formatProgramContext(context: string[], separator = ' > '): string {
  return context.join(separator);
}
