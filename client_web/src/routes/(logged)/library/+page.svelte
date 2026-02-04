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
	<div class="mt-8 flex flex-col justify-center items-center gap-2.5">
		<h2>Error loading Song</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !list}
	<p class="mt-8 text-center">Loading...</p>
{:else}
	<!-- page top -->
	<div class="mt-8 flex flex-row items-center justify-center">
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
	</div>
	<div>
		{#each list as item}
			<div class="flex items-center gap-2">
				<p>{item.title}</p>
				<div class="flex gap-1">
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
{/if}
