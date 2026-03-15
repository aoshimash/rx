import type { Log, Plan } from '@/types/api';

interface ExportOptions {
  includeHeaders?: boolean;
  delimiter?: string;
}

function arrayToCSV(data: string[][], options: ExportOptions = {}): string {
  const { delimiter = ',' } = options;

  return data
    .map((row) =>
      row
        .map((cell) => {
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

function downloadCSV(filename: string, csv: string): void {
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

export function exportLogsToCSV(
  logs: Log[],
  plan?: Plan,
  _options: { scope: 'all' | 'current-week' } = { scope: 'all' }
): void {
  const headers = ['Date', 'Time', 'Plan', 'Exercise', 'Sets', 'Reps', 'Load (kg)', 'RPE', 'Notes'];

  const rows: string[][] = [headers];

  for (const log of logs) {
    const date = new Date(log.performed_at);
    const dateStr = date.toLocaleDateString('en-US');
    const timeStr = date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
    const planName = plan?.name || 'No Plan';

    for (const entry of log.entries) {
      rows.push([
        dateStr,
        timeStr,
        planName,
        entry.exercise_name,
        String(entry.sets ?? ''),
        String(entry.reps ?? ''),
        String(entry.load_kg ?? ''),
        String(entry.rpe ?? ''),
        entry.notes || '',
      ]);
    }
  }

  const csv = arrayToCSV(rows);
  const filename = `rx-logs-${new Date().toISOString().split('T')[0]}.csv`;
  downloadCSV(filename, csv);
}

export function exportPlanToCSV(plan: Plan): void {
  const headers = ['Week', 'Day', 'Exercise', 'Sets', 'Reps', 'Load (kg)', 'RPE', 'Notes'];
  const rows: string[][] = [headers];

  for (const entry of plan.entries || []) {
    rows.push([
      (entry.metadata?.week as string) || '',
      (entry.metadata?.day as string) || '',
      entry.exercise_name,
      String(entry.sets ?? ''),
      String(entry.reps ?? ''),
      String(entry.load_kg ?? ''),
      String(entry.rpe ?? ''),
      entry.notes || '',
    ]);
  }

  const csv = arrayToCSV(rows);
  const filename = `rx-plan-${plan.name.toLowerCase().replace(/\s+/g, '-')}-${new Date().toISOString().split('T')[0]}.csv`;
  downloadCSV(filename, csv);
}
