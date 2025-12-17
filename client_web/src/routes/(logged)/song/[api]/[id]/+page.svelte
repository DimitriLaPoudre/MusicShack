<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";
	import { apiFetch } from "$lib/functions/apiFetch";
	import { downloadSong } from "$lib/functions/download";

	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let song = $state<any | null>(null);

	afterNavigate(async () => {
		try {
			const res = await apiFetch(
				`/api/song/${page.params.api}/${page.params.id}`,
			);
			song = await res.json();
			if (!res.ok) {
				throw new Error(song.error || "Failed to fetch song");
			}
			song.Duration = `${Math.floor(song.Duration / 60)}:${(song.Duration % 60).toString().padStart(2, "0")}`;
			isLoading = false;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load song";
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
				<img class="cover" src={song.Album.CoverUrl} alt={song.Title} />
				<div class="data">
					<h1>{song.Title}</h1>
					<a href="/album/{page.params.api}/{song.Album.Id}">
						{song.Album.Title}
					</a>
					<div class="artists">
						{#each song.Artists as artist}
							<a href="/artist/{page.params.api}/{artist.Id}">
								{artist.Name}
							</a>
						{/each}
					</div>
					<br />
					<p>{song.Duration}</p>
					<p>{song.AudioQuality}</p>
				</div>
			</div>
		</div>
		<button
			class="download"
			onclick={async () => {
				error = await downloadSong(page.params.api!, song.Id);
			}}
		>
			Download Song
		</button>
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
</style>
