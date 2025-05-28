<script lang="ts">
	import type { PageProps } from "./$types";

	let { data }: PageProps = $props();
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>

<a href="/signin/debug">Sign In (Debug)</a>
<a href="/signin/github">Sign In (GitHub)</a>
<a href="/signin/google">Sign In (Google)</a>

<ul>
	{#each data.articles as article}
		<li>
			{article.title}
		</li>
	{/each}
</ul>
