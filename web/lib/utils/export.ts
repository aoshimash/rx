import type { Log } from '@/types/api';

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
  _options: { scope: 'all' | 'current-week' } = { scope: 'all' }
): void {
  const headers = [
    'Date',
    'Time',
    'Session',
    'Exercise',
    'Sets',
    'Reps',
    'Load (kg)',
    'RPE',
    'Notes',
  ];

  const rows: string[][] = [headers];

  for (const log of logs) {
    const date = new Date(log.performed_at);
    const dateStr = date.toLocaleDateString('en-US');
    const timeStr = date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
    const sessionName = log.session_name || '';

    for (const entry of log.entries) {
      rows.push([
        dateStr,
        timeStr,
        sessionName,
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
