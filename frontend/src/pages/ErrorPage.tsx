import { isRouteErrorResponse, useRouteError } from "react-router-dom";
import { findGeneralError, isStructuredErrorsResponse } from "../errors";
import ButtonLink from "../components/ButtonLink";

export default function ErrorPage() {
	let statusCode = 500;
	let statusText = "Internal server error";
	let message = "Sorry, something went wrong.";

	// If the error is an ErrorResponse (react-router-dom)...
	const resp = useRouteError();
	if (isRouteErrorResponse(resp)) {
		// Update the status and status text.
		statusCode = resp.status;
		statusText = resp.statusText;
		// If the error is a StructuredErrorsResponse (bloggulus)...
		if (isStructuredErrorsResponse(resp.data)) {
			// Find the first general error and use it for the message.
			const generalError = findGeneralError(resp.data.errors);
			if (generalError) {
				message = generalError;
			}
		}
	}

	return (
		<div className="min-h-screen flex items-center justify-center">
			<div className="text-center">
				<p className="font-semibold text-gray-700">{statusCode}</p>
				<h1 className="mt-4 text-3xl font-bold text-gray-900">{statusText}</h1>
				<p className="mt-6">{message}</p>
				<div className="mt-10">
					<ButtonLink to="/">Go back home</ButtonLink>
				</div>
			</div>
		</div>
	);
}
