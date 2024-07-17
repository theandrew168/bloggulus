import type { Article } from "./types";
import TagBadge from "./TagBadge";

export type Props = {
	article: Article;
};

export default function ArticleCard({ article }: Props) {
	return (
		<div className="max-w-3xl mx-auto bg-white overflow-hidden shadow-md rounded-lg mb-6 p-6" key={article.url}>
			<div className="flex justify-between items-center mb-2">
				<span className="text-sm font-light text-gray-600">{new Date(article.publishedAt).toLocaleDateString()}</span>
				<div className="flex items-center gap-x-2">
					{article.tags.map((tag) => (
						<TagBadge tag={tag} />
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
	);
}
