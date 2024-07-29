import { Form, redirect, useActionData, type ActionFunctionArgs } from "react-router-dom";

import { findSpecificErrors, type StructuredErrorsResponse } from "../errors";
import FormInput from "../components/FormInput";
import Button from "../components/Button";

export async function registerPageAction({ request }: ActionFunctionArgs) {
	const form = await request.formData();
	const username = form.get("username");
	const password = form.get("password");
	const resp = await fetch(`/api/v1/accounts`, {
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

	return redirect("/login");
}

export default function RegisterPage() {
	// https://reactrouter.com/en/main/hooks/use-action-data
	const errors = useActionData() as StructuredErrorsResponse | undefined;
	const specificErrors = findSpecificErrors(errors?.errors ?? []);

	return (
		<div className="h-full flex items-center justify-center">
			<Form method="POST" className="bg-white p-8 shadow rounded-md flex flex-col gap-6">
				<FormInput name="username" label="Username" error={specificErrors["username"]} />
				<FormInput name="password" label="Password" type="password" error={specificErrors["password"]} />
				<Button type="submit">Register</Button>
			</Form>
		</div>
	);
}
