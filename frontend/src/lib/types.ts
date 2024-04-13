// reader domain types

export type Post = {
	title: string;
	url: string;
	blogTitle: string;
	blogURL: string;
	publishedAt: string;
	tags: string[];
};

export type PostsResponse = {
	posts: Post[];
};
