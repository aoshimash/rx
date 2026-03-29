import { z } from 'zod';

// ============================================================================
// Log Entry Forms
// ============================================================================

export const logEntrySchema = z.object({
  exercise_name: z.string().min(1, 'Exercise name is required').max(200, 'Name too long'),
  fields: z.record(z.unknown()).optional(),
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
