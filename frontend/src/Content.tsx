import { Link } from "react-router-dom";

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
					<div className="max-w-3xl mx-auto bg-white overflow-hidden shadow-md rounded-lg mb-6 p-6" key={article.url}>
						<div className="flex justify-between items-center mb-2">
							<span className="text-sm font-light text-gray-600">
								{new Date(article.publishedAt).toLocaleDateString()}
							</span>
							<div className="flex items-center gap-x-2">
								{article.tags.map((tag) => (
									<Link
										to={`/?q=${tag}`}
										key={tag}
										className="text-sm font-bold px-3 py-1 bg-gray-600 text-gray-100 rounded hover:bg-gray-500"
									>
										{tag}
									</Link>
								))}
							</div>
						</div>

						<a href={article.url} className="text-2xl text-gray-700 font-bold hover:underline block mb-2">
							{article.title}
						</a>

						<a href={article.blogURL} className="text-gray-700 font-bold hover:underline block">
							{article.blogTitle}
						</a>
					</div>
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
