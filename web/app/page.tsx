'use client';

import { WeekView } from '@/components/week-view/WeekView';
import { usePrograms } from '@/lib/hooks/usePrograms';
import { useWorkouts } from '@/lib/hooks/useWorkouts';
import { Skeleton } from '@/components/ui/skeleton';

export default function Home() {
  // Fetch programs and workouts
  const { data: programsData, isLoading: programsLoading } = usePrograms();
  const { data: workoutsData, isLoading: workoutsLoading } = useWorkouts();

  // For PoC, use first program (program selection is Phase 7)
  const program = programsData?.data[0] || null;
  const workouts = workoutsData?.data || [];

  if (programsLoading || workoutsLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <Skeleton className="h-[400px] w-full" />
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Training Week</h1>
        {program && (
          <p className="text-muted-foreground mt-1">
            Program: {program.name}
          </p>
        )}
      </div>
      <WeekView program={program} workouts={workouts} />
    </main>
  );
}

