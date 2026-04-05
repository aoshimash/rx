'use client';

import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { useEffect, useRef, useState } from 'react';

interface InlineEditProps {
  value: string;
  onSave: (value: string) => void;
  placeholder?: string;
  type?: 'text' | 'number' | 'date';
  className?: string;
  inputClassName?: string;
  emptyDisplay?: string;
}

export function InlineEdit({
  value,
  onSave,
  placeholder,
  type = 'text',
  className,
  inputClassName,
  emptyDisplay = '—',
}: InlineEditProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [draft, setDraft] = useState(value);
  const inputRef = useRef<HTMLInputElement>(null);
  const skipBlurRef = useRef(false);

  useEffect(() => {
    setDraft(value);
  }, [value]);

  useEffect(() => {
    if (isEditing) {
      inputRef.current?.focus();
      if (type !== 'date') {
        inputRef.current?.select();
      }
    }
  }, [isEditing, type]);

  const commit = () => {
    setIsEditing(false);
    const trimmed = draft.trim();
    if (trimmed !== value) {
      onSave(trimmed);
    }
  };

  const cancel = () => {
    setDraft(value);
    setIsEditing(false);
  };

  if (isEditing) {
    return (
      <Input
        ref={inputRef}
        type={type}
        value={draft}
        onChange={(e) => setDraft(e.target.value)}
        onBlur={() => {
          if (skipBlurRef.current) {
            skipBlurRef.current = false;
            return;
          }
          commit();
        }}
        onKeyDown={(e) => {
          if (e.key === 'Enter') {
            e.preventDefault();
            skipBlurRef.current = true;
            commit();
          }
          if (e.key === 'Escape') {
            e.preventDefault();
            skipBlurRef.current = true;
            cancel();
          }
        }}
        placeholder={placeholder}
        className={cn('h-7 px-1', inputClassName)}
      />
    );
  }

  const displayValue = value || emptyDisplay;

  return (
    <span
      className={cn(
        'cursor-pointer rounded px-1 py-0.5 hover:bg-muted transition-colors',
        !value && 'text-muted-foreground italic',
        className
      )}
      onClick={() => setIsEditing(true)}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          setIsEditing(true);
        }
      }}
      tabIndex={0}
      role="button"
    >
      {displayValue}
    </span>
  );
}
