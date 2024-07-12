<script lang="ts">
	import { goto } from "$app/navigation";

	// TODO: Add proper validation / error handling
	async function register(event: Event) {
		const form = new FormData(event.target as HTMLFormElement);
		const username = form.get("username");
		const password = form.get("password");
		const resp = await fetch(`/api/v1/accounts`, {
			method: "POST",
			body: JSON.stringify({ username, password }),
		});

		if (resp.ok) {
			await goto("/login");
		}
	}
</script>

<div class="container">
	<h1>Register</h1>
	<form on:submit|preventDefault={register}>
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
		<button type="submit">Register</button>
	</form>
</div>

<style>
	h1 {
		font-size: 24px;
		font-weight: 600;
		margin-top: 1.5rem;
		margin-bottom: 1rem;
	}
</style>
