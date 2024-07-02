import { error } from "@sveltejs/kit";

import type { PageLoad } from "./$types";
import type { BlogResponse, PostsResponse } from "$lib/types";

// TODO: Load these two reqs at the same time.
export const load: PageLoad = async ({ params, fetch }) => {
	const id = params.id;

	const blogResp = await fetch(`/api/v1/blogs/${id}`);
	if (!blogResp.ok) {
		error(blogResp.status, await blogResp.text());
	}

	const blog: BlogResponse = await blogResp.json();

	const postsResp = await fetch(`/api/v1/blogs/${id}/posts`);
	if (!postsResp.ok) {
		error(postsResp.status, await postsResp.text());
	}

	const posts: PostsResponse = await postsResp.json();

	return { ...blog, ...posts };
};
