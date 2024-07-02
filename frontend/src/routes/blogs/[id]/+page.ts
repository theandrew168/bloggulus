import { error } from "@sveltejs/kit";

import type { PageLoad } from "./$types";
import type { BlogResponse } from "$lib/types";

export const load: PageLoad = async ({ params, fetch }) => {
	const id = params.id;
	const resp = await fetch(`/api/v1/blogs/${id}`);
	if (!resp.ok) {
		error(resp.status, await resp.text());
	}

	const blog: BlogResponse = await resp.json();
	return blog;
};
