<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";

	export let data;

	// TODO: Add proper validation / error handling
	async function deletePost() {
		const token = localStorage.getItem("token");
		if (!token) {
			await goto("/login");
		}

		const params = $page.params;
		await fetch(`/api/v1/blogs/${params.blogID}/posts/${params.postID}`, {
			method: "DELETE",
			headers: {
				Authorization: `Bearer ${token}`,
			},
		});

		await goto(`/blogs/${params.blogID}`);
	}
</script>

<div class="container">
	<h1>{data.post?.title}</h1>
	<div class="links">
		<a href={data.post.url}>(URL)</a>
		<a href="/blogs/{data.post.blogID}">(Blog)</a>
	</div>
	<div class="times">
		<h2>Updated</h2>
		<div>{new Date(data.post.publishedAt).toLocaleString()}</div>
	</div>
	<div class="actions">
		<h2>Actions</h2>
		<div class="buttons">
			<form on:submit|preventDefault={deletePost}>
				<button type="submit">Delete</button>
			</form>
		</div>
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
</style>
