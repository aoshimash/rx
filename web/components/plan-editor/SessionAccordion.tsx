import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { PlanEntryCreate } from '@/types/api';
import { X } from 'lucide-react';
import { useState } from 'react';
import { DeleteConfirmDialog } from './DeleteConfirmDialog';
import { ExerciseTable } from './ExerciseTable';
import type { SessionGroup } from './types';

interface SessionAccordionProps {
  sessions: SessionGroup[];
  onChange: (index: number, updated: SessionGroup) => void;
  onDelete: (index: number) => void;
}

export function SessionAccordion({ sessions, onChange, onDelete }: SessionAccordionProps) {
  const [deleteTarget, setDeleteTarget] = useState<number | null>(null);

  const handleExerciseNameChange = (sessionIdx: number, exIdx: number, name: string) => {
    const session = sessions[sessionIdx];
    if (!session) return;
    const exercises = [...session.exercises];
    const ex = exercises[exIdx];
    if (!ex) return;
    exercises[exIdx] = { ...ex, exercise_name: name };
    onChange(sessionIdx, { ...session, exercises });
  };

  const handleSetsChange = (sessionIdx: number, exIdx: number, value: number) => {
    const session = sessions[sessionIdx];
    if (!session) return;
    const exercises = [...session.exercises];
    const ex = exercises[exIdx];
    if (!ex) return;
    exercises[exIdx] = { ...ex, sets: value };
    onChange(sessionIdx, { ...session, exercises });
  };

  const handleRepsChange = (sessionIdx: number, exIdx: number, value: number) => {
    const session = sessions[sessionIdx];
    if (!session) return;
    const exercises = [...session.exercises];
    const ex = exercises[exIdx];
    if (!ex) return;
    exercises[exIdx] = { ...ex, reps: value };
    onChange(sessionIdx, { ...session, exercises });
  };

  const handleLoadKgChange = (sessionIdx: number, exIdx: number, value: number) => {
    const session = sessions[sessionIdx];
    if (!session) return;
    const exercises = [...session.exercises];
    const ex = exercises[exIdx];
    if (!ex) return;
    exercises[exIdx] = { ...ex, load_kg: value };
    onChange(sessionIdx, { ...session, exercises });
  };

  const handleRpeChange = (sessionIdx: number, exIdx: number, value: number) => {
    const session = sessions[sessionIdx];
    if (!session) return;
    const exercises = [...session.exercises];
    const ex = exercises[exIdx];
    if (!ex) return;
    exercises[exIdx] = { ...ex, rpe: value };
    onChange(sessionIdx, { ...session, exercises });
  };

  const handleRemoveExercise = (sessionIdx: number, exIdx: number) => {
    const session = sessions[sessionIdx];
    if (!session) return;
    const exercises = session.exercises
      .filter((_, idx) => idx !== exIdx)
      .map((e, idx) => ({ ...e, order: idx }));
    onChange(sessionIdx, { ...session, exercises });
  };

  const handleAddExercise = (sessionIdx: number) => {
    const session = sessions[sessionIdx];
    if (!session) return;
    const newExercise: PlanEntryCreate = {
      exercise_name: 'Exercise',
      order: session.exercises.length,
      sets: 3,
      reps: 10,
      rpe: 7,
    };
    onChange(sessionIdx, { ...session, exercises: [...session.exercises, newExercise] });
  };

  return (
    <>
      <Accordion type="multiple" className="w-full">
        {sessions.map((session, sessionIdx) => (
          <AccordionItem key={sessionIdx} value={`session-${sessionIdx}`}>
            <AccordionTrigger className="hover:no-underline">
              <div className="flex items-center justify-between w-full pr-4">
                <Input
                  value={session.name}
                  onChange={(e) => onChange(sessionIdx, { ...session, name: e.target.value })}
                  onClick={(e) => e.stopPropagation()}
                  placeholder="e.g., Block1 Week2 Day3, Week1 Day2"
                  className="font-semibold border-none shadow-none p-0 h-auto focus-visible:ring-0 bg-transparent"
                />
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={(e) => {
                    e.stopPropagation();
                    setDeleteTarget(sessionIdx);
                  }}
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
            </AccordionTrigger>
            <AccordionContent>
              <div className="space-y-4 pt-4">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label>Date (optional)</Label>
                    <Input
                      type="date"
                      value={session.date || ''}
                      onChange={(e) =>
                        onChange(sessionIdx, {
                          ...session,
                          date: e.target.value || undefined,
                        })
                      }
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label>Exercises</Label>
                  <ExerciseTable
                    exercises={session.exercises}
                    onExerciseNameChange={(idx, name) =>
                      handleExerciseNameChange(sessionIdx, idx, name)
                    }
                    onSetsChange={(idx, value) => handleSetsChange(sessionIdx, idx, value)}
                    onRepsChange={(idx, value) => handleRepsChange(sessionIdx, idx, value)}
                    onLoadKgChange={(idx, value) => handleLoadKgChange(sessionIdx, idx, value)}
                    onRpeChange={(idx, value) => handleRpeChange(sessionIdx, idx, value)}
                    onRemove={(idx) => handleRemoveExercise(sessionIdx, idx)}
                    onAdd={() => handleAddExercise(sessionIdx)}
                  />
                </div>
              </div>
            </AccordionContent>
          </AccordionItem>
        ))}
      </Accordion>

      <DeleteConfirmDialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null);
        }}
        onConfirm={() => {
          if (deleteTarget !== null) {
            onDelete(deleteTarget);
            setDeleteTarget(null);
          }
        }}
        title="Delete Session?"
        description="This will remove this session and all its exercises. This action cannot be undone."
      />
    </>
  );
}
