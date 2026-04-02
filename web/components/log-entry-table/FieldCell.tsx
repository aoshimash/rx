'use client';

import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { videosApi } from '@/lib/api/videos';
import type { FieldDef } from '@/types/api';
import { CheckCircle2, Loader2, Upload, XCircle } from 'lucide-react';
import { useRef, useState } from 'react';
import { EditableCell } from './EditableCell';

interface FieldCellProps {
  fieldDef: FieldDef;
  value: unknown;
  onChange: (value: unknown) => void;
  onVideoUploaded?: (objectKey: string | undefined) => void;
  videoObjectKey?: string;
  /** When true, show plan value as read-only until clicked */
  defaultReadOnly?: boolean;
  planValue?: unknown;
}

export function FieldCell({
  fieldDef,
  value,
  onChange,
  onVideoUploaded,
  videoObjectKey,
  defaultReadOnly = false,
  planValue,
}: FieldCellProps) {
  if (fieldDef.type === 'number') {
    return (
      <EditableCell
        value={value as number | undefined}
        onChange={onChange}
        defaultReadOnly={defaultReadOnly}
        isEdited={value !== planValue}
        displayText={planValue != null ? String(planValue) : undefined}
        min={0}
        step={0.5}
      />
    );
  }

  if (fieldDef.type === 'text') {
    return (
      <Input
        value={(value as string | undefined) ?? ''}
        onChange={(e) => onChange(e.target.value || undefined)}
        placeholder={planValue != null ? String(planValue) : fieldDef.name}
        className="h-8 text-sm"
      />
    );
  }

  if (fieldDef.type === 'select') {
    const options = fieldDef.options ?? [];
    return (
      <Select
        value={(value as string | undefined) ?? ''}
        onValueChange={(v) => onChange(v || undefined)}
      >
        <SelectTrigger className="h-8 text-sm">
          <SelectValue placeholder="—" />
        </SelectTrigger>
        <SelectContent>
          {options.map((opt) => (
            <SelectItem key={opt} value={opt}>
              {opt}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    );
  }

  if (fieldDef.type === 'video') {
    return (
      <VideoUploadCell objectKey={videoObjectKey} onUploaded={onVideoUploaded ?? (() => {})} />
    );
  }

  return null;
}

// ============================================================================
// VideoUploadCell — handles file selection, upload, and status display
// ============================================================================

type UploadStatus = 'idle' | 'uploading' | 'done' | 'error';

function VideoUploadCell({
  objectKey,
  onUploaded,
}: {
  objectKey: string | undefined;
  onUploaded: (objectKey: string | undefined) => void;
}) {
  const [status, setStatus] = useState<UploadStatus>(objectKey ? 'done' : 'idle');
  const [error, setError] = useState<string | undefined>();
  const inputRef = useRef<HTMLInputElement>(null);

  const handleFile = async (file: File) => {
    setStatus('uploading');
    setError(undefined);
    try {
      const { upload_url, object_key } = await videosApi.getUploadUrl({
        content_type: file.type,
      });
      await videosApi.uploadToPresignedUrl(upload_url, file);
      onUploaded(object_key);
      setStatus('done');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed');
      setStatus('error');
      onUploaded(undefined);
    }
  };

  if (status === 'uploading') {
    return (
      <div className="flex items-center gap-1 text-xs text-muted-foreground">
        <Loader2 className="h-3 w-3 animate-spin" />
        Uploading…
      </div>
    );
  }

  if (status === 'done' && objectKey) {
    return (
      <div className="flex items-center gap-1 text-xs text-muted-foreground">
        <CheckCircle2 className="h-3 w-3 text-green-500" />
        <button
          type="button"
          className="underline hover:text-foreground cursor-pointer"
          onClick={() => {
            onUploaded(undefined);
            setStatus('idle');
            if (inputRef.current) inputRef.current.value = '';
          }}
        >
          Replace
        </button>
      </div>
    );
  }

  if (status === 'error') {
    return (
      <div className="flex flex-col gap-0.5">
        <div className="flex items-center gap-1 text-xs text-destructive">
          <XCircle className="h-3 w-3" />
          {error ?? 'Error'}
        </div>
        <button
          type="button"
          className="text-xs underline text-muted-foreground hover:text-foreground cursor-pointer"
          onClick={() => {
            setStatus('idle');
            setError(undefined);
          }}
        >
          Try again
        </button>
      </div>
    );
  }

  return (
    <label className="flex items-center gap-1 cursor-pointer text-xs text-muted-foreground hover:text-foreground">
      <Upload className="h-3 w-3" />
      Upload
      <input
        ref={inputRef}
        type="file"
        accept="video/*"
        className="sr-only"
        onChange={(e) => {
          const file = e.target.files?.[0];
          if (file) handleFile(file);
        }}
      />
    </label>
  );
}
