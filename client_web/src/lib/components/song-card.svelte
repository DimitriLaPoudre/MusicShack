<script lang="ts">
	import { goto } from "$app/navigation";
	import { ApiError, downloadSong, listDownloads } from "$lib/api";
	import { Button } from "$lib/components/ui/button/index.js";
	import Explicit from "$lib/components/explicit.svelte";
	import Quality from "$lib/components/quality.svelte";
	import { Disc, Download } from "@lucide/svelte";
	import { downloadList } from "$lib/stores";
	import type { SongInfo } from "$lib/types";

	let {
		provider,
		song,
		onerror,
	}: {
		provider: string;
		song: SongInfo;
		onerror?: (message: string | null) => void;
	} = $props();

	async function startDownload() {
		try {
			await downloadSong({ provider, id: song.id });
			onerror?.(null);
		} catch (e) {
			onerror?.(e instanceof ApiError ? e.message : "Failed to download song");
		}
		try {
			downloadList.set(await listDownloads());
		} catch {
			// keep cached list on refresh failure
		}
	}
</script>

<div class="w-50 h-auto">
	<Button
		variant="hover-full"
		class="flex w-50 h-auto flex-col items-center gap-2 py-5 overflow-hidden shadow-[inset_0_1px_0_var(--foreground),inset_1px_0_0_var(--foreground),inset_-1px_0_0_var(--foreground)]"
		onclick={(e) => {
			if (e.target instanceof Element && e.target.closest("a")) return;
			goto(`/song/${provider}/${song.id}`);
		}}
	>
		<div class="size-40">
			{#if song.album.cover_url}
				<img
					class="size-full object-cover"
					src={song.album.cover_url}
					alt={song.title}
				/>
			{:else}
				<div class="grid size-full place-items-center bg-muted">
					<Disc class="size-24" />
				</div>
			{/if}
		</div>
		<p class="flex flex-row flex-wrap items-center justify-center gap-2 break-words font-extrabold">
			{song.title}
			{#if song.explicit}
				<Explicit />
			{/if}
		</p>
		<nav class="flex flex-col gap-y-[0.2rem] gap-x-4 italic">
			{#each song.artists as artist}
				<a href="/artist/{provider}/{artist.id}">{artist.name}</a>
			{/each}
		</nav>
		<Quality quality={song.audio_quality} />
	</Button>
	<Button
		variant="hover-full"
		class="w-full p-3 shadow-[inset_0_-1px_0_var(--foreground),inset_1px_0_0_var(--foreground),inset_-1px_0_0_var(--foreground)]"
		onclick={startDownload}
	>
		<Download />
	</Button>
</div>
