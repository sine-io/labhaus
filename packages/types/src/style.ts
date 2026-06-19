import { z } from 'zod';

export const styleSchema = z.object({
  id: z.string().uuid(),
  case_id: z.number().int().positive(),
  title: z.string().min(1),
  prompt: z.string().min(1),
  prompt_preview: z.string().nullable(),
  category: z.string().min(1),
  styles: z.array(z.string()),
  scenes: z.array(z.string()),
  image_url: z.string().url().nullable(),
  source_label: z.string().nullable(),
  source_url: z.string().url().nullable(),
  github_url: z.string().url().nullable(),
  featured: z.boolean(),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
});

export const styleQuerySchema = z.object({
  category: z.string().optional(),
  style: z.string().optional(),
  scene: z.string().optional(),
  featured: z.boolean().optional(),
  search: z.string().optional(),
  page: z.coerce.number().int().positive().default(1),
  limit: z.coerce.number().int().positive().max(100).default(20),
});

export type Style = z.infer<typeof styleSchema>;
export type StyleQuery = z.infer<typeof styleQuerySchema>;
