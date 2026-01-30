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
import { X } from 'lucide-react';
import { ExerciseTable } from './ExerciseTable';

interface DayAccordionProps {
  days: ProgramNodeCreate[];
  availableExercises: Exercise[];
  onDayNameChange: (dayIndex: number, name: string) => void;
  onExerciseChange: (dayIndex: number, exerciseIndex: number, exerciseId: string) => void;
  onSetsChange: (dayIndex: number, exerciseIndex: number, value: number) => void;
  onRepsChange: (dayIndex: number, exerciseIndex: number, value: number) => void;
  onRpeChange: (dayIndex: number, exerciseIndex: number, value: number) => void;
  onRemoveExercise: (dayIndex: number, exerciseIndex: number) => void;
  onAddExercise: (dayIndex: number) => void;
  onRemoveDay: (dayIndex: number) => void;
}

/**
 * Collapsible day accordion for program editor
 */
export function DayAccordion({
  days,
  availableExercises,
  onDayNameChange,
  onExerciseChange,
  onSetsChange,
  onRepsChange,
  onRpeChange,
  onRemoveExercise,
  onAddExercise,
  onRemoveDay,
}: DayAccordionProps) {
  return (
    <Accordion type="multiple" className="w-full">
      {days.map((day, dayIdx) => (
        <AccordionItem key={dayIdx} value={`day-${dayIdx}`}>
          <AccordionTrigger className="hover:no-underline">
            <div className="flex items-center justify-between w-full pr-4">
              <span className="font-semibold">{day.name || `Day ${dayIdx + 1}`}</span>
              <Button
                variant="ghost"
                size="sm"
                onClick={(e) => {
                  e.stopPropagation();
                  onRemoveDay(dayIdx);
                }}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          </AccordionTrigger>
          <AccordionContent>
            <div className="space-y-4 pt-4">
              <div className="space-y-2">
                <Label>Day Name</Label>
                <Input
                  value={day.name}
                  onChange={(e) => onDayNameChange(dayIdx, e.target.value)}
                  placeholder="e.g., Push Day, Lower Body"
                />
              </div>

              <div className="space-y-2">
                <Label>Exercises</Label>
                <ExerciseTable
                  exercises={day.children || []}
                  availableExercises={availableExercises}
                  onExerciseChange={(exIdx, id) => onExerciseChange(dayIdx, exIdx, id)}
                  onSetsChange={(exIdx, value) => onSetsChange(dayIdx, exIdx, value)}
                  onRepsChange={(exIdx, value) => onRepsChange(dayIdx, exIdx, value)}
                  onRpeChange={(exIdx, value) => onRpeChange(dayIdx, exIdx, value)}
                  onRemove={(exIdx) => onRemoveExercise(dayIdx, exIdx)}
                  onAdd={() => onAddExercise(dayIdx)}
                />
              </div>
            </div>
          </AccordionContent>
        </AccordionItem>
      ))}
    </Accordion>
  );
}
