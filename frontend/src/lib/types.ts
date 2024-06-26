export type Article = {
	title: string;
	url: string;
	blogTitle: string;
	blogURL: string;
	publishedAt: string;
	tags: string[];
};

export type ArticlesResponse = {
	count: number;
	articles: Article[];
};

export type Blog = {
	id: string;
	feedURL: string;
	siteURL: string;
	title: string;
	syncedAt: string;
};

export type BlogResponse = {
	blog: Blog;
};

export type BlogsResponse = {
	blogs: Blog[];
};

export type Post = {
	id: string;
	blogID: string;
	url: string;
	title: string;
	publishedAt: string;
};

export type PostResponse = {
	post: Post;
};

export type PostsResponse = {
	posts: Post[];
};

export type Tag = {
	id: string;
	name: string;
};

export type TagResponse = {
	tag: Tag;
};

export type TagsResponse = {
	tags: Tag[];
};
