'use client';

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { useGenerateProgram } from '@/lib/hooks/useProgramTemplates';
import { useCreateProgram } from '@/lib/hooks/usePrograms';
import type { GenerateProgramRequest, ProgramTemplate } from '@/types/api';
import { useState } from 'react';
import { ChoiceStep, type CreationMethod } from './create-program/ChoiceStep';
import { ImportStep } from './create-program/ImportStep';
import { ScratchStep } from './create-program/ScratchStep';
import { TemplateConfigStep } from './create-program/TemplateConfigStep';
import { TemplateSelectStep } from './create-program/TemplateSelectStep';
import type { ProgramImport } from './create-program/importSchema';

type Step =
  | { type: 'choice' }
  | { type: 'template-select' }
  | { type: 'template-config'; template: ProgramTemplate }
  | { type: 'import' }
  | { type: 'scratch' };

interface CreateProgramDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const STEP_TITLES: Record<Step['type'], string> = {
  choice: 'Create Program',
  'template-select': 'Create from Template',
  'template-config': 'Create from Template',
  import: 'Import Program',
  scratch: 'Create from Scratch',
};

const STEP_DESCRIPTIONS: Record<Step['type'], string> = {
  choice: 'How would you like to create it?',
  'template-select': 'Select a template to generate from.',
  'template-config': 'Enter your weights to calculate working loads.',
  import: 'Paste or drop a JSON file exported from this app.',
  scratch: 'Define the program name and session names.',
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
  const generateProgram = useGenerateProgram();

  const handleOpenChange = (value: boolean) => {
    onOpenChange(value);
    if (!value) {
      setStep({ type: 'choice' });
      setNameError(undefined);
    }
  };

  const handleMethodSelect = (method: CreationMethod) => {
    if (method === 'template') setStep({ type: 'template-select' });
    else if (method === 'import') setStep({ type: 'import' });
    else setStep({ type: 'scratch' });
    setNameError(undefined);
  };

  const handleScratchSubmit = async (data: {
    name: string;
    notes?: string;
    sessions: Parameters<typeof createProgram.mutateAsync>[0]['sessions'];
  }) => {
    try {
      await createProgram.mutateAsync(data);
      handleOpenChange(false);
    } catch (err) {
      setNameError(conflictError(err) ?? 'Failed to create program');
    }
  };

  const handleGenerateSubmit = async (templateId: string, data: GenerateProgramRequest) => {
    try {
      await generateProgram.mutateAsync({ id: templateId, data });
      handleOpenChange(false);
    } catch (err) {
      setNameError(conflictError(err) ?? 'Failed to generate program');
    }
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

  const isPending = createProgram.isPending || generateProgram.isPending;

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{STEP_TITLES[step.type]}</DialogTitle>
          <DialogDescription>{STEP_DESCRIPTIONS[step.type]}</DialogDescription>
        </DialogHeader>

        {step.type === 'choice' && <ChoiceStep onSelect={handleMethodSelect} />}
        {step.type === 'template-select' && (
          <TemplateSelectStep
            onBack={() => setStep({ type: 'choice' })}
            onSelect={(t) => {
              setStep({ type: 'template-config', template: t });
              setNameError(undefined);
            }}
          />
        )}
        {step.type === 'template-config' && (
          <TemplateConfigStep
            template={step.template}
            onBack={() => setStep({ type: 'template-select' })}
            onSubmit={handleGenerateSubmit}
            isPending={isPending}
            nameError={nameError}
          />
        )}
        {step.type === 'import' && (
          <ImportStep
            onBack={() => setStep({ type: 'choice' })}
            onSubmit={handleImportSubmit}
            isPending={isPending}
            nameError={nameError}
          />
        )}
        {step.type === 'scratch' && (
          <ScratchStep
            onBack={() => setStep({ type: 'choice' })}
            onSubmit={handleScratchSubmit}
            isPending={isPending}
            nameError={nameError}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}
