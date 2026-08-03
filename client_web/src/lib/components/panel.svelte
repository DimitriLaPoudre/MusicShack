<script lang="ts">
	import { afterNavigate } from "$app/navigation";
	import Search from "./panels/search.svelte";
	import Follow from "./panels/follow.svelte";
	import Download from "./panels/download.svelte";
	import Settings from "./panels/settings.svelte";
	import Admin from "./panels/admin.svelte";
	import { openPanel } from "$lib/stores";

	let panelEl: HTMLElement | undefined = $state();

	afterNavigate(() => {
		openPanel.set(null);
	});
</script>

<div class="flex" bind:this={panelEl}>
	<Search anchor={panelEl} />
	<Follow anchor={panelEl} />
	<Download anchor={panelEl} />
	<Settings anchor={panelEl} />
	<Admin anchor={panelEl} />
</div>
