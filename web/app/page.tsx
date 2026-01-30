'use client';

import { WeekView } from '@/components/week-view/WeekView';
import { WorkoutModal } from '@/components/workout-input/WorkoutModal';
import { ExportButton } from '@/components/export/ExportButton';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { usePrograms, useProgram } from '@/lib/hooks/usePrograms';
import { useWorkouts, useCreateWorkout } from '@/lib/hooks/useWorkouts';
import { useProgramStore } from '@/stores/program';
import { Skeleton } from '@/components/ui/skeleton';
import { Plus } from 'lucide-react';
import { useState, useMemo } from 'react';
import Link from 'next/link';
import type { WorkoutEntryCreate, ProgramNode } from '@/types/api';

export default function Home() {
  const [workoutModalOpen, setWorkoutModalOpen] = useState(false);
  const [selectedDayNode, setSelectedDayNode] = useState<ProgramNode | undefined>();
  const { selectedProgramId, setSelectedProgram } = useProgramStore();

  // Fetch programs and selected program details
  const { data: programsData, isLoading: programsLoading } = usePrograms();
  const { data: selectedProgram, isLoading: selectedProgramLoading } = useProgram(
    selectedProgramId
  );
  const { data: workoutsData, isLoading: workoutsLoading } = useWorkouts();
  const createWorkout = useCreateWorkout();

  // Use selected program or fall back to first program
  const programs = programsData?.data || [];
  const program = selectedProgram || (programs.length > 0 ? programs[0] : null);
  const workouts = workoutsData?.data || [];

  // Get first day node for quick workout recording
  const firstDayNode = useMemo(() => {
    if (!program) return undefined;
    const firstWeek = program.root_nodes?.find((n) => n.node_type === 'week');
    return firstWeek?.children?.find((n) => n.node_type === 'day');
  }, [program]);

  const handleSaveWorkout = async (entries: WorkoutEntryCreate[], notes: string) => {
    await createWorkout.mutateAsync({
      timestamp: new Date().toISOString(),
      notes,
      entries,
      program_node_id: selectedDayNode?.id,
    });
  };

  const handleProgramChange = (programId: string) => {
    setSelectedProgram(programId);
  };

  const handleOpenWorkoutModal = () => {
    setSelectedDayNode(firstDayNode);
    setWorkoutModalOpen(true);
  };

  if (programsLoading || selectedProgramLoading || workoutsLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-12 w-[300px]" />
        <Skeleton className="h-[400px] w-full" />
      </main>
    );
  }

  if (programs.length === 0) {
    return (
      <main className="container mx-auto p-6">
        <div className="text-center py-12">
          <h1 className="text-3xl font-bold mb-4">Welcome to OPTel Workout</h1>
          <p className="text-muted-foreground mb-6">
            Create your first training program to get started.
          </p>
          <Link href="/programs/new">
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              Create Program
            </Button>
          </Link>
        </div>
      </main>
    );
  }

  return (
    <main className="container mx-auto p-6">
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <h1 className="text-3xl font-bold">Training Week</h1>
          {programs.length > 1 && (
            <Select value={program?.id || ''} onValueChange={handleProgramChange}>
              <SelectTrigger className="w-[250px]">
                <SelectValue placeholder="Select program" />
              </SelectTrigger>
              <SelectContent>
                {programs.map((p) => (
                  <SelectItem key={p.id} value={p.id}>
                    {p.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
          {programs.length === 1 && program && (
            <p className="text-muted-foreground">Program: {program.name}</p>
          )}
        </div>
        <div className="flex items-center gap-2">
          <ExportButton workouts={workouts} program={program} />
          <Button onClick={handleOpenWorkoutModal}>
            <Plus className="h-4 w-4 mr-2" />
            Record Workout
          </Button>
        </div>
      </div>

      <WeekView program={program} workouts={workouts} />

      <WorkoutModal
        open={workoutModalOpen}
        onOpenChange={setWorkoutModalOpen}
        dayNode={selectedDayNode}
        onSave={handleSaveWorkout}
      />
    </main>
  );
}

