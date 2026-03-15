import type { ProgramEntryCreate } from '@/types/api';
import { z } from 'zod';

// ============================================================================
// Exercise Forms
// ============================================================================

export const exerciseCreateSchema = z.object({
  name: z.string().min(1, 'Name is required').max(200, 'Name too long'),
  description: z.string().max(2000, 'Description too long').optional(),
  aliases: z.array(z.string()).optional(),
  muscle_groups: z.array(z.string()).optional(),
  load_increment: z.number().positive('Must be positive').optional(),
});

export type ExerciseFormData = z.infer<typeof exerciseCreateSchema>;

// ============================================================================
// Workout Entry Forms
// ============================================================================

export const workoutEntrySchema = z.object({
  exercise_id: z.string().uuid('Invalid exercise ID'),
  display_name: z.string().max(200, 'Display name too long').optional(),
  entry_type: z.string().max(50, 'Entry type too long').nullable().optional(),
  sets: z.number().int().min(1, 'Must be at least 1'),
  reps: z.number().int().min(1, 'Must be at least 1'),
  load_kg: z.number().min(0, 'Must be non-negative'),
  rpe: z.number().int().min(1, 'Must be 1-10').max(10, 'Must be 1-10'),
  notes: z.string().max(2000, 'Notes too long').optional(),
  program_node_id: z.string().uuid().optional(),
});

export type WorkoutEntryFormData = z.infer<typeof workoutEntrySchema>;

// ============================================================================
// Workout Forms
// ============================================================================

export const workoutCreateSchema = z.object({
  timestamp: z.string().datetime('Invalid datetime format'),
  session_start: z.string().datetime().optional(),
  session_end: z.string().datetime().optional(),
  body_weight_kg: z.number().positive('Must be positive').optional(),
  fatigue_level: z.number().int().min(1).max(5).optional(),
  sleep_hours: z.number().min(0).max(24).optional(),
  condition_notes: z.string().max(2000, 'Notes too long').optional(),
  program_node_id: z.string().uuid().optional(),
  program_context: z.array(z.string()).optional(),
  notes: z.string().max(5000, 'Notes too long').optional(),
  entries: z.array(workoutEntrySchema).min(1, 'At least one entry required'),
});

export type WorkoutFormData = z.infer<typeof workoutCreateSchema>;

// ============================================================================
// Program Entry Forms
// ============================================================================

export const programEntrySchema: z.ZodType<ProgramEntryCreate> = z.object({
  name: z.string().min(1, 'Name is required').max(200, 'Name too long'),
  order: z.number().int().min(0, 'Order must be non-negative'),
  metadata: z.record(z.string(), z.unknown()).optional(),
  exercise_id: z.string().uuid('Invalid exercise ID').optional(),
  target_sets: z.number().int().min(1, 'Must be at least 1').optional(),
  target_reps: z.number().int().min(1, 'Must be at least 1').optional(),
  target_rpe: z.number().int().min(1).max(10).optional(),
  percent_1rm: z.number().min(0).max(1, 'Must be 0-1').optional(),
  planned_rest_seconds: z.number().int().min(0).optional(),
  muscle_groups: z.array(z.string()).optional(),
  notes: z.string().max(2000, 'Notes too long').optional(),
});

export type ProgramEntryFormData = z.infer<typeof programEntrySchema>;

// ============================================================================
// Program Forms
// ============================================================================

export const programCreateSchema = z.object({
  name: z.string().min(1, 'Name is required').max(200, 'Name too long'),
  description: z.string().max(2000, 'Description too long').optional(),
  entries: z.array(programEntrySchema).max(1000, 'Too many entries').optional(),
});

export type ProgramFormData = z.infer<typeof programCreateSchema>;

// ============================================================================
// Schedule Forms
// ============================================================================

export const scheduleSettingsSchema = z.object({
  programId: z.string().uuid('Invalid program ID'),
  startDate: z.string().datetime('Invalid start date'),
  skipWeekends: z.boolean(),
  avoidConsecutive: z.boolean(),
});

export type ScheduleSettingsFormData = z.infer<typeof scheduleSettingsSchema>;
