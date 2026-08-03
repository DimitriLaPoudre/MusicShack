<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { ApiError, search } from "$lib/api";
	import AlbumCard from "$lib/components/album-card.svelte";
	import ArtistCard from "$lib/components/artist-card.svelte";
	import PlaylistCard from "$lib/components/playlist-card.svelte";
	import SongCard from "$lib/components/song-card.svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import type { SearchResult } from "$lib/types";

	type SearchType = "songs" | "albums" | "artists" | "playlists";

	let error = $state<string | null>(null);
	let provider = $state<string>("");
	let type = $state<SearchType>("songs");
	let result = $state<SearchResult | null>(null);

	const searchData = $derived(page.url.searchParams.get("q"));
	const providerResult = $derived(result?.provider_result ?? null);

	async function fetchData(q: string | null) {
		result = null;
		error = null;
		try {
			if (!q) throw new Error("No Search");
			const data = await search(q);
			if (data.item?.type && data.item?.data) {
				const item = data.item.data as { provider: string; id: string };
				await goto(`/${data.item.type}/${item.provider}/${item.id}`);
			} else if (data.provider_result) {
				if (Object.keys(data.provider_result).length === 0) {
					throw new Error("instances missing");
				}
				provider = Object.keys(data.provider_result)[0];
				result = data;
			} else {
				throw new Error("No result");
			}
		} catch (e) {
			error =
				e instanceof ApiError || e instanceof Error
					? e.message
					: "Failed to search";
		}
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
{:else if !providerResult}
	<p class="mt-6 text-center">Searching...</p>
{:else}
	<div class="flex flex-col gap-2 py-4 items-center">
		<div class="flex flex-row gap-2">
			{#each Object.keys(providerResult) as key}
				<Button
					variant="hover-full"
					class={provider === key
						? "bg-foreground text-background shadow-none"
						: ""}
					onclick={() => (provider = key)}
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
					onclick={() => (type = tab.key as SearchType)}
				>
					{tab.label}
				</Button>
			{/each}
		</div>
	</div>
	<div class="grid grid-cols-[repeat(auto-fit,200px)] justify-center gap-4">
		{#if type === "songs"}
			{#if providerResult[provider].songs.items.length === 0}
				<p class="flex justify-center">No song found</p>
			{/if}
			{#each providerResult[provider].songs.items as song}
				<SongCard {provider} {song} onerror={(msg) => (error = msg)} />
			{/each}
		{:else if type === "albums"}
			{#if providerResult[provider].albums.items.length === 0}
				<p class="flex justify-center">No album found</p>
			{/if}
			{#each providerResult[provider].albums.items as album}
				<AlbumCard
					{provider}
					{album}
					onerror={(msg) => (error = msg)}
				/>
			{/each}
		{:else if type === "artists"}
			{#if providerResult[provider].artists.items.length === 0}
				<p class="flex justify-center">No artist found</p>
			{/if}
			{#each providerResult[provider].artists.items as artist}
				<ArtistCard
					{provider}
					{artist}
					onerror={(msg) => (error = msg)}
				/>
			{/each}
		{:else if type === "playlists"}
			{#if providerResult[provider].playlists.items.length === 0}
				<p class="flex justify-center">No playlist found</p>
			{/if}
			{#each providerResult[provider].playlists.items as playlist}
				<PlaylistCard
					{provider}
					{playlist}
					onerror={(msg) => (error = msg)}
				/>
			{/each}
		{/if}
	</div>
{/if}
