import { Link } from "react-router-dom";

export type Props = {
	tag: string;
};

export default function TagBadge({ tag }: Props) {
	return (
		<Link
			to={`/?q=${tag}`}
			key={tag}
			className="text-sm font-bold px-3 py-1 bg-gray-600 text-gray-100 rounded hover:bg-gray-500"
		>
			{tag}
		</Link>
	);
}
