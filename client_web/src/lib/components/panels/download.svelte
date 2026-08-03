<script lang="ts">
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Error } from "$lib/components/ui/error/index.js";
	import {
		ApiError,
		cancelDownload,
		listDownloads,
		removeDoneDownloads,
		removeDownload,
		retryAllDownloads,
		retryDownload,
	} from "$lib/api";
	import {
		CircleAlert,
		CircleCheck,
		CircleDashed,
		CircleX,
		Disc,
		Loader,
		RotateCcw,
		Trash,
	} from "@lucide/svelte";
	import { onMount } from "svelte";
	import type { DownloadTask } from "$lib/types";
	import { downloadList, openPanel } from "$lib/stores";

	let {
		anchor = undefined,
	}: {
		anchor?: HTMLElement | null;
	} = $props();

	let error = $state<string | null>(null);

	function message(e: unknown, fallback: string) {
		return e instanceof ApiError ? e.message : fallback;
	}

	async function load() {
		try {
			downloadList.set(await listDownloads());
			error = null;
		} catch (e) {
			error = message(e, "Failed to load downloads");
		}
	}

	async function run(fn: () => Promise<void>, fallback: string) {
		try {
			await fn();
			error = null;
		} catch (e) {
			error = message(e, fallback);
		}
		await load();
	}

	onMount(() => {
		if ($downloadList === null) load();
	});

	const done = (t: DownloadTask) => t.status === "done";
	const failed = (t: DownloadTask) => t.status === "failed" || t.status === "cancel";
</script>

<Popover.Root
	open={$openPanel === "download"}
	onOpenChange={(v) => openPanel.set(v ? "download" : null)}
>
	<Popover.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="hover-full" size="icon-panel">
				<Disc />
			</Button>
		{/snippet}
	</Popover.Trigger>
	<Popover.Content customAnchor={anchor} class="w-full">
		<h1 class="text-xs font-semibold tracking-widest uppercase">
			Downloads Queue
		</h1>
		<Error text={error} />
		{#if $downloadList === null}
			<p class="py-2 text-center text-xs text-muted-foreground">
				Loading...
			</p>
		{:else}
			{#if $downloadList.some((t) => done(t) || failed(t))}
				<div class="flex gap-1">
					{#if $downloadList.some(failed)}
						<Button
							variant="hover-full"
							size="icon"
							class="flex-1"
							onclick={() => run(retryAllDownloads, "Failed to retry downloads")}
						>
							<RotateCcw />
						</Button>
					{/if}
					{#if $downloadList.some(done)}
						<Button
							variant="hover-full"
							size="icon"
							class="flex-1"
							onclick={() => run(removeDoneDownloads, "Failed to remove done downloads")}
						>
							<CircleCheck class="fill-current" />
						</Button>
					{/if}
				</div>
			{/if}
			{#if $downloadList.length === 0}
				<p class="py-2 text-center text-xs text-muted-foreground">
					No downloads
				</p>
			{:else}
				<div class="flex flex-col gap-1">
					{#each $downloadList as task}
						<div class="grid grid-cols-[auto_1fr_auto] gap-2">
							{#if task.cover_url}
								<img
									src={task.cover_url}
									alt={task.title}
									class="size-14.5 self-center object-cover"
								/>
							{:else}
								<div
									class="size-14.5 grid place-items-center self-center bg-muted"
								>
									<Disc class="size-6" />
								</div>
							{/if}
							<div class="flex min-w-0 flex-col justify-center pl-2">
								<p class="truncate font-semibold">{task.title}</p>
								{#if task.artists.length > 0}
									<p class="truncate italic text-muted-foreground">
										{task.artists.join(", ")}
									</p>
								{/if}
							</div>
							<div class="grid grid-cols-1 gap-1">
								{#if task.status === "done"}
									<div class="flex items-center justify-center p-2">
										<CircleCheck class="fill-current" />
									</div>
								{:else if task.status === "pending" || task.status === "running"}
									<Button
										variant="hover-full"
										size="icon"
										class="group/cancel relative"
										onclick={() =>
											run(() => cancelDownload(task.id), "Failed to cancel download")
										}
									>
										<span class="block group-hover/cancel:hidden">
											{#if task.status === "pending"}
												<CircleDashed />
											{:else}
												<Loader class="animate-spin" />
											{/if}
										</span>
										<span class="hidden group-hover/cancel:block">
											<CircleX />
										</span>
									</Button>
								{:else if failed(task)}
									<Button
										variant="hover-full"
										size="icon"
										class="group/retry relative"
										onclick={() =>
											run(() => retryDownload(task.id), "Failed to retry download")
										}
									>
										<span class="block group-hover/retry:hidden">
											<CircleAlert />
										</span>
										<span class="hidden group-hover/retry:block">
											<RotateCcw />
										</span>
									</Button>
								{/if}
								<Button
									variant="hover-full"
									size="icon"
									onclick={() =>
										run(() => removeDownload(task.id), "Failed to remove download")
									}
								>
									<Trash />
								</Button>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		{/if}
	</Popover.Content>
</Popover.Root>
