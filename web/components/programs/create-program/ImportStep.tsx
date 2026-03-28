'use client';

import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Upload } from 'lucide-react';
import { useCallback, useRef, useState } from 'react';
import { type ProgramImport, programImportSchema } from './importSchema';

interface ImportStepProps {
  onBack: () => void;
  onSubmit: (data: ProgramImport) => void;
  isPending: boolean;
  nameError?: string;
}

export function ImportStep({ onBack, onSubmit, isPending, nameError }: ImportStepProps) {
  const [text, setText] = useState('');
  const [parseError, setParseError] = useState<string | null>(null);
  const [parsed, setParsed] = useState<ProgramImport | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const parseAndStore = useCallback((raw: string) => {
    if (!raw.trim()) {
      setParseError(null);
      setParsed(null);
      return;
    }
    try {
      const json = JSON.parse(raw);
      const result = programImportSchema.safeParse(json);
      if (!result.success) {
        const firstError = result.error.issues[0];
        setParseError(
          `Invalid format: ${firstError?.path.join('.') ?? ''} — ${firstError?.message}`
        );
        setParsed(null);
        return;
      }
      setParseError(null);
      setParsed(result.data);
    } catch {
      setParseError('Invalid JSON: could not parse');
      setParsed(null);
    }
  }, []);

  const handleTextChange = (value: string) => {
    setText(value);
    parseAndStore(value);
  };

  const loadFile = useCallback(
    (file: File) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const content = e.target?.result as string;
        setText(content);
        parseAndStore(content);
      };
      reader.readAsText(file);
    },
    [parseAndStore]
  );

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    const file = e.dataTransfer.files[0];
    if (file) loadFile(file);
  };

  const canImport = parsed !== null && !isPending;

  return (
    <div className="space-y-3">
      <div
        onDragOver={(e) => {
          e.preventDefault();
          setIsDragging(true);
        }}
        onDragLeave={() => setIsDragging(false)}
        onDrop={handleDrop}
        className={`rounded-lg border-2 border-dashed p-6 text-center transition-colors ${
          isDragging ? 'border-primary bg-accent' : 'border-muted-foreground/30'
        }`}
      >
        <Upload className="h-6 w-6 mx-auto mb-2 text-muted-foreground" />
        <p className="text-sm text-muted-foreground">Drop JSON file here or paste below</p>
      </div>

      <Textarea
        value={text}
        onChange={(e) => handleTextChange(e.target.value)}
        placeholder={'{\n  "rx_version": "1",\n  "name": "...",\n  "sessions": []\n}'}
        className="font-mono text-xs h-32"
      />

      {(parseError || nameError) && (
        <p className="text-sm text-destructive">{nameError ?? parseError}</p>
      )}

      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() => fileInputRef.current?.click()}
          type="button"
        >
          Browse file
        </Button>
        <input
          ref={fileInputRef}
          type="file"
          accept=".json"
          className="hidden"
          onChange={(e) => {
            const f = e.target.files?.[0];
            if (f) loadFile(f);
          }}
        />
      </div>

      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onBack}>
          Back
        </Button>
        <Button
          onClick={() => {
            if (parsed) onSubmit(parsed);
          }}
          disabled={!canImport}
        >
          {isPending ? 'Importing...' : 'Import'}
        </Button>
      </div>
    </div>
  );
}
