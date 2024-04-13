<script lang="ts">
	import { page } from "$app/stores";

	$: q = $page.url.searchParams.get("q") ?? "";
</script>

<header class="shadow">
	<div class="container">
		<div class="items">
			<a class="home" href="/">Bloggulus</a>
			<form data-sveltekit-keepfocus>
				<input value={q} name="q" type="text" placeholder="Search" />
			</form>
		</div>
	</div>
</header>

<style>
	header {
		background-color: var(--light-color);
		padding: 1rem 0;
	}
	.items {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	.home {
		color: var(--dark-color);
		font-size: 1.5rem;
		text-decoration: none;
	}
	input {
		padding: 0.5rem;
		border: 1px solid black;
		border-radius: 0.375rem;
	}
</style>
