<script lang="ts">
	import { Download } from "lucide-svelte";
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";

	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let album = $state<any | null>(null);

	afterNavigate(async () => {
		try {
			const res = await fetch(
				`http://localhost:8080/api/album/${page.params.api}/${page.params.id}`,
				{
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			album = await res.json();
			if (!res.ok) {
				throw new Error(album.error || "Failed to fetch album");
			}
			album.Duration = `${Math.floor(album.Duration / 60)}:${(album.Duration % 60).toString().padStart(2, "0")}`;
			isLoading = false;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load album";
			isLoading = false;
		}
	});

	async function download(api: string, id: string) {
		try {
			const res = await fetch(
				`http://localhost:8080/api/users/downloads/${api}/${id}`,
				{
					method: "POST",
					credentials: "include",
				},
			);

			if (res.status === 401) {
				goto("/login");
				return;
			}
			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || "Failed to download song");
			}
		} catch (e) {
			error =
				e instanceof Error ? e.message : "Failed to load download song";
		}
	}
</script>

<svelte:head>
	<title
		>{album?.Title || "Album"} | {album?.Artist?.Name || "Artist"} - MusicShack</title
	>
</svelte:head>

{#if isLoading}
	<p>Loading...</p>
{:else if error}
	<h2>Error Loading Album</h2>
	<p>{error}</p>
	<a href="/">Go to Home</a>
{:else}
	<!-- page top -->
	<div
		style="display: flex; flex-direction: column; align-items: center; gap: 10px"
	>
		<div
			style="display: flex; flex-direction: row; gap: 10px; justify-content: center;"
		>
			<img
				src={album.CoverUrl}
				alt={album.Title}
				style="width:200px; height:auto;"
			/>
			<div style="display: flex; flex-direction: column; gap: 10px">
				<p>[{album.Type}]</p>
				<h1>{album.Title}</h1>
				<div style="display: flex; gap: 10px;">
					{#each album.Artists as artist}
						<a href="/artist/{page.params.api}/{artist.Id}">
							{artist.Name}
						</a>
					{/each}
				</div>
				<br />
				<p>{album.Duration}</p>
				<div style="display: flex; flex-direction: row; gap: 10px">
					<p>{album.ReleaseDate}</p>
					<p>{album.AudioQuality}</p>
				</div>
			</div>
		</div>
		<button
			style="display: flex; flex-direction: row; gap: 10px; padding: 10px 10px; "
			onclick={() => {}}
		>
			<Download size="24" />
			<p>Download Album</p>
		</button>
	</div>
	<!-- page body -->
	<div style="display: grid; gap: 10px; padding: 20px 20px;">
		{#each album.Songs as song}
			<div
				style="display: flex; flex-direction: row; justify-content: space-between; gap: 16px"
			>
				<p>
					{song.TrackNumber}
				</p>
				<a href="/song/{page.params.api}/{song.Id}">
					{song.Title}
				</a>
				<div style="display: flex; flex-direction: row; gap: 1em">
					{#each song.Artists as artist}
						<a href="/artist/{page.params.api}/{artist.Id}">
							{artist.Name}
						</a>
					{/each}
				</div>
				<p style="margin-left: auto;">
					{`${Math.floor(song.Duration / 60)}:${(song.Duration % 60).toString().padStart(2, "0")}`}
				</p>
				<button
					onclick={() => {
						download(page.params.api!, song.Id);
					}}
				>
					<Download size="28" />
				</button>
			</div>
		{/each}
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
