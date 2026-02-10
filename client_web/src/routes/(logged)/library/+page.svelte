<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import Explicit from "$lib/components/explicit.svelte";
	import {
		deleteSong,
		loadLibrary,
		syncLibrary,
	} from "$lib/functions/library";
	import {
		ChevronLeft,
		ChevronRight,
		SearchIcon,
		Trash,
	} from "lucide-svelte";
	import { onMount } from "svelte";
	import { Pagination, AlertDialog } from "bits-ui";
	import { libraryPage } from "$lib/stores/panel/library";
	import { page } from "$app/state";
	import type { ErrorResponse, ResponseSong } from "$lib/types/response";

	let error = $state<null | string>(null);

	let uploadDialog = $state(false);
	let errorUploadDialog = $state<null | string>(null);

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
		await syncLibrary();
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
	<p class="mt-3"></p>
	<div class="flex items-center justify-center">
		<button class="">
			<SearchIcon />
		</button>

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
				onclick={() => (errorUploadDialog = null)}
			>
				Upload
			</AlertDialog.Trigger>
		</div>
		<AlertDialog.Portal>
			<AlertDialog.Overlay class="fixed inset-0 z-50 bg-black/80" />
			<AlertDialog.Content
				class="rounded-card-lg bg-bg shadow-popover outline-hidden fixed left-[50%] top-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 border p-7 sm:max-w-lg md:w-full "
			>
				<form
					onsubmit={async (e) => {
						try {
							e.preventDefault();

							const form = e.currentTarget;
							const fd = new FormData(form);

							const res = await fetch("/api/library", {
								method: "POST",
								credentials: "include",
								body: fd,
							});
							if (res.status === 401) {
								await goto("/login");
							}
							if (!res.ok) {
								throw new Error(
									((await res.json()) as ErrorResponse)
										.error || "Failed to upload song",
								);
							}
							error = await loadLibrary(search, limit, offset);
							uploadDialog = false;
						} catch (e) {
							errorUploadDialog =
								e instanceof Error
									? e.message
									: "Failed to upload song";
						}
					}}
				>
					<div class="flex flex-col gap-4">
						<AlertDialog.Title class="text-lg font-semibold">
							Upload a song
						</AlertDialog.Title>
						<AlertDialog.Description
							class="flex flex-col text-foreground-alt text-sm gap-4"
						>
							{#if errorUploadDialog}
								<p>{errorUploadDialog}</p>
							{/if}
							<input
								required
								type="file"
								name="file"
								accept="audio/*"
							/>
							<input
								required
								type="text"
								name="title"
								placeholder="Title"
							/>
							<input
								required
								type="text"
								name="album"
								placeholder="Album"
							/>
							<input
								required
								type="text"
								name="albumArtists"
								placeholder="Album Artists (eg: thaHomey, Skuna)"
							/>

							<input
								required
								type="text"
								name="artists"
								placeholder="Artists (eg: thaHomey, LaFève)"
							/>
							<button>Optionnal</button>
							<input
								type="number"
								name="trackNumber"
								placeholder="Track Number"
								min="1"
								step="1"
							/>
							<input
								type="number"
								name="volumeNumber"
								placeholder="Volume Number"
								min="1"
								step="1"
							/>
							<input
								type="text"
								name="isrc"
								placeholder="ISRC (eg: FR5R00909899)"
								minlength="12"
								maxlength="12"
								pattern="^[A-Z]{2}[A-Z0-9]{3}\d{2}\d{5}$"
							/>
							<input type="date" name="releaseDate" />
							<label>
								Explicit
								<input type="checkbox" name="explicit" />
							</label>
						</AlertDialog.Description>
					</div>
					<div class="flex w-full items-center justify-center gap-2">
						<AlertDialog.Cancel
							class="h-input rounded-input bg-muted shadow-mini hover:bg-dark-10 focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex w-full items-center justify-center text-[15px] font-medium transition-all focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
						>
							Cancel
						</AlertDialog.Cancel>
						<AlertDialog.Action
							class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex w-full items-center justify-center text-[15px] font-semibold transition-all focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
						>
							Save
						</AlertDialog.Action>
					</div>
				</form>
			</AlertDialog.Content>
		</AlertDialog.Portal>
	</AlertDialog.Root>
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
			<div class="w-[200px] h-auto">
				<button
					class="hover-full flex flex-col items-center w-[200px] h-auto overflow-hidden gap-3 shadow-[inset_0_1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
				>
					<div class="w-[160px] h-[160px]">
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
				<button
					class="hover-full w-full p-4 shadow-[inset_0_-1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
					onclick={() => (deletedItem = item)}
				>
					<Trash />
				</button>
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
{/if}
{#if deletedItem}
	<AlertDialog.Root bind:open={deleteDialog}>
		<AlertDialog.Portal>
			<AlertDialog.Overlay class="fixed inset-0 z-50 bg-black/80" />
			<AlertDialog.Content
				class="rounded-card-lg bg-bg shadow-popover outline-hidden fixed left-[50%] top-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 border p-7 sm:max-w-lg md:w-full "
			>
				<div class="flex flex-col gap-4 pb-6">
					<AlertDialog.Title class="text-lg font-semibold">
						Delete {deletedItem.title}
						{deletedItem.album}
						{deletedItem.artists}
					</AlertDialog.Title>
					<AlertDialog.Description
						class="text-foreground-alt text-sm"
					>
						Are you sure you want to delete this song?
					</AlertDialog.Description>
				</div>
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
