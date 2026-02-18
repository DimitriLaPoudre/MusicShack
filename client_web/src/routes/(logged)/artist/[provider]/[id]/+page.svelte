<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { Download, HeartIcon, HeartOff } from "lucide-svelte";
	import { addFollow, removeFollow } from "$lib/functions/follow";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import type { ArtistData, ArtistDataAlbum } from "$lib/types/response";
	import Quality from "$lib/components/Quality.svelte";
	import Explicit from "$lib/components/Explicit.svelte";
	import Owned from "$lib/components/Owned.svelte";

	let error = $state<null | string>(null);
	let artist = $state<null | ArtistData>(null);
	let albums = $state<Record<"Albums" | "EP" | "Singles", ArtistDataAlbum[]>>(
		{
			Albums: [],
			EP: [],
			Singles: [],
		},
	);
	let followInProgress = false;

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
			artist = data;
			albums["Albums"] = artist.albums;
			albums["EP"] = artist.ep;
			albums["Singles"] = artist.singles;
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load artist";
		}
	}

	$effect(() => {
		if (provider && id) {
			fetchData(provider, id);
		}
	});
</script>

<svelte:head>
	<title>{artist?.name || "Artist"} - MusicShack</title>
</svelte:head>

{#if error}
	<div class="mt-4 flex flex-col justify-center items-center gap-2.5">
		<h2>Error loading Artist</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !artist}
	<p class="mt-6 text-center">Loading...</p>
{:else}
	<!-- page top -->
	<div class="mt-1 mx-auto table border-separate border-spacing-y-2.5">
		<div class="table-row">
			<div class="flex flex-row flex-wrap justify-center gap-2.5">
				<img
					class="w-70 h-70"
					src={artist.pictureUrl}
					alt={artist.name}
				/>
				<h1 class="font-extrabold m-auto">{artist.name}</h1>
			</div>
		</div>
		<div class="table-row">
			<div class="grid grid-cols-2 gap-2">
				<button
					class="hover-full flex-col"
					onclick={async () => {
						if (artist!.followed) {
							if (followInProgress) {
								artist!.followed = 0;
								return;
							}

							const follow = artist!.followed;
							artist!.followed = 0;
							followInProgress = true;

							const error = await removeFollow(follow);
							if (error) {
								artist!.followed = follow;
							} else {
								artist!.followed = 0;
							}
							followInProgress = false;
						} else {
							if (followInProgress) {
								artist!.followed = -1;
								return;
							}

							artist!.followed = -1;
							followInProgress = true;
							const { follow, error } = await addFollow({
								provider: page.params.provider!,
								id: artist!.id,
							});
							if (error || !follow) {
								artist!.followed = 0;
							} else {
								artist!.followed = follow.id;
							}
							followInProgress = false;
						}
					}}
				>
					{#if artist!.followed}
						<p>Unfollow</p>
						<HeartOff />
					{:else}
						<p>Follow</p>
						<HeartIcon />
					{/if}
				</button>
				<button
					class="hover-full flex-col"
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
	<div class="flex flex-col gap-8 mt-4">
		{#each Object.entries(albums) as [type, list]}
			{#if list && list.length > 0}
				<div class="flex flex-col gap-2">
					<h2 class="text-center font-bold">{type}</h2>
					<div
						class="grid grid-cols-[repeat(auto-fit,200px)] justify-center gap-4"
					>
						{#each list as album}
							<div class="w-50 h-auto">
								<button
									class="hover-full flex flex-col items-center w-50 h-auto overflow-hidden gap-2 shadow-[inset_0_1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
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
										class="w-40 h-40"
										src={album.coverUrl}
										alt={album.title}
									/>
									<p
										class="flex flex-row items-center justify-center gap-2 font-extrabold"
									>
										{#if album.downloaded}
											<Owned />
										{/if}
										{album.title}
										{#if album.explicit}
											<Explicit />
										{/if}
									</p>
									<nav
										class="flex flex-col gap-y-[0.2rem] gap-x-4 italic"
									>
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
									class="hover-full w-full p-3 shadow-[inset_0_-1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
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
