<script lang="ts">
	import { page } from "$app/state";
	import {
		ApiError,
		downloadAlbum,
		getAlbumInfo,
		getAlbumSongs,
		listDownloads,
	} from "$lib/api";
	import { Button } from "$lib/components/ui/button/index.js";
	import Explicit from "$lib/components/explicit.svelte";
	import LoadMore from "$lib/components/load-more.svelte";
	import Quality from "$lib/components/quality.svelte";
	import SongRow from "$lib/components/song-row.svelte";
	import { Clock4, Download } from "@lucide/svelte";
	import { downloadList } from "$lib/stores";
	import type { AlbumInfo, SongInfo } from "$lib/types";

	let album = $state<AlbumInfo | null>(null);
	let songs = $state<SongInfo[]>([]);
	let total = $state(0);
	let loadingMore = $state(false);
	let error = $state<string | null>(null);

	const provider = $derived(page.params.provider);
	const id = $derived(page.params.id);

	const discs = $derived.by(() => {
		const d: SongInfo[][] = [];
		for (const song of songs) {
			if (!d[song.volume_number]) d[song.volume_number] = [];
			d[song.volume_number][song.track_number] = song;
		}
		return d;
	});

	const hasMore = $derived(songs.length < total);

	async function load() {
		if (!provider || !id) return;
		album = null;
		songs = [];
		total = 0;
		error = null;
		try {
			const [albumData, songsData] = await Promise.all([
				getAlbumInfo(provider, id),
				getAlbumSongs(provider, id, { limit: 100 }),
			]);
			album = albumData;
			songs = songsData.items;
			total = songsData.total_number_of_items;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to load album";
		}
	}

	async function loadMore() {
		if (!provider || !id || loadingMore) return;
		loadingMore = true;
		try {
			const songsData = await getAlbumSongs(provider, id, {
				limit: 100,
				offset: songs.length,
			});
			songs = [...songs, ...songsData.items];
			total = songsData.total_number_of_items;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to load album";
		}
		loadingMore = false;
	}

	$effect(() => {
		if (provider && id) load();
	});

	function formatDuration(seconds: number) {
		return `${Math.floor(seconds / 60)}:${(seconds % 60).toString().padStart(2, "0")}`;
	}

	async function startDownload() {
		if (!provider || !id || !album) return;
		try {
			await downloadAlbum({ provider, id });
			error = null;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to download album";
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
		>{album?.title || "Album"} | {album?.artists
			.map((artist) => artist.name)
			.join(" ") || "Artist"} - MusicShack</title
	>
</svelte:head>

{#if error}
	<div class="flex flex-col items-center justify-center gap-2 py-4">
		<h1 class="text-xs font-semibold tracking-widest uppercase">
			Error loading Album
		</h1>
		<p class="text-sm text-muted-foreground">{error}</p>
		<a class="underline" href="/dashboard">Go to Home</a>
	</div>
{:else if !album}
	<p class="py-4 text-center text-sm text-muted-foreground">Loading...</p>
{:else}
	<div class="mt-1 mx-auto table border-separate border-spacing-y-2.5">
		<div class="table-row">
			<div class="flex flex-row flex-wrap justify-center gap-2.5">
				{#if album.cover_url}
					<img
						class="size-70 object-cover"
						src={album.cover_url}
						alt={album.title}
					/>
				{:else}
					<div class="size-70 bg-muted"></div>
				{/if}
				<div class="flex flex-col gap-1.5">
					<h1 class="flex flex-row items-center gap-2 font-extrabold">
					{album.title}
					{#if album.explicit}
						<Explicit />
					{/if}
					</h1>
					<div class="italic flex flex-wrap gap-x-2 gap-y-1">
						{#each album.artists as artist}
							<a href="/artist/{provider}/{artist.id}">
								{artist.name}
							</a>
						{/each}
					</div>
					<br />
					<p>{album.number_tracks} Tracks</p>
					{#if album.number_volumes > 1}
						<p>{album.number_volumes} Discs</p>
					{/if}
					<div class="flex items-center gap-1">
						<Clock4 class="size-4" />
						<p>{formatDuration(album.duration)}</p>
					</div>
					<p>{album.release_date.slice(0, 10)}</p>
				<Quality quality={album.audio_quality} />
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
			Download Album
		</Button>
	</div>

	<div class="flex flex-col gap-8 pt-8">
		{#each discs as disc, i}
			{#if disc}
				<div class="grid gap-2.5 pl-1.5">
					{#if album.number_volumes > 1}
						<h2 class="text-center font-bold">Disc {i}</h2>
					{/if}
					{#each disc as song}
						{#if song}
							<SongRow
								provider={provider!}
								song={song}
								onerror={(msg) => (error = msg)}
							/>
						{/if}
					{/each}
				</div>
			{/if}
		{/each}
		<LoadMore
			hasMore={hasMore}
			loading={loadingMore}
			onclick={loadMore}
		/>
	</div>
{/if}
