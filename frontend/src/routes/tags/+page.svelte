<script lang="ts">
	import { superForm, setMessage, setError } from "sveltekit-superforms";
	import { zod } from "sveltekit-superforms/adapters";

	import type { TagResponse } from "$lib/types.js";
	import { TagSchema } from "$lib/schemas.js";

	export let data;

	const { form, errors, message, constraints, enhance } = superForm(data.form, {
		SPA: true,
		validators: zod(TagSchema),
		async onUpdate({ form }) {
			// Form validation
			if (!form.data.name) {
				setError(form, "name", "Must not be empty.");
			} else if (form.valid) {
				const resp = await fetch(`/api/v1/tags`, {
					method: "POST",
					body: JSON.stringify({ name }),
				});
				const tagResp: TagResponse = await resp.json();
				data.tags = [tagResp.tag, ...data.tags];

				// TODO: Call an external API with form.data, await the result and update form
				setMessage(form, "Valid data!");
			}
		},
	});

	let name = "";

	// async function createTag(name: string) {
	// 	const resp = await fetch(`/api/v1/tags`, {
	// 		method: "POST",
	// 		body: JSON.stringify({ name }),
	// 	});
	// 	const tagResp: TagResponse = await resp.json();
	// 	data.tags = [tagResp.tag, ...data.tags];
	// }

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
		<form method="POST" use:enhance>
			<label>
				Add tag
				<input placeholder="Name" bind:value={$form.name} {...$constraints.name} />
			</label>
			{#if $errors.name}<span>{$errors.name}</span>{/if}
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
