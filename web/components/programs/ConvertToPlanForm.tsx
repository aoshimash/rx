'use client';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useConvertProgramToPlan } from '@/lib/hooks/usePrograms';
import type { Program } from '@/types/api';
import { useRouter } from 'next/navigation';
import { useMemo, useState } from 'react';

interface ConvertToPlanFormProps {
  program: Program;
}

export function ConvertToPlanForm({ program }: ConvertToPlanFormProps) {
  const router = useRouter();
  const convertMutation = useConvertProgramToPlan();

  const [planName, setPlanName] = useState(program.name);

  const exerciseNames = useMemo(() => {
    const names = new Set<string>();
    for (const entry of program.entries || []) {
      names.add(entry.exercise_name);
    }
    return Array.from(names);
  }, [program.entries]);

  const [targetWeights, setTargetWeights] = useState<Record<string, string>>(() =>
    Object.fromEntries(exerciseNames.map((name) => [name, '']))
  );

  const [loadIncrements, setLoadIncrements] = useState<Record<string, string>>(() =>
    Object.fromEntries(exerciseNames.map((name) => [name, '']))
  );

  const handleSubmit = async () => {
    const weights: Record<string, number> = {};
    for (const [name, value] of Object.entries(targetWeights)) {
      const num = Number(value);
      if (value && !Number.isNaN(num)) {
        weights[name] = num;
      }
    }

    const increments: Record<string, number> = {};
    for (const [name, value] of Object.entries(loadIncrements)) {
      const num = Number(value);
      if (value && !Number.isNaN(num)) {
        increments[name] = num;
      }
    }

    await convertMutation.mutateAsync({
      program_id: program.id,
      name: planName || undefined,
      target_weights: weights,
      load_increments: Object.keys(increments).length > 0 ? increments : undefined,
    });

    router.push('/plans');
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Plan Details</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="plan-name">Plan Name</Label>
            <Input
              id="plan-name"
              value={planName}
              onChange={(e) => setPlanName(e.target.value)}
              placeholder={program.name}
            />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Target Weights & Load Increments</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {exerciseNames.map((name) => (
            <div key={name} className="space-y-2 border rounded-lg p-4">
              <Label className="font-semibold">{name}</Label>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-1">
                  <Label className="text-xs text-muted-foreground">Target Weight (kg)</Label>
                  <Input
                    type="number"
                    step="0.5"
                    min={0}
                    value={targetWeights[name] ?? ''}
                    onChange={(e) =>
                      setTargetWeights((prev) => ({
                        ...prev,
                        [name]: e.target.value,
                      }))
                    }
                    placeholder="0"
                  />
                </div>
                <div className="space-y-1">
                  <Label className="text-xs text-muted-foreground">Load Increment (kg)</Label>
                  <Input
                    type="number"
                    step="0.5"
                    min={0}
                    value={loadIncrements[name] ?? ''}
                    onChange={(e) =>
                      setLoadIncrements((prev) => ({
                        ...prev,
                        [name]: e.target.value,
                      }))
                    }
                    placeholder="2.5"
                  />
                </div>
              </div>
            </div>
          ))}

          {exerciseNames.length === 0 && (
            <p className="text-sm text-muted-foreground">No exercises found in this program.</p>
          )}
        </CardContent>
      </Card>

      <div className="flex justify-end">
        <Button
          onClick={handleSubmit}
          disabled={convertMutation.isPending || exerciseNames.length === 0}
        >
          {convertMutation.isPending ? 'Converting...' : 'Convert to Plan'}
        </Button>
      </div>
    </div>
  );
}
