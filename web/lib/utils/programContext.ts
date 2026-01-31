import type { Program, ProgramNode } from '@/types/api';

/**
 * Generate program context (hierarchy path) from a day node
 *
 * Traverses the program tree to find the path from root to the target day node.
 * Supports any hierarchy depth (e.g., cycle > block > week > day).
 *
 * @param dayNodeId - The ID of the day node
 * @param program - The program containing the node
 * @returns Array of node names from root to day (e.g., ["Week 1", "Day 1"])
 */
export function generateProgramContext(dayNodeId: string, program: Program): string[] {
  if (!program.root_nodes) return [];

  const path: string[] = [];

  function findPath(nodes: ProgramNode[], targetId: string): boolean {
    for (const node of nodes) {
      path.push(node.name);

      if (node.id === targetId) {
        return true;
      }

      if (node.children && node.children.length > 0) {
        if (findPath(node.children, targetId)) {
          return true;
        }
      }

      path.pop();
    }
    return false;
  }

  findPath(program.root_nodes, dayNodeId);
  return path;
}

/**
 * Format program context as a breadcrumb string
 *
 * @param context - Array of node names
 * @param separator - Separator between levels (default: " > ")
 * @returns Formatted breadcrumb string (e.g., "Week 1 > Day 1")
 */
export function formatProgramContext(context: string[], separator = ' > '): string {
  return context.join(separator);
}
