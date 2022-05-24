<script lang="ts">
	import type { Service } from './types/service'

	export let services: Service[] = []

	function getStatusString(service: Service): string {
		if (service.status === 'offline') {
			return 'Offline'
		} else if (service.status === 'online') {
			return 'Online'
		} else if (service.status === 'pending') {
			return 'Pending'
		}
		return 'Unknown'
	}
</script>

<main>
	<h1>Service Pinger</h1>
	{#each services as service}
		<article>
			<header>{service.name}</header>
			<aside>{service.address}</aside>
			<section class:offline={service.status==='offline'}>{getStatusString(service)}</section>
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
		grid-template-columns: auto auto minmax(0, 1fr);
		grid-template-rows: minmax(0, 1fr);
	}

	header {
		font-weight: bold;
	}

	aside {
		opacity: 0.75;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>
