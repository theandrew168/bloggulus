<script lang="ts">
	import { goto } from "$app/navigation";

	export let data;

	// TODO: Add proper validation / error handling
	async function deleteBlog(event: Event) {
		const form = new FormData(event.target as HTMLFormElement);
		const id = form.get("id");
		await fetch(`/api/v1/blogs/${id}`, {
			method: "DELETE",
		});

		await goto("/blogs");
	}
</script>

<div class="container">
	<h1>{data.blog?.title}</h1>
	<div class="links">
		<a href={data.blog.siteURL}>(Site URL)</a>
		<a href={data.blog.feedURL}>(Feed URL)</a>
	</div>
	<div class="times">
		<h2>Synced</h2>
		<div>{new Date(data.blog.syncedAt).toLocaleString()}</div>
	</div>

	<div class="actions">
		<h2>Actions</h2>
		<div class="buttons">
			<form on:submit|preventDefault={deleteBlog}>
				<input type="hidden" name="id" value={data.blog.id} />
				<button type="submit">Delete</button>
			</form>
		</div>
	</div>

	<div class="posts">
		<h2>{data.posts.length} Posts</h2>
		{#each data.posts as post}
			<div class="post">
				<a href="/admin/posts/{post.id}">{post.title}</a>
				<span>{new Date(post.publishedAt).toDateString()}</span>
			</div>
		{/each}
	</div>
</div>

<style>
	h1 {
		font-size: 1.5rem;
		font-weight: 500;
		margin-top: 1.5rem;
		margin-bottom: 0.5rem;
	}
	h2 {
		font-size: 1.5rem;
	}
	.links {
		margin-bottom: 1rem;
	}
	.times {
		margin-bottom: 1rem;
	}
	.actions {
		margin-bottom: 1rem;
	}
	.buttons {
		display: flex;
		gap: 0.5rem;
	}
	.posts {
		margin-bottom: 1rem;
	}
	.post {
		display: flex;
		justify-content: space-between;
	}
</style>
