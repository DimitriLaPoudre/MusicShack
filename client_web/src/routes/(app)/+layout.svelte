<script lang="ts">
	import Logo from "$lib/components/logo.svelte";
	import Panel from "$lib/components/panel.svelte";

	let { children } = $props();
</script>

<header
	class="fixed top-0 left-0 z-50 flex w-full flex-col items-center justify-between bg-background text-foreground flex-wrap"
>
	<Logo />
	<Panel />
</header>

<main class="mx-auto w-[clamp(320px,80vw,1200px)] pt-35">
	{@render children?.()}
	<div class="h-8"></div>
</main>
