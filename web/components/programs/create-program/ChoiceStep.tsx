'use client';

import { FileJson, Pencil } from 'lucide-react';

export type CreationMethod = 'import' | 'scratch';

const options: {
  method: CreationMethod;
  icon: React.ReactNode;
  title: string;
  description: string;
}[] = [
  {
    method: 'import',
    icon: <FileJson className="h-6 w-6" />,
    title: 'Import',
    description: 'Import from JSON (shared by coach or another device)',
  },
  {
    method: 'scratch',
    icon: <Pencil className="h-6 w-6" />,
    title: 'From Scratch',
    description: 'Create manually with custom sessions',
  },
];

interface ChoiceStepProps {
  onSelect: (method: CreationMethod) => void;
}

export function ChoiceStep({ onSelect }: ChoiceStepProps) {
  return (
    <div className="grid gap-3">
      {options.map(({ method, icon, title, description }) => (
        <button
          key={method}
          type="button"
          onClick={() => onSelect(method)}
          className="flex items-center gap-4 rounded-lg border p-4 text-left hover:bg-accent transition-colors"
        >
          <div className="text-muted-foreground">{icon}</div>
          <div>
            <p className="font-medium">{title}</p>
            <p className="text-sm text-muted-foreground">{description}</p>
          </div>
        </button>
      ))}
    </div>
  );
}
