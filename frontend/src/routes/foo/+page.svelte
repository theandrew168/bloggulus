<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<h1>Page Foo!</h1>

{#each data.posts as post}
	<a href={post.url}>{post.title}</a>
	<a href={post.blogURL}>{post.blogTitle}</a>
	<p>{post.publishedAt}</p>
	<p>{post.tags}</p>
{/each}
