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

	async function downloadSong(api: string, id: string) {
		try {
			const res = await fetch(
				`http://localhost:8080/api/users/downloads/song/${api}/${id}`,
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

	async function downloadAlbum(api: string, id: string) {
		try {
			const res = await fetch(
				`http://localhost:8080/api/users/downloads/album/${api}/${id}`,
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
				throw new Error(data.error || "Failed to download album");
			}
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Failed to load download album";
		}
	}
</script>

<svelte:head>
	<title
		>{album?.Title || "Album"} | {album?.Artist?.Name || "Artist"} - MusicShack</title
	>
</svelte:head>

{#if isLoading}
	<p class="loading">Loading...</p>
{:else if error}
	<div class="error">
		<h2>Error Loading Song</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else}
	<!-- page top -->
	<div class="header">
		<div class="top">
			<div class="top-data">
				<img class="cover" src={album.CoverUrl} alt={album.Title} />
				<div class="data">
					<h1>{album.Title}</h1>
					<div class="artists">
						{#each album.Artists as artist}
							<a href="/artist/{page.params.api}/{artist.Id}">
								{artist.Name}
							</a>
						{/each}
					</div>
					<br />
					<p>{album.Duration}</p>
					<p>{album.AudioQuality}</p>
				</div>
			</div>
		</div>
		<button
			class="download"
			onclick={() => {
				downloadAlbum(page.params.api!, album.Id);
			}}>Download Album</button
		>
	</div>

	<!-- page body -->
	<div class="body">
		{#each album.Songs as song}
			<div class="song">
				<p class="number">{song.TrackNumber}</p>
				<a class="title" href="/song/{page.params.api}/{song.Id}">
					{song.Title}
				</a>
				<p class="split">|</p>
				<div class="wrap-artists">
					<div class="artists">
						{#each song.Artists as artist}
							<a href="/artist/{page.params.api}/{artist.Id}">
								{artist.Name}
							</a>
						{/each}
					</div>
				</div>
				<p class="duration">
					{`${Math.floor(song.Duration / 60)}:${(song.Duration % 60).toString().padStart(2, "0")}`}
				</p>
				<button
					onclick={() => {
						downloadSong(page.params.api!, song.Id);
					}}
				>
					<Download />
				</button>
			</div>
		{/each}
	</div>
{/if}

<style>
	.loading {
		margin-top: 30px;
		text-align: center;
	}
	.error {
		margin-top: 30px;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 10px;

		* {
			margin: 0;
		}
	}

	.header {
		display: table;
		margin: 0 auto;
		border-spacing: 0 10px;
	}
	.top {
		display: table-row;
	}

	.top-data {
		display: flex;
		flex-direction: row;
		flex-wrap: wrap;
		justify-content: center;
		gap: 10px;
	}

	.cover {
		width: 160px;
		height: 160px;
	}

	.data {
		display: flex;
		flex-direction: column;
		gap: 7px;

		* {
			margin: 0;
		}

		.artists {
			display: flex;
			flex-wrap: wrap;
			gap: 0px 0.5rem;
		}
	}

	.download {
		display: table-row;
		width: 100%;
	}

	.body {
		display: grid;
		gap: 10px;
		padding: 0 0 0 5px;
	}

	.song {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		gap: 16px;
		overflow-x: hidden;

		.number {
			margin: 0;
		}

		.title {
			margin: 0 0 0 auto;
		}

		.split {
			margin: 0;
		}

		.wrap-artists {
			overflow-x: hidden;
			.artists {
				margin: 0;
				display: flex;
				flex-direction: row;
				flex-wrap: nowrap;
				flex-shrink: 0;
				overflow-x: hidden;
				gap: 1em;
			}
		}

		.duration {
			margin: 0;
			margin-left: auto;
			overflow-x: hidden;
		}

		button {
			aspect-ratio: 1/1;
		}
	}
</style>
