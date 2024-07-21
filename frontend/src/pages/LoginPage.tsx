import { Form, redirect, type ActionFunctionArgs } from "react-router-dom";

import type { TokenResponse } from "../types";
import LabeledInput from "../components/LabeledInput";
import Button from "../components/Button";

// TODO: Handle input validation and errors.
export async function loginPageAction({ request }: ActionFunctionArgs) {
	const form = await request.formData();
	const username = form.get("username");
	const password = form.get("password");
	const resp = await fetch(`/api/v1/tokens`, {
		method: "POST",
		body: JSON.stringify({ username, password }),
	});

	if (resp.ok) {
		const token: TokenResponse = await resp.json();
		localStorage.setItem("token", token.token.value);
		return redirect("/");
	}

	return null;
}

export default function LoginPage() {
	return (
		<div className="h-full flex items-center justify-center">
			<Form method="POST" className="max-w-xl bg-white p-8 shadow rounded-md flex flex-col gap-6">
				<LabeledInput name="username" label="Username" />
				<LabeledInput name="password" label="Password" type="password" />
				<Button type="submit">Login</Button>
			</Form>
		</div>
	);
}
