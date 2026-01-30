import type { Program, Workout } from '@/types/api';

/**
 * CSV export utility with UTF-8 BOM
 */

interface ExportOptions {
  includeHeaders?: boolean;
  delimiter?: string;
}

/**
 * Convert data to CSV format
 */
function arrayToCSV(data: string[][], options: ExportOptions = {}): string {
  const { delimiter = ',' } = options;

  return data
    .map((row) =>
      row
        .map((cell) => {
          // Escape quotes and wrap in quotes if contains delimiter, quotes, or newlines
          const cellStr = String(cell);
          if (cellStr.includes(delimiter) || cellStr.includes('"') || cellStr.includes('\n')) {
            return `"${cellStr.replace(/"/g, '""')}"`;
          }
          return cellStr;
        })
        .join(delimiter)
    )
    .join('\n');
}

/**
 * Download CSV file with UTF-8 BOM
 */
function downloadCSV(filename: string, csv: string): void {
  // Add UTF-8 BOM for Excel compatibility
  const BOM = '\uFEFF';
  const blob = new Blob([BOM + csv], { type: 'text/csv;charset=utf-8;' });
  const link = document.createElement('a');

  if (link.download !== undefined) {
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', filename);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }
}

/**
 * Export workouts to CSV
 */
export function exportWorkoutsToCSV(
  workouts: Workout[],
  program?: Program,
  _options: { scope: 'all' | 'current-week' } = { scope: 'all' }
): void {
  const headers = [
    'Date',
    'Time',
    'Program',
    'Exercise',
    'Sets',
    'Reps',
    'Load (kg)',
    'RPE',
    'Entry Type',
    'Notes',
  ];

  const rows: string[][] = [headers];

  for (const workout of workouts) {
    const date = new Date(workout.timestamp);
    const dateStr = date.toLocaleDateString('en-US');
    const timeStr = date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
    const programName = program?.name || 'No Program';

    for (const entry of workout.entries) {
      rows.push([
        dateStr,
        timeStr,
        programName,
        entry.display_name || 'Unknown Exercise',
        String(entry.sets),
        String(entry.reps),
        String(entry.load_kg),
        String(entry.rpe),
        entry.entry_type || '',
        entry.notes || '',
      ]);
    }
  }

  const csv = arrayToCSV(rows);
  const filename = `optel-workouts-${new Date().toISOString().split('T')[0]}.csv`;
  downloadCSV(filename, csv);
}

/**
 * Export program structure to CSV
 */
export function exportProgramToCSV(program: Program): void {
  const headers = ['Week', 'Day', 'Exercise', 'Sets', 'Reps', 'RPE', 'Notes'];
  const rows: string[][] = [headers];

  const weeks = program.root_nodes?.filter((node) => node.node_type === 'week') || [];

  for (const week of weeks) {
    const days = week.children?.filter((child) => child.node_type === 'day') || [];

    for (const day of days) {
      const exercises = day.children?.filter((child) => child.node_type === 'exercise') || [];

      for (const exercise of exercises) {
        rows.push([
          week.name,
          day.name,
          exercise.name,
          String(exercise.target_sets || ''),
          String(exercise.target_reps || ''),
          String(exercise.target_rpe || ''),
          exercise.notes || '',
        ]);
      }
    }
  }

  const csv = arrayToCSV(rows);
  const filename = `optel-program-${program.name.toLowerCase().replace(/\s+/g, '-')}-${new Date().toISOString().split('T')[0]}.csv`;
  downloadCSV(filename, csv);
}
