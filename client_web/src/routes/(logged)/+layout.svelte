<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import "../../app.css";
	import {
		DownloadIcon,
		Heart,
		Search as SearchIcon,
		SettingsIcon,
	} from "lucide-svelte";
	import Follow from "$lib/components/panel/Follow.svelte";
	import Downloads from "$lib/components/panel/Downloads.svelte";
	import Search from "$lib/components/panel/Search.svelte";
	import Settings from "$lib/components/panel/Settings.svelte";

	let { children } = $props();
	let barState = $state<null | string>(null);

	afterNavigate(() => {
		barState = null;
	});
</script>

<header>
	<button
		class="logo"
		onclick={() => {
			barState = null;
			goto("/dashboard");
		}}
	>
		MusicShack
	</button>
	<div class="bar">
		<button
			onclick={() => {
				barState = barState === "search" ? null : "search";
			}}
			class:active={barState === "search"}
		>
			<SearchIcon />
		</button>
		<button
			onclick={() => {
				barState = barState === "follow" ? null : "follow";
			}}
			class:active={barState === "follow"}
		>
			<Heart />
		</button>
		<button
			onclick={() => {
				barState = barState === "download" ? null : "download";
			}}
			class:active={barState === "download"}
		>
			<DownloadIcon />
		</button>
		<button
			onclick={() => {
				barState = barState === "settings" ? null : "settings";
			}}
			class:active={barState === "settings"}
		>
			<SettingsIcon />
		</button>
	</div>

	{#if barState === "search"}
		<div class="panel-search">
			<Search />
		</div>
	{:else if barState === "follow"}
		<div class="panel-default">
			<Follow />
		</div>
	{:else if barState === "download"}
		<div class="panel-default">
			<Downloads />
		</div>
	{:else if barState === "settings"}
		<div class="panel-default">
			<Settings />
		</div>
	{/if}
</header>

<main>
	{@render children?.()}
</main>

<style>
	header {
		flex-wrap: wrap;
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		background-color: #0e0e0e;
		color: #fff;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		align-items: center;
		z-index: 1000;

		.logo {
			text-transform: uppercase;
			font-size: 3rem;
			border: none;
			padding: 10px;
		}
		.logo:hover {
			background-color: inherit;
			color: inherit;
		}
		.logo:active {
			background-color: inherit;
			color: inherit;
			outline: none;
		}

		.bar {
			display: flex;
			flex-direction: row;
			gap: 0px 10px;
			button {
				aspect-ratio: 1/1;
			}
		}
		.panel-search {
			width: clamp(320px, 70vw + 20px, 1200px);
			max-height: calc(95vh - 135px);
			overflow-y: auto;
			outline: 1px solid #ffffff;
			display: flex;
			flex-direction: column;
			gap: 8px;
		}
		.panel-default {
			width: clamp(320px, 70vw + 20px, 1200px);
			max-height: calc(95vh - 135px);
			overflow-y: auto;
			outline: 1px solid #ffffff;
			display: flex;
			flex-direction: column;
			gap: 8px;
			padding: 8px;
		}
	}

	main {
		padding-top: 130px;
		width: clamp(320px, 70vw + 20px, 1200px);
		margin: 0 auto;
	}
</style>
