import { Badge } from '@/components/ui/badge';
import type { EntryType } from '@/types/api';

interface EntryTypeBadgeProps {
  entryType?: EntryType;
  className?: string;
}

/**
 * Get badge variant based on entry type
 *
 * Common entry types have specific colors:
 * - Top: primary (most important set)
 * - Main: secondary (main working sets)
 * - Backoff: outline (lighter weight sets)
 * - Others: default outline style
 */
function getVariant(
  entryType: string
): 'default' | 'secondary' | 'outline' | 'destructive' | 'ghost' | 'link' {
  const normalized = entryType.toLowerCase();
  switch (normalized) {
    case 'top':
      return 'default';
    case 'main':
      return 'secondary';
    case 'backoff':
      return 'outline';
    default:
      return 'outline';
  }
}

/**
 * Badge component for displaying entry type
 *
 * Displays the entry type with color-coded variants for common types.
 * Returns null if entry type is null/undefined/empty.
 */
export function EntryTypeBadge({ entryType, className }: EntryTypeBadgeProps) {
  if (!entryType) {
    return null;
  }

  return (
    <Badge variant={getVariant(entryType)} className={className}>
      {entryType}
    </Badge>
  );
}
