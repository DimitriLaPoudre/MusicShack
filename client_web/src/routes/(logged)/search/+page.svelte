<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";

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
				<div
					class="song"
					onclick={() => goto(`/song/${api}/${song.Id}`)}
				>
					<img src={song.CoverUrl} alt={song.Title} />
					<p>{song.Title}</p>
					{#each song.Artists as artist}
						<p>{artist.Name}</p>
					{/each}
				</div>
			{/each}
		{:else if type === "albums"}
			{#each result[api].Albums as album}
				<div
					class="album"
					onclick={() => goto(`/album/${api}/${album.Id}`)}
				>
					<img src={album.CoverUrl} alt={album.Title} />
					<p>{album.Title}</p>
					{#each album.Artists as artist}
						<p>{artist.Name}</p>
					{/each}
				</div>
			{/each}
		{:else}
			{#each result[api].Artists as artist}
				<div
					class="artist"
					onclick={() => goto(`/artist/${api}/${artist.Id}`)}
				>
					<img src={artist.PictureUrl} alt={artist.Name} />
					<p>{artist.Name}</p>
				</div>
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
	.album {
		width: 200px;
		height: auto;
	}
	.artist {
		width: 200px;
		height: auto;
	}
</style>
