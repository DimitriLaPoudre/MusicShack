<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";

	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let song = $state<any | null>(null);

	afterNavigate(async () => {
		try {
			const res = await fetch(
				`http://localhost:8080/api/song/${page.params.api}/${page.params.id}`,
				{
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			if (!res.ok) {
				throw new Error("Failed to fetch album");
			}
			song = await res.json();
			song.Duration = `${Math.floor(song.Duration / 60)}:${(song.Duration % 60).toString().padStart(2, "0")}`;
			isLoading = false;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load album";
			isLoading = false;
		}
	});
</script>

<svelte:head>
	<title
		>{song?.Title || "Song"} | {song?.Artist?.Name || "Artist"} - MusicShack</title
	>
</svelte:head>

{#if isLoading}
	<p>Loading...</p>
{:else if error}
	<h2>Error Loading Song</h2>
	<p>{error}</p>
	<a href="/dashboard"> Go to Dashboard </a>
{:else}
	<!-- page top -->
	<div style="display: flex; flex-direction: row; gap: 10px;">
		<img
			src={song.Album.CoverUrl}
			alt={song.title}
			style="width:200px; height:auto;"
		/>
		<div style="display: flex; flex-direction: column; gap: 10px">
			<h1>{song.Title}</h1>
			<a
				href="/album/{page.params.api}/{song.Album.Id}"
				style="display: block;"
			>
				{song.Album.Title}
			</a>
			<div style="display: flex; gap: 10px;">
				{#each song.Artists as artist}
					<a href="/artist/{page.params.api}/{artist.Id}">
						{artist.Name}
					</a>
				{/each}
			</div>
			<br />
			<p>{song.Duration}</p>
			<p>{song.AudioQuality}</p>
			<button onclick={() => {}}>Download song</button>
		</div>
	</div>
{/if}

<style>
	h1 {
		margin: 0;
	}
	p {
		margin: 0;
	}
</style>
