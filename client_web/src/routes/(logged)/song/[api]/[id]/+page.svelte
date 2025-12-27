<script lang="ts">
	import { afterNavigate } from "$app/navigation";
	import { page } from "$app/state";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import type { SongData } from "$lib/types/response";
	import { quality } from "$lib/types/quality";

	let error = $state<null | string>(null);
	let song = $state<null | SongData>(null);

	afterNavigate(async () => {
		try {
			const data = await apiFetch<SongData>(
				`/song/${page.params.api}/${page.params.id}`,
			);
			if ("error" in data) {
				throw new Error(data.error || "Failed to fetch song");
			}
			song = data;
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load song";
		}
	});
</script>

<svelte:head>
	<title
		>{song?.title || "Song"} | {song?.artists[0]?.name || "Artist"} - MusicShack</title
	>
</svelte:head>

{#if error}
	<div class="error">
		<h2>Error loading Song</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !song}
	<p class="loading">Loading...</p>
{:else}
	<!-- page top -->
	<div class="header">
		<div class="top">
			<div class="top-data">
				<img class="cover" src={song.album.coverUrl} alt={song.title} />
				<div class="data">
					<h1>{song.title}</h1>
					<a href="/album/{page.params.api}/{song.album.id}">
						{song.album.title}
					</a>
					<div class="artists">
						{#each song.artists as artist}
							<a href="/artist/{page.params.api}/{artist.id}">
								{artist.name}
							</a>
						{/each}
					</div>
					<br />
					<p>
						{`${Math.floor(song.duration / 60)}:${(song.duration % 60).toString().padStart(2, "0")}`}
					</p>
					<p>{quality[song.audioQuality]}</p>
				</div>
			</div>
		</div>
		<button
			class="download"
			onclick={async () => {
				error = await download({
					api: page.params.api!,
					type: "song",
					id: song!.id,
					quality: "",
				});
			}}
		>
			Download Song
		</button>
	</div>
{/if}

<style>
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
</style>
