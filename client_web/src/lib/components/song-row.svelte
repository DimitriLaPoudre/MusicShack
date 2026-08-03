<script lang="ts">
	import { goto } from "$app/navigation";
	import { ApiError, downloadSong, listDownloads } from "$lib/api";
	import { Button } from "$lib/components/ui/button/index.js";
	import Explicit from "$lib/components/explicit.svelte";
	import { Clock4, Download } from "@lucide/svelte";
	import { downloadList } from "$lib/stores";
	import type { SongInfo } from "$lib/types";

	let {
		provider,
		song,
		number: leading = song.track_number,
		onerror,
	}: {
		provider: string;
		song: SongInfo;
		number?: number;
		onerror?: (message: string | null) => void;
	} = $props();

	function formatDuration(seconds: number) {
		return `${Math.floor(seconds / 60)}:${(seconds % 60).toString().padStart(2, "0")}`;
	}

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

<div class="grid grid-cols-[1fr_auto] gap-2">
	<Button
		variant="hover-soft"
		class="grid w-full grid-cols-[auto_1fr_auto] items-center gap-2"
		onclick={(e) => {
			if (e.target instanceof Element && e.target.closest("a")) return;
			goto(`/song/${provider}/${song.id}`);
		}}
	>
		<p>{leading}</p>
		<div class="flex w-full flex-row flex-wrap items-center justify-center gap-2">
			<p class="flex flex-row flex-wrap items-center justify-center gap-2 font-extrabold">
				{song.title}
				{#if song.explicit}
					<Explicit />
				{/if}
			</p>
			<nav class="italic flex flex-wrap items-center justify-center gap-x-2 gap-y-1">
				{#each song.artists as artist}
					<a href="/artist/{provider}/{artist.id}">
						{artist.name}
					</a>
				{/each}
			</nav>
		</div>
		<div class="flex items-center gap-1">
			<Clock4 class="size-4" />
			<p>{formatDuration(song.duration)}</p>
		</div>
	</Button>
	<Button
		variant="hover-full"
		size="icon"
		class="self-center"
		onclick={startDownload}
	>
		<Download />
	</Button>
</div>
