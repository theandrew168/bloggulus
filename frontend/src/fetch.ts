import { redirect } from "react-router-dom";

export type FetchParams = {
	method?: string;
	body?: string;
	authRequired?: boolean;
};

/**
 * Perform a fetch request to the backend API using the "token" found in the
 * browser's local storage. If a token isn't found (or it is expired) but auth
 * is required, the user will be redirected to the login page.
 */
export async function fetchAPI(url: string, params?: FetchParams): Promise<Response> {
	const method = params?.method ?? "GET";
	const body = params?.body ?? null;
	const authRequired = params?.authRequired ?? false;

	// If no token is found but auth is required, redirect to the login page.
	const token = localStorage.getItem("token");
	if (!token && authRequired) {
		throw redirect("/login");
	}

	// Construct the request headers, optionally including the auth token.
	const headers: Record<string, string> = {};
	if (token) {
		headers["Authorization"] = `Bearer ${token}`;
	}

	// Make the fetch request, optionally including the body.
	const resp = await fetch(url, {
		method,
		headers,
		body: body ?? null,
	});

	// Check for an expired token and redirect if found.
	if (resp.status === 401) {
		localStorage.removeItem("token");
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
