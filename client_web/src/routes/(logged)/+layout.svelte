<script lang="ts">
	import { afterNavigate } from "$app/navigation";
	import "../../app.css";
	import {
		DownloadIcon,
		Heart,
		Search as SearchIcon,
		SettingsIcon,
	} from "lucide-svelte";
	import Follow from "$lib/components/panel/Follow.svelte";
	import Download from "$lib/components/panel/Download.svelte";
	import Search from "$lib/components/panel/Search.svelte";
	import Setting from "$lib/components/panel/Setting.svelte";

	let { children } = $props();
	let barState = $state<null | string>(null);

	afterNavigate(() => {
		barState = null;
	});
</script>

<header class="flex-wrap fixed top-0 left-0 w-full bg-bg text-fg flex flex-col justify-between items-center z-[1000]">
	<a
		href="/dashboard"
		class="uppercase no-underline text-5xl px-2.5 py-4"
		onclick={() => {
			barState = null;
		}}
	>
		MusicShack
	</a>
	<div class="flex flex-row gap-x-2.5">
		<button
			class="hover-full w-14 h-15"
			onclick={() => {
				barState = barState === "search" ? null : "search";
			}}
			class:active={barState === "search"}
		>
			<SearchIcon />
		</button>
		<button
			class="hover-full w-14 h-15"
			onclick={() => {
				barState = barState === "follow" ? null : "follow";
			}}
			class:active={barState === "follow"}
		>
			<Heart />
		</button>
		<button
			class="hover-full w-14 h-15"
			onclick={() => {
				barState = barState === "download" ? null : "download";
			}}
			class:active={barState === "download"}
		>
			<DownloadIcon />
		</button>
		<button
			class="hover-full w-14 h-15"
			class:active={barState === "settings"}
			onclick={() => {
				barState = barState === "settings" ? null : "settings";
			}}
		>
			<SettingsIcon />
		</button>
	</div>
</header>

<main class="pt-[140px] w-[clamp(320px,70vw+20px,1200px)] mx-auto">
	{#if barState}
		<div class="fixed z-[1000] top-[140px] flex flex-col items-center bg-bg">
			{#if barState === "search"}
				<div class="w-[clamp(320px,70vw+20px,1200px)] max-h-[calc(95vh-135px)] overflow-y-auto outline outline-1 outline-fg">
					<Search />
				</div>
			{:else if barState === "follow"}
				<div class="p-3 w-[clamp(320px,70vw+20px,1200px)] max-h-[calc(95vh-135px)] overflow-y-auto outline outline-1 outline-fg">
					<Follow />
				</div>
			{:else if barState === "download"}
				<div class="p-3 w-[clamp(320px,70vw+20px,1200px)] max-h-[calc(95vh-135px)] overflow-y-auto outline outline-1 outline-fg">
					<Download />
				</div>
			{:else if barState === "settings"}
				<div class="p-3 w-[clamp(320px,70vw+20px,1200px)] max-h-[calc(95vh-135px)] overflow-y-auto outline outline-1 outline-fg">
					<Setting />
				</div>
			{/if}
		</div>
	{/if}

	{@render children?.()}
</main>
