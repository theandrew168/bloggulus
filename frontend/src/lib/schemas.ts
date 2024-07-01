import { z } from "zod";

export const CreateTagSchema = z.object({
	name: z.string(),
});

export const DeleteTagSchema = z.object({});
