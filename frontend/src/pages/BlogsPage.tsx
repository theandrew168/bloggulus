import {
	Form,
	Link,
	redirect,
	useLoaderData,
	type ActionFunctionArgs,
	type LoaderFunctionArgs,
} from "react-router-dom";

import type { BlogsResponse } from "../types";
import Button from "../components/Button";

export async function blogsPageLoader({ request }: LoaderFunctionArgs) {
	const token = localStorage.getItem("token");
	if (!token) {
		return redirect("/login");
	}

	const search = new URLSearchParams();
	const url = new URL(request.url);

	const p = url.searchParams.get("p");
	if (p) {
		const page = parseInt(p);
		if (!Number.isNaN(page)) {
			search.set("page", page.toString());
		}
	}

	const resp = await fetch("/api/v1/blogs?" + search, {
		headers: {
			Authorization: `Bearer ${token}`,
		},
	});
	const blogs: BlogsResponse = await resp.json();
	return blogs;
}

export async function blogsPageAction({ request }: ActionFunctionArgs) {
	const token = localStorage.getItem("token");
	if (!token) {
		return redirect("/login");
	}

	const form = await request.formData();
	const feedURL = form.get("feedURL");
	const resp = await fetch(`/api/v1/blogs`, {
		method: "POST",
		headers: {
			Authorization: `Bearer ${token}`,
		},
		body: JSON.stringify({ feedURL }),
	});

	// If the input wasn't valid, return the errors back to the form.
	if (resp.status === 422) {
		return resp.json();
	}

	// For other errors (not related to input validation), throw to the nearest boundary.
	if (!resp.ok) {
		throw resp;
	}

	return null;
}

export default function BlogsPage() {
	const { blogs } = useLoaderData() as BlogsResponse;
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
			<div className="mb-4">
				{blogs.map((blog) => (
					<div>
						<Link to={`/blogs/${blog.id}`}>{blog.title}</Link>
					</div>
				))}
			</div>
		</div>
	);
}
