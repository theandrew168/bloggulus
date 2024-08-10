import { Form, Link, useLoaderData, type ActionFunctionArgs, type LoaderFunctionArgs } from "react-router-dom";

import { fetchAPI } from "../fetch";
import type { BlogsResponse } from "../types";
import Button from "../components/Button";

type BlogsAllAndFollowing = {
	all: BlogsResponse;
	following: BlogsResponse;
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

	const allResp = await fetchAPI("/api/v1/blogs?" + search, { authRequired: true });
	const allBlogs: BlogsResponse = await allResp.json();

	const followingResp = await fetchAPI("/api/v1/blogs/following?" + search, { authRequired: true });
	const followingBlogs: BlogsResponse = await followingResp.json();
	return { all: allBlogs, following: followingBlogs };
}

export async function blogsPageAction({ request }: ActionFunctionArgs) {
	const form = await request.formData();
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

	return null;
}

export default function BlogsPage() {
	const { all, following } = useLoaderData() as BlogsAllAndFollowing;
	return (
		<div className="container mx-auto">
			<h1 className="text-lg font-semibold mt-6 mb-2">Blogs</h1>
			<div className="mb-4">
				<Form method="POST" className="flex flex-row gap-4">
					<input
						id="feedURL"
						name="feedURL"
						placeholder="Add Feed"
						className="block rounded-md border-0 py-1.5 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-gray-800"
					/>
					<Button type="submit">Add</Button>
				</Form>
			</div>
			<h2 className="text-lg font-semibold">Following Blogs</h2>
			<div className="mb-4">
				{following.blogs.map((blog) => (
					<div key={blog.id}>
						<Link to={`/blogs/${blog.id}`}>{blog.title}</Link>
					</div>
				))}
			</div>
			<h2 className="text-lg font-semibold">All Blogs</h2>
			<div className="mb-4">
				{all.blogs.map((blog) => (
					<div key={blog.id}>
						<Link to={`/blogs/${blog.id}`}>{blog.title}</Link>
					</div>
				))}
			</div>
		</div>
	);
}
