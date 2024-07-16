export default function Footer() {
	return (
		<footer className="text-gray-100 bg-gray-800">
			<div className="max-w-3xl mx-auto py-4">
				<div className="flex flex-col items-center justify-between md:flex-row space-y-1 md:space-y-0">
					<a href="/" className="text-xl font-bold text-gray-100 hover:text-gray-400">
						Bloggulus
					</a>
					<a href="https://shallowbrooksoftware.com" className="text-gray-100 hover:text-gray-400">
						Shallow Brook Software
					</a>
				</div>
			</div>
		</footer>
	);
}
