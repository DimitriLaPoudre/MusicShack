<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";
	import { Disc, DiscAlbum, Download, User } from "lucide-svelte";

	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let result = $state<any | null>(null);
	let searchData = $state<string | null>(null);
	let api = $state<string>("");
	let type = $state<string>("songs");

	afterNavigate(async () => {
		try {
			searchData = page.url.searchParams.get("q");
			if (!searchData) {
				throw new Error("No Search");
			}
			const res = await fetch(
				`http://localhost:8080/api/search?q=${searchData}`,
				{
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			result = await res.json();
			if (!res.ok) {
				throw new Error(result.error || "Failed to fetch search");
			}
			api = Object.keys(result)[0];
			isLoading = false;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load song";
			isLoading = false;
		}
	});
</script>

<svelte:head>
	<title>{"Search"} | {searchData} - MusicShack</title>
</svelte:head>

{#if isLoading}
	<p>Loading...</p>
{:else if error}
	<h2>Error Loading Song</h2>
	<p>{error}</p>
	<a href="/">Go to Home</a>
{:else}
	<div class="api">
		{#each Object.entries(result as Record<string, any>) as [key, _]}
			<button class="api-btn">{key}</button>
		{/each}
	</div>
	<div class="type">
		<button class="type-btn" onclick={() => (type = "songs")}>Songs</button>
		<button class="type-btn" onclick={() => (type = "albums")}
			>Albums</button
		>
		<button class="type-btn" onclick={() => (type = "artists")}
			>Artists</button
		>
	</div>
	<div class="items">
		{#if type === "songs"}
			{#each result[api].Songs as song}
				<a class="song" href="/song/{api}/{song.Id}">
					{#if song.CoverUrl !== ""}
						<img src={song.CoverUrl} alt={song.CoverUrl} />
					{:else}
						<Disc style="width: 160px; height: 160px;" />
					{/if}
					<div class="song-detail">
						<div>
							<p>{song.Title}</p>
							{#each song.Artists as artist}
								<a href="/artist/{api}/{artist.Id}"
									>{artist.Name}</a
								>
							{/each}
						</div>
						<button><Download /></button>
					</div>
				</a>
			{/each}
		{:else if type === "albums"}
			{#each result[api].Albums as album}
				<a class="album" href="/album/{api}/{album.Id}">
					{#if album.CoverUrl !== ""}
						<img src={album.CoverUrl} alt={album.Title} />
					{:else}
						<DiscAlbum style="width: 160px; height: 160px;" />
					{/if}

					<div class="album-detail">
						<div>
							<p>{album.Title}</p>
							{#each album.Artists as artist}
								<a href="/artist/{api}/{artist.Id}"
									>{artist.Name}</a
								>
							{/each}
						</div>
						<button><Download /></button>
					</div>
				</a>
			{/each}
		{:else}
			{#each result[api].Artists as artist}
				<a class="artist" href="/artist/{api}/{artist.Id}">
					{#if artist.PictureUrl !== ""}
						<img src={artist.PictureUrl} alt={artist.PictureUrl} />
					{:else}
						<User style="width: 160px; height: 160px;" />
					{/if}
					<p>{artist.Name}</p>
				</a>
			{/each}
		{/if}
	</div>
{/if}

<style>
	.api {
		display: flex;
		flex-direction: row;
		background-color: var(--color-primary);
		gap: 10px;
		padding: 10px;
	}
	.api-btn {
		padding: 10px;
		background-color: var(--color-secondary);
	}

	.type {
		display: flex;
		flex-direction: row;
		background-color: var(--color-secondary);
		gap: 10px;
		padding: 10px;
	}
	.type-btn {
		padding: 10px;
		background-color: var(--color-primary);
	}

	.items {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;
	}

	.song {
		width: 200px;
		height: auto;
	}
	.song-detail {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
	}

	.album {
		width: 200px;
		height: auto;
	}
	.album-detail {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
	}

	.artist {
		width: 200px;
		height: auto;
	}
</style>
