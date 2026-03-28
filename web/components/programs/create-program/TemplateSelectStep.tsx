'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Skeleton } from '@/components/ui/skeleton';
import { useProgramTemplates } from '@/lib/hooks/useProgramTemplates';
import type { ProgramTemplate } from '@/types/api';
import { Search } from 'lucide-react';
import { useState } from 'react';

interface TemplateSelectStepProps {
  onBack: () => void;
  onSelect: (template: ProgramTemplate) => void;
}

export function TemplateSelectStep({ onBack, onSelect }: TemplateSelectStepProps) {
  const { data, isLoading } = useProgramTemplates(false);
  const [query, setQuery] = useState('');

  const templates = (data?.data ?? []).filter((t) =>
    t.name.toLowerCase().includes(query.toLowerCase())
  );

  return (
    <div className="space-y-4">
      <div className="relative">
        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
        <Input
          className="pl-8"
          placeholder="Search templates..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
      </div>

      {isLoading ? (
        <div className="space-y-2">
          <Skeleton className="h-16" />
          <Skeleton className="h-16" />
        </div>
      ) : templates.length === 0 ? (
        <p className="text-sm text-muted-foreground text-center py-6">
          {query ? 'No templates match your search.' : 'No templates yet.'}
        </p>
      ) : (
        <div className="space-y-2 max-h-64 overflow-y-auto">
          {templates.map((template) => {
            const exerciseNames = [
              ...new Set((template.entries ?? []).map((e) => e.exercise_name)),
            ].slice(0, 3);
            const sessionCount = new Set(
              (template.entries ?? []).map((e) => e.metadata?.session as string)
            ).size;
            return (
              <button
                key={template.id}
                type="button"
                onClick={() => onSelect(template)}
                className="w-full flex flex-col gap-1 rounded-lg border p-3 text-left hover:bg-accent transition-colors"
              >
                <p className="font-medium">{template.name}</p>
                <p className="text-xs text-muted-foreground">
                  {sessionCount} session
                  {sessionCount !== 1 ? 's' : ''}
                  {exerciseNames.length > 0 &&
                    ` · ${exerciseNames.join(', ')}${(template.entries?.length ?? 0) > 3 ? '...' : ''}`}
                </p>
              </button>
            );
          })}
        </div>
      )}

      <div className="flex justify-start">
        <Button variant="outline" onClick={onBack}>
          Back
        </Button>
      </div>
    </div>
  );
}
