<script lang="ts">
	import { goto, invalidateAll } from "$app/navigation";

	export let data;

	// TODO: Add proper validation / error handling
	async function createTag(event: Event) {
		const token = localStorage.getItem("token");
		if (!token) {
			await goto("/login");
		}

		const form = new FormData(event.target as HTMLFormElement);
		const name = form.get("name");
		await fetch(`/api/v1/tags`, {
			method: "POST",
			headers: {
				Authorization: `Bearer ${token}`,
			},
			body: JSON.stringify({ name }),
		});

		await invalidateAll();
	}

	// TODO: Add proper validation / error handling
	async function deleteTag(event: Event) {
		const token = localStorage.getItem("token");
		if (!token) {
			await goto("/login");
		}

		const form = new FormData(event.target as HTMLFormElement);
		const id = form.get("id");
		await fetch(`/api/v1/tags/${id}`, {
			method: "DELETE",
			headers: {
				Authorization: `Bearer ${token}`,
			},
		});

		await invalidateAll();
	}
</script>

<div class="container">
	<h1>Tags</h1>
	<div class="add">
		<form on:submit|preventDefault={createTag}>
			<label>
				Add tag:&nbsp;
				<input name="name" placeholder="Name" />
			</label>
			<button type="submit">Add</button>
		</form>
	</div>
	<div class="tags">
		{#each data.tags as tag}
			<div class="tag">
				<form on:submit|preventDefault={deleteTag}>
					<input type="hidden" name="id" value={tag.id} />
					<button type="submit">Delete</button>
				</form>
				<p>{tag.name}</p>
			</div>
		{/each}
	</div>
</div>

<style>
	h1 {
		font-size: 24px;
		font-weight: 600;
		margin-top: 1.5rem;
		margin-bottom: 1rem;
	}
	.add {
		margin-bottom: 1rem;
	}
	.tags {
		margin-bottom: 1rem;
	}
	.tag {
		display: flex;
		gap: 1rem;
	}
</style>
