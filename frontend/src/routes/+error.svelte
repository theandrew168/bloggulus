<script lang="ts">
	import { page } from "$app/stores";
</script>

<div class="container">
	<h1>{$page.status}: {$page.error?.message}</h1>
</div>

<style>
	h1 {
		margin-top: 1.5rem;
	}
</style>
