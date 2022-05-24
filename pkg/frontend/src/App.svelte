<script lang="ts">
	import { onMount } from 'svelte'
	import moment from 'moment'

	import type { Service } from './types/service'

	let services: Service[] = []

	let title: string = 'pingmon'
	let lastDate: Date = new Date()
	let lastMoment: string = ''

	async function poll() {
		let s = await fetch("/api/services")
		services = await s.json()
		lastDate = new Date()
		refresh()
	}
	
	function refresh() {
		lastMoment = moment(lastDate).fromNow()
	}

	// Poll the API every 30 seconds.
	setTimeout(poll, 30000)
	setTimeout(refresh, 500)

	onMount(async () => {
		poll()
		refresh()

		let s = await (await fetch("/api/title")).json()
		title = s.Name
		document.title = `${s.Prefix}${s.Name}${s.Suffix}`
	})
</script>

<main>
	<h1>{title}</h1>
	{#each services as service}
		<article>
			<header>{service.address}</header>
			<section class={service.status}>{service.status}</section>
			<aside>{service.delay/1000}ms</aside>
		</article>
	{/each}
	<h4>updated {lastMoment}</h4>
</main>
<footer><a href="https://github.com/kettek/pingmon">pingmon</a> copyright © 2022 <a href="https://kettek.net">Ketchetwahmeegwun T. Southall</a>. Licensed under the <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">GPLv3</a></footer>

<style>
	main {
		text-align: center;
		padding: 1em;
		margin: 0 auto;
		max-width: 35em;
	}

	h1 {
		color: #003eff;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	h4 {
		font-weight: 100;
	}

	article {
		display: grid;
		grid-template-columns: minmax(0, 1fr) 10em minmax(0, 1fr);
		grid-template-rows: minmax(0, 1fr);
	}

	header {
		font-weight: bold;
		text-align: left;
	}

	section.offline {
		color: rgb(255, 0, 0);
	}
	section.online {
		color: rgba(0, 255, 0);
	}
	section.pending {
		color: rgb(255, 255, 0);
	}

	aside {
		opacity: 0.75;
		text-align: right;
	}

	footer {
		text-align: center;
		opacity: 0.75;
		font-size: 80%;
	}
</style>
