<script lang="ts">
	import { onMount } from 'svelte'

	import type { Service } from './types/service'

	let services: Service[] = []

	let title: {
		Prefix: string
		Suffix: string
		Name: string
	} = {Prefix: '', Suffix: '', Name: 'pingmon'}
	let lastDate: number = Date.now()
	let lastMoment: string = ''

	async function notifyStatus(service: Service) {
		if (!('Notification' in window)) {
			console.log('notifications unsupported')
			return
		}
		if (Notification.permission !== 'granted') {
			if (await Notification.requestPermission() !== 'granted') {
				console.log('notifications rejected')
				return
			}
		}
		new Notification(`${service.address} ${service.status}`, {
			icon: '/favicon.svg',
			tag: `${service.address}:${service.port}`,
			requireInteraction: true,
			body: `${service.method}:${service.address}:${service.port} is now ${service.status} @ ${service.delay/1000}ms`,
			timestamp: Date.now(),
		})
	}

	async function poll() {
		try {
			let s = await (await fetch("/api/services")).json()
			if (services.length > 0 ) {
				for (let service of services) {
					let t = s.targets.find(v=>v.method === service.method && v.address === service.address && v.port === service.port)
					if (t) {
						if (service.status !== t.status) {
							notifyStatus(t)
						}
					}
				}
			}
			services = s.targets || []
			lastDate = Date.now() - s.elapsed
			updateTitle()
			refresh()
		} catch(err) {
			services = []
			console.error(err)
		}
	}

	const rtf = new Intl.RelativeTimeFormat(undefined, { style: 'long'})
	
	function refresh() {
		lastMoment = rtf.format(-Math.round((Date.now() - lastDate)/1000), 'seconds')
	}

	function updateTitle() {
		document.title = `(${services.filter(v=>v.status==='online').length}/${services.length}) ${title.Prefix}${title.Name}${title.Suffix}`
	}

	// Poll the API every 10 seconds.
	setInterval(poll, 10000)
	setInterval(refresh, 500)

	onMount(async () => {
		poll()
		refresh()

		let s = await (await fetch("/api/title")).json()
		title = s
		updateTitle()
	})
</script>

<main>
	<h1>{title.Name}</h1>
	{#if services.length}
		{#each services as service}
			<article>
				<header>{service.address}</header>
				<section class={service.status}>{service.status}</section>
				<aside>{service.delay/1000}ms</aside>
			</article>
		{/each}
	{:else}
		<article>
			<header></header>
			<section>no services</section>
			<aside></aside>
		</article>
	{/if}
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
