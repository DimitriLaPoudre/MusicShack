<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import {
		ApiError,
		searchAlbum,
		searchArtist,
		searchPlaylist,
		searchSetup,
		searchSong,
	} from "$lib/api";
	import AlbumCard from "$lib/components/album-card.svelte";
	import ArtistCard from "$lib/components/artist-card.svelte";
	import LoadMore from "$lib/components/load-more.svelte";
	import PlaylistCard from "$lib/components/playlist-card.svelte";
	import SongCard from "$lib/components/song-card.svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import type { AlbumInfo, ArtistInfo, PlaylistInfo, SongInfo } from "$lib/types";

	type SearchType = "songs" | "albums" | "artists" | "playlists";

	const PAGE_SIZE = 25;

	let error = $state<string | null>(null);
	let providers = $state<string[]>([]);
	let provider = $state<string>("");
	let type = $state<SearchType>("songs");
	let songs = $state<SongInfo[]>([]);
	let albums = $state<AlbumInfo[]>([]);
	let artists = $state<ArtistInfo[]>([]);
	let playlists = $state<PlaylistInfo[]>([]);
	let total = $state(0);
	let loadingMore = $state(false);

	const searchData = $derived(page.url.searchParams.get("q"));

	const currentItems = $derived.by(() => {
		switch (type) {
			case "songs":
				return songs;
			case "albums":
				return albums;
			case "artists":
				return artists;
			case "playlists":
				return playlists;
		}
	});

	const hasMore = $derived(currentItems.length < total);

	async function fetchPage(
		q: string,
		prov: string,
		t: SearchType,
		offset: number,
	) {
		switch (t) {
			case "songs":
				return searchSong(q, prov, { limit: PAGE_SIZE, offset });
			case "albums":
				return searchAlbum(q, prov, { limit: PAGE_SIZE, offset });
			case "artists":
				return searchArtist(q, prov, { limit: PAGE_SIZE, offset });
			case "playlists":
				return searchPlaylist(q, prov, { limit: PAGE_SIZE, offset });
		}
	}

	async function loadType(q: string, prov: string, t: SearchType) {
		switch (t) {
			case "songs":
				songs = [];
				break;
			case "albums":
				albums = [];
				break;
			case "artists":
				artists = [];
				break;
			case "playlists":
				playlists = [];
				break;
		}
		total = 0;
		const data = await fetchPage(q, prov, t, 0);
		total = data.total_number_of_items;
		switch (t) {
			case "songs":
				songs = data.items as SongInfo[];
				break;
			case "albums":
				albums = data.items as AlbumInfo[];
				break;
			case "artists":
				artists = data.items as ArtistInfo[];
				break;
			case "playlists":
				playlists = data.items as PlaylistInfo[];
				break;
		}
	}

	async function fetchData(q: string) {
		error = null;
		songs = [];
		albums = [];
		artists = [];
		playlists = [];
		total = 0;
		providers = [];
		provider = "";
		try {
			if (!q) throw new Error("No Search");
			const setup = await searchSetup(q);
			if (setup.item?.type && setup.item?.data) {
				const item = setup.item.data as { provider: string; id: string };
				await goto(`/${setup.item.type}/${item.provider}/${item.id}`);
				return;
			}
			const keys = Object.keys(setup.provider_result ?? {});
			if (keys.length === 0) {
				throw new Error("instances missing");
			}
			providers = keys;
			provider = keys[0];
			await loadType(q, provider, type);
		} catch (e) {
			error =
				e instanceof ApiError || e instanceof Error
					? e.message
					: "Failed to search";
		}
	}

	async function loadMore() {
		if (!searchData || !provider || loadingMore || !hasMore) return;
		loadingMore = true;
		try {
			const data = await fetchPage(searchData, provider, type, currentItems.length);
			total = data.total_number_of_items;
			switch (type) {
				case "songs":
					songs = [...songs, ...(data.items as SongInfo[])];
					break;
				case "albums":
					albums = [...albums, ...(data.items as AlbumInfo[])];
					break;
				case "artists":
					artists = [...artists, ...(data.items as ArtistInfo[])];
					break;
				case "playlists":
					playlists = [...playlists, ...(data.items as PlaylistInfo[])];
					break;
			}
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to load more";
		}
		loadingMore = false;
	}

	async function switchProvider(p: string) {
		provider = p;
		if (searchData) await loadType(searchData, p, type);
	}

	async function switchType(t: SearchType) {
		type = t;
		if (searchData && provider) await loadType(searchData, provider, t);
	}

	$effect(() => {
		if (searchData) fetchData(searchData);
	});
</script>

<svelte:head>
	<title>Search | {searchData} - MusicShack</title>
</svelte:head>

<h1 class="mt-4 text-center">"{searchData}"</h1>
{#if error}
	<div class="mt-4 flex flex-col justify-center items-center gap-2.5">
		<h2>Error loading Search result</h2>
		<p>{error}</p>
		<a class="underline" href="/dashboard">Go to Home</a>
	</div>
{:else if !provider}
	<p class="mt-6 text-center">Searching...</p>
{:else}
	<div class="flex flex-col gap-2 py-4 items-center">
		<div class="flex flex-row gap-2">
			{#each providers as key}
				<Button
					variant="hover-full"
					class={provider === key
						? "bg-foreground text-background shadow-none"
						: ""}
					onclick={() => switchProvider(key)}
				>
					{key}
				</Button>
			{/each}
		</div>
		<div class="flex flex-row gap-2">
			{#each [{ key: "songs", label: "Songs" }, { key: "albums", label: "Albums" }, { key: "artists", label: "Artists" }, { key: "playlists", label: "Playlists" }] as tab}
				<Button
					variant="hover-full"
					class={type === tab.key
						? "bg-foreground text-background shadow-none"
						: ""}
					onclick={() => switchType(tab.key as SearchType)}
				>
					{tab.label}
				</Button>
			{/each}
		</div>
	</div>
	<div class="grid grid-cols-[repeat(auto-fit,200px)] justify-center gap-4">
		{#if type === "songs"}
			{#if currentItems.length === 0}
				<p class="flex justify-center">No song found</p>
			{/if}
			{#each songs as song}
				<SongCard {provider} {song} onerror={(msg) => (error = msg)} />
			{/each}
		{:else if type === "albums"}
			{#if currentItems.length === 0}
				<p class="flex justify-center">No album found</p>
			{/if}
			{#each albums as album}
				<AlbumCard {provider} {album} onerror={(msg) => (error = msg)} />
			{/each}
		{:else if type === "artists"}
			{#if currentItems.length === 0}
				<p class="flex justify-center">No artist found</p>
			{/if}
			{#each artists as artist}
				<ArtistCard {provider} {artist} onerror={(msg) => (error = msg)} />
			{/each}
		{:else if type === "playlists"}
			{#if currentItems.length === 0}
				<p class="flex justify-center">No playlist found</p>
			{/if}
			{#each playlists as playlist}
				<PlaylistCard
					{provider}
					{playlist}
					onerror={(msg) => (error = msg)}
				/>
			{/each}
		{/if}
	</div>
	<div class="pt-4">
		<LoadMore
			hasMore={hasMore}
			loading={loadingMore}
			onclick={loadMore}
		/>
	</div>
{/if}
