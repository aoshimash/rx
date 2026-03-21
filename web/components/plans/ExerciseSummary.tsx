import { Badge } from '@/components/ui/badge';
import type { PlanEntryCreate } from '@/types/api';

interface ExerciseGroup {
  name: string;
  entries: PlanEntryCreate[];
}

function groupByExercise(entries: PlanEntryCreate[]): ExerciseGroup[] {
  const groups: ExerciseGroup[] = [];
  const map = new Map<string, ExerciseGroup>();
  const sorted = [...entries].sort((a, b) => a.order - b.order);
  for (const entry of sorted) {
    if (!map.has(entry.exercise_name)) {
      const g: ExerciseGroup = { name: entry.exercise_name, entries: [] };
      groups.push(g);
      map.set(entry.exercise_name, g);
    }
    map.get(entry.exercise_name)?.entries.push(entry);
  }
  return groups;
}

function formatEntryText(entry: PlanEntryCreate): string {
  const parts: string[] = [];
  if (entry.rpe != null) parts.push(`RPE${entry.rpe}`);
  if (entry.reps != null) parts.push(`${entry.reps}reps`);
  if (entry.sets != null) parts.push(`${entry.sets}sets`);
  if (entry.load_kg != null) parts.push(`${entry.load_kg}kg`);
  return parts.join(' ');
}

interface ExerciseSummaryProps {
  exercises: PlanEntryCreate[];
}

export function ExerciseSummary({ exercises }: ExerciseSummaryProps) {
  const groups = groupByExercise(exercises);

  return (
    <div className="divide-y text-sm">
      {groups.map((group) => (
        <div key={group.name} className="flex items-baseline gap-3 py-1.5 first:pt-0 last:pb-0">
          <span className="w-40 shrink-0 font-medium truncate">{group.name}</span>
          <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-muted-foreground">
            {group.entries.map((entry, idx) => {
              const label = (entry.metadata?.label ?? entry.metadata?.set_type) as
                | string
                | undefined;
              return (
                <span key={`${group.name}-${idx}`} className="inline-flex items-center gap-1.5">
                  {label && <Badge variant="outline">{label}</Badge>}
                  {formatEntryText(entry)}
                </span>
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
}
