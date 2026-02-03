<script lang="ts">
	import {
		deleteSong,
		loadLibrary,
		syncLibrary,
	} from "$lib/functions/library";
	import type { ResponseSong } from "$lib/types/response";
	import { Trash } from "lucide-svelte";
	import { onMount } from "svelte";

	let error = $state<null | string>(null);
	let list = $state<null | ResponseSong[]>(null);

	onMount(async () => {
		({ list, error } = await loadLibrary());
	});
</script>

<svelte:head>
	<title>Library - MusicShack</title>
</svelte:head>

{#if error}
	<div class="error">
		<h2>Error loading Song</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !list}
	<p class="loading">Loading...</p>
{:else}
	<!-- page top -->
	<div class="header">
		<button
			class="hover-full"
			onclick={async () => {
				error = await syncLibrary();
				if (!error) {
					({ list, error } = await loadLibrary());
				}
			}}
		>
			Sync
		</button>
		<div class="wrap-item">
			{#each list as item}
				<div class="item">
					<p>{item.title}</p>
					<div class="artists">
						{#each item.artists as artist}
							<p>{artist}</p>
						{/each}
					</div>
				</div>
				<button
					class="hover-full"
					onclick={async () => {
						error = await deleteSong(item.id);
						if (!error) {
							({ list, error } = await loadLibrary());
						}
					}}
				>
					<Trash />
				</button>
			{/each}
		</div>
	</div>
{/if}

<style>
	.loading {
		margin-top: 30px;
		text-align: center;
	}

	.error {
		margin-top: 30px;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 10px;
	}

	.header {
		margin: 15px auto 0;
		display: table;
		border-spacing: 0 10px;
	}
</style>
