<script lang="ts">
	import { goto, invalidateAll } from "$app/navigation";

	export let data;

	// TODO: Add proper validation / error handling
	async function createBlog(event: Event) {
		const token = localStorage.getItem("token");
		if (!token) {
			await goto("/login");
		}

		const form = new FormData(event.target as HTMLFormElement);
		const feedURL = form.get("feedURL");
		await fetch(`/api/v1/blogs`, {
			method: "POST",
			headers: {
				Authorization: `Bearer ${token}`,
			},
			body: JSON.stringify({ feedURL }),
		});

		await invalidateAll();
	}
</script>

<div class="container">
	<h1>Blogs</h1>
	<div class="add">
		<form on:submit|preventDefault={createBlog}>
			<input name="feedURL" placeholder="RSS / Atom Feed URL" />
			<button type="submit">Add</button>
		</form>
	</div>
	<div class="blogs">
		{#each data.blogs as blog}
			<div>
				<a href="/blogs/{blog.id}">{blog.title}</a>
			</div>
		{/each}
	</div>
</div>

<style>
	h1 {
		font-size: 24px;
		font-weight: 600;
		margin-top: 1.5rem;
		margin-bottom: 0.5rem;
	}
	.add {
		margin-bottom: 1rem;
	}
	.blogs {
		margin-bottom: 1rem;
	}
</style>
