import { error } from "@sveltejs/kit";

import type { PageLoad } from "./$types";
import type { PostResponse } from "$lib/types";
import { goto } from "$app/navigation";

export const load: PageLoad = async ({ fetch, params }) => {
	const token = localStorage.getItem("token");
	if (!token) {
		await goto("/login");
	}

	const blogID = params.blogID;
	const postID = params.postID;

	const resp = await fetch(`/api/v1/blogs/${blogID}/posts/${postID}`, {
		headers: {
			Authorization: `Bearer ${token}`,
		},
	});
	if (!resp.ok) {
		error(resp.status, await resp.text());
	}

	const post: PostResponse = await resp.json();
	return post;
};
