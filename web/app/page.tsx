'use client';

import { WeekView } from '@/components/week-view/WeekView';
import { WorkoutModal } from '@/components/workout-input/WorkoutModal';
import { Button } from '@/components/ui/button';
import { usePrograms } from '@/lib/hooks/usePrograms';
import { useWorkouts, useCreateWorkout } from '@/lib/hooks/useWorkouts';
import { Skeleton } from '@/components/ui/skeleton';
import { Plus } from 'lucide-react';
import { useState } from 'react';
import type { WorkoutEntryCreate } from '@/types/api';

export default function Home() {
  const [workoutModalOpen, setWorkoutModalOpen] = useState(false);

  // Fetch programs and workouts
  const { data: programsData, isLoading: programsLoading } = usePrograms();
  const { data: workoutsData, isLoading: workoutsLoading } = useWorkouts();
  const createWorkout = useCreateWorkout();

  // For PoC, use first program (program selection is Phase 7)
  const program = programsData?.data[0] || null;
  const workouts = workoutsData?.data || [];

  const handleSaveWorkout = async (entries: WorkoutEntryCreate[], notes: string) => {
    await createWorkout.mutateAsync({
      timestamp: new Date().toISOString(),
      notes,
      entries,
    });
  };

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
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Training Week</h1>
          {program && (
            <p className="text-muted-foreground mt-1">
              Program: {program.name}
            </p>
          )}
        </div>
        <Button onClick={() => setWorkoutModalOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Record Workout
        </Button>
      </div>

      <WeekView program={program} workouts={workouts} />

      <WorkoutModal
        open={workoutModalOpen}
        onOpenChange={setWorkoutModalOpen}
        onSave={handleSaveWorkout}
      />
    </main>
  );
}

