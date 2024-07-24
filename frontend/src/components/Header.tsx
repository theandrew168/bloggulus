import { Link } from "react-router-dom";
import SearchInput from "./SearchInput";

export default function Header() {
	return (
		<nav className="bg-white shadow">
			<div className="max-w-3xl mx-auto py-3 px-6 md:px-0 flex justify-between items-center gap-x-2">
				<Link to="/" className="text-gray-800 text-xl md:text-2xl hover:text-gray-500">
					Bloggulus
				</Link>

				<SearchInput />
			</div>
		</nav>
	);
}
