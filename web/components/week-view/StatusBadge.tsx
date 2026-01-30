import { Badge } from '@/components/ui/badge';
import { type DiffStatus, getStatusIcon, getStatusVariant } from '@/lib/utils/diff';

interface StatusBadgeProps {
  status: DiffStatus;
  differences?: string[];
  className?: string;
}

/**
 * Status indicator badge for Plan vs Actual comparison
 *
 * Icons:
 * - ✓ (match): Plan and actual match
 * - ≠ (diff): Plan and actual differ
 * - ○ (pending): Planned but not executed yet
 * - 📝 (unplanned): No plan exists
 */
export function StatusBadge({ status, differences, className }: StatusBadgeProps) {
  const icon = getStatusIcon(status);
  const variant = getStatusVariant(status);

  return (
    <Badge variant={variant} className={className}>
      <span className="mr-1">{icon}</span>
      <span className="capitalize">{status}</span>
      {differences && differences.length > 0 && (
        <span className="ml-1 text-xs">({differences.join(', ')})</span>
      )}
    </Badge>
  );
}
