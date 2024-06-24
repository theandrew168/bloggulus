<script lang="ts">
	import { page } from "$app/stores";
	import Article from "$lib/components/Article.svelte";

	export let data;

	$: q = $page.url.searchParams.get("q") ?? "";
	$: p = parseInt($page.url.searchParams.get("p") ?? "1") || 1;
	$: moreLink = `/?p=${p + 1}` + (q ? `&q=${q}` : "");
	$: hasMorePages = p * 20 < data.count;
</script>

<div class="container">
	{#if q === ""}
		<h1>Recent Articles</h1>
	{:else}
		<h1>Relevant Articles</h1>
	{/if}
	<div class="articles">
		{#each data.articles as article}
			<Article {article} />
		{/each}
	</div>
	{#if hasMorePages}
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
	.articles {
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
