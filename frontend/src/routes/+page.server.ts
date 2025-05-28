import type { PageServerLoad } from "./$types";

type Article = {
	title: string;
	url: string;
	blogTitle: string;
	blogURL: string;
	publishedAt: string;
	tags: string[];
};

type ArticlesResponse = {
	articles: Article[];
};

export const load: PageServerLoad = async () => {
	const articlesResp = await fetch("http://localhost:5000/api/v1/articles");
	const articles: ArticlesResponse = await articlesResp.json();
	return { articles: articles.articles };
};