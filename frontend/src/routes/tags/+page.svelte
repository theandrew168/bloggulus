<script lang="ts">
	import type { TagResponse } from "$lib/types.js";

	export let data;

	let name = "";

	async function createTag(name: string) {
		const resp = await fetch(`/api/v1/tags`, {
			method: "POST",
			body: JSON.stringify({ name }),
		});
		const tagResp: TagResponse = await resp.json();
		data.tags = [tagResp.tag, ...data.tags];
	}

	async function deleteTag(id: string) {
		await fetch(`/api/v1/tags/${id}`, {
			method: "DELETE",
		});
		data.tags = data.tags.filter((tag) => tag.id !== id);
	}
</script>

<div class="container">
	<h1>Tags</h1>
	<div class="add">
		<form on:submit|preventDefault={async () => await createTag(name)}>
			<input bind:value={name} placeholder="Name" />
			<button type="submit">Add</button>
		</form>
	</div>
	<div class="tags">
		{#each data.tags as tag}
			<div class="tag">
				<button on:click={async () => await deleteTag(tag.id)}>Delete</button>
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
