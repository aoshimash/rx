import type { ProgramTemplateEntryCreate } from '@/types/api';
import { z } from 'zod';

// ============================================================================
// Log Entry Forms
// ============================================================================

export const logEntrySchema = z.object({
  exercise_name: z.string().min(1, 'Exercise name is required').max(200, 'Name too long'),
  sets: z.number().int().min(1, 'Must be at least 1').optional(),
  reps: z.number().int().min(1, 'Must be at least 1').optional(),
  load_kg: z.number().min(0, 'Must be non-negative').optional(),
  rpe: z.number().int().min(1, 'Must be 1-10').max(10, 'Must be 1-10').optional(),
  notes: z.string().max(2000, 'Notes too long').optional(),
});

export type LogEntryFormData = z.infer<typeof logEntrySchema>;

// ============================================================================
// Log Forms
// ============================================================================

export const logCreateSchema = z.object({
  performed_at: z.string().datetime('Invalid datetime format'),
  program_id: z.string().uuid().optional(),
  session_name: z.string().optional(),
  notes: z.string().max(5000, 'Notes too long').optional(),
  entries: z.array(logEntrySchema).min(1, 'At least one entry required'),
});

export type LogFormData = z.infer<typeof logCreateSchema>;

// ============================================================================
// Program Template Entry Forms
// ============================================================================

export const programTemplateEntrySchema: z.ZodType<ProgramTemplateEntryCreate> = z.object({
  exercise_name: z.string().min(1, 'Exercise name is required').max(200, 'Name too long'),
  order: z.number().int().min(0, 'Order must be non-negative'),
  sets: z.number().int().min(1, 'Must be at least 1').optional(),
  reps: z.number().int().min(1, 'Must be at least 1').optional(),
  rpe: z.number().int().min(1).max(10).optional(),
  percent_1rm: z.number().min(0).max(1).optional(),
  notes: z.string().max(2000, 'Notes too long').optional(),
  metadata: z.record(z.string(), z.unknown()).optional(),
});

export type ProgramTemplateEntryFormData = z.infer<typeof programTemplateEntrySchema>;

// ============================================================================
// Program Template Forms
// ============================================================================

export const programTemplateCreateSchema = z.object({
  name: z.string().min(1, 'Name is required').max(200, 'Name too long'),
  description: z.string().max(2000, 'Description too long').optional(),
  notes: z.string().max(5000, 'Notes too long').optional(),
  entries: z.array(programTemplateEntrySchema).max(1000, 'Too many entries').optional(),
});

export type ProgramTemplateFormData = z.infer<typeof programTemplateCreateSchema>;
