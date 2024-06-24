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
