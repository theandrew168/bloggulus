import type { PageLoad } from "./$types";

export type Post = {
	title: string;
	url: string;
	blogTitle: string;
	blogURL: string;
	publishedAt: Date;
	tags: string[];
};

export type PostsResponse = {
	posts: Post[];
};

export const load: PageLoad = async () => {
	const resp = await fetch("/api/v1/posts");
	const posts: PostsResponse = await resp.json();
	return posts;
};
