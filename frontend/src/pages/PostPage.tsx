import {
	Form,
	Link,
	redirect,
	useLoaderData,
	type ActionFunctionArgs,
	type LoaderFunctionArgs,
} from "react-router-dom";

import { fetchAPI } from "../fetch";
import type { PostResponse } from "../types";
import Button from "../components/Button";

export async function postPageLoader({ params }: LoaderFunctionArgs) {
	// The blogID and postID URL params must exist for this route to match.
	const blogID = params["blogID"]!;
	const postID = params["postID"]!;

	const resp = await fetchAPI(`/api/v1/blogs/${blogID}/posts/${postID}`, { authRequired: true });
	return resp.json();
}

export async function postPageAction({ request }: ActionFunctionArgs) {
	const form = await request.formData();
	const blogID = form.get("blogID");
	const postID = form.get("postID");
	await fetchAPI(`/api/v1/blogs/${blogID}/posts/${postID}`, {
		method: "DELETE",
		authRequired: true,
	});

	return redirect(`/blogs/${blogID}`);
}

export default function PostPage() {
	const { post } = useLoaderData() as PostResponse;
	return (
		<div className="container mx-auto">
			<h1 className="text-2xl mt-6 mb-2">{post.title}</h1>
			<div className="mb-4">
				<a href={post.url}>(URL)</a>
				<Link to={`/blogs/${post.blogID}`}>(Blog)</Link>
			</div>
			<div className="mb-4">
				<h2 className="text-2xl">Updated</h2>
				<div>{new Date(post.publishedAt).toLocaleString()}</div>
			</div>
			<div className="mb-4">
				<h2 className="text-2xl">Actions</h2>
				<div className="buttons">
					<Form method="DELETE">
						<input type="hidden" name="blogID" value={post.blogID} />
						<input type="hidden" name="postID" value={post.id} />
						<Button type="submit">Delete</Button>
					</Form>
				</div>
			</div>
		</div>
	);
}
