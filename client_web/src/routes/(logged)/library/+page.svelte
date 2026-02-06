<script lang="ts">
	import {
		deleteSong,
		loadLibrary,
		syncLibrary,
	} from "$lib/functions/library";
	import type { ResponseLibrary } from "$lib/types/response";
	import { Trash } from "lucide-svelte";
	import { onMount } from "svelte";

	let error = $state<null | string>(null);
	let page = $state<null | ResponseLibrary>(null);

	onMount(async () => {
		await syncLibrary();
		({ page, error } = await loadLibrary());
	});
</script>

<svelte:head>
	<title>Library - MusicShack</title>
</svelte:head>

{#if error}
	<div class="mt-4 flex flex-col justify-center items-center gap-2.5">
		<h2>Error loading Song</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !page}
	<p class="mt-6 text-center">Loading...</p>
{:else}
	<!-- page top -->
	<div class="mt-4 flex flex-row items-center justify-center"></div>
	<div class="flex flex-col gap-2 items-center">
		{#each page.items as item}
			<div class="grid grid-cols-[1fr_auto] gap-2">
				<div class="flex items-center gap-2 w-full">
					<p class="font-extrabold">{item.title}</p>
					<div class="flex gap-1 italic">
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
							({ page, error } = await loadLibrary());
						}
					}}
				>
					<Trash />
				</button>
			</div>
		{/each}
	</div>
{/if}
