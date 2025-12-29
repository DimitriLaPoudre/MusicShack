<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { page } from "$app/state";
	import { Download } from "lucide-svelte";
	import { addFollow } from "$lib/functions/follow";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import type { ArtistData } from "$lib/types/response";
	import { quality } from "$lib/types/quality";

	let error = $state<null | string>(null);
	let artist = $state<null | ArtistData>(null);

	afterNavigate(async () => {
		try {
			const data = await apiFetch<ArtistData>(
				`/artist/${page.params.api}/${page.params.id}`,
			);
			if ("error" in data) {
				throw new Error(data.error || "Failed to fetch artist");
			}
			artist = data;
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load artist";
		}
	});
</script>

<svelte:head>
	<title>{artist?.name || "Artist"} - MusicShack</title>
</svelte:head>

{#if error}
	<div class="error">
		<h2>Error loading Artist</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !artist}
	<p class="loading">Loading...</p>
{:else}
	<!-- page top -->
	<div class="header">
		<div class="top">
			<div class="top-data">
				<img
					class="picture"
					src={artist.pictureUrl}
					alt={artist.name}
				/>
				<div class="data">
					<h1>{artist.name}</h1>
				</div>
			</div>
		</div>
		<div class="bottom">
			<div class="bottom-data">
				<button
					onclick={() => {
						addFollow({ api: page.params.api!, id: artist!.id });
					}}
				>
					Follow
				</button>
				<button
					onclick={async () => {
						error = await download({
							api: page.params.api!,
							type: "artist",
							id: artist!.id,
							quality: "",
						});
					}}
				>
					Download Discography
				</button>
			</div>
		</div>
	</div>

	<!-- page body -->
	<div>
		<div>
			<h2>Albums</h2>
			<div class="container">
				{#each artist.albums as album}
					<div class="wrap-item">
						<button
							class="item"
							onclick={(e) => {
								if (
									e.target instanceof Element &&
									e.target.closest("a")
								)
									return;
								goto(`/album/${page.params.api}/${album.id}`);
							}}
						>
							<img src={album.coverUrl} alt={album.title} />
							<p>{album.title}</p>
							<nav>
								{#each album.artists as artist}
									<a
										href="/artist/{page.params
											.api}/{artist.id}"
									>
										{artist.name}
									</a>
								{/each}
							</nav>
							<p>{quality[album.audioQuality]}</p>
						</button>
						<button
							class="download"
							onclick={async () =>
								(error = await download({
									api: page.params.api!,
									type: "album",
									id: album!.id,
									quality: "",
								}))}
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
				{#each artist.ep as ep}
					<button
						class="item"
						onclick={() => {
							goto(`/album/${page.params.api}/${ep.id}`);
						}}
					>
						<img src={ep.coverUrl} alt={ep.title} />
						<p>{ep.title}</p>
						<div class="list">
							{#each ep.artists as artist}
								<a href="/artist/{page.params.api}/{artist.id}">
									{artist.name}
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
				{#each artist.singles as single}
					<button
						class="item"
						onclick={() => {
							goto(`/album/${page.params.api}/${single.id}`);
						}}
					>
						<img src={single.coverUrl} alt={single.title} />
						<p>{single.title}</p>

						<div class="list">
							{#each single.artists as artist}
								<a href="/artist/{page.params.api}/{artist.id}">
									{artist.name}
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
		margin: 15px auto 0;
		border-spacing: 0 10px;

		.top {
			display: table-row;

			.top-data {
				display: flex;
				flex-direction: row;
				flex-wrap: wrap;
				justify-content: center;
				gap: 10px;

				.picture {
					width: 160px;
					height: 160px;
				}
				.data {
					margin: auto;
				}
			}
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
				gap: 8px;

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
