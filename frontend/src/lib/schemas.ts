import { z } from "zod";

export const NewTagSchema = z.object({
	name: z.string(),
});
