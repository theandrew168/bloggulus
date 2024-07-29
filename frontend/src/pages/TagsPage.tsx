import { Form, useLoaderData, type ActionFunctionArgs, type LoaderFunctionArgs } from "react-router-dom";

import { authenticatedFetch } from "../utils";
import type { TagsResponse } from "../types";
import Button from "../components/Button";

export async function tagsPageLoader({ request }: LoaderFunctionArgs) {
	const search = new URLSearchParams();
	const url = new URL(request.url);

	const p = url.searchParams.get("p");
	if (p) {
		const page = parseInt(p);
		if (!Number.isNaN(page)) {
			search.set("page", page.toString());
		}
	}

	const resp = await authenticatedFetch("/api/v1/tags?" + search);
	const tags: TagsResponse = await resp.json();
	return tags;
}

export async function tagsPageAction({ request }: ActionFunctionArgs) {
	const form = await request.formData();
	const name = form.get("name");

	const resp = await authenticatedFetch(`/api/v1/tags`, "POST", JSON.stringify({ name }));

	// If the input wasn't valid, return the errors back to the form.
	if (resp.status === 422) {
		return resp.json();
	}

	return null;
}

export default function TagsPage() {
	const { tags } = useLoaderData() as TagsResponse;
	return (
		<div className="container mx-auto">
			<h1 className="text-lg font-semibold mt-6 mb-2">Tags</h1>
			<div className="mb-4">
				<Form method="POST" className="flex flex-row gap-4">
					<input
						id="name"
						name="name"
						placeholder="Add Tag"
						className="block rounded-md border-0 py-1.5 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-gray-800"
					/>
					<Button type="submit">Add</Button>
				</Form>
			</div>
			<div className="mb-4">
				{tags.map((tag) => (
					<div key={tag.id}>{tag.name}</div>
				))}
			</div>
		</div>
	);
}
