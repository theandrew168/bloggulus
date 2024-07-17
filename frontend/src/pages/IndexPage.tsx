import { Link, useLoaderData, useLocation } from "react-router-dom";
import type { LoaderFunctionArgs } from "react-router-dom";

import type { ArticlesResponse } from "../types";
import Footer from "../components/Footer";
import Header from "../components/Header";
import ArticleCard from "../components/ArticleCard";

export async function loader({ request }: LoaderFunctionArgs) {
	const search = new URLSearchParams();

	const url = new URL(request.url);
	const q = url.searchParams.get("q");
	if (q) {
		search.set("q", q);
	}

	const p = url.searchParams.get("p");
	if (p) {
		const page = parseInt(p);
		if (!Number.isNaN(page)) {
			search.set("page", page.toString());
		}
	}

	const resp = await fetch("/api/v1/articles?" + search);
	const articles: ArticlesResponse = await resp.json();
	return articles;
}

export default function IndexPage() {
	const { articles, count } = useLoaderData() as ArticlesResponse;

	const location = useLocation();
	const search = new URLSearchParams(location.search);
	const q = search.get("q") ?? "";
	const p = parseInt(search.get("p") ?? "1") ?? 1;
	const moreLink = `/?p=${p + 1}` + (q ? `&q=${q}` : "");
	const hasMorePages = p * 20 < count;

	return (
		<div className="bg-gray-100 sans-serif flex flex-col min-h-screen">
			<Header q={q} />
			<main className="flex-grow">
				<div className="max-w-3xl mx-auto flex justify-start items-center my-6 px-6 md:px-0">
					<h1 className="text-xl font-bold text-gray-700 md:text-2xl">Recent Articles</h1>
				</div>

				<div className="px-6 md:px-0">
					{articles.map((article) => (
						<ArticleCard article={article} />
					))}
				</div>

				{hasMorePages && (
					<div className="mx-auto mb-6 px-16 md:px-0 flex justify-center items-center gap-x-4">
						<Link to={moreLink} className="bg-white text-gray-700 font-bold shadow hover:shadow-md rounded px-6 py-2">
							See More
						</Link>
					</div>
				)}
			</main>
			<Footer />
		</div>
	);
}
