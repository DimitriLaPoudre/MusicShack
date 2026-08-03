<script lang="ts">
	import { page } from "$app/state";
	import { ApiError, downloadSong, getSongInfo, listDownloads } from "$lib/api";
	import { Button } from "$lib/components/ui/button/index.js";
	import Explicit from "$lib/components/explicit.svelte";
	import Quality from "$lib/components/quality.svelte";
	import { Clock4, Download } from "@lucide/svelte";
	import { downloadList } from "$lib/stores";
	import type { SongInfo } from "$lib/types";

	let song = $state<SongInfo | null>(null);
	let error = $state<string | null>(null);

	const provider = $derived(page.params.provider);
	const id = $derived(page.params.id);

	async function load() {
		if (!provider || !id) return;
		song = null;
		error = null;
		try {
			song = await getSongInfo(provider, id);
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to load song";
		}
	}

	$effect(() => {
		if (provider && id) load();
	});

	function formatDuration(seconds: number) {
		return `${Math.floor(seconds / 60)}:${(seconds % 60).toString().padStart(2, "0")}`;
	}

	async function startDownload() {
		if (!provider || !id) return;
		try {
			await downloadSong({ provider, id });
			error = null;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to download song";
		}
		try {
			downloadList.set(await listDownloads());
		} catch {
			// keep cached list on refresh failure
		}
	}
</script>

<svelte:head>
	<title
		>{song?.title || "Song"} | {song?.artists.map((artist) => artist.name).join(" ") ||
			"Artist"} - MusicShack</title
	>
</svelte:head>

{#if error}
	<div class="flex flex-col items-center justify-center gap-2 py-4">
		<h1 class="text-xs font-semibold tracking-widest uppercase">
			Error loading Song
		</h1>
		<p class="text-sm text-muted-foreground">{error}</p>
		<a class="underline" href="/dashboard">Go to Home</a>
	</div>
{:else if !song}
	<p class="py-4 text-center text-sm text-muted-foreground">Loading...</p>
{:else}
	<div class="mt-1 mx-auto table border-separate border-spacing-y-2.5">
		<div class="table-row">
			<div class="flex flex-row flex-wrap justify-center gap-2.5">
				{#if song.album.cover_url}
					<img
						class="size-70 object-cover"
						src={song.album.cover_url}
						alt={song.title}
					/>
				{:else}
					<div class="size-70 bg-muted"></div>
				{/if}
				<div class="flex flex-col gap-1.5">
					<h1 class="flex flex-row items-center gap-2 font-extrabold">
					{song.title}
					{#if song.explicit}
						<Explicit />
					{/if}
					</h1>
					<a
						class="font-extrabold italic no-underline"
						href="/album/{provider}/{song.album.id}"
					>
						{song.album.title}
					</a>
					<nav class="italic flex flex-wrap gap-x-2 gap-y-1">
						{#each song.artists as artist}
							<a href="/artist/{provider}/{artist.id}">
								{artist.name}
							</a>
						{/each}
					</nav>
					<br />
					<div class="flex items-center gap-1">
						<Clock4 class="size-4" />
						<p>{formatDuration(song.duration)}</p>
					</div>
				<Quality quality={song.audio_quality} />
				</div>
			</div>
		</div>
		<Button
			variant="hover-full"
			size="default"
			class="w-full"
			onclick={startDownload}
		>
			<Download />
			Download Song
		</Button>
	</div>
{/if}
