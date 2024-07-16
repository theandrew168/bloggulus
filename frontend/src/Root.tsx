import { useLoaderData, useLocation } from "react-router-dom";
import type { LoaderFunctionArgs } from "react-router-dom";

import Content from "./Content";
import Footer from "./Footer";
import Header from "./Header";
import type { ArticlesResponse } from "./types";

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

export default function Root() {
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
				<Content articles={articles} moreLink={moreLink} hasMorePages={hasMorePages} />
			</main>
			<Footer />
		</div>
	);
}
