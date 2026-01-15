<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { Download, HeartIcon, HeartOff } from "lucide-svelte";
	import {
		addFollow,
		loadFollows,
		removeFollow,
	} from "$lib/functions/follow";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import type { ArtistData, ArtistDataAlbum } from "$lib/types/response";
	import Quality from "$lib/components/quality.svelte";
	import Explicit from "$lib/components/explicit.svelte";

	let error = $state<null | string>(null);
	let artist = $state<null | ArtistData>(null);
	let albums = $state<Record<"Albums" | "EP" | "Singles", ArtistDataAlbum[]>>(
		{
			Albums: [],
			EP: [],
			Singles: [],
		},
	);
	let followed = $state<null | number>(null);

	const provider = $derived(page.params.provider);
	const id = $derived(page.params.id);

	async function fetchData(
		provider: string | undefined,
		id: string | undefined,
	) {
		artist = null;
		try {
			const data = await apiFetch<ArtistData>(
				`/artist/${provider}/${id}`,
			);
			if ("error" in data) {
				throw new Error(data.error || "Failed to fetch artist");
			}
			artist = data;
			albums["Albums"] = artist.albums;
			albums["EP"] = artist.ep;
			albums["Singles"] = artist.singles;
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load artist";
		}
		setFollowButton();
	}

	$effect(() => {
		if (provider && id) {
			fetchData(provider, id);
		}
	});

	async function setFollowButton() {
		const { list, error } = await loadFollows();
		if (error) {
			return;
		}
		const follow = list?.find(
			(item) =>
				item.provider === page.params.provider &&
				item.artistId === page.params.id,
		);
		if (follow) {
			followed = follow.id;
		} else {
			followed = null;
		}
	}
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
				<h1 class="name">{artist.name}</h1>
			</div>
		</div>
		<div class="bottom">
			<div class="bottom-data">
				<button
					class="hover-full"
					onclick={async () => {
						if (followed) {
							await removeFollow(followed);
						} else {
							await addFollow({
								provider: page.params.provider!,
								id: artist!.id,
							});
						}
						await setFollowButton();
					}}
				>
					{#if followed}
						<p>Unfollow</p>
						<HeartOff />
					{:else}
						<p>Follow</p>
						<HeartIcon />
					{/if}
				</button>
				<button
					class="hover-full"
					onclick={async () => {
						error = await download({
							provider: page.params.provider!,
							type: "artist",
							id: artist!.id,
						});
					}}
				>
					<p>Download Discography</p>
					<Download />
				</button>
			</div>
		</div>
	</div>

	<!-- page body -->
	<div>
		{#each Object.entries(albums) as [type, list]}
			{#if list && list.length > 0}
				<div class="wrap-container">
					<h2>{type}</h2>
					<div class="container">
						{#each list as album}
							<div class="wrap-item">
								<button
									class="item hover-full"
									onclick={(e) => {
										if (
											e.target instanceof Element &&
											e.target.closest("a")
										)
											return;
										goto(
											`/album/${page.params.provider}/${album.id}`,
										);
									}}
								>
									<img
										class="cover"
										src={album.coverUrl}
										alt={album.title}
									/>
									<p class="title">
										{album.title}
										{#if album.explicit}
											<Explicit />
										{/if}
									</p>
									<nav class="artists">
										{#each album.artists as artist}
											<a
												href="/artist/{page.params
													.provider}/{artist.id}"
											>
												{artist.name}
											</a>
										{/each}
									</nav>
									<Quality quality={album.audioQuality} />
								</button>
								<button
									class="download hover-full"
									onclick={async () =>
										(error = await download({
											provider: page.params.provider!,
											type: "album",
											id: album!.id,
										}))}
								>
									<Download />
								</button>
							</div>
						{/each}
					</div>
				</div>
			{/if}
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
					width: 280px;
					height: 280px;
				}
				.name {
					font-weight: bolder;
					margin: auto;
				}
			}
		}
		.bottom {
			display: table-row;

			.bottom-data {
				display: grid;
				grid-template-columns: 1fr 1fr;
				gap: 0.5rem;
			}
		}
	}

	.wrap-container {
		padding-top: 2rem;
		h2 {
			text-align: center;
			font-weight: bold;
		}

		.container {
			display: grid;
			grid-template-columns: repeat(auto-fit, 200px);
			justify-content: center;
			gap: 1rem;

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
					box-shadow:
						inset 0 1px 0 var(--fg),
						inset 1px 0 0 var(--fg),
						inset -1px 0 0 var(--fg);
					gap: 0.5rem;

					.cover {
						width: 160px;
						height: 160px;
					}
					.title {
						display: flex;
						flex-direction: row;
						align-items: center;
						justify-content: center;
						gap: 0.5rem;
						font-weight: bolder;
					}
					.artists {
						display: flex;
						flex-direction: column;
						gap: 0.2rem 1rem;
						font-style: italic;
					}
				}
				.download {
					width: 100%;
					box-shadow:
						inset 0 -1px 0 var(--fg),
						inset 1px 0 0 var(--fg),
						inset -1px 0 0 var(--fg);
					padding: 0.75rem;
				}
			}
		}
	}
</style>
