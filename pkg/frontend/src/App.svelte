<script lang="ts">
	export let name: string;
</script>

<main>
	<h1>Service Pinger</h1>
	{#each services as service}
		<article>
			<header>{service.name}</header>
			<section class:offline={service.offline}>{service.offline?'Offline':'Online'}</section>
		</article>
	{/each}
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	article {
		display: grid;
		grid-template-columns: auto minmax(0, 1fr);
		grid-template-rows: minmax(0, 1fr);
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>
