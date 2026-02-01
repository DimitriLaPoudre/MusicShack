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
	import Download from "$lib/components/panel/Download.svelte";
	import Search from "$lib/components/panel/Search.svelte";
	import Setting from "$lib/components/panel/Setting.svelte";

	let { children } = $props();
	let barState = $state<null | string>(null);

	afterNavigate(() => {
		barState = null;
	});
</script>

<header>
	<a
		href="/dashboard"
		class="logo"
		onclick={() => {
			barState = null;
		}}
	>
		MusicShack
	</a>
	<div class="bar">
		<button
			class="hover-full"
			onclick={() => {
				barState = barState === "search" ? null : "search";
			}}
			class:active={barState === "search"}
		>
			<SearchIcon />
		</button>
		<button
			class="hover-full"
			onclick={() => {
				barState = barState === "follow" ? null : "follow";
			}}
			class:active={barState === "follow"}
		>
			<Heart />
		</button>
		<button
			class="hover-full"
			onclick={() => {
				barState = barState === "download" ? null : "download";
			}}
			class:active={barState === "download"}
		>
			<DownloadIcon />
		</button>
		<button
			class="hover-full"
			class:active={barState === "settings"}
			onclick={() => {
				barState = barState === "settings" ? null : "settings";
			}}
		>
			<SettingsIcon />
		</button>
	</div>
</header>

<main>
	{#if barState}
		<div class="wrap-panel">
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
					<Download />
				</div>
			{:else if barState === "settings"}
				<div class="panel-default">
					<Setting />
				</div>
			{/if}
		</div>
	{/if}

	{@render children?.()}
</main>

<style>
	header {
		flex-wrap: wrap;
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		background-color: var(--bg);
		color: var(--fg);
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		align-items: center;
		z-index: 1000;

		.logo {
			text-transform: uppercase;
			text-decoration: none;
			font-size: 3rem;
			padding: 10px;
		}

		.bar {
			display: flex;
			flex-direction: row;
			gap: 0 10px;
			button {
				aspect-ratio: 1/1;
			}
		}
	}

	main {
		.wrap-panel {
			position: fixed;
			z-index: 1000;
			top: 144px;
			display: flex;
			flex-direction: column;
			align-items: center;
			background-color: var(--bg);

			.panel-search {
				width: clamp(320px, 70vw + 20px, 1200px);
				max-height: calc(95vh - 135px);
				overflow-y: auto;
				outline: 1px solid var(--fg);
			}
			.panel-default {
				padding: 0.75rem;
				width: clamp(320px, 70vw + 20px, 1200px);
				max-height: calc(95vh - 135px);
				overflow-y: auto;
				outline: 1px solid var(--fg);
			}
		}
		padding-top: 130px;
		width: clamp(320px, 70vw + 20px, 1200px);
		margin: 0 auto;
	}
</style>
