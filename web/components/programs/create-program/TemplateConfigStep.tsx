'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { GenerateProgramRequest, ProgramTemplate } from '@/types/api';
import { useState } from 'react';

interface ExerciseConfig {
  name: string;
  hasPercentRm: boolean;
  weight: string;
  increment: string;
}

function buildExerciseConfigs(template: ProgramTemplate): ExerciseConfig[] {
  const seen = new Map<string, boolean>();
  for (const entry of template.entries ?? []) {
    if (!seen.has(entry.exercise_name)) {
      seen.set(entry.exercise_name, entry.percent_1rm != null);
    } else if (entry.percent_1rm != null) {
      seen.set(entry.exercise_name, true);
    }
  }
  return [...seen.entries()].map(([name, hasPercentRm]) => ({
    name,
    hasPercentRm,
    weight: '',
    increment: '2.5',
  }));
}

interface TemplateConfigStepProps {
  template: ProgramTemplate;
  onBack: () => void;
  onSubmit: (templateId: string, data: GenerateProgramRequest) => void;
  isPending: boolean;
  nameError?: string;
}

export function TemplateConfigStep({
  template,
  onBack,
  onSubmit,
  isPending,
  nameError,
}: TemplateConfigStepProps) {
  const [programName, setProgramName] = useState(template.name);
  const [exercises, setExercises] = useState<ExerciseConfig[]>(() =>
    buildExerciseConfigs(template)
  );

  const updateExercise = (idx: number, field: 'weight' | 'increment', value: string) => {
    setExercises((prev) => prev.map((ex, i) => (i === idx ? { ...ex, [field]: value } : ex)));
  };

  const allFilled = exercises.every((ex) => {
    const w = Number.parseFloat(ex.weight);
    return !Number.isNaN(w) && w > 0;
  });
  const hasExercises = exercises.length > 0;

  const handleSubmit = () => {
    const target_weights: Record<string, number> = {};
    const load_increments: Record<string, number> = {};
    for (const ex of exercises) {
      target_weights[ex.name] = Number.parseFloat(ex.weight);
      const inc = Number.parseFloat(ex.increment);
      if (!Number.isNaN(inc) && inc > 0) load_increments[ex.name] = inc;
    }
    onSubmit(template.id, {
      name: programName,
      target_weights,
      load_increments,
    });
  };

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="template-program-name">Program Name</Label>
        <Input
          id="template-program-name"
          value={programName}
          onChange={(e) => setProgramName(e.target.value)}
        />
        {nameError && <p className="text-sm text-destructive">{nameError}</p>}
      </div>

      {!hasExercises ? (
        <p className="text-sm text-muted-foreground">This template has no exercises.</p>
      ) : (
        <div className="space-y-3 max-h-72 overflow-y-auto">
          {exercises.map((ex, idx) => (
            <div key={ex.name} className="rounded-lg border p-3 space-y-2">
              <div className="flex items-center justify-between">
                <p className="font-medium text-sm">{ex.name}</p>
                <span className="text-xs text-muted-foreground">
                  {ex.hasPercentRm ? '% 1RM' : 'RPE only'}
                </span>
              </div>
              <div className="grid grid-cols-2 gap-2">
                <div className="space-y-1">
                  <Label className="text-xs">{ex.hasPercentRm ? '1RM (kg)' : 'Weight (kg)'}</Label>
                  <Input
                    type="number"
                    min={0.1}
                    step={0.5}
                    value={ex.weight}
                    onChange={(e) => updateExercise(idx, 'weight', e.target.value)}
                    placeholder="e.g., 140"
                  />
                </div>
                <div className="space-y-1">
                  <Label className="text-xs">Increment (kg)</Label>
                  <Input
                    type="number"
                    min={0.5}
                    step={0.5}
                    value={ex.increment}
                    onChange={(e) => updateExercise(idx, 'increment', e.target.value)}
                  />
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onBack}>
          Back
        </Button>
        <Button
          onClick={handleSubmit}
          disabled={isPending || !programName.trim() || !allFilled || !hasExercises}
        >
          {isPending ? 'Generating...' : 'Generate Program'}
        </Button>
      </div>
    </div>
  );
}
