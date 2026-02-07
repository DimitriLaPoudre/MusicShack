<script lang="ts">
	import { afterNavigate } from "$app/navigation";
	import Explicit from "$lib/components/explicit.svelte";
	import {
		deleteSong,
		loadLibrary,
		syncLibrary,
	} from "$lib/functions/library";
	import { ChevronLeft, ChevronRight, Trash } from "lucide-svelte";
	import { onMount } from "svelte";
	import { Pagination } from "bits-ui";
	import { libraryPage } from "$lib/stores/panel/library";

	let error = $state<null | string>(null);

	let currentPage = $state(1);
	const limit = 10;
	let offset = $derived((currentPage - 1) * limit);

	onMount(async () => {
		await syncLibrary();
		error = await loadLibrary(limit, offset);
	});

	afterNavigate(async () => {
		await syncLibrary();
		error = await loadLibrary(limit, offset);
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
	<!-- page top -->
	<div class="flex flex-col items-center justify-center">
		<Pagination.Root
			bind:page={currentPage}
			count={$libraryPage.total}
			perPage={limit}
			onPageChange={async () => {
				error = await loadLibrary(limit, offset);
			}}
		>
			{#snippet children({ pages })}
				<div class="my-8 flex items-center">
					<Pagination.PrevButton
						class="mr-[26px] inline-flex size-10 items-center justify-center active:text-bg active:bg-fg disabled:cursor-default disabled:text-bg hover:disabled:bg-transparent"
					>
						<ChevronLeft />
					</Pagination.PrevButton>
					<div class="flex items-center gap-2.5">
						{#each pages as page (page.key)}
							{#if page.type === "ellipsis"}
								<div
									class="text-foreground-alt select-none text-[15px] font-medium"
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
						class="mr-[26px] inline-flex size-10 items-center justify-center active:text-bg active:bg-fg disabled:cursor-default disabled:text-bg hover:disabled:bg-transparent"
					>
						<ChevronRight />
					</Pagination.NextButton>
				</div>
			{/snippet}
		</Pagination.Root>
	</div>
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
					class="hover-full w-full p-3 shadow-[inset_0_-1px_0_var(--fg),inset_1px_0_0_var(--fg),inset_-1px_0_0_var(--fg)]"
					onclick={async () => {
						error = await deleteSong(item.id);
						if (!error) {
							error = await loadLibrary(limit, offset);
						}
					}}
				>
					<Trash />
				</button>
			</div>
		{/each}
	</div>
{/if}
