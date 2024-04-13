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
		const np = Number(p);
		if (!Number.isNaN(np)) {
			const limit = 20;
			const offset = (np - 1) * limit;
			search.set("limit", limit.toString());
			search.set("offset", offset.toString());
		}
	}

	const resp = await fetch("/api/v1/posts?" + search);
	const posts: PostsResponse = await resp.json();
	return posts;
};
