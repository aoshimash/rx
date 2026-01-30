import { Button } from '@/components/ui/button';
import { Plus } from 'lucide-react';

interface AddExerciseButtonProps {
  onClick: () => void;
}

/**
 * Button to add unplanned exercises to workout
 */
export function AddExerciseButton({ onClick }: AddExerciseButtonProps) {
  return (
    <Button variant="outline" onClick={onClick} className="w-full">
      <Plus className="h-4 w-4 mr-2" />
      Add Exercise
    </Button>
  );
}
