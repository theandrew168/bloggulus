import { error } from "@sveltejs/kit";

import type { PageLoad } from "./$types";
import type { PostResponse } from "$lib/types";

export const load: PageLoad = async ({ fetch, params }) => {
	const blogID = params.blogID;
	const postID = params.postID;

	const resp = await fetch(`/api/v1/blogs/${blogID}/posts/${postID}`);
	if (!resp.ok) {
		error(resp.status, await resp.text());
	}

	const post: PostResponse = await resp.json();
	return post;
};
