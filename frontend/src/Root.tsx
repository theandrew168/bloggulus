import Content from "./Content";
import Footer from "./Footer";
import Header from "./Header";

export default function Root() {
	return (
		<div className="h-full bg-gray-100 sans-serif flex flex-col">
			<Header />
			<main className="flex-grow">
				<Content />
			</main>
			<Footer />
		</div>
	);
}
