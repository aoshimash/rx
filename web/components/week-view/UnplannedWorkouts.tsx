import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ExerciseTable } from './ExerciseTable';
import type { Workout } from '@/types/api';

interface UnplannedWorkoutsProps {
  workouts: Workout[];
}

/**
 * Section showing workouts without program association
 * 
 * Displays workouts that were not linked to any program node
 */
export function UnplannedWorkouts({ workouts }: UnplannedWorkoutsProps) {
  if (workouts.length === 0) {
    return null;
  }

  return (
    <Card className="mt-6">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <span>📝</span>
          <span>Unplanned Workouts</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-6">
          {workouts.map((workout) => {
            const exercises = workout.entries.map((entry) => ({
              exerciseName: entry.display_name || 'Unknown Exercise',
              plan: null,
              actual: {
                sets: entry.sets,
                reps: entry.reps,
                load: entry.load_kg,
                rpe: entry.rpe,
              },
              entry,
            }));

            return (
              <div key={workout.id} className="space-y-2">
                <div className="text-sm text-muted-foreground">
                  {new Date(workout.timestamp).toLocaleString('en-US', {
                    month: 'short',
                    day: 'numeric',
                    hour: 'numeric',
                    minute: '2-digit',
                  })}
                  {workout.notes && (
                    <span className="ml-2">- {workout.notes}</span>
                  )}
                </div>
                <ExerciseTable exercises={exercises} />
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}
