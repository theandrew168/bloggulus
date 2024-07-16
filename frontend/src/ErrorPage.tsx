import { isRouteErrorResponse, useRouteError } from "react-router-dom";

export default function ErrorPage() {
	let statusCode = 500;
	let statusText = "Internal server error";
	let message = "Sorry, something went wrong.";

	const error = useRouteError();
	if (isRouteErrorResponse(error)) {
		statusCode = error.status;
		statusText = error.statusText;
		if (error.data?.message) {
			message = error.data.message + ".";
		}
	}

	return (
		<div className="container mx-auto">
			<h1 className="mt-4">
				<p>
					{statusCode}: {statusText}
				</p>
				<p>{message}</p>
			</h1>
		</div>
	);
}
