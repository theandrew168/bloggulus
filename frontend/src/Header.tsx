export default function Header() {
	return (
		<nav className="bg-white shadow">
			<div className="max-w-3xl mx-auto py-3 px-6 md:px-0 flex justify-between items-center gap-x-2">
				<div className="flex flex-wrap items-baseline gap-2">
					<a href="/" className="text-gray-800 text-xl md:text-2xl hover:text-gray-500">
						Bloggulus
					</a>
					<a href="/api/v1/" className="text-gray-800 text-base md:text-lg hover:text-gray-500">
						[API]
					</a>
				</div>

				<form method="GET" action="/" className="block relative shrink-0">
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
						placeholder="Search"
						className="w-full py-2 pl-10 pr-4 text-gray-700 bg-white border border-gray-300 rounded-md focus:border-blue-500 focus:outline-none focus:ring"
					/>
				</form>
			</div>
		</nav>
	);
}
