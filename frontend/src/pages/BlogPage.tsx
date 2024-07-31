import {
	Form,
	Link,
	redirect,
	useLoaderData,
	type ActionFunctionArgs,
	type LoaderFunctionArgs,
} from "react-router-dom";

import { fetchAPI } from "../fetch";
import type { BlogResponse, PostsResponse } from "../types";
import Button from "../components/Button";

export async function blogPageLoader({ request, params }: LoaderFunctionArgs) {
	// The blogID URL param must exist for this route to match.
	const blogID = params["blogID"]!;

	const [blogResp, postsResp] = await Promise.all([
		fetchAPI(`/api/v1/blogs/${blogID}`, { authRequired: true }),
		fetchAPI(`/api/v1/blogs/${blogID}/posts`, { authRequired: true }),
	]);

	const [blog, posts] = await Promise.all([blogResp.json(), postsResp.json()]);
	return { ...blog, ...posts };
}

export async function blogPageAction({ request }: ActionFunctionArgs) {
	const form = await request.formData();
	const blogID = form.get("blogID");
	await fetchAPI(`/api/v1/blogs/${blogID}`, {
		method: "DELETE",
		authRequired: true,
	});

	return redirect("/blogs");
}

export default function BlogPage() {
	const { blog, posts } = useLoaderData() as BlogResponse & PostsResponse;
	return (
		<div className="container mx-auto">
			<h1 className="text-2xl mt-6 mb-2">{blog.title}</h1>
			<div className="mb-4">
				<a href={blog.siteURL}>(Site URL)</a>
				<a href={blog.feedURL}>(Feed URL)</a>
			</div>
			<div className="times">
				<h2>Synced at:</h2>
				<div>{new Date(blog.syncedAt).toLocaleString()}</div>
			</div>

			<div className="mb-4">
				<h2 className="text-2xl">Actions</h2>
				<div className="flex gap-2">
					<Form method="DELETE">
						<input type="hidden" name="blogID" value={blog.id} />
						<Button type="submit">Delete</Button>
					</Form>
				</div>
			</div>

			<div className="mb-4">
				<h2 className="text-2xl">{posts.length} Posts</h2>
				{posts.map((post) => (
					<div key={post.id} className="flex justify-between">
						<Link to={`/blogs/${blog.id}/posts/${post.id}`}>{post.title}</Link>
						<span>{new Date(post.publishedAt).toDateString()}</span>
					</div>
				))}
			</div>
		</div>
	);
}
