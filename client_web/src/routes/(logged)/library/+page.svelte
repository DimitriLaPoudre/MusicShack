<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import {
		deleteSong,
		loadLibrary,
		syncLibrary,
	} from "$lib/functions/library";
	import {
		ChevronLeft,
		ChevronRight,
		ChevronDown,
		ChevronUp,
		Pencil,
		SearchIcon,
		Trash,
		UploadIcon,
	} from "lucide-svelte";
	import { onMount } from "svelte";
	import { Pagination, AlertDialog } from "bits-ui";
	import { libraryPage } from "$lib/stores/panel/library";
	import { page } from "$app/state";
	import type { ResponseSong, StatusResponse } from "$lib/types/response";
	import { apiFetchFormData } from "$lib/functions/fetch";
	import Explicit from "$lib/components/Explicit.svelte";

	let error = $state<null | string>(null);

	let uploadArtists = $state<string>("");
	let uploadArtistsList = $state<string[]>([]);
	let uploadAlbumArtists = $state<string>("");
	let uploadAlbumArtistsList = $state<string[]>([]);
	let uploadDialog = $state(false);
	let uploadOptionnal = $state(false);
	let errorUploadDialog = $state<null | string>(null);

	let editItem = $state<null | ResponseSong>(null);
	let editArtists = $state<string>("");
	let editArtistsList = $state<string[]>([]);
	let editAlbumArtists = $state<string>("");
	let editAlbumArtistsList = $state<string[]>([]);
	let errorEditDialog = $state<null | string>(null);
	let editDialog = $derived(editItem !== null);

	let deletedItem = $state<null | ResponseSong>(null);
	let deleteDialog = $derived(deletedItem !== null);

	let search = $derived(page.url.searchParams.get("q") || "");
	let currentPage = $derived(Number(page.url.searchParams.get("page")) || 1);
	const limit = 10;
	let offset = $derived((currentPage - 1) * limit);

	onMount(async () => {
		await syncLibrary();
		error = await loadLibrary(search, limit, offset);
	});

	afterNavigate(async () => {
		error = await loadLibrary(search, limit, offset);
	});
</script>

<svelte:head>
	<title>Library - MusicShack</title>
</svelte:head>

{#if error}
	<div class="mt-4 flex flex-col justify-center items-center gap-2.5">
		<h2>Error loading Song</h2>
		<p>{error}</p>
		<a href="/">Go to Home</a>
	</div>
{:else if !$libraryPage}
	<p class="mt-6 text-center">Loading...</p>
{:else}
	<div class="flex flex-col gap-4 py-4">
		<div class="flex items-center justify-center gap-4">
			<div>
				<SearchIcon />
			</div>

			<input
				bind:value={search}
				oninput={() => {
					let query = new URLSearchParams(
						page.url.searchParams.toString(),
					);
					query.set("q", search);
					goto(`?${query.toString()}`, {
						keepFocus: true,
					});
				}}
				placeholder="Search"
			/>
		</div>

		<AlertDialog.Root bind:open={uploadDialog}>
			<div class="flex justify-center items-center">
				<AlertDialog.Trigger
					class="hover-full"
					onclick={() => {
						errorUploadDialog = null;
						uploadOptionnal = false;
						uploadAlbumArtistsList = [];
						uploadArtistsList = [];
					}}
				>
					<div class="flex gap-4">
						<UploadIcon />
						<p>Upload</p>
					</div>
				</AlertDialog.Trigger>
			</div>
			<AlertDialog.Portal>
				<AlertDialog.Overlay class="fixed inset-0 z-50 bg-black/80" />
				<AlertDialog.Content
					class="fixed inset-0 z-50 flex items-center justify-center"
				>
					<div
						class="bg-bg border p-7 gap-3 grid grid-rows-[auto_1fr_auto] max-h-[calc(100%-2rem)] w-full max-w-[calc(100%-2rem)] sm:max-w-lg md:w-full"
					>
						<AlertDialog.Title class="text-2xl font-semibold">
							Upload a song
						</AlertDialog.Title>
						<form
							class="flex flex-col gap-2 overflow-y-auto p-4"
							id="upload-form"
							onsubmit={async (e) => {
								try {
									e.preventDefault();

									const form = e.currentTarget;
									const fd = new FormData(form);

									uploadAlbumArtistsList.forEach((t) =>
										fd.append("albumArtists", t),
									);
									uploadArtistsList.forEach((t) =>
										fd.append("artists", t),
									);

									const emptyKeys = [];
									for (let [key, value] of fd.entries()) {
										if (value === "") emptyKeys.push(key);
										if (
											value instanceof File &&
											value.size === 0
										) {
											emptyKeys.push(key);
										}
									}
									for (let key of emptyKeys) {
										fd.delete(key);
									}

									await apiFetchFormData<StatusResponse>(
										"/library",
										fd,
									);

									error = await loadLibrary(
										search,
										limit,
										offset,
									);
									uploadDialog = false;
								} catch (e) {
									errorUploadDialog =
										e instanceof Error
											? e.message
											: "Failed to upload song";
								}
							}}
						>
							{#if errorUploadDialog}
								<p class="text-center bg-err p-2">
									{errorUploadDialog}
								</p>
							{/if}
							<label
								class="grid grid-cols-[minmax(0,1fr)_minmax(0,3fr)] gap-2 items-center"
							>
								Cover
								<input
									type="file"
									name="cover"
									accept="image/*"
								/>
							</label>
							<label
								class="grid grid-cols-[minmax(0,1fr)_minmax(0,3fr)] gap-2 items-center"
							>
								File
								<input
									required
									type="file"
									name="file"
									accept="audio/*"
								/>
							</label>
							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								Title
								<input
									type="text"
									name="title"
									placeholder="Shawty"
								/>
							</label>
							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								Album
								<input
									type="text"
									name="album"
									placeholder="Magma Road"
								/>
							</label>

							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								<div class="flex gap-2 flex-wrap">
									<span>Album Artists: </span>
									{#each uploadAlbumArtistsList as artist, index}
										<div class="flex gap-1">
											<span class="underline">
												{artist}
											</span>
											<button
												class="p-0 m-0 border-0 bg-transparent text-inherit font-extrabold"
												type="button"
												onclick={() =>
													uploadAlbumArtistsList.splice(
														index,
														1,
													)}
											>
												x
											</button>
										</div>
									{/each}
								</div>
								<input
									bind:value={uploadAlbumArtists}
									onkeydown={(e) => {
										const albumArtist =
											uploadAlbumArtists.trim();
										if (
											e.key === "Enter" &&
											albumArtist !== ""
										) {
											e.preventDefault();
											uploadAlbumArtistsList = [
												...uploadAlbumArtistsList,
												albumArtist,
											];
											uploadAlbumArtists = "";
										}
									}}
									type="text"
									placeholder="Album Artists"
								/>
							</label>

							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								<div class="flex gap-2 flex-wrap">
									<span>Artists: </span>
									{#each uploadArtistsList as artist, index}
										<div class="flex gap-1">
											<span class="underline">
												{artist}
											</span>
											<button
												class="p-0 m-0 border-0 bg-transparent text-inherit font-extrabold"
												type="button"
												onclick={() =>
													uploadArtistsList.splice(
														index,
														1,
													)}
											>
												x
											</button>
										</div>
									{/each}
								</div>
								<input
									bind:value={uploadArtists}
									onkeydown={(e) => {
										const artist = uploadArtists.trim();
										if (
											e.key === "Enter" &&
											artist !== ""
										) {
											e.preventDefault();
											uploadArtistsList = [
												...uploadArtistsList,
												artist,
											];
											uploadArtists = "";
										}
									}}
									type="text"
									placeholder="Artists"
								/>
							</label>
							{#if !uploadOptionnal}
								<button
									type="button"
									class="flex gap-4 py-2"
									onclick={() =>
										(uploadOptionnal = !uploadOptionnal)}
								>
									<p>Optionnal</p>
									<ChevronUp />
								</button>
							{:else}
								<button
									type="button"
									class="flex gap-4 py-2"
									onclick={() =>
										(uploadOptionnal = !uploadOptionnal)}
								>
									<p>Optionnal</p>
									<ChevronDown />
								</button>

								<label
									class="grid grid-cols-[1fr_auto] gap-2 items-center"
								>
									Track Number
									<input
										type="number"
										name="trackNumber"
										placeholder="5"
										min="1"
										step="1"
									/>
								</label>
								<label
									class="grid grid-cols-[1fr_auto] gap-2 items-center"
								>
									Volume Number
									<input
										type="number"
										name="volumeNumber"
										placeholder="1"
										min="1"
										step="1"
									/>
								</label>
								<label
									class="grid grid-cols-[1fr_auto] gap-2 items-center"
								>
									ISRC
									<input
										type="text"
										name="isrc"
										placeholder="FR5R00909899"
									/>
								</label>
								<label
									class="grid grid-cols-[1fr_auto] gap-2 items-center"
								>
									Release Date
									<input type="date" name="releaseDate" />
								</label>
								<label
									class="grid grid-cols-[1fr_auto] items-center"
								>
									Explicit
									<input type="checkbox" name="explicit" />
								</label>
							{/if}
						</form>
						<div
							class="grid grid-cols-2 items-center justify-center gap-2 pt-1"
						>
							<AlertDialog.Cancel
								type="button"
								class="hover-full"
							>
								Cancel
							</AlertDialog.Cancel>
							<AlertDialog.Action
								type="submit"
								form="upload-form"
								class="hover-full"
							>
								Upload
							</AlertDialog.Action>
						</div>
					</div>
				</AlertDialog.Content>
			</AlertDialog.Portal>
		</AlertDialog.Root>
	</div>
	<Pagination.Root
		bind:page={currentPage}
		count={$libraryPage.total}
		perPage={limit}
		onPageChange={async (pageIndex: number) => {
			let query = new URLSearchParams(page.url.searchParams.toString());
			query.set("page", String(pageIndex));
			goto(`?${query.toString()}`);
		}}
	>
		{#snippet children({ pages })}
			<div class="flex items-center justify-center mb-4 mt-2">
				<Pagination.PrevButton
					class="inline-flex size-10 items-center justify-center active:text-bg active:bg-fg disabled:cursor-default disabled:text-bg hover:disabled:bg-transparent"
				>
					<ChevronLeft size={24} />
				</Pagination.PrevButton>
				<div class="flex items-center gap-2.5">
					{#each pages as page (page.key)}
						{#if page.type === "ellipsis"}
							<div
								class="text-foreground-alt select-none text-[10px] font-medium"
							>
								...
							</div>
						{:else}
							<Pagination.Page
								{page}
								class="data-selected:bg-fg data-selected:text-bg inline-flex size-10 select-none items-center justify-center text-[15px] font-medium disabled:cursor-not-allowed disabled:opacity-50"
							>
								{page.value}
							</Pagination.Page>
						{/if}
					{/each}
				</div>
				<Pagination.NextButton
					class="inline-flex size-10 items-center justify-center active:text-bg active:bg-fg disabled:cursor-default disabled:text-bg hover:disabled:bg-transparent"
				>
					<ChevronRight size={24} />
				</Pagination.NextButton>
			</div>
		{/snippet}
	</Pagination.Root>
	<div class="grid grid-cols-[repeat(auto-fit,200px)] justify-center gap-4">
		{#each $libraryPage.items as item}
			<div class="w-50 h-auto">
				<button
					class="hover-full flex flex-col items-center w-50 h-auto overflow-hidden gap-3 shadow-[inset_0_1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
				>
					<div class="w-40 h-40">
						<img
							src="/api/library/{item.id}/img"
							alt={item.title}
						/>
					</div>
					<p
						class="flex flex-row items-center justify-center gap-2 font-extrabold"
					>
						{item.title}
						{#if item.explicit}
							<Explicit />
						{/if}
					</p>
					<nav class="flex flex-col gap-y-[0.2rem] gap-x-4 italic">
						{#each item.artists as artist}
							<span>
								{artist}
							</span>
						{/each}
					</nav>
				</button>
				<div class="flex">
					<button
						class="hover-full w-full p-4 shadow-[inset_0_-1px_0_var(--fg),inset_1px_0_0_var(--fg)]"
						onclick={() => {
							editItem = item;
							editArtistsList = editItem.artists;
							editAlbumArtistsList = editItem.albumArtists;
							errorEditDialog = null;
						}}
					>
						<Pencil />
					</button>
					<button
						class="hover-full w-full p-4 shadow-[inset_0_-1px_0_var(--fg),inset_-1px_0_0_var(--fg)]"
						onclick={() => (deletedItem = item)}
					>
						<Trash />
					</button>
				</div>
			</div>
		{/each}
	</div>
	<Pagination.Root
		bind:page={currentPage}
		count={$libraryPage.total}
		perPage={limit}
		onPageChange={async (pageIndex: number) => {
			let query = new URLSearchParams(page.url.searchParams.toString());
			query.set("page", String(pageIndex));
			goto(`?${query.toString()}`);
		}}
	>
		{#snippet children({ pages })}
			<div class="flex items-center justify-center mb-4 mt-2">
				<Pagination.PrevButton
					class="inline-flex size-10 items-center justify-center active:text-bg active:bg-fg disabled:cursor-default disabled:text-bg hover:disabled:bg-transparent"
				>
					<ChevronLeft size={24} />
				</Pagination.PrevButton>
				<div class="flex items-center gap-2.5">
					{#each pages as page (page.key)}
						{#if page.type === "ellipsis"}
							<div
								class="text-foreground-alt select-none text-[10px] font-medium"
							>
								...
							</div>
						{:else}
							<Pagination.Page
								{page}
								class="data-selected:bg-fg data-selected:text-bg inline-flex size-10 select-none items-center justify-center text-[15px] font-medium disabled:cursor-not-allowed disabled:opacity-50"
							>
								{page.value}
							</Pagination.Page>
						{/if}
					{/each}
				</div>
				<Pagination.NextButton
					class="inline-flex size-10 items-center justify-center active:text-bg active:bg-fg disabled:cursor-default disabled:text-bg hover:disabled:bg-transparent"
				>
					<ChevronRight size={24} />
				</Pagination.NextButton>
			</div>
		{/snippet}
	</Pagination.Root>
	{#if editItem}
		<AlertDialog.Root bind:open={editDialog}>
			<AlertDialog.Portal>
				<AlertDialog.Overlay class="fixed inset-0 z-50 bg-black/80" />
				<AlertDialog.Content
					class="fixed inset-0 z-50 flex items-center justify-center"
				>
					<div
						class="bg-bg border p-7 gap-3 grid grid-rows-[auto_1fr_auto] max-h-[calc(100%-2rem)] w-full max-w-[calc(100%-2rem)] sm:max-w-lg md:w-full"
					>
						<AlertDialog.Title class="text-2xl font-semibold">
							Edit
						</AlertDialog.Title>
						<form
							class="flex flex-col gap-2 overflow-y-auto p-4"
							id="edit-form"
							onsubmit={async (e) => {
								try {
									e.preventDefault();

									const form = e.currentTarget;
									const fd = new FormData(form);

									editAlbumArtistsList.forEach((t) =>
										fd.append("albumArtists", t),
									);
									editArtistsList.forEach((t) =>
										fd.append("artists", t),
									);

									const emptyKeys = [];
									for (let [key, value] of fd.entries()) {
										if (value === "") emptyKeys.push(key);
										if (
											value instanceof File &&
											value.size === 0
										) {
											emptyKeys.push(key);
										}
									}
									for (let key of emptyKeys) {
										fd.delete(key);
									}

									await apiFetchFormData<StatusResponse>(
										`/library/${editItem!.id}`,
										fd,
										"PUT",
									);

									error = await loadLibrary(
										search,
										limit,
										offset,
									);

									editItem = null;
								} catch (e) {
									errorEditDialog =
										e instanceof Error
											? e.message
											: "Failed to upload song";
								}
							}}
						>
							{#if errorEditDialog}
								<p class="text-center bg-err p-2">
									{errorEditDialog}
								</p>
							{/if}
							<label
								class="grid grid-cols-[minmax(0,1fr)_minmax(0,3fr)] gap-2 items-center"
							>
								Cover
								<input
									type="file"
									name="cover"
									accept="image/*"
								/>
							</label>
							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								Title
								<input
									type="text"
									name="title"
									placeholder={editItem.title}
									value={editItem.title}
								/>
							</label>
							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								Album
								<input
									type="text"
									name="album"
									value={editItem.album}
									placeholder={editItem.album}
								/>
							</label>

							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								<div class="flex gap-2 flex-wrap">
									<span>Album Artists: </span>
									{#each editAlbumArtistsList as artist, index}
										<div class="flex gap-1">
											<span class="underline">
												{artist}
											</span>
											<button
												class="p-0 m-0 border-0 bg-transparent text-inherit font-extrabold"
												type="button"
												onclick={() =>
													editAlbumArtistsList.splice(
														index,
														1,
													)}
											>
												x
											</button>
										</div>
									{/each}
								</div>
								<input
									bind:value={editAlbumArtists}
									onkeydown={(e) => {
										const albumArtist =
											editAlbumArtists.trim();
										if (
											e.key === "Enter" &&
											albumArtist !== ""
										) {
											e.preventDefault();
											editAlbumArtistsList = [
												...editAlbumArtistsList,
												albumArtist,
											];
											editAlbumArtists = "";
										}
									}}
									type="text"
									placeholder="Album Artists"
								/>
							</label>

							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								<div class="flex gap-2 flex-wrap">
									<span>Artists: </span>
									{#each editArtistsList as artist, index}
										<div class="flex gap-1">
											<span class="underline">
												{artist}
											</span>
											<button
												class="p-0 m-0 border-0 bg-transparent text-inherit font-extrabold"
												type="button"
												onclick={() =>
													editArtistsList.splice(
														index,
														1,
													)}
											>
												x
											</button>
										</div>
									{/each}
								</div>
								<input
									bind:value={editArtists}
									onkeydown={(e) => {
										const artist = editArtists.trim();
										if (
											e.key === "Enter" &&
											artist !== ""
										) {
											e.preventDefault();
											editArtistsList = [
												...editArtistsList,
												artist,
											];
											editArtists = "";
										}
									}}
									type="text"
									placeholder="Artists"
								/>
							</label>

							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								Track Number
								<input
									type="number"
									name="trackNumber"
									value={editItem.trackNumber}
									placeholder={editItem.trackNumber.toString()}
									min="1"
									step="1"
								/>
							</label>
							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								Volume Number
								<input
									type="number"
									name="volumeNumber"
									value={editItem.volumeNumber}
									placeholder={editItem.volumeNumber.toString()}
									min="1"
									step="1"
								/>
							</label>
							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								ISRC
								<input
									type="text"
									name="isrc"
									value={editItem.isrc}
									placeholder={editItem.isrc}
								/>
							</label>
							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								Release Date
								<input
									type="date"
									name="releaseDate"
									value={editItem.releaseDate}
								/>
							</label>
							<label
								class="grid grid-cols-[1fr_auto] items-center"
							>
								Explicit
								<input
									type="checkbox"
									name="explicit"
									checked={editItem.explicit}
								/>
							</label>

							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								Album Gain
								<input
									type="number"
									name="albumGain"
									step="any"
									value={editItem.albumGain}
									placeholder={editItem.albumGain.toString()}
								/>
							</label>
							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								Album Peak
								<input
									type="number"
									name="albumPeak"
									step="any"
									value={editItem.albumPeak}
									placeholder={editItem.albumPeak.toString()}
								/>
							</label>
							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								Track Gain
								<input
									type="number"
									name="trackGain"
									step="any"
									value={editItem.trackGain}
									placeholder={editItem.trackGain.toString()}
								/>
							</label>
							<label
								class="grid grid-cols-[1fr_auto] gap-2 items-center"
							>
								Track Peak
								<input
									type="number"
									name="trackPeak"
									step="any"
									value={editItem.trackPeak}
									placeholder={editItem.trackPeak.toString()}
								/>
							</label>
						</form>
						<div
							class="grid grid-cols-2 items-center justify-center gap-2 pt-1"
						>
							<AlertDialog.Cancel
								type="button"
								class="hover-full"
							>
								Cancel
							</AlertDialog.Cancel>
							<AlertDialog.Action
								type="submit"
								form="edit-form"
								class="hover-full"
							>
								Edit
							</AlertDialog.Action>
						</div>
					</div>
				</AlertDialog.Content>
			</AlertDialog.Portal>
		</AlertDialog.Root>
	{/if}
	{#if deletedItem}
		<AlertDialog.Root bind:open={deleteDialog}>
			<AlertDialog.Portal>
				<AlertDialog.Overlay class="fixed inset-0 z-50 bg-black/80" />
				<AlertDialog.Content
					class="rounded-card-lg bg-bg shadow-popover outline-hidden fixed left-[50%] top-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 border p-7 sm:max-w-lg md:w-full "
				>
					<AlertDialog.Title class="text-lg font-semibold">
						Delete
						<span class="font-extrabold">
							{deletedItem.title}
						</span>
						<span class="italic font-extrabold">
							{deletedItem.album}
						</span>
						<span class="italic">
							{deletedItem.artists}
						</span>
					</AlertDialog.Title>
					<AlertDialog.Description
						class="text-foreground-alt text-sm"
					>
						Are you sure you want to delete this song?
					</AlertDialog.Description>
					<div class="flex w-full items-center justify-center gap-2">
						<AlertDialog.Cancel
							class="h-input rounded-input bg-muted shadow-mini hover:bg-dark-10 focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex w-full items-center justify-center text-[15px] font-medium transition-all focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
						>
							Cancel
						</AlertDialog.Cancel>
						<AlertDialog.Action
							class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex w-full items-center justify-center text-[15px] font-semibold transition-all focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
							onclick={async () => {
								if (deletedItem) {
									error = await deleteSong(deletedItem.id);
									if (!error) {
										error = await loadLibrary(
											search,
											limit,
											offset,
										);
									}
								}
								deletedItem = null;
							}}
						>
							Confirm
						</AlertDialog.Action>
					</div>
				</AlertDialog.Content>
			</AlertDialog.Portal>
		</AlertDialog.Root>
	{/if}
{/if}
