<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>

{#each data.posts as post}
	<a href={post.url}>{post.title}</a>
	<a href={post.blogURL}>{post.blogTitle}</a>
	<p>{post.publishedAt}</p>
	<p>{post.tags}</p>
{/each}
