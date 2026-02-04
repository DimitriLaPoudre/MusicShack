<script lang="ts">
	import { loadFollows, removeFollow } from "$lib/functions/follow";
	import { followList } from "$lib/stores/panel/follow";
	import { HeartIcon, HeartOff } from "lucide-svelte";
	import { onMount } from "svelte";

	let error = $state<null | string>(null);

	onMount(() => {
		async function intervalFunc() {
			error = await loadFollows();
		}
		intervalFunc();
		const interval = setInterval(intervalFunc, 500);
		return () => {
			clearInterval(interval);
		};
	});
</script>

<div class="flex flex-col gap-3">
	<h1 class="font-extrabold">Followed Artist</h1>
	{#if error}
		<p class="text-center bg-err p-2 m-0">{error}</p>
	{/if}
	{#if !$followList}
		<p class="text-center">Loading...</p>
	{:else}
		<div class="flex flex-col gap-1">
			{#each $followList as item}
				<div class="grid grid-cols-[1fr_auto] gap-2 group">
					<a
						class="grid grid-cols-[auto_1fr] items-stretch gap-2 group/data"
						href="/artist/{item.provider}/{item.artistId}"
					>
						<img
							class="w-[58px] h-[58px] self-center"
							src={item.artistPictureUrl}
							alt={item.artistName}
						/>
						<p class="italic pl-3 flex items-center group-hover/data:outline group-hover/data:outline-1 group-hover/data:outline-fg group-hover/data:-outline-offset-1">
							{item.artistName}
						</p>
					</a>
					<button
						class="hover-full"
						onclick={async () => {
							error = await removeFollow(item.id);
							if (!error) {
								error = await loadFollows();
							}
						}}
					>
						<div class="block group-hover:hidden">
							<HeartIcon />
						</div>
						<div class="hidden group-hover:block">
							<HeartOff />
						</div>
					</button>
				</div>
			{/each}
		</div>
	{/if}
</div>
