'use client';

import { Button } from '@/components/ui/button';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { cn } from '@/lib/utils';
import type { EntryType } from '@/types/api';
import { Check, ChevronsUpDown, X } from 'lucide-react';
import { useState } from 'react';

/** Default entry type suggestions */
const DEFAULT_SUGGESTIONS = ['Top', 'Main', 'Backoff'];

interface EntryTypeComboboxProps {
  value: EntryType;
  onChange: (value: EntryType) => void;
  className?: string;
}

/**
 * Combobox for entry type selection
 *
 * Features:
 * - Default suggestions (Top, Main, Backoff)
 * - Free text input for custom types
 * - Support for null/empty entry type
 */
export function EntryTypeCombobox({ value, onChange, className }: EntryTypeComboboxProps) {
  const [open, setOpen] = useState(false);
  const [inputValue, setInputValue] = useState('');

  const handleSelect = (selectedValue: string) => {
    onChange(selectedValue || null);
    setOpen(false);
    setInputValue('');
  };

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation();
    onChange(null);
    setInputValue('');
  };

  // Filter suggestions based on input
  const filteredSuggestions = DEFAULT_SUGGESTIONS.filter((suggestion) =>
    suggestion.toLowerCase().includes(inputValue.toLowerCase())
  );

  // Show custom input option if input doesn't match any suggestion exactly
  const showCustomOption =
    inputValue &&
    !DEFAULT_SUGGESTIONS.some((s) => s.toLowerCase() === inputValue.toLowerCase()) &&
    inputValue !== value;

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className={cn('w-full justify-between', className)}
        >
          <span className={cn(!value && 'text-muted-foreground')}>{value || 'Select type...'}</span>
          <div className="flex items-center gap-1">
            {value && (
              <X className="h-4 w-4 shrink-0 opacity-50 hover:opacity-100" onClick={handleClear} />
            )}
            <ChevronsUpDown className="h-4 w-4 shrink-0 opacity-50" />
          </div>
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0" align="start">
        <Command>
          <CommandInput
            placeholder="Search or enter custom..."
            value={inputValue}
            onValueChange={setInputValue}
          />
          <CommandList>
            <CommandEmpty>
              {inputValue ? (
                <button
                  type="button"
                  className="w-full px-2 py-1.5 text-sm text-left hover:bg-accent rounded"
                  onClick={() => handleSelect(inputValue)}
                >
                  Use &quot;{inputValue}&quot;
                </button>
              ) : (
                'Type to add custom entry type'
              )}
            </CommandEmpty>
            <CommandGroup heading="Suggestions">
              {filteredSuggestions.map((suggestion) => (
                <CommandItem
                  key={suggestion}
                  value={suggestion}
                  onSelect={() => handleSelect(suggestion)}
                >
                  <Check
                    className={cn(
                      'mr-2 h-4 w-4',
                      value?.toLowerCase() === suggestion.toLowerCase()
                        ? 'opacity-100'
                        : 'opacity-0'
                    )}
                  />
                  {suggestion}
                </CommandItem>
              ))}
              {showCustomOption && (
                <CommandItem value={inputValue} onSelect={() => handleSelect(inputValue)}>
                  <Check className="mr-2 h-4 w-4 opacity-0" />
                  Use &quot;{inputValue}&quot;
                </CommandItem>
              )}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
