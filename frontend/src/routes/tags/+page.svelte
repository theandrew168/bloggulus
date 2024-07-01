<script lang="ts">
	import { superForm } from "sveltekit-superforms";
	import { zod } from "sveltekit-superforms/adapters";

	import { NewTagSchema } from "$lib/schemas.js";

	export let data;

	const { form, errors, enhance } = superForm(data.form, {
		SPA: true,
		validators: zod(NewTagSchema),
		async onUpdate({ form }) {
			// If the form isn't valid on the client-side, don't submit.
			if (!form.valid) {
				return;
			}

			await fetch(`/api/v1/tags`, {
				method: "POST",
				body: JSON.stringify({ name: form.data.name }),
			});
		},
	});

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
		<form use:enhance>
			<label>
				<input bind:value={$form.name} placeholder="Name" />
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
