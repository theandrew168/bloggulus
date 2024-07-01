import { error } from "@sveltejs/kit";
import { superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";

import type { PageLoad } from "./$types";
import type { TagsResponse } from "$lib/types";
import { NewTagSchema } from "$lib/schemas";

export const load: PageLoad = async ({ fetch }) => {
	const resp = await fetch("/api/v1/tags");
	if (!resp.ok) {
		error(resp.status, await resp.text());
	}

	const tags: TagsResponse = await resp.json();

	const form = await superValidate(zod(NewTagSchema));
	return { tags: tags.tags, form };
};
