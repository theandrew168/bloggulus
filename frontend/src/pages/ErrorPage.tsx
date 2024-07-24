import { isRouteErrorResponse, useRouteError } from "react-router-dom";
import { isStructuredErrorsResponse, type StructuredError } from "../errors";

export default function ErrorPage() {
	let statusCode = 500;
	let statusText = "Internal server error";
	const errors: StructuredError[] = [];

	// If the error is an ErrorResponse (react-router-dom)...
	const resp = useRouteError();
	if (isRouteErrorResponse(resp)) {
		// Update the status and status text.
		statusCode = resp.status;
		statusText = resp.statusText;
		// If the error is a StructuredErrorsResponse (bloggulus)...
		if (isStructuredErrorsResponse(resp.data)) {
			// Update the list of specific errors.
			errors.push(...resp.data.errors);
		}
	}

	return (
		<div className="container mx-auto">
			<h1 className="mt-4">
				<p>
					{statusCode}: {statusText}
				</p>
				<ul>
					{errors.map((e, i) => (
						<li key={i}>
							{e.message} - {e.field}
						</li>
					))}
				</ul>
			</h1>
		</div>
	);
}
