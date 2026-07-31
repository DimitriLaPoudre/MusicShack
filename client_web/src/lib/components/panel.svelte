<script lang="ts">
	import Search from "./panels/search.svelte";

	let panelEl: HTMLElement | undefined = $state();
</script>

<div class="flex" bind:this={panelEl}>
	<Search anchor={panelEl} />
</div>
