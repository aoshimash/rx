import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { Exercise, ProgramNodeCreate } from '@/types/api';
import { Plus, X } from 'lucide-react';
import { DayAccordion } from './DayAccordion';

interface WeekAccordionProps {
  weeks: ProgramNodeCreate[];
  availableExercises: Exercise[];
  onWeekNameChange: (weekIndex: number, name: string) => void;
  onDayNameChange: (weekIndex: number, dayIndex: number, name: string) => void;
  onExerciseChange: (
    weekIndex: number,
    dayIndex: number,
    exerciseIndex: number,
    exerciseId: string
  ) => void;
  onSetsChange: (weekIndex: number, dayIndex: number, exerciseIndex: number, value: number) => void;
  onRepsChange: (weekIndex: number, dayIndex: number, exerciseIndex: number, value: number) => void;
  onRpeChange: (weekIndex: number, dayIndex: number, exerciseIndex: number, value: number) => void;
  onRemoveExercise: (weekIndex: number, dayIndex: number, exerciseIndex: number) => void;
  onAddExercise: (weekIndex: number, dayIndex: number) => void;
  onRemoveDay: (weekIndex: number, dayIndex: number) => void;
  onAddDay: (weekIndex: number) => void;
  onRemoveWeek: (weekIndex: number) => void;
}

/**
 * Collapsible week accordion for program editor
 */
export function WeekAccordion({
  weeks,
  availableExercises,
  onWeekNameChange,
  onDayNameChange,
  onExerciseChange,
  onSetsChange,
  onRepsChange,
  onRpeChange,
  onRemoveExercise,
  onAddExercise,
  onRemoveDay,
  onAddDay,
  onRemoveWeek,
}: WeekAccordionProps) {
  return (
    <Accordion type="multiple" className="w-full">
      {weeks.map((week, weekIdx) => (
        <AccordionItem key={weekIdx} value={`week-${weekIdx}`}>
          <AccordionTrigger className="hover:no-underline">
            <div className="flex items-center justify-between w-full pr-4">
              <span className="font-semibold">{week.name || `Week ${weekIdx + 1}`}</span>
              <Button
                variant="ghost"
                size="sm"
                onClick={(e) => {
                  e.stopPropagation();
                  onRemoveWeek(weekIdx);
                }}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          </AccordionTrigger>
          <AccordionContent>
            <div className="space-y-4 pt-4">
              <div className="space-y-2">
                <Label>Week Name</Label>
                <Input
                  value={week.name}
                  onChange={(e) => onWeekNameChange(weekIdx, e.target.value)}
                  placeholder="e.g., Week 1, Deload Week"
                />
              </div>

              <div className="space-y-2">
                <Label>Days</Label>
                <DayAccordion
                  days={week.children || []}
                  availableExercises={availableExercises}
                  onDayNameChange={(dayIdx, name) => onDayNameChange(weekIdx, dayIdx, name)}
                  onExerciseChange={(dayIdx, exIdx, id) =>
                    onExerciseChange(weekIdx, dayIdx, exIdx, id)
                  }
                  onSetsChange={(dayIdx, exIdx, value) =>
                    onSetsChange(weekIdx, dayIdx, exIdx, value)
                  }
                  onRepsChange={(dayIdx, exIdx, value) =>
                    onRepsChange(weekIdx, dayIdx, exIdx, value)
                  }
                  onRpeChange={(dayIdx, exIdx, value) => onRpeChange(weekIdx, dayIdx, exIdx, value)}
                  onRemoveExercise={(dayIdx, exIdx) => onRemoveExercise(weekIdx, dayIdx, exIdx)}
                  onAddExercise={(dayIdx) => onAddExercise(weekIdx, dayIdx)}
                  onRemoveDay={(dayIdx) => onRemoveDay(weekIdx, dayIdx)}
                />
                <Button variant="outline" onClick={() => onAddDay(weekIdx)} className="w-full">
                  <Plus className="h-4 w-4 mr-2" />
                  Add Day
                </Button>
              </div>
            </div>
          </AccordionContent>
        </AccordionItem>
      ))}
    </Accordion>
  );
}
