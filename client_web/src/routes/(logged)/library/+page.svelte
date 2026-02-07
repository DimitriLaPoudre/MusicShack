<script lang="ts">
	import { afterNavigate } from "$app/navigation";
	import Explicit from "$lib/components/explicit.svelte";
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

	afterNavigate(async () => {
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
	<div class="grid grid-cols-[repeat(auto-fit,200px)] justify-center gap-4">
		{#each page.items as item}
			<div class="w-[200px] h-auto">
				<button
					class="hover-full flex flex-col items-center w-[200px] h-auto overflow-hidden gap-3 shadow-[inset_0_1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
				>
					<div class="w-[160px] h-[160px]">
						<img
							src="/api/library/{item.id}/img"
							alt={item.title}
						/>
					</div>
					<p
						class="flex flex-row items-center justify-center gap-2 font-extrabold"
					>
						{item.title}
						{#if item.explicit}
							<Explicit />
						{/if}
					</p>
					<nav class="flex flex-col gap-y-[0.2rem] gap-x-4 italic">
						{#each item.artists as artist}
							<span>
								{artist}
							</span>
						{/each}
					</nav>
				</button>
				<button
					class="hover-full w-full p-3 shadow-[inset_0_-1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
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
