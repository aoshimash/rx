'use client';

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { useCreateProgram } from '@/lib/hooks/usePrograms';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { ChoiceStep, type CreationMethod } from './create-program/ChoiceStep';
import { ImportStep } from './create-program/ImportStep';
import type { ProgramImport } from './create-program/importSchema';

type Step = { type: 'choice' } | { type: 'import' };

interface CreateProgramDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const STEP_TITLES: Record<Step['type'], string> = {
  choice: 'Create Program',
  import: 'Import Program',
};

const STEP_DESCRIPTIONS: Record<Step['type'], string> = {
  choice: 'How would you like to create it?',
  import: 'Paste or drop a JSON file exported from this app.',
};

function conflictError(err: unknown): string | undefined {
  if (
    err &&
    typeof err === 'object' &&
    'response' in err &&
    typeof (err as Record<string, unknown>).response === 'object'
  ) {
    const response = (err as { response: { status?: number } }).response;
    if (response?.status === 409) {
      return 'A program with this name already exists';
    }
  }
  return undefined;
}

export function CreateProgramDialog({ open, onOpenChange }: CreateProgramDialogProps) {
  const [step, setStep] = useState<Step>({ type: 'choice' });
  const [nameError, setNameError] = useState<string | undefined>();
  const createProgram = useCreateProgram();
  const router = useRouter();

  const handleOpenChange = (value: boolean) => {
    onOpenChange(value);
    if (!value) {
      setStep({ type: 'choice' });
      setNameError(undefined);
    }
  };

  const handleMethodSelect = (method: CreationMethod) => {
    if (method === 'import') setStep({ type: 'import' });
    else {
      handleOpenChange(false);
      router.push('/programs/new');
    }
    setNameError(undefined);
  };

  const handleImportSubmit = async (importData: ProgramImport) => {
    const { rx_version: _v, ...payload } = importData;
    try {
      await createProgram.mutateAsync(payload);
      handleOpenChange(false);
    } catch (err) {
      setNameError(conflictError(err) ?? 'Failed to import program');
    }
  };

  const isPending = createProgram.isPending;

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{STEP_TITLES[step.type]}</DialogTitle>
          <DialogDescription>{STEP_DESCRIPTIONS[step.type]}</DialogDescription>
        </DialogHeader>

        {step.type === 'choice' && <ChoiceStep onSelect={handleMethodSelect} />}
        {step.type === 'import' && (
          <ImportStep
            onBack={() => setStep({ type: 'choice' })}
            onSubmit={handleImportSubmit}
            isPending={isPending}
            nameError={nameError}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}
