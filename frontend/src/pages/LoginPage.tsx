import { Form } from "react-router-dom";

export default function LoginPage() {
	return (
		<div className="h-full flex items-center justify-center">
			<Form className="max-w-xl bg-white p-8 shadow rounded-md flex flex-col gap-4">
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
					Login
				</button>
			</Form>
		</div>
	);
}
