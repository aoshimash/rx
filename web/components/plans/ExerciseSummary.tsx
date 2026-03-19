import type { PlanEntryCreate } from '@/types/api';

const SET_TYPE_LABELS: Record<string, string> = {
  top: 'Top',
  main: 'Main',
  backoff: 'Backoff',
};

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

function formatEntry(entry: PlanEntryCreate): string {
  const parts: string[] = [];
  if (entry.sets != null && entry.reps != null) {
    parts.push(`${entry.sets}x${entry.reps}`);
  }
  if (entry.load_kg != null) {
    parts.push(`${entry.load_kg}kg`);
  }
  if (entry.rpe != null) {
    parts.push(`@${entry.rpe}`);
  }
  return parts.join('  ');
}

interface ExerciseSummaryProps {
  exercises: PlanEntryCreate[];
}

export function ExerciseSummary({ exercises }: ExerciseSummaryProps) {
  const groups = groupByExercise(exercises);

  return (
    <div className="space-y-2 text-sm">
      {groups.map((group) => (
        <div key={group.name}>
          <div className="font-medium">{group.name}</div>
          <div className="ml-3 space-y-0.5 text-muted-foreground">
            {group.entries.map((entry, idx) => {
              const setType = entry.metadata?.set_type as string | undefined;
              const label = setType ? (SET_TYPE_LABELS[setType] ?? setType) : null;
              return (
                <div key={`${group.name}-${idx}`} className="flex gap-2">
                  {label && <span className="w-14 text-xs shrink-0">{label}</span>}
                  <span>{formatEntry(entry)}</span>
                </div>
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
}
