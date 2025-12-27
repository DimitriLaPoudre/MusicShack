<script lang="ts">
	import { Download } from "lucide-svelte";
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import type { AlbumData } from "$lib/types/response";
	import { quality } from "$lib/types/quality";

	let error = $state<null | string>(null);
	let album = $state<null | AlbumData>(null);

	afterNavigate(async () => {
		try {
			const data = await apiFetch<AlbumData>(
				`/album/${page.params.api}/${page.params.id}`,
			);
			if ("error" in data) {
				throw new Error(data.error || "Failed to fetch album");
			}
			album = data;
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load album";
		}
	});
</script>

<svelte:head>
	<title
		>{album?.title || "Album"} | {album?.artists[0]?.name || "Artist"} - MusicShack</title
	>
</svelte:head>

{#if error}
	<div class="error">
		<h2>Error loading Album</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !album}
	<p class="loading">Loading...</p>
{:else}
	<!-- page top -->
	<div class="header">
		<div class="top">
			<div class="top-data">
				<img class="cover" src={album.coverUrl} alt={album.title} />
				<div class="data">
					<h1>{album.title}</h1>
					<div class="artists">
						{#each album.artists as artist}
							<a href="/artist/{page.params.api}/{artist.id}">
								{artist.name}
							</a>
						{/each}
					</div>
					<br />
					<p>
						{`${Math.floor(album.duration / 60)}:${(album.duration % 60).toString().padStart(2, "0")}`}
					</p>
					<p>{quality[album.audioQuality]}</p>
				</div>
			</div>
		</div>
		<button
			class="download"
			onclick={async () => {
				error = await download({
					api: page.params.api!,
					type: "album",
					id: album!.id,
					quality: "",
				});
			}}
		>
			Download Album
		</button>
	</div>

	<!-- page body -->
	<div class="body">
		{#each album.songs as song}
			<div class="song">
				<button
					class="data"
					onclick={(e) => {
						if (
							e.target instanceof Element &&
							e.target.closest("a")
						)
							return;
						goto(`/song/${page.params.api}/${song.id}`);
					}}
				>
					<p class="number">{song.trackNumber}</p>
					<p>{song.title}</p>
					<div class="artists">
						{#each song.artists as artist}
							<a href="/artist/{page.params.api}/{artist.id}">
								{artist.name}
							</a>
						{/each}
					</div>
					<p class="duration">
						{`${Math.floor(song.duration / 60)}:${(song.duration % 60).toString().padStart(2, "0")}`}
					</p>
				</button>
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
					<Download />
				</button>
			</div>
		{/each}
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
