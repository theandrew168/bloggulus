<script lang="ts">
	import { page } from "$app/stores";
	import Post from "$lib/components/Post.svelte";

	export let data;

	$: q = $page.url.searchParams.get("q") ?? "";
	$: p = parseInt($page.url.searchParams.get("p") ?? "1") || 1;
	$: moreLink = `/?p=${p + 1}` + (q ? `&q=${q}` : "");
</script>

<div class="container">
	<h1>Recent Posts</h1>
	<div class="posts">
		{#each data.posts as post}
			<Post {post} />
		{/each}
	</div>
	{#if data.posts.length === 15}
		<div class="more">
			<a class="shadow" href={moreLink}>See More</a>
		</div>
	{/if}
</div>

<style>
	h1 {
		color: var(--dark-color);
		font-size: 24px;
		font-weight: 600;
		margin-top: 1.5rem;
		margin-bottom: 1.5rem;
	}
	.posts {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
		margin-bottom: 1.5rem;
	}
	.more {
		display: flex;
		justify-content: center;
		align-items: center;
		margin-bottom: 1.5rem;
	}
	.more a {
		padding: 0.5rem 1.5rem;
		border-radius: 0.25rem;

		font-weight: 600;
		text-decoration: none;
		color: var(--dark-color);
		background-color: var(--light-color);
	}
</style>
