import { Outlet, useLocation } from "react-router-dom";

import Header from "../components/Header";
import Footer from "../components/Footer";

export default function SiteLayout() {
	const location = useLocation();
	const search = new URLSearchParams(location.search);
	const q = search.get("q") ?? "";

	return (
		<div className="bg-gray-100 sans-serif flex flex-col items-stretch min-h-screen">
			<Header q={q} />
			<main className="flex-grow basis-0">
				<Outlet />
			</main>
			<Footer />
		</div>
	);
}
