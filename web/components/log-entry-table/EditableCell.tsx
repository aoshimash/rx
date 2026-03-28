'use client';

import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { useEffect, useRef, useState } from 'react';

interface EditableCellProps {
  value: number | undefined;
  onChange: (value: number | undefined) => void;
  /** When true, cell starts as read-only text and becomes input on click */
  defaultReadOnly?: boolean;
  /** Display text when in read-only mode (e.g., plan value) */
  displayText?: string;
  /** Whether the cell has been manually edited */
  isEdited?: boolean;
  placeholder?: string;
  min?: number;
  max?: number;
  step?: number;
  className?: string;
}

export function EditableCell({
  value,
  onChange,
  defaultReadOnly = false,
  displayText,
  isEdited = false,
  placeholder,
  min,
  max,
  step,
  className,
}: EditableCellProps) {
  const [isEditing, setIsEditing] = useState(!defaultReadOnly);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [isEditing]);

  if (defaultReadOnly && !isEditing && !isEdited) {
    return (
      <button
        type="button"
        className={cn(
          'w-full text-right tabular-nums text-muted-foreground hover:text-foreground cursor-pointer px-2 py-1 rounded hover:bg-muted/50',
          className
        )}
        onClick={() => setIsEditing(true)}
      >
        {displayText ?? '—'}
      </button>
    );
  }

  return (
    <Input
      ref={inputRef}
      type="number"
      value={value ?? ''}
      onChange={(e) => {
        const raw = e.target.value;
        onChange(raw === '' ? undefined : Number(raw));
      }}
      onBlur={() => {
        if (defaultReadOnly && !isEdited && value === undefined) {
          setIsEditing(false);
        }
      }}
      placeholder={placeholder}
      min={min}
      max={max}
      step={step}
      className={cn('h-8 w-full tabular-nums text-right', className)}
    />
  );
}
