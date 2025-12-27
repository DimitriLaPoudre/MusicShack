<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import { Disc, DiscAlbum, Download, User } from "lucide-svelte";
	import type { ErrorResponse, SearchResponse } from "$lib/types/response";
	import { quality } from "$lib/types/quality";

	let error = $state<null | string>(null);
	let searchData = $state<null | string>(null);
	let api = $state<string>("");
	let type = $state<string>("songs");
	let result = $state<SearchResponse | null>(null);

	afterNavigate(async () => {
		try {
			searchData = page.url.searchParams.get("q");
			if (!searchData) {
				throw new Error("No Search");
			}
			const data = await apiFetch<SearchResponse>(
				`/search?q=${searchData}`,
			);
			if ("error" in data) {
				throw new Error(
					(data as ErrorResponse).error || "Failed to fetch search",
				);
			}
			if (Object.keys(data).length === 0) {
				throw new Error("instances missing");
			}
			result = data;
			api = Object.keys(data)[0];
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load song";
		}
	});
</script>

<svelte:head>
	<title>{"Search"} | {searchData} - MusicShack</title>
</svelte:head>

<h1 class="research">"{searchData}"</h1>
{#if error}
	<div class="error">
		<h2>Error loading Search result</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !result}
	<p class="loading">Searching...</p>
{:else}
	<div class="top">
		<div class="section">
			{#each Object.entries(result as Record<string, any>) as [key, _]}
				<button onclick={() => (api = key)} class:active={api === key}>
					{key}</button
				>
			{/each}
		</div>
		<div class="section">
			<button
				onclick={() => (type = "songs")}
				class:active={type === "songs"}
			>
				Songs</button
			>
			<button
				onclick={() => (type = "albums")}
				class:active={type === "albums"}>Albums</button
			>
			<button
				onclick={() => (type = "artists")}
				class:active={type === "artists"}>Artists</button
			>
		</div>
	</div>
	<div class="items">
		{#if type === "songs"}
			{#each result[api].songs as song}
				<div class="wrap-item">
					<button
						class="item"
						onclick={(e) => {
							if (
								e.target instanceof Element &&
								e.target.closest("a")
							)
								return;
							goto(`/song/${api}/${song.id}`);
						}}
					>
						<div class="cover">
							{#if song.album.coverUrl !== ""}
								<img
									src={song.album.coverUrl}
									alt={song.title}
								/>
							{:else}
								<Disc size={140} />
							{/if}
						</div>
						<p>{song.title}</p>
						<nav>
							{#each song.artists as artist}
								<a href="/artist/{api}/{artist.id}">
									{artist.name}
								</a>
							{/each}
						</nav>
						<p>{quality[song.audioQuality]}</p>
					</button>
					<button
						class="download"
						onclick={async () => {
							error = await download({
								api: api,
								type: "song",
								id: song!.id,
								quality: "",
							});
						}}
					>
						<Download />
					</button>
				</div>
			{/each}
		{:else if type === "albums"}
			{#each result[api].albums as album}
				<div class="wrap-item">
					<button
						class="item"
						onclick={(e) => {
							if (
								e.target instanceof Element &&
								e.target.closest("a")
							)
								return;
							goto(`/album/${api}/${album.id}`);
						}}
					>
						<div class="cover">
							{#if album.coverUrl !== ""}
								<img src={album.coverUrl} alt={album.title} />
							{:else}
								<DiscAlbum size={140} />
							{/if}
						</div>
						<p>{album.title}</p>
						<nav>
							{#each album.artists as artist}
								<a href="/artist/{api}/{artist.id}">
									{artist.name}
								</a>
							{/each}
						</nav>
						<p>{quality[album.audioQuality]}</p>
					</button>
					<button
						class="download"
						onclick={async () => {
							error = await download({
								api: api,
								type: "album",
								id: album!.id,
								quality: "",
							});
						}}
					>
						<Download />
					</button>
				</div>
			{/each}
		{:else}
			{#each result[api].artists as artist}
				<button
					class="artist"
					onclick={() => goto(`/artist/${api}/${artist.id}`)}
				>
					<div class="picture">
						{#if artist.pictureUrl !== ""}
							<img src={artist.pictureUrl} alt={artist.name} />
						{:else}
							<User size={140} />
						{/if}
					</div>
					<p>{artist.name}</p>
				</button>
			{/each}
		{/if}
	</div>
{/if}

<style>
	.research {
		margin-top: 15px;
		text-align: center;
	}

	.loading {
		margin-top: 15px;
		text-align: center;
	}
	.error {
		margin-top: 15px;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 10px;
	}

	.top {
		display: flex;
		flex-direction: column;
		gap: 10px;
		padding: 10px 0;
		.section {
			display: flex;
			flex-direction: row;
			gap: 10px;
			button {
				padding: 10px;
			}
		}
	}

	.items {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;

		.wrap-item {
			width: 200px;
			height: auto;
			.item {
				display: flex;
				flex-direction: column;
				align-items: center;
				width: 200px;
				height: auto;
				overflow: hidden;
				border-bottom: none;
				gap: 8px;

				.cover {
					width: 160px;
					height: 160px;
				}

				nav {
					display: flex;
					flex-direction: column;
					gap: 0.2rem 1rem;
				}
			}
			.download {
				width: 100%;
				border-top: none;
			}
		}
	}

	.artist {
		width: 200px;
		height: auto;
		.picture {
			width: 160px;
			height: 160px;
			border-radius: 50%;
		}
	}
</style>
