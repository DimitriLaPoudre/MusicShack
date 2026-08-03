<script lang="ts">
	import { page } from "$app/state";
	import {
		ApiError,
		addFollow,
		deleteFollow,
		downloadArtist,
		getArtistAlbums,
		getArtistInfo,
		listDownloads,
		listFollows,
	} from "$lib/api";
	import { Button } from "$lib/components/ui/button/index.js";
	import {
		Accordion,
		AccordionContent,
		AccordionHeader,
		AccordionItem,
		AccordionTrigger,
	} from "$lib/components/ui/accordion/index.js";
	import AlbumCard from "$lib/components/album-card.svelte";
	import LoadMore from "$lib/components/load-more.svelte";
	import { ChevronDown, Download, Heart } from "@lucide/svelte";
	import { downloadList, followList } from "$lib/stores";
	import type { AlbumInfo, ArtistInfo } from "$lib/types";

	let artist = $state<ArtistInfo | null>(null);
	let albums = $state<AlbumInfo[]>([]);
	let eps = $state<AlbumInfo[]>([]);
	let singles = $state<AlbumInfo[]>([]);
	let totalAlbums = $state(0);
	let loadingMore = $state(false);
	let openSections = $state<string[]>([]);
	let error = $state<string | null>(null);
	let followInProgress = $state(false);

	const provider = $derived(page.params.provider);
	const id = $derived(page.params.id);

	const sections = $derived([
		{ label: "Albums", items: albums },
		{ label: "EP", items: eps },
		{ label: "Singles", items: singles },
	]);

	const hasMore = $derived(albums.length < totalAlbums);

	async function load() {
		if (!provider || !id) return;
		artist = null;
		albums = [];
		eps = [];
		singles = [];
		totalAlbums = 0;
		error = null;
		try {
			const [info, albumsData] = await Promise.all([
				getArtistInfo(provider, id),
				getArtistAlbums(provider, id, { limit: 100 }),
			]);
			artist = info;
			albums = albumsData.albums.items;
			eps = albumsData.ep.items;
			singles = albumsData.singles.items;
			totalAlbums = albumsData.albums.total_number_of_items;
			openSections = sections
				.filter((section) => section.items.length > 0)
				.map((section) => section.label);
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to load artist";
		}
	}

	async function loadMore() {
		if (!provider || !id || loadingMore) return;
		loadingMore = true;
		try {
			const albumsData = await getArtistAlbums(provider, id, {
				limit: 100,
				offset: albums.length,
			});
			albums = [...albums, ...albumsData.albums.items];
			totalAlbums = albumsData.albums.total_number_of_items;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to load artist";
		}
		loadingMore = false;
	}

	$effect(() => {
		if (provider && id) load();
	});

	async function refreshFollows() {
		try {
			followList.set(await listFollows());
		} catch {
			// keep cached list on refresh failure
		}
	}

	async function toggleFollow() {
		if (!artist || followInProgress) return;
		followInProgress = true;
		try {
			if (artist.followed) {
				await deleteFollow(artist.followed);
				artist.followed = null;
			} else {
				const follow = await addFollow({
					provider: provider!,
					artist_id: artist.id,
					featuring: false,
				});
				artist.followed = follow.id;
			}
			error = null;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to toggle follow";
		}
		followInProgress = false;
		await refreshFollows();
	}

	async function startDownload() {
		if (!provider || !id || !artist) return;
		try {
			await downloadArtist({ provider, id });
			error = null;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to download discography";
		}
		try {
			downloadList.set(await listDownloads());
		} catch {
			// keep cached list on refresh failure
		}
	}
</script>

<svelte:head>
	<title>{artist?.name || "Artist"} - MusicShack</title>
</svelte:head>

{#if error}
	<div class="flex flex-col items-center justify-center gap-2 py-4">
		<h1 class="text-xs font-semibold tracking-widest uppercase">
			Error loading Artist
		</h1>
		<p class="text-sm text-muted-foreground">{error}</p>
		<a class="underline" href="/dashboard">Go to Home</a>
	</div>
{:else if !artist}
	<p class="py-4 text-center text-sm text-muted-foreground">Loading...</p>
{:else}
	<div class="mt-1 mx-auto table border-separate border-spacing-y-2.5">
		<div class="table-row">
			<div class="flex flex-row flex-wrap justify-center gap-2.5">
				{#if artist.picture_url}
					<img
						class="size-70 rounded-full object-cover"
						src={artist.picture_url}
						alt={artist.name}
					/>
				{:else}
					<div class="size-70 rounded-full bg-muted"></div>
				{/if}
				<h1 class="m-auto font-extrabold">{artist.name}</h1>
			</div>
		</div>
		<div class="table-row">
			<div class="grid grid-cols-2 gap-2">
				<Button
					variant="hover-full"
					class="flex-col"
					onclick={toggleFollow}
				>
					{#if artist.followed}
						<p>Followed</p>
						<Heart class="fill-current" />
					{:else}
						<p>Follow</p>
						<Heart />
					{/if}
				</Button>
				<Button
					variant="hover-full"
					class="flex-col"
					onclick={startDownload}
				>
					<p>Download Discography</p>
					<Download />
				</Button>
			</div>
		</div>
	</div>

	<Accordion
		type="multiple"
		value={openSections}
		onValueChange={(value: string[]) => (openSections = value)}
		class="flex w-full flex-col items-center gap-2"
	>
		{#each sections as section}
			{#if section.items.length > 0}
				<AccordionItem value={section.label} class="w-full">
					<AccordionHeader>
						<AccordionTrigger>
							{section.label}
							<ChevronDown
								class="size-4 transition-transform group-data-[state=open]:rotate-180"
							/>
						</AccordionTrigger>
					</AccordionHeader>
					<AccordionContent>
						<div
							class="grid grid-cols-[repeat(auto-fit,200px)] justify-center gap-4"
						>
							{#each section.items as album}
								<AlbumCard
									provider={provider!}
									album={album}
									onerror={(msg) => (error = msg)}
								/>
							{/each}
						</div>
					</AccordionContent>
				</AccordionItem>
			{/if}
		{/each}
	</Accordion>
	<div class="pt-8">
		<LoadMore
			hasMore={hasMore}
			loading={loadingMore}
			onclick={loadMore}
		/>
	</div>
{/if}
