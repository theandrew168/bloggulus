import type { PostsResponse } from "$lib/types";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const resp = await fetch("/api/v1/posts?" + url.searchParams);
	const posts: PostsResponse = await resp.json();
	return posts;
};
