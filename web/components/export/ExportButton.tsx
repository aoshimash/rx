import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Download } from 'lucide-react';
import { exportWorkoutsToCSV, exportProgramToCSV } from '@/lib/utils/export';
import type { Workout, Program } from '@/types/api';

interface ExportButtonProps {
  workouts: Workout[];
  program: Program | null;
}

/**
 * Export button with options dropdown
 */
export function ExportButton({ workouts, program }: ExportButtonProps) {
  const handleExportWorkouts = (scope: 'all' | 'current-week') => {
    exportWorkoutsToCSV(workouts, program || undefined, { scope });
  };

  const handleExportProgram = () => {
    if (program) {
      exportProgramToCSV(program);
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
        <DropdownMenuItem onClick={() => handleExportWorkouts('all')}>
          All Workouts (CSV)
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => handleExportWorkouts('current-week')}>
          Current Week (CSV)
        </DropdownMenuItem>
        {program && (
          <>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleExportProgram}>
              Program Structure (CSV)
            </DropdownMenuItem>
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
