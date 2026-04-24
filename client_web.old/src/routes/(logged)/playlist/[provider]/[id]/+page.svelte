<script lang="ts">
	import { Clock4, Download } from "lucide-svelte";
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import Explicit from "$lib/components/Explicit.svelte";
	import Owned from "$lib/components/Owned.svelte";
	import type { PlaylistData } from "$lib/types/response";

	let error = $state<null | string>(null);
	let playlist = $state<null | PlaylistData>(null);

	const provider = $derived(page.params.provider);
	const id = $derived(page.params.id);

	async function fetchData(
		provider: string | undefined,
		id: string | undefined,
	) {
		try {
			const data = await apiFetch<PlaylistData>(
				`/playlist/${provider}/${id}`,
			);
			playlist = data;
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load playlist";
		}
	}

	$effect(() => {
		if (provider && id) {
			fetchData(provider, id);
		}
	});
</script>

<svelte:head>
	<title>{playlist?.title || "Playlist"} - MusicShack</title>
</svelte:head>

{#if error}
	<div class="mt-4 flex flex-col justify-center items-center gap-2.5">
		<h2>Error loading playlist</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !playlist}
	<p class="mt-6 text-center">Loading...</p>
{:else}
	<!-- page top -->
	<div class="mt-1 mx-auto table border-separate border-spacing-y-2.5">
		<div class="table-row">
			<div class="flex flex-row flex-wrap justify-center gap-2.5">
				<img
					class="w-70 h-70"
					src={playlist.coverUrl}
					alt={playlist.title}
				/>
				<div class="flex flex-col gap-1.75">
					<h1 class="flex flex-row items-center gap-2 font-extrabold">
						{#if playlist.downloaded}
							<Owned />
						{/if}
						{playlist.title}
					</h1>
					<h2>
						{playlist.description}
					</h2>
					<br />
					<p>{playlist.numberOfTracks} Tracks</p>
					<div class="flex items-center gap-1">
						<Clock4 size={16} />
						<p>
							{`${Math.floor(playlist.duration / 60)}:${(playlist.duration % 60).toString().padStart(2, "0")}`}
						</p>
					</div>
					<p>{playlist.lastUpdated}</p>
				</div>
			</div>
		</div>
		<button
			class="hover-full flex flex-row gap-2 w-full items-center justify-center"
			onclick={async () => {
				error = await download({
					provider: page.params.provider!,
					type: "playlist",
					id: playlist!.id,
				});
			}}
		>
			<Download />
			Download all song separately
		</button>
	</div>

	<!-- page body -->
	<div class="pt-8 flex flex-col gap-8">
		{#each playlist.songs as song, index}
			{#if song}
				<div class="grid grid-cols-[1fr_auto] gap-2">
					<button
						class="hover-soft grid grid-cols-[auto_1fr_auto] items-center gap-2 border-none"
						onclick={(e) => {
							if (
								e.target instanceof Element &&
								e.target.closest("a")
							)
								return;
							goto(`/song/${page.params.provider}/${song.id}`);
						}}
					>
						<p>{index + 1}</p>
						<div
							class="flex flex-row flex-wrap justify-center items-center w-full gap-2"
						>
							<p
								class="flex flex-row items-center justify-center gap-2 font-extrabold wrap-break-words"
							>
								{#if song.downloaded}
									<Owned />
								{/if}
								{song.title}
								{#if song.explicit}
									<Explicit />
								{/if}
							</p>
							<nav
								class="italic flex flex-wrap justify-center items-center gap-x-2 gap-y-1 wrap-break-words"
							>
								{#each song.artists as artist}
									<a
										href="/artist/{page.params
											.provider}/{artist.id}"
									>
										{artist.name}
									</a>
								{/each}
							</nav>
						</div>
						<div class="flex items-center gap-1">
							<Clock4 size={16} />
							<p>
								{`${Math.floor(song.duration / 60)}:${(song.duration % 60).toString().padStart(2, "0")}`}
							</p>
						</div>
					</button>
					<button
						class="hover-full h-auto"
						onclick={async () => {
							error = await download({
								provider: page.params.provider!,
								type: "song",
								id: song!.id,
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
