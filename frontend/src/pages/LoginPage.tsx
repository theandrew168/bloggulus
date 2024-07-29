import { Form, redirect, useActionData, type ActionFunctionArgs } from "react-router-dom";

import { findGeneralError, findSpecificErrors, type StructuredErrorsResponse } from "../errors";
import type { TokenResponse } from "../types";
import FormInput from "../components/FormInput";
import Button from "../components/Button";

export async function loginPageAction({ request }: ActionFunctionArgs) {
	const form = await request.formData();
	const username = form.get("username");
	const password = form.get("password");
	const resp = await fetch(`/api/v1/tokens`, {
		method: "POST",
		body: JSON.stringify({ username, password }),
	});

	// If the input wasn't valid, return the errors back to the form.
	if (resp.status === 422) {
		return resp.json();
	}

	// For other errors (not related to input validation), throw to the nearest boundary.
	if (!resp.ok) {
		throw resp;
	}

	const token: TokenResponse = await resp.json();
	localStorage.setItem("token", token.token.value);
	return redirect("/");
}

export default function LoginPage() {
	// https://reactrouter.com/en/main/hooks/use-action-data
	const errors = useActionData() as StructuredErrorsResponse | undefined;
	const generalError = findGeneralError(errors?.errors ?? []);
	const specificErrors = findSpecificErrors(errors?.errors ?? []);

	return (
		<div className="h-full flex items-center justify-center">
			<Form method="POST" className="bg-white p-8 shadow rounded-md flex flex-col gap-6">
				{generalError && <div className="text-sm text-red-500">{generalError}</div>}
				<FormInput name="username" label="Username" error={specificErrors["username"]} />
				<FormInput name="password" label="Password" type="password" error={specificErrors["password"]} />
				<Button type="submit">Login</Button>
			</Form>
		</div>
	);
}
