import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { Log } from '@/types/api';
import { ExerciseTable } from './ExerciseTable';

interface UnplannedLogsProps {
  logs: Log[];
}

export function UnplannedLogs({ logs }: UnplannedLogsProps) {
  if (logs.length === 0) {
    return null;
  }

  return (
    <Card className="mt-6">
      <CardHeader>
        <CardTitle>Unplanned Logs</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-6">
          {logs.map((log) => {
            const exercises = log.entries.map((entry) => ({
              exerciseName: entry.exercise_name,
              plan: null,
              actual: {
                sets: entry.sets,
                reps: entry.reps,
                load_kg: entry.load_kg,
                rpe: entry.rpe,
              },
              logEntry: entry,
            }));

            return (
              <div key={log.id} className="space-y-2">
                <div className="text-sm text-muted-foreground">
                  {new Date(log.performed_at).toLocaleString('en-US', {
                    month: 'short',
                    day: 'numeric',
                    hour: 'numeric',
                    minute: '2-digit',
                  })}
                  {log.notes && <span className="ml-2">- {log.notes}</span>}
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
