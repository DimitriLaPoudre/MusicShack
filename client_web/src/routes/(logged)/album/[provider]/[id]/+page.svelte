<script lang="ts">
	import { Clock4, Download } from "lucide-svelte";
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import type { AlbumData, AlbumDataSong } from "$lib/types/response";
	import Quality from "$lib/components/quality.svelte";
	import Explicit from "$lib/components/explicit.svelte";

	let error = $state<null | string>(null);
	let album = $state<null | AlbumData>(null);
	let discs = $state<null | AlbumDataSong[][]>(null);

	const provider = $derived(page.params.provider);
	const id = $derived(page.params.id);

	async function fetchData(
		provider: string | undefined,
		id: string | undefined,
	) {
		album = null;
		try {
			const data = await apiFetch<AlbumData>(`/album/${provider}/${id}`);
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
	}

	$effect(() => {
		if (provider && id) {
			fetchData(provider, id);
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
	<div class="mt-4 flex flex-col justify-center items-center gap-2.5">
		<h2>Error loading Album</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !album}
	<p class="mt-6 text-center">Loading...</p>
{:else}
	<!-- page top -->
	<div class="mt-1 mx-auto table border-separate border-spacing-y-2.5">
		<div class="table-row">
			<div class="flex flex-row flex-wrap justify-center gap-2.5">
				<img class="w-70 h-70" src={album.coverUrl} alt={album.title} />
				<div class="flex flex-col gap-1.75">
					<h1 class="flex flex-row items-center gap-2 font-extrabold">
						{album.title}
						{#if album.explicit}
							<Explicit />
						{/if}
					</h1>
					<div class="italic flex flex-wrap gap-x-2 gap-y-1">
						{#each album.artists as artist}
							<a
								href="/artist/{page.params
									.provider}/{artist.id}"
							>
								{artist.name}
							</a>
						{/each}
					</div>
					<br />
					<p>{album.numberTracks} Tracks</p>
					{#if album.numberVolumes > 1}
						<p>{album.numberVolumes} Discs</p>
					{/if}
					<div class="flex items-center gap-1">
						<Clock4 size={16} />
						<p>
							{`${Math.floor(album.duration / 60)}:${(album.duration % 60).toString().padStart(2, "0")}`}
						</p>
					</div>
					<p>{album.releaseDate}</p>
					<Quality quality={album.audioQuality} />
				</div>
			</div>
		</div>
		<button
			class="hover-full flex flex-row gap-2 w-full items-center justify-center"
			onclick={async () => {
				error = await download({
					provider: page.params.provider!,
					type: "album",
					id: album!.id,
				});
			}}
		>
			<Download />
			Download Album
		</button>
	</div>

	<!-- page body -->
	<div class="pt-8 flex flex-col gap-8">
		{#each discs as disc, i}
			{#if disc}
				<div class="grid gap-2.5 pl-1.25">
					{#if album.numberVolumes > 1}
						<h2 class="text-center font-bold">Disc {i}</h2>
					{/if}
					{#each disc as song}
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
										goto(
											`/song/${page.params.provider}/${song.id}`,
										);
									}}
								>
									<p>{song.trackNumber}</p>
									<div
										class="flex flex-row flex-wrap justify-center items-center w-full gap-2"
									>
										<p
											class="flex flex-row items-center justify-center gap-2 font-extrabold wrap-break-words"
										>
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
		{/each}
	</div>
{/if}
