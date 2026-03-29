import { Button } from '@/components/ui/button';
import { Plus } from 'lucide-react';

interface EmptyStateProps {
  onAddSession: () => void;
}

export function EmptyState({ onAddSession }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-16">
      <p className="text-lg font-medium mb-1">No sessions planned yet</p>
      <p className="text-sm text-muted-foreground mb-6">Add a session to get started</p>
      <Button onClick={onAddSession}>
        <Plus className="h-4 w-4 mr-2" />
        Add Session
      </Button>
    </div>
  );
}
