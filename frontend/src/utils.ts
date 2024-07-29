import { redirect } from "react-router-dom";

// TODO: Make params an object, add "authRequired" field

/**
 * Perform an authenticated fetch request using the "token" found in the
 * browser's local storage. If a token isn't found (or it is expired), the
 * user will be redirected to the login page.
 */
export async function authenticatedFetch(url: string, method: string = "GET", body?: string): Promise<Response> {
	const token = localStorage.getItem("token");
	if (!token) {
		throw redirect("/login");
	}

	const resp = await fetch(url, {
		method,
		headers: {
			Authorization: `Bearer ${token}`,
		},
		body: body ?? null,
	});

	// Check for an expired token and redirect if found.
	if (resp.status === 401) {
		throw redirect("/login");
	}

	// If the input wasn't valid, return the errors back to the form.
	if (resp.status === 422) {
		return resp;
	}

	// For all other errors, throw to the nearest boundary.
	if (!resp.ok) {
		throw resp;
	}

	return resp;
}
