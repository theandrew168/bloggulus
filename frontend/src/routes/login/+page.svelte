<script lang="ts">
	import { goto } from "$app/navigation";
	import type { TokenResponse } from "$lib/types";

	// TODO: Add proper validation / error handling
	async function login(event: Event) {
		const form = new FormData(event.target as HTMLFormElement);
		const username = form.get("username");
		const password = form.get("password");
		const resp = await fetch(`/api/v1/tokens`, {
			method: "POST",
			body: JSON.stringify({ username, password }),
		});

		if (resp.ok) {
			const token: TokenResponse = await resp.json();
			localStorage.setItem("token", token.token.value);
			await goto("/");
		}
	}
</script>

<div class="container">
	<h1>Login</h1>
	<form on:submit|preventDefault={login}>
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
		<button type="submit">Login</button>
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
