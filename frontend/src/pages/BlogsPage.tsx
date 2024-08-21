import React from "react";
import {
	Await,
	defer,
	Link,
	useFetcher,
	useLoaderData,
	type ActionFunctionArgs,
	type LoaderFunctionArgs,
} from "react-router-dom";

import { fetchAPI } from "../fetch";
import type { Blog, BlogsResponse } from "../types";
import Button from "../components/Button";
import Toggle from "../components/Toggle";
import DisabledToggle from "../components/DisabledToggle";

type BlogsWithDeferredFollowing = {
	count: number;
	blogs: Blog[];
	following: Record<string, Promise<Response>>;
};

export async function blogsPageLoader({ request }: LoaderFunctionArgs) {
	const search = new URLSearchParams();
	const url = new URL(request.url);

	const p = url.searchParams.get("p");
	if (p) {
		const page = parseInt(p);
		if (!Number.isNaN(page)) {
			search.set("page", page.toString());
		}
	}

	// Fetch the current page of blogs.
	const blogsResp = await fetchAPI("/api/v1/blogs?" + search, { authRequired: true });
	const blogs: BlogsResponse = await blogsResp.json();

	// Build a mapping of blog IDs to their deferred following status (204 vs 404).
	const following = blogs.blogs.reduce(
		(acc, blog) => {
			return {
				...acc,
				[blog.id]: fetchAPI(`/api/v1/blogs/${blog.id}/following`, {
					authRequired: true,
					ignoreNotFound: true,
				}),
			};
		},
		{} as Record<string, Promise<Response>>,
	);

	return defer({ count: blogs.count, blogs: blogs.blogs, following });
}

export async function blogsPageAction({ request }: ActionFunctionArgs) {
	const form = await request.formData();
	const intent = form.get("intent");

	if (intent === "add") {
		const feedURL = form.get("feedURL");
		const resp = await fetchAPI(`/api/v1/blogs`, {
			method: "POST",
			body: JSON.stringify({ feedURL }),
			authRequired: true,
		});

		// If the input wasn't valid, return the errors back to the form.
		if (resp.status === 422) {
			return resp.json();
		}
	}

	if (intent === "follow") {
		const id = form.get("id");

		await fetchAPI(`/api/v1/blogs/${id}/follow`, {
			method: "POST",
			authRequired: true,
		});
	}

	if (intent === "unfollow") {
		const id = form.get("id");

		await fetchAPI(`/api/v1/blogs/${id}/unfollow`, {
			method: "POST",
			authRequired: true,
		});
	}

	return null;
}

export default function BlogsPage() {
	const fetcher = useFetcher();
	const { blogs, following } = useLoaderData() as BlogsWithDeferredFollowing;
	return (
		<div className="max-w-3xl px-6 md:px-0 mx-auto">
			<h1 className="text-lg font-semibold mt-6 mb-2">Blogs</h1>
			<div className="mb-4">
				<fetcher.Form method="POST" className="flex flex-row gap-4">
					<input
						id="feedURL"
						name="feedURL"
						placeholder="Add Feed"
						className="block rounded-md border-0 py-1.5 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-gray-800"
					/>
					<Button type="submit" name="intent" value="add">
						Add
					</Button>
				</fetcher.Form>
			</div>
			<div className="divide-y divide-gray-200">
				{blogs.map((blog) => (
					<div key={blog.id} className="py-2 flex gap-4 items-center justify-between">
						<Link className="hover:underline" to={`/blogs/${blog.id}`}>
							{blog.title}
						</Link>
						<React.Suspense fallback={<DisabledToggle />}>
							<Await resolve={following[blog.id]}>
								{(followingResponse: Response) => (
									<fetcher.Form method="POST">
										<input type="hidden" name="id" value={blog.id} />
										<Toggle
											initiallyEnabled={followingResponse.ok}
											onToggle={(enabled) => {
												fetcher.submit(
													{
														id: blog.id,
														intent: enabled ? "follow" : "unfollow",
													},
													{ method: "POST" },
												);
											}}
										/>
									</fetcher.Form>
								)}
							</Await>
						</React.Suspense>
					</div>
				))}
			</div>
		</div>
	);
}
