<script lang="ts">
	export let name: string;
</script>

<a href="?q={name}">{name}</a>

<style>
	a {
		color: var(--light-color);
		background-color: var(--mid-color);
		font-weight: 600;
		padding: 0.25rem 0.75rem;
		border-radius: 0.25rem;
		text-decoration: none;
	}
</style>
