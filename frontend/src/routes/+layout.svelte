<script lang="ts">
	import "$lib/styles/reset.css";
	import "$lib/styles/fonts.css";
	import "$lib/styles/style.css";
	import Header from "$lib/components/Header.svelte";
	import Footer from "$lib/components/Footer.svelte";
</script>

<div class="sans-serif">
	<Header />
	<main>
		<slot />
	</main>
	<Footer />
</div>

<style>
	div {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}
	main {
		flex-grow: 1;
	}
</style>
