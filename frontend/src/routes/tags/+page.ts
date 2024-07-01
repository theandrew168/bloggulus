import { error } from "@sveltejs/kit";

import type { PageLoad } from "./$types";
import type { TagsResponse } from "$lib/types";

export const load: PageLoad = async ({ fetch }) => {
	const resp = await fetch("/api/v1/tags");
	if (!resp.ok) {
		error(resp.status, await resp.text());
	}

	const tags: TagsResponse = await resp.json();
	return tags;
};
