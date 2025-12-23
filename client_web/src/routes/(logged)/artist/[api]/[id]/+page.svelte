<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";
	import { Download } from "lucide-svelte";
	import { addFollow } from "$lib/functions/follow";
	import { apiFetch } from "$lib/functions/fetch";
	import { downloadAlbum, downloadArtist } from "$lib/functions/download";

	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let artist = $state<any | null>(null);

	afterNavigate(async () => {
		try {
			const res = await apiFetch(
				`/artist/${page.params.api}/${page.params.id}`,
			);
			artist = await res.json();
			if (!res.ok) {
				throw new Error(artist.error || "Failed to fetch artist");
			}
			isLoading = false;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load artist";
			isLoading = false;
		}
	});
</script>

<svelte:head>
	<title>{artist?.Name || "Artist"} - MusicShack</title>
</svelte:head>

{#if isLoading}
	<p class="loading">Loading...</p>
{:else if error}
	<div class="error">
		<h2>Error loading Artist</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else}
	<!-- page top -->
	<div class="header">
		<div class="top">
			<div class="top-data">
				<img
					class="picture"
					src={artist.PictureUrl}
					alt={artist.Name}
				/>
				<div class="data">
					<h1>{artist.Name}</h1>
				</div>
			</div>
		</div>
		<div class="bottom">
			<div class="bottom-data">
				<button
					onclick={() => {
						addFollow({ api: page.params.api!, id: artist.Id });
					}}>Favorite</button
				>
				<button
					onclick={async () => {
						error = await downloadArtist(
							page.params.api!,
							artist.Id,
						);
					}}
				>
					Download Song
				</button>
			</div>
		</div>
	</div>

	<!-- page body -->
	<div>
		<div>
			<h2>Albums</h2>
			<div class="container">
				{#each artist.Albums as album}
					<div class="wrap-item">
						<button
							class="item"
							onclick={(e) => {
								if (
									e.target instanceof Element &&
									e.target.closest("a")
								)
									return;
								goto(`/album/${page.params.api}/${album.Id}`);
							}}
						>
							<img src={album.CoverUrl} alt={album.Title} />
							<p>{album.Title}</p>
							<nav>
								{#each album.Artists as artist}
									<a
										href="/artist/{page.params
											.api}/{artist.Id}"
									>
										{artist.Name}
									</a>
								{/each}
							</nav>
						</button>
						<button
							class="download"
							onclick={async () =>
								(error = await downloadAlbum(
									page.params.api!,
									album.Id,
								))}
						>
							<Download />
						</button>
					</div>
				{/each}
			</div>
		</div>
		<div>
			<h2>EPs</h2>
			<div class="container">
				{#each artist.Ep as ep}
					<button
						class="item"
						onclick={() => {
							goto(`/album/${page.params.api}/${ep.Id}`);
						}}
					>
						<img src={ep.CoverUrl} alt={ep.Title} />
						<p>{ep.Title}</p>
						<div class="list">
							{#each ep.Artists as artist}
								<a href="/artist/{page.params.api}/{artist.Id}">
									{artist.Name}
								</a>
							{/each}
						</div>
					</button>
				{/each}
			</div>
		</div>
		<div>
			<h2>Singles</h2>
			<div class="container">
				{#each artist.Singles as single}
					<button
						class="item"
						onclick={() => {
							goto(`/album/${page.params.api}/${single.Id}`);
						}}
					>
						<img src={single.CoverUrl} alt={single.Title} />
						<p>{single.Title}</p>

						<div class="list">
							{#each single.Artists as artist}
								<a href="/artist/{page.params.api}/{artist.Id}">
									{artist.Name}
								</a>
							{/each}
						</div>
					</button>
				{/each}
			</div>
		</div>
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

	.picture {
		width: 160px;
		height: 160px;
	}

	.data {
		margin: auto;
	}

	.bottom {
		display: table-row;

		.bottom-data {
			display: flex;
			flex-direction: row;
			flex-wrap: wrap;
			gap: 10px;
			button {
				flex: 1 1 calc(50% - 5px);
			}
		}
	}

	.container {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;
		.wrap-item {
			width: 200px;
			height: auto;
			.item {
				display: flex;
				flex-direction: column;
				align-items: center;
				width: 200px;
				height: auto;
				overflow: hidden;
				border-bottom: none;

				img {
					width: 160px;
					height: 160px;
				}

				nav {
					display: flex;
					flex-direction: column;
					gap: 0.2rem 1rem;
				}
			}
			.download {
				width: 100%;
				border-top: none;
			}
		}
	}
</style>
