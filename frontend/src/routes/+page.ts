import { error } from "@sveltejs/kit";

import type { PostsResponse } from "$lib/types";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const search = new URLSearchParams();

	const q = url.searchParams.get("q");
	if (q) {
		search.set("q", q);
	}

	const p = url.searchParams.get("p");
	if (p) {
		const page = parseInt(p);
		if (!Number.isNaN(page)) {
			search.set("page", page.toString());
		}
	}

	const resp = await fetch("/api/v1/posts?" + search);
	if (!resp.ok) {
		error(resp.status, await resp.text());
	}

	const posts: PostsResponse = await resp.json();
	return posts;
};
