import { Form, redirect, type ActionFunctionArgs } from "react-router-dom";

import FormInput from "../components/FormInput";
import Button from "../components/Button";

// TODO: Handle input validation and errors.
export async function registerPageAction({ request }: ActionFunctionArgs) {
	const form = await request.formData();
	const username = form.get("username");
	const password = form.get("password");
	const resp = await fetch(`/api/v1/accounts`, {
		method: "POST",
		body: JSON.stringify({ username, password }),
	});

	if (resp.ok) {
		return redirect("/login");
	}

	return null;
}

export default function RegisterPage() {
	return (
		<div className="h-full flex items-center justify-center">
			<Form method="POST" className="max-w-xl bg-white p-8 shadow rounded-md flex flex-col gap-6">
				<FormInput name="username" label="Username" />
				<FormInput name="password" label="Password" type="password" />
				<Button type="submit">Register</Button>
			</Form>
		</div>
	);
}
