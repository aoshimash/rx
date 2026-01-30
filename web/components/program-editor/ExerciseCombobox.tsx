'use client';

import { useState } from 'react';
import { Check, ChevronsUpDown } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import type { Exercise } from '@/types/api';

interface ExerciseComboboxProps {
  exercises: Exercise[];
  value?: string;
  onSelect: (exerciseId: string) => void;
  placeholder?: string;
}

/**
 * Autocomplete combobox for exercise selection
 */
export function ExerciseCombobox({
  exercises,
  value,
  onSelect,
  placeholder = 'Select exercise...',
}: ExerciseComboboxProps) {
  const [open, setOpen] = useState(false);

  const selectedExercise = exercises.find((ex) => ex.id === value);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between"
        >
          {selectedExercise ? selectedExercise.name : placeholder}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[400px] p-0">
        <Command>
          <CommandInput placeholder="Search exercise..." />
          <CommandList>
            <CommandEmpty>No exercise found.</CommandEmpty>
            <CommandGroup>
              {exercises.map((exercise) => (
                <CommandItem
                  key={exercise.id}
                  value={exercise.name}
                  onSelect={() => {
                    onSelect(exercise.id);
                    setOpen(false);
                  }}
                >
                  <Check
                    className={cn(
                      'mr-2 h-4 w-4',
                      value === exercise.id ? 'opacity-100' : 'opacity-0'
                    )}
                  />
                  {exercise.name}
                  {exercise.muscle_groups && exercise.muscle_groups.length > 0 && (
                    <span className="ml-2 text-xs text-muted-foreground">
                      ({exercise.muscle_groups.join(', ')})
                    </span>
                  )}
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
