import { error } from "@sveltejs/kit";

import type { PageLoad } from "./$types";
import type { TagsResponse } from "$lib/types";
import { goto } from "$app/navigation";

export const load: PageLoad = async ({ fetch }) => {
	const token = localStorage.getItem("token");
	if (!token) {
		await goto("/login");
	}

	const resp = await fetch("/api/v1/tags", {
		headers: {
			Authorization: `Bearer ${token}`,
		},
	});
	if (!resp.ok) {
		error(resp.status, await resp.text());
	}

	const tags: TagsResponse = await resp.json();
	return tags;
};
