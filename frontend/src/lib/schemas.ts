import { z } from 'zod';

export const TagSchema = z.object({
	name: z.string().min(2),
});
