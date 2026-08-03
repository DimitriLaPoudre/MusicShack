<script lang="ts">
	import { page } from "$app/state";
	import {
		ApiError,
		downloadPlaylist,
		getPlaylistInfo,
		getPlaylistSongs,
		listDownloads,
	} from "$lib/api";
	import { Button } from "$lib/components/ui/button/index.js";
	import LoadMore from "$lib/components/load-more.svelte";
	import SongRow from "$lib/components/song-row.svelte";
	import { Clock4, Download } from "@lucide/svelte";
	import { downloadList } from "$lib/stores";
	import type { PlaylistInfo, SongInfo } from "$lib/types";

	let playlist = $state<PlaylistInfo | null>(null);
	let songs = $state<SongInfo[]>([]);
	let total = $state(0);
	let loadingMore = $state(false);
	let error = $state<string | null>(null);

	const provider = $derived(page.params.provider);
	const id = $derived(page.params.id);

	const hasMore = $derived(songs.length < total);

	async function load() {
		if (!provider || !id) return;
		playlist = null;
		songs = [];
		total = 0;
		error = null;
		try {
			const [playlistData, songsData] = await Promise.all([
				getPlaylistInfo(provider, id),
				getPlaylistSongs(provider, id, { limit: 100 }),
			]);
			playlist = playlistData;
			songs = songsData.items;
			total = songsData.total_number_of_items;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to load playlist";
		}
	}

	async function loadMore() {
		if (!provider || !id || loadingMore) return;
		loadingMore = true;
		try {
			const songsData = await getPlaylistSongs(provider, id, {
				limit: 100,
				offset: songs.length,
			});
			songs = [...songs, ...songsData.items];
			total = songsData.total_number_of_items;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to load playlist";
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
		if (!provider || !id || !playlist) return;
		try {
			await downloadPlaylist({ provider, id });
			error = null;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to download playlist";
		}
		try {
			downloadList.set(await listDownloads());
		} catch {
			// keep cached list on refresh failure
		}
	}
</script>

<svelte:head>
	<title>{playlist?.title || "Playlist"} - MusicShack</title>
</svelte:head>

{#if error}
	<div class="flex flex-col items-center justify-center gap-2 py-4">
		<h1 class="text-xs font-semibold tracking-widest uppercase">
			Error loading Playlist
		</h1>
		<p class="text-sm text-muted-foreground">{error}</p>
		<a class="underline" href="/dashboard">Go to Home</a>
	</div>
{:else if !playlist}
	<p class="py-4 text-center text-sm text-muted-foreground">Loading...</p>
{:else}
	<div class="mt-1 mx-auto table border-separate border-spacing-y-2.5">
		<div class="table-row">
			<div class="flex flex-row flex-wrap justify-center gap-2.5">
				{#if playlist.cover_url}
					<img
						class="size-70 object-cover"
						src={playlist.cover_url}
						alt={playlist.title}
					/>
				{:else}
					<div class="size-70 bg-muted"></div>
				{/if}
				<div class="flex flex-col gap-1.5">
					<h1 class="flex flex-row items-center gap-2 font-extrabold">
						{playlist.title}
					</h1>
					<h2>{playlist.description}</h2>
					<br />
					<p>{playlist.number_of_tracks} Tracks</p>
					<div class="flex items-center gap-1">
						<Clock4 class="size-4" />
						<p>{formatDuration(playlist.duration)}</p>
					</div>
					<p>{playlist.last_updated.slice(0, 10)}</p>
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
			Download all song separately
		</Button>
	</div>

	<div class="flex flex-col gap-8 pt-8">
		<div class="grid gap-2.5 pl-1.5">
			{#each songs as song, index}
				{#if song}
					<SongRow
						provider={provider!}
						song={song}
						number={index + 1}
						onerror={(msg) => (error = msg)}
					/>
				{/if}
			{/each}
		</div>
		<LoadMore
			hasMore={hasMore}
			loading={loadingMore}
			onclick={loadMore}
		/>
	</div>
{/if}
