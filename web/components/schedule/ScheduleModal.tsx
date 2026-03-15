'use client';

import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { cn } from '@/lib/utils';
import type { DaySchedule } from '@/lib/utils/schedule';
import { type ScheduleOptions, generateSchedule } from '@/lib/utils/schedule';
import type { Program } from '@/types/api';
import { format } from 'date-fns';
import { CalendarIcon } from 'lucide-react';
import { useState } from 'react';

interface ScheduleModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  program: Program;
  onScheduleGenerated: (schedule: DaySchedule[]) => void;
}

/**
 * Modal for configuring schedule settings
 *
 * Allows setting:
 * - Start date
 * - Skip weekends option
 * - Avoid consecutive days option
 */
export function ScheduleModal({
  open,
  onOpenChange,
  program,
  onScheduleGenerated,
}: ScheduleModalProps) {
  const [startDate, setStartDate] = useState<Date>(new Date());
  const [skipWeekends, setSkipWeekends] = useState(true);
  const [avoidConsecutive, setAvoidConsecutive] = useState(false);

  // Count unique training days using composite week::day key to handle multi-week programs
  const entries = program.entries || [];
  const totalDays =
    new Set(
      entries
        .map((e) => {
          const w = e.metadata?.week;
          const d = e.metadata?.day;
          return w !== undefined && d !== undefined ? `${w}::${d}` : undefined;
        })
        .filter((k): k is string => k !== undefined)
    ).size || entries.length;

  const handleGenerate = () => {
    const options: ScheduleOptions = {
      startDate: startDate.toISOString(),
      skipWeekends,
      avoidConsecutive,
    };

    const schedule = generateSchedule(totalDays, options);
    onScheduleGenerated(schedule);
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Schedule Settings</DialogTitle>
        </DialogHeader>

        <div className="space-y-6">
          <div className="space-y-2">
            <Label>Start Date</Label>
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  className={cn(
                    'w-full justify-start text-left font-normal',
                    !startDate && 'text-muted-foreground'
                  )}
                >
                  <CalendarIcon className="mr-2 h-4 w-4" />
                  {startDate ? format(startDate, 'PPP') : <span>Pick a date</span>}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0">
                <Calendar
                  mode="single"
                  selected={startDate}
                  onSelect={(date) => date && setStartDate(date)}
                  initialFocus
                />
              </PopoverContent>
            </Popover>
          </div>

          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label>Skip Weekends</Label>
                <div className="text-sm text-muted-foreground">
                  Avoid scheduling workouts on Saturday/Sunday
                </div>
              </div>
              <input
                type="checkbox"
                checked={skipWeekends}
                onChange={(e) => setSkipWeekends(e.target.checked)}
                className="h-4 w-4"
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label>Avoid Consecutive Days</Label>
                <div className="text-sm text-muted-foreground">Add rest day between workouts</div>
              </div>
              <input
                type="checkbox"
                checked={avoidConsecutive}
                onChange={(e) => setAvoidConsecutive(e.target.checked)}
                className="h-4 w-4"
              />
            </div>
          </div>

          <div className="text-sm text-muted-foreground border-t pt-4">
            This will schedule {totalDays} training days starting from{' '}
            {format(startDate, 'MMM d, yyyy')}.
          </div>

          <div className="flex justify-end gap-2">
            <Button variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button onClick={handleGenerate}>Generate Schedule</Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
