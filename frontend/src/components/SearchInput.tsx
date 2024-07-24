import { useEffect, useState } from "react";
import { Form, useLocation } from "react-router-dom";

export default function SearchInput() {
	const location = useLocation();
	const searchParams = new URLSearchParams(location.search);
	const q = searchParams.get("q") ?? "";

	const [search, setSearch] = useState(q);
	useEffect(() => {
		setSearch(q);
	}, [q]);

	return (
		<Form method="GET" action="/" className="block relative shrink-0">
			<span className="absolute inset-y-0 left-0 flex items-center pl-3">
				<svg className="w-5 h-5 text-gray-400" viewBox="0 0 24 24" fill="none">
					<path
						d="M21 21L15 15M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z"
						stroke="currentColor"
						strokeWidth="2"
						strokeLinecap="round"
						strokeLinejoin="round"
					></path>
				</svg>
			</span>

			<input
				name="q"
				type="text"
				value={search}
				onChange={(e) => setSearch(e.currentTarget.value)}
				placeholder="Search"
				className="w-full py-2 pl-10 pr-4 text-gray-700 bg-white border border-gray-300 rounded-md ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-gray-800"
			/>
		</Form>
	);
}
