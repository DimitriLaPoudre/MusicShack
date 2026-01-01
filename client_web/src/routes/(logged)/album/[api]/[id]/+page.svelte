<script lang="ts">
	import { Clock4, Download } from "lucide-svelte";
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import type { AlbumData, AlbumDataSong } from "$lib/types/response";
	import { quality } from "$lib/types/quality";

	let error = $state<null | string>(null);
	let album = $state<null | AlbumData>(null);
	let discs = $state<null | AlbumDataSong[][]>(null);

	afterNavigate(async () => {
		try {
			const data = await apiFetch<AlbumData>(
				`/album/${page.params.api}/${page.params.id}`,
			);
			if ("error" in data) {
				throw new Error(data.error || "Failed to fetch album");
			}
			album = data;
			discs = [];
			for (const song of album.songs) {
				if (!(song.volumeNumber in discs))
					discs[song.volumeNumber] = [];
				discs[song.volumeNumber][song.trackNumber] = song;
			}
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load album";
		}
	});
</script>

<svelte:head>
	<title
		>{album?.title || "Album"} | {album?.artists
			.map((artist) => artist.name)
			.join(" ") || "Artist"} - MusicShack</title
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
					<h1 class="title">{album.title}</h1>
					<div class="artists">
						{#each album.artists as artist}
							<a href="/artist/{page.params.api}/{artist.id}">
								{artist.name}
							</a>
						{/each}
					</div>
					<br />
					<div class="duration">
						<Clock4 size={16} />
						<p>
							{`${Math.floor(album.duration / 60)}:${(album.duration % 60).toString().padStart(2, "0")}`}
						</p>
					</div>
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
	<div class="discs">
		{#each discs as disc, i}
			{#if disc !== undefined}
				<div class="disc">
					<h1>Disc {i}</h1>
					{#each disc as song}
						{#if song !== undefined}
							<div class="item">
								<button
									class="song"
									onclick={(e) => {
										if (
											e.target instanceof Element &&
											e.target.closest("a")
										)
											return;
										goto(
											`/song/${page.params.api}/${song.id}`,
										);
									}}
								>
									<p class="number">{song.trackNumber}</p>
									<div class="data">
										<p class="title">{song.title}</p>
										<div class="artists">
											{#each song.artists as artist}
												<a
													href="/artist/{page.params
														.api}/{artist.id}"
												>
													{artist.name}
												</a>
											{/each}
										</div>
									</div>
									<div class="duration">
										<Clock4 size={16} />
										<p>
											{`${Math.floor(song.duration / 60)}:${(song.duration % 60).toString().padStart(2, "0")}`}
										</p>
									</div>
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
						{/if}
					{/each}
				</div>
			{/if}
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
					width: 160px;
					height: 160px;
				}
				.data {
					display: flex;
					flex-direction: column;
					gap: 7px;

					.title {
						font-weight: bolder;
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
			display: table-row;
			width: 100%;
		}
	}

	.discs {
		display: flex;
		flex-direction: column;
		gap: 2rem;

		.disc {
			display: grid;
			gap: 10px;
			padding: 0 0 0 5px;

			.item {
				display: grid;
				grid-template-columns: 1fr auto;
				gap: 8px;

				.song {
					display: grid;
					grid-template-columns: auto 1fr auto;
					align-items: center;
					gap: 8px;
					border: none;

					.data {
						display: flex;
						flex-direction: row;
						flex-wrap: wrap;
						justify-content: center;
						align-items: center;
						width: 100%;
						gap: 1rem;

						.title {
							font-weight: bolder;
							flex: 0 0 49%;
						}
						.artists {
							font-style: italic;
							flex: 0 0 49%;
							display: flex;
							flex-wrap: wrap;
							justify-content: center;
							align-items: center;
							gap: 0.25rem 0.5rem;
						}
					}
					.duration {
						display: flex;
						align-items: center;
						gap: 0.25rem;
					}
				}
				.song:hover {
					outline: 1px solid #ffffff;
					outline-offset: -1px;
					background-color: inherit;
					color: inherit;
				}
				.download {
					height: auto;
				}
			}
		}
	}
</style>
