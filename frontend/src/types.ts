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
	count: number;
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
	count: number;
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
	count: number;
	tags: Tag[];
};

export type Account = {
	id: string;
	username: string;
};

export type AccountResponse = {
	account: Account;
};

export type Token = {
	id: string;
	value: string;
	expiresAt: string;
};

export type TokenResponse = {
	token: Token;
};
