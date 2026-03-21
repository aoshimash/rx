import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { exportLogsToCSV } from '@/lib/utils/export';
import type { Log } from '@/types/api';
import { Download } from 'lucide-react';

interface ExportButtonProps {
  logs: Log[];
  plan?: null;
}

export function ExportButton({ logs }: ExportButtonProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm">
          <Download className="h-4 w-4 mr-2" />
          Export
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuLabel>Export Data</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={() => exportLogsToCSV(logs, { scope: 'all' })}>
          All Logs (CSV)
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => exportLogsToCSV(logs, { scope: 'current-week' })}>
          Current Week (CSV)
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
