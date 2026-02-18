<script lang="ts">
	import { goto } from "$app/navigation";
	import {
		cancelDownload,
		deleteDownload,
		doneDownload,
		loadDownloads,
		retryAllDownload,
		retryDownload,
	} from "$lib/functions/download";
	import {
		CircleAlert,
		CircleCheck,
		CircleDashed,
		CircleX,
		Disc,
		LoaderIcon,
		RotateCcw,
		Trash,
	} from "lucide-svelte";
	import { onMount } from "svelte";
	import Explicit from "../Explicit.svelte";
	import { downloadList } from "$lib/stores/panel/download";
	import Owned from "../Owned.svelte";

	let error = $state<null | string>(null);
	let buttonHover = $state<null | number>(null);

	onMount(() => {
		async function intervalFunc() {
			error = await loadDownloads();
		}
		intervalFunc();
		const interval = setInterval(intervalFunc, 500);
		return () => {
			clearInterval(interval);
		};
	});
</script>

<div class="flex flex-col gap-3">
	<h1 class="font-extrabold">Downloads Queue</h1>
	{#if error}
		<p class="text-center bg-err p-2">{error}</p>
	{/if}
	{#if !$downloadList}
		<p class="text-center">Loading...</p>
	{:else}
		{#if $downloadList.some((task) => task.status === "failed" || task.status === "cancel" || task.status === "done")}
			<div class="flex flex-row justify-center items-center gap-2">
				{#if $downloadList.some((task) => task.status === "failed" || task.status === "cancel")}
					<button
						class="hover-full w-full py-3"
						onclick={async () => {
							error = await retryAllDownload();
							if (!error) {
								error = await loadDownloads();
							}
						}}
					>
						<RotateCcw />
					</button>
				{/if}
				{#if $downloadList.some((task) => task.status === "done")}
					<button
						class="hover-full w-full py-3"
						onclick={async () => {
							error = await doneDownload();
							if (!error) {
								error = await loadDownloads();
							}
						}}
					>
						<CircleCheck />
					</button>
				{/if}
			</div>
		{/if}
		<div class="flex flex-col gap-1">
			{#each $downloadList as download, index}
				<div
					class="grid grid-cols-[auto_1fr_auto] gap-2 items-stretch @container"
				>
					{#if download.data.id === ""}
						<div class="w-14.5 h-14.5 self-center">
							<Disc />
						</div>
						<button
							class="hover-soft self-center grid grid-cols-2 @max-[520px]:grid-cols-1 items-center justify-items-center border-none gap-2 h-full"
						>
							<p>Unreleased</p>
							<p>Unknown</p>
						</button>
					{:else}
						<div class="w-14.5 h-14.5 self-center">
							{#if download.data.album.coverUrl !== ""}
								<img
									src={download.data.album.coverUrl}
									alt={download.data.album.coverUrl}
								/>
							{:else}
								<Disc />
							{/if}
						</div>
						<button
							class="hover-soft self-center grid grid-cols-2 @max-[520px]:grid-cols-1 items-center justify-items-center border-none gap-2 h-full"
							onclick={(e) => {
								if (
									e.target instanceof Element &&
									e.target.closest("a")
								)
									return;
								goto(
									`/song/${download.provider}/${download.data.id}`,
								);
							}}
						>
							<p
								class="flex flex-row items-center justify-center gap-2 font-extrabold"
							>
								{#if download.data.downloaded}
									<Owned />
								{/if}
								{download.data.title}
								{#if download.data.explicit}
									<Explicit />
								{/if}
							</p>
							<a
								class="italic"
								href="/artist/{download.provider}/{download.data
									.artists[0].id}"
								>{download.data.artists[0].name}</a
							>
						</button>
					{/if}
					<div
						class="grid grid-cols-2 @max-[520px]:grid-cols-1 gap-1"
					>
						{#if download.status === "done"}
							<div class="p-4 flex items-center justify-center">
								<CircleCheck />
							</div>
						{:else if download.status === "pending" || download.status === "running"}
							<button
								class="hover-full"
								onmouseenter={() => (buttonHover = index)}
								onmouseleave={() => (buttonHover = null)}
								onclick={async () => {
									error = await cancelDownload(download.id);
									if (!error) {
										error = await loadDownloads();
									}
								}}
							>
								{#if buttonHover != null && buttonHover === index}
									<CircleX />
								{:else if download.status === "pending"}
									<CircleDashed />
								{:else}
									<LoaderIcon />
								{/if}
							</button>
						{:else if download.status === "failed" || download.status === "cancel"}
							<button
								class="hover-full"
								onmouseenter={() => (buttonHover = index)}
								onmouseleave={() => (buttonHover = null)}
								onclick={async () => {
									error = await retryDownload(download.id);
									if (!error) {
										error = await loadDownloads();
									}
								}}
							>
								{#if buttonHover != null && buttonHover === index}
									<RotateCcw />
								{:else}
									<CircleAlert />
								{/if}
							</button>
						{/if}
						<button
							class="hover-full"
							onclick={async () => {
								error = await deleteDownload(download.id);
								if (!error) {
									error = await loadDownloads();
								}
							}}
						>
							<Trash />
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
