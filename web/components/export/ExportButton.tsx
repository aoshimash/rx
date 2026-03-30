import { Button } from '@/components/ui/button';
import { exportLogsToCSV } from '@/lib/utils/export';
import type { Log } from '@/types/api';
import { Download } from 'lucide-react';

interface ExportButtonProps {
  logs: Log[];
}

export function ExportButton({ logs }: ExportButtonProps) {
  return (
    <Button variant="outline" size="sm" onClick={() => exportLogsToCSV(logs)}>
      <Download className="h-4 w-4 mr-2" />
      Export CSV
    </Button>
  );
}
