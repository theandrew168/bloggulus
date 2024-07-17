import { Link } from "react-router-dom";

import ArticleCard from "./ArticleCard";
import type { Article } from "./types";

export type Props = {
	articles: Article[];
	moreLink: string;
	hasMorePages: boolean;
};

export default function Content({ articles, moreLink, hasMorePages }: Props) {
	return (
		<>
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
		</>
	);
}
