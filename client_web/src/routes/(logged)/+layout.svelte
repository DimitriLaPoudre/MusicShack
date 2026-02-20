<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import "../../app.css";
	import {
		DownloadIcon,
		Heart,
		Library,
		Search as SearchIcon,
		SettingsIcon,
	} from "lucide-svelte";
	import Follow from "$lib/components/panel/Follow.svelte";
	import Download from "$lib/components/panel/Download.svelte";
	import Search from "$lib/components/panel/Search.svelte";
	import Setting from "$lib/components/panel/Setting.svelte";
	import { page } from "$app/state";
	import { syncLibrary } from "$lib/functions/library";

	let { children } = $props();
	let barState = $state<null | string>(null);

	afterNavigate(() => {
		barState = null;
	});
</script>

<header
	class="flex-wrap fixed top-0 left-0 w-full bg-bg text-fg flex flex-col justify-between items-center z-50"
>
	<a
		href="/dashboard"
		class="no-underline text-5xl px-2.5 py-4 flex"
		onclick={() => {
			barState = null;
		}}
	>
		<span class="italic font-extrabold">MUSIC</span>
		<span class="italic font-thin">SHACK</span>
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
				goto("/library");
			}}
			class:active={page.url.pathname == "/library"}
		>
			<Library />
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

<main class="pt-35 w-[clamp(320px,70vw+20px,1200px)] mx-auto">
	{#if barState}
		<div class="fixed z-1000 top-35 flex flex-col items-center bg-bg">
			{#if barState === "search"}
				<div
					class="w-[clamp(320px,70vw+20px,1200px)] max-h-[calc(100vh-140px-3rem)] overflow-y-auto outline outline-fg"
				>
					<Search />
				</div>
			{:else if barState === "follow"}
				<div
					class="p-3 w-[clamp(320px,70vw+20px,1200px)] max-h-[calc(100vh-140px-3rem)] overflow-y-auto outline outline-fg"
				>
					<Follow />
				</div>
			{:else if barState === "download"}
				<div
					class="p-3 w-[clamp(320px,70vw+20px,1200px)] max-h-[calc(100vh-140px-3rem)] overflow-y-auto outline outline-fg"
				>
					<Download />
				</div>
			{:else if barState === "settings"}
				<div
					class="p-3 w-[clamp(320px,70vw+20px,1200px)] max-h-[calc(100vh-140px-3rem)] overflow-y-auto outline outline-fg"
				>
					<Setting />
				</div>
			{/if}
		</div>
	{/if}

	{@render children?.()}
	<div class="h-8"></div>
</main>
