<script lang="ts">
	import { afterNavigate } from "$app/navigation";
	import { page } from "$app/state";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import type { SongData } from "$lib/types/response";
	import { quality } from "$lib/types/quality";
	import { Clock4, Download } from "lucide-svelte";
	import Quality from "$lib/components/quality.svelte";

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
		>{song?.title || "Song"} | {song?.artists
			.map((artist) => artist.name)
			.join(" ") || "Artist"} - MusicShack</title
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
					<h1 class="title">{song.title}</h1>
					<a
						class="album"
						href="/album/{page.params.api}/{song.album.id}"
					>
						{song.album.title}
					</a>
					<nav class="artists">
						{#each song.artists as artist}
							<a href="/artist/{page.params.api}/{artist.id}">
								{artist.name}
							</a>
						{/each}
					</nav>
					<br />
					<div class="duration">
						<Clock4 size={16} />
						<p>
							{`${Math.floor(song.duration / 60)}:${(song.duration % 60).toString().padStart(2, "0")}`}
						</p>
					</div>
					<Quality quality={quality[song.audioQuality]} />
				</div>
			</div>
		</div>
		<button
			class="download hover-full"
			onclick={async () => {
				error = await download({
					api: page.params.api!,
					type: "song",
					id: song!.id,
					quality: "",
				});
			}}
		>
			<Download />
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
	}

	.header {
		margin: 15px auto 0;
		display: table;
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
					width: 280px;
					height: 280px;
				}
				.data {
					display: flex;
					flex-direction: column;
					gap: 7px;

					.title {
						font-weight: bolder;
					}
					.album {
						font-weight: bolder;
						font-style: italic;
						text-decoration: none;
					}
					.artists {
						font-style: italic;
						display: flex;
						flex-wrap: wrap;
						gap: 0.25rem 0.5rem;
					}
					.duration {
						display: flex;
						align-items: center;
						gap: 0.25rem;
					}
				}
			}
		}
		.download {
			display: flex;
			flex-direction: row;
			gap: 0.5rem;
			width: 100%;
			align-items: center;
			justify-content: center;
		}
	}
</style>
