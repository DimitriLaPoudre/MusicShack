<script lang="ts">
	import { Download } from "lucide-svelte";
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";
	import { apiFetch } from "$lib/functions/fetch";
	import { downloadAlbum, downloadSong } from "$lib/functions/download";

	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let album = $state<any | null>(null);

	afterNavigate(async () => {
		try {
			const res = await apiFetch(
				`/album/${page.params.api}/${page.params.id}`,
			);
			album = await res.json();
			if (!res.ok) {
				throw new Error(album.error || "Failed to fetch album");
			}
			album.Duration = `${Math.floor(album.Duration / 60)}:${(album.Duration % 60).toString().padStart(2, "0")}`;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load album";
		}
		isLoading = false;
	});
</script>

<svelte:head>
	<title
		>{album?.Title || "Album"} | {album?.Artist?.Name || "Artist"} - MusicShack</title
	>
</svelte:head>

{#if error}
	<div class="error">
		<h2>Error loading Album</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if isLoading}
	<p class="loading">Loading...</p>
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
			onclick={async () => {
				error = await downloadAlbum(page.params.api!, album.Id);
			}}
		>
			Download Album
		</button>
	</div>

	<!-- page body -->
	<div class="body">
		{#each album.Songs as song}
			<div class="song">
				<button
					class="data"
					onclick={(e) => {
						if (
							e.target instanceof Element &&
							e.target.closest("a")
						)
							return;
						goto(`/song/${page.params.api}/${song.Id}`);
					}}
				>
					<p class="number">{song.TrackNumber}</p>
					<p>{song.Title}</p>
					<div class="artists">
						{#each song.Artists as artist}
							<a href="/artist/{page.params.api}/{artist.Id}">
								{artist.Name}
							</a>
						{/each}
					</div>
					<p class="duration">
						{`${Math.floor(song.Duration / 60)}:${(song.Duration % 60).toString().padStart(2, "0")}`}
					</p>
				</button>
				<button
					class="download"
					onclick={async () => {
						error = await downloadSong(page.params.api!, song.Id);
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
		.top {
			display: table-row;

			.top-data {
				display: flex;
				flex-direction: row;
				flex-wrap: wrap;
				justify-content: center;
				gap: 10px;

				.cover {
					width: 160px;
					height: 160px;
				}

				.data {
					display: flex;
					flex-direction: column;
					gap: 7px;

					.artists {
						display: flex;
						flex-wrap: wrap;
						gap: 0px 0.5rem;
					}
				}
			}
		}
		.download {
			display: table-row;
			width: 100%;
		}
	}

	.body {
		display: grid;
		gap: 10px;
		padding: 0 0 0 5px;

		.song {
			display: grid;
			grid-template-columns: 1fr auto;
			gap: 8px;
			.data {
				display: grid;
				grid-template-columns: auto 1fr 1fr auto;
				align-items: center;
				gap: 8px;
				border: none;

				.artists {
					display: flex;
					gap: 1rem;
					overflow: hidden;
					a {
						white-space: nowrap;
						overflow-x: hidden;
						text-overflow: ellipsis;
					}
				}
			}
			.data:hover {
				outline: 1px solid #ffffff;
				outline-offset: -1px;
				background-color: inherit;
				color: inherit;
			}

			.download {
				aspect-ratio: 1/1;
			}
		}
	}
</style>
