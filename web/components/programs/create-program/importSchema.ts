import { z } from 'zod';

const programSessionEntryImportSchema = z.object({
  exercise_name: z.string().min(1),
  order: z.number().int().min(0),
  sets: z.number().int().min(1).optional(),
  reps: z.number().int().min(1).optional(),
  load_kg: z.number().positive().optional(),
  rpe: z.number().int().min(1).max(10).optional(),
  notes: z.string().optional(),
});

const programSessionImportSchema = z.object({
  session_name: z.string().min(1),
  order: z.number().int().min(0),
  date: z.string().optional(),
  entries: z.array(programSessionEntryImportSchema).optional(),
});

export const programImportSchema = z.object({
  rx_version: z.literal('1'),
  name: z.string().min(1),
  notes: z.string().optional(),
  sessions: z.array(programSessionImportSchema),
});

export type ProgramImport = z.infer<typeof programImportSchema>;
