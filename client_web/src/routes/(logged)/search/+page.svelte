<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import { Disc, DiscAlbum, Download, HeartIcon, User } from "lucide-svelte";
	import type { SearchResponse, SearchResult } from "$lib/types/response";
	import Quality from "$lib/components/Quality.svelte";
	import Explicit from "$lib/components/Explicit.svelte";
	import { onMount } from "svelte";
	import Owned from "$lib/components/Owned.svelte";
	import { addFollow, removeFollow } from "$lib/functions/follow";

	let error = $state<null | string>(null);
	let provider = $state<string>("");
	let type = $state<string>("songs");
	let result = $state<SearchResult | null>(null);

	const searchData = $derived(page.url.searchParams.get("q"));

	async function fetchData(searchData: string | null) {
		result = null;
		try {
			if (!searchData) {
				throw new Error("No Search");
			}
			const data = await apiFetch<SearchResponse>(
				`/search?q=${searchData}`,
			);

			if ("url" in data) {
				goto(`/${data.url.type}/${data.url.provider}/${data.url.id}`);
			} else {
				if (Object.keys(data.result).length === 0) {
					throw new Error("instances missing");
				}
				result = data.result;
				provider = Object.keys(data.result)[0];
				error = null;
			}
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load song";
		}
	}

	$effect(() => {
		if (searchData) {
			fetchData(searchData);
		}
	});

	onMount(async () => {});
</script>

<svelte:head>
	<title>Search | {searchData} - MusicShack</title>
</svelte:head>

<h1 class="mt-4 text-center">"{searchData}"</h1>
{#if error}
	<div class="mt-4 flex flex-col justify-center items-center gap-2.5">
		<h2>Error loading Search result</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !result}
	<p class="mt-6 text-center">Searching...</p>
{:else}
	<div class="flex flex-col gap-2 py-4 items-center">
		<div class="flex flex-row gap-2">
			{#each Object.entries(result as Record<string, any>) as [key, _]}
				<button
					class="hover-full p-4"
					onclick={() => (provider = key)}
					class:active={provider === key}
				>
					{key}</button
				>
			{/each}
		</div>
		<div class="flex flex-row gap-2">
			<button
				class="hover-full p-4"
				onclick={() => (type = "songs")}
				class:active={type === "songs"}
			>
				Songs</button
			>
			<button
				class="hover-full p-4"
				onclick={() => (type = "albums")}
				class:active={type === "albums"}>Albums</button
			>
			<button
				class="hover-full p-4"
				onclick={() => (type = "artists")}
				class:active={type === "artists"}>Artists</button
			>
			<button
				class="hover-full p-4"
				onclick={() => (type = "playlists")}
				class:active={type === "playlists"}>Playlists</button
			>
		</div>
	</div>
	<div class="grid grid-cols-[repeat(auto-fit,200px)] justify-center gap-4">
		{#if type === "songs"}
			{#if result[provider].songs.length === 0}
				<p class="flex justify-center">No song found</p>
			{/if}
			{#each result[provider].songs as song}
				<div class="w-50 h-auto">
					<button
						class="hover-full flex flex-col items-center w-50 h-auto overflow-hidden gap-3 shadow-[inset_0_1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
						onclick={(e) => {
							if (
								e.target instanceof Element &&
								e.target.closest("a")
							)
								return;
							goto(`/song/${provider}/${song.id}`);
						}}
					>
						<div class="w-40 h-40">
							{#if song.album.coverUrl !== ""}
								<img
									src={song.album.coverUrl}
									alt={song.title}
								/>
							{:else}
								<Disc size={140} />
							{/if}
						</div>
						<p
							class="flex flex-row items-center justify-center gap-2 font-extrabold"
						>
							{#if song.downloaded}
								<Owned />
							{/if}
							{song.title}
							{#if song.explicit}
								<Explicit />
							{/if}
						</p>
						<nav
							class="flex flex-col gap-y-[0.2rem] gap-x-4 italic"
						>
							{#each song.artists as artist}
								<a href="/artist/{provider}/{artist.id}">
									{artist.name}
								</a>
							{/each}
						</nav>
						<Quality quality={song.audioQuality} />
					</button>
					<button
						class="hover-full w-full p-3 shadow-[inset_0_-1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
						onclick={async () => {
							error = await download({
								provider: provider,
								type: "song",
								id: song!.id,
							});
						}}
					>
						<Download />
					</button>
				</div>
			{/each}
		{:else if type === "albums"}
			{#if result[provider].albums.length === 0}
				<p class="flex justify-center">No album found</p>
			{/if}
			{#each result[provider].albums as album}
				<div class="w-50 h-auto">
					<button
						class="hover-full flex flex-col items-center w-50 h-auto overflow-hidden gap-3 shadow-[inset_0_1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
						onclick={(e) => {
							if (
								e.target instanceof Element &&
								e.target.closest("a")
							)
								return;
							goto(`/album/${provider}/${album.id}`);
						}}
					>
						<div class="w-40 h-40">
							{#if album.coverUrl !== ""}
								<img src={album.coverUrl} alt={album.title} />
							{:else}
								<DiscAlbum size={140} />
							{/if}
						</div>
						<p
							class="flex flex-row items-center justify-center gap-2 font-extrabold"
						>
							{#if album.downloaded}
								<Owned />
							{/if}
							{album.title}
							{#if album.explicit}
								<Explicit />
							{/if}
						</p>
						<nav
							class="flex flex-col gap-y-[0.2rem] gap-x-4 italic"
						>
							{#each album.artists as artist}
								<a href="/artist/{provider}/{artist.id}">
									{artist.name}
								</a>
							{/each}
						</nav>
						<Quality quality={album.audioQuality} />
					</button>
					<button
						class="hover-full w-full p-3 shadow-[inset_0_-1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
						onclick={async () => {
							error = await download({
								provider: provider,
								type: "album",
								id: album!.id,
							});
						}}
					>
						<Download />
					</button>
				</div>
			{/each}
		{:else if type === "artists"}
			{#if result[provider].artists.length === 0}
				<p class="flex justify-center">No artist found</p>
			{/if}
			{#each result[provider].artists as artist}
				<div class="w-50 h-auto">
					<button
						class="hover-full flex flex-col items-center w-50 h-auto overflow-hidden gap-3 shadow-[inset_0_1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
						onclick={() => goto(`/artist/${provider}/${artist.id}`)}
					>
						<div class="w-40 h-40 flex items-center justify-center">
							{#if artist.pictureUrl !== ""}
								<img
									class="rounded-full"
									src={artist.pictureUrl}
									alt={artist.name}
								/>
							{:else}
								<User size={140} />
							{/if}
						</div>
						<p>{artist.name}</p>
					</button>
					<button
						class="hover-full w-full p-3 shadow-[inset_0_-1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
						onclick={async () => {
							if (
								artist!.followed == -1 ||
								artist!.followed == -2
							) {
								return;
							}
							if (artist!.followed) {
								const follow = artist!.followed;
								artist!.followed = -1;

								const error = await removeFollow(follow);
								if (error) {
									artist!.followed = follow;
								} else {
									artist!.followed = 0;
								}
							} else {
								artist!.followed = -2;
								const { follow, error } = await addFollow({
									provider: provider,
									id: artist!.id,
								});
								if (error || !follow) {
									artist!.followed = 0;
								} else {
									artist!.followed = follow.id;
								}
							}
						}}
					>
						{#if artist!.followed != -1 && artist!.followed != 0}
							<HeartIcon color="#0F0" fill="#0F0" />
						{:else}
							<HeartIcon />
						{/if}
					</button>
				</div>
			{/each}
		{:else if type === "playlists"}
			{#if result[provider].playlists.length === 0}
				<p class="flex justify-center">No playlist found</p>
			{/if}
			{#each result[provider].playlists as playlist}
				<div class="w-50 h-auto">
					<button
						class="hover-full flex flex-col items-center w-50 h-auto overflow-hidden gap-3 shadow-[inset_0_1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
						onclick={(e) => {
							if (
								e.target instanceof Element &&
								e.target.closest("a")
							)
								return;
							goto(`/playlist/${provider}/${playlist.id}`);
						}}
					>
						<div class="w-40 h-40">
							{#if playlist.coverUrl !== ""}
								<img
									src={playlist.coverUrl}
									alt={playlist.title}
								/>
							{:else}
								<DiscAlbum size={140} />
							{/if}
						</div>
						<p
							class="flex flex-row items-center justify-center gap-2 font-extrabold"
						>
							{#if playlist.downloaded}
								<Owned />
							{/if}
							{playlist.title}
						</p>
					</button>
					<button
						class="hover-full w-full p-3 shadow-[inset_0_-1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
						onclick={async () => {
							error = await download({
								provider: provider,
								type: "playlist",
								id: playlist!.id,
							});
						}}
					>
						<Download />
					</button>
				</div>
			{/each}
		{/if}
	</div>
{/if}
