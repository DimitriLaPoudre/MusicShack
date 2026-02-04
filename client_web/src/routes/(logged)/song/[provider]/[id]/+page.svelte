<script lang="ts">
	import { page } from "$app/state";
	import { apiFetch } from "$lib/functions/fetch";
	import { download } from "$lib/functions/download";
	import type { SongData } from "$lib/types/response";
	import { Clock4, Download } from "lucide-svelte";
	import Quality from "$lib/components/quality.svelte";
	import Explicit from "$lib/components/explicit.svelte";

	let error = $state<null | string>(null);
	let song = $state<null | SongData>(null);

	const provider = $derived(page.params.provider);
	const id = $derived(page.params.id);

	async function fetchData(
		provider: string | undefined,
		id: string | undefined,
	) {
		song = null;
		try {
			const data = await apiFetch<SongData>(`/song/${provider}/${id}`);
			if ("error" in data) {
				throw new Error(data.error || "Failed to fetch song");
			}
			song = data;
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to load song";
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
		>{song?.title || "Song"} | {song?.artists
			.map((artist) => artist.name)
			.join(" ") || "Artist"} - MusicShack</title
	>
</svelte:head>

{#if error}
	<div class="mt-4 flex flex-col justify-center items-center gap-2.5">
		<h2>Error loading Song</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !song}
	<p class="mt-6 text-center">Loading...</p>
{:else}
	<!-- page top -->
	<div class="mt-1 mx-auto table border-separate border-spacing-y-2.5">
		<div class="table-row">
			<div class="flex flex-row flex-wrap justify-center gap-2.5">
				<img class="w-[280px] h-[280px]" src={song.album.coverUrl} alt={song.title} />
				<div class="flex flex-col gap-[7px]">
					<h1 class="flex flex-row items-center gap-2 font-extrabold">
						{song.title}
						{#if song.explicit}
							<Explicit />
						{/if}
					</h1>
					<a
						class="font-extrabold italic no-underline"
						href="/album/{page.params.provider}/{song.album.id}"
					>
						{song.album.title}
					</a>
					<nav class="italic flex flex-wrap gap-x-2 gap-y-1">
						{#each song.artists as artist}
							<a
								href="/artist/{page.params
									.provider}/{artist.id}"
							>
								{artist.name}
							</a>
						{/each}
					</nav>
					<br />
					<div class="flex items-center gap-1">
						<Clock4 size={16} />
						<p>
							{`${Math.floor(song.duration / 60)}:${(song.duration % 60).toString().padStart(2, "0")}`}
						</p>
					</div>
					<Quality quality={song.audioQuality} />
				</div>
			</div>
		</div>
		<button
			class="hover-full flex flex-row gap-2 w-full items-center justify-center"
			onclick={async () => {
				error = await download({
					provider: page.params.provider!,
					type: "song",
					id: song!.id,
				});
			}}
		>
			<Download />
			Download Song
		</button>
	</div>
{/if}
