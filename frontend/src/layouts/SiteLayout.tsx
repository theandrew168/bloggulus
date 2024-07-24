import { Outlet } from "react-router-dom";

import Header from "../components/Header";
import Footer from "../components/Footer";

export default function SiteLayout() {
	return (
		<div className="bg-gray-100 sans-serif flex flex-col items-stretch min-h-screen">
			<Header />
			<main className="flex-grow basis-0">
				<Outlet />
			</main>
			<Footer />
		</div>
	);
}
