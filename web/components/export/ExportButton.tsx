import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { exportLogsToCSV, exportPlanToCSV } from '@/lib/utils/export';
import type { Log, Plan } from '@/types/api';
import { Download } from 'lucide-react';

interface ExportButtonProps {
  logs: Log[];
  plan: Plan | null;
}

export function ExportButton({ logs, plan }: ExportButtonProps) {
  const handleExportLogs = (scope: 'all' | 'current-week') => {
    exportLogsToCSV(logs, plan || undefined, { scope });
  };

  const handleExportPlan = () => {
    if (plan) {
      exportPlanToCSV(plan);
    }
  };

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
        <DropdownMenuItem onClick={() => handleExportLogs('all')}>All Logs (CSV)</DropdownMenuItem>
        <DropdownMenuItem onClick={() => handleExportLogs('current-week')}>
          Current Week (CSV)
        </DropdownMenuItem>
        {plan && (
          <>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleExportPlan}>Plan Structure (CSV)</DropdownMenuItem>
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
