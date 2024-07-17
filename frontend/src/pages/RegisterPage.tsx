import { Form, redirect, type ActionFunctionArgs } from "react-router-dom";

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
			<Form method="POST" className="max-w-xl bg-white p-8 shadow rounded-md flex flex-col gap-4">
				<label>
					Username:
					<br />
					<input name="username" placeholder="Username" />
				</label>
				<br />
				<label>
					Password:
					<br />
					<input name="password" placeholder="Password" type="password" />
				</label>
				<br />
				<button
					className="text-sm font-bold px-3 py-1 bg-gray-700 text-gray-100 rounded hover:bg-gray-500"
					type="submit"
				>
					Register
				</button>
			</Form>
		</div>
	);
}
