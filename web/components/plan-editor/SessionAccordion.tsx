import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
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

  return (
    <>
      <Accordion type="multiple" className="w-full space-y-3">
        {sessions.map((session, sessionIdx) => (
          <div key={sessionIdx} className="border rounded-lg px-4">
            <AccordionItem value={`session-${sessionIdx}`} className="border-0">
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
                      onChange={(exercises) => onChange(sessionIdx, { ...session, exercises })}
                    />
                  </div>
                </div>
              </AccordionContent>
            </AccordionItem>
          </div>
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
