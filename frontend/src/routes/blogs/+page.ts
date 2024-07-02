import { error } from "@sveltejs/kit";

import type { PageLoad } from "./$types";
import type { BlogsResponse } from "$lib/types";

export const load: PageLoad = async ({ fetch }) => {
	const resp = await fetch("/api/v1/blogs");
	if (!resp.ok) {
		error(resp.status, await resp.text());
	}

	const blogs: BlogsResponse = await resp.json();
	return blogs;
};
