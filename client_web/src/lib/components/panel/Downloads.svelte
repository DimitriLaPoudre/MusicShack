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
	import type { DownloadListResponse } from "$lib/types/response";
	import {
		CircleAlert,
		CircleCheck,
		CircleDashed,
		CircleX,
		Disc,
		LoaderCircleIcon,
		RotateCcw,
		Trash,
	} from "lucide-svelte";
	import { onMount } from "svelte";

	let list = $state<null | DownloadListResponse>(null);
	let error = $state<null | string>(null);
	let buttonHover = $state<null | number>(null);

	onMount(() => {
		async function firstInterval() {
			({ list, error } = await loadDownloads());
		}
		firstInterval();

		const interval = setInterval(async () => {
			({ list, error } = await loadDownloads());
		}, 500);
		return () => clearInterval(interval);
	});
</script>

<div class="body">
	<h1>Downloads Queue</h1>
	{#if error}
		<p class="error">{error}</p>
	{/if}
	{#if !list}
		<p class="loading">Loading...</p>
	{:else}
		{#if list.some((task) => task.status === "failed" || task.status === "cancel" || task.status === "done")}
			<div class="all">
				{#if list.some((task) => task.status === "failed" || task.status === "cancel")}
					<button
						onclick={async () => {
							error = await retryAllDownload();
							if (!error) {
								({ list, error } = await loadDownloads());
							}
						}}
					>
						<RotateCcw />
					</button>
				{/if}
				{#if list.some((task) => task.status === "done")}
					<button
						onclick={async () => {
							error = await doneDownload();
							if (!error) {
								({ list, error } = await loadDownloads());
							}
						}}
					>
						<CircleCheck />
					</button>
				{/if}
			</div>
		{/if}
		<div class="items">
			{#each list as download, index}
				<div class="item">
					{#if download.data.id === ""}
						<div class="img">
							<Disc />
						</div>
						<button class="item-data">
							<p>Unreleased</p>
							<p>Unknown</p>
						</button>
					{:else}
						<div class="img">
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
							class="item-data"
							onclick={(e) => {
								if (
									e.target instanceof Element &&
									e.target.closest("a")
								)
									return;
								goto(
									`/song/${download.api}/${download.data.id}`,
								);
							}}
						>
							<p class="title">{download.data.title}</p>
							<a
								class="artist"
								href="/artist/{download.api}/{download.data
									.artists[0].id}"
								>{download.data.artists[0].name}</a
							>
						</button>
					{/if}
					<div class="item-btn">
						{#if download.status === "done"}
							<button>
								<CircleCheck />
							</button>
						{:else if download.status === "pending" || download.status === "running"}
							<button
								onmouseenter={() => (buttonHover = index)}
								onmouseleave={() => (buttonHover = null)}
								onclick={async () => {
									error = await cancelDownload(download.id);
									if (!error) {
										({ list, error } =
											await loadDownloads());
									}
								}}
							>
								{#if buttonHover != null && buttonHover === index}
									<CircleX />
								{:else if download.status === "pending"}
									<CircleDashed />
								{:else}
									<LoaderCircleIcon />
								{/if}
							</button>
						{:else if download.status === "failed" || download.status === "cancel"}
							<button
								onmouseenter={() => (buttonHover = index)}
								onmouseleave={() => (buttonHover = null)}
								onclick={async () => {
									error = await retryDownload(download.id);
									if (!error) {
										({ list, error } =
											await loadDownloads());
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
							onclick={async () => {
								error = await deleteDownload(download.id);
								if (!error) {
									({ list, error } = await loadDownloads());
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

<style>
	h1 {
		font-weight: bolder;
	}
	.loading {
		text-align: center;
	}
	.error {
		text-align: center;
		background-color: var(--err);
		padding: 0.5rem;
	}

	.body {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;

		.all {
			display: flex;
			flex-direction: row;
			justify-content: center;
			align-items: center;
			gap: 0.5rem;

			button {
				width: 100%;
				padding: 0.5rem 0;
			}
		}
		.items {
			display: flex;
			flex-direction: column;
			gap: 0.25rem;
			.item {
				display: grid;
				grid-template-columns: auto 1fr auto;
				gap: 0.75rem;
				align-items: stretch;
				container-type: inline-size;

				.img {
					width: 58px;
					height: 58px;
					align-self: center;
				}

				.item-data {
					align-self: center;
					display: grid;
					grid-template-columns: 1fr 1fr;
					align-items: center;
					justify-items: center;
					border: none;
					gap: 0.5rem 0.5rem;
					height: 100%;

					.title {
						font-weight: bolder;
					}
					.artist {
						font-style: italic;
					}
				}
				.item-data:hover {
					outline: 1px solid #ffffff;
					outline-offset: -1px;
					background-color: inherit;
					color: inherit;
				}

				.item-btn {
					display: grid;
					grid-template-columns: 1fr 1fr;
					gap: 0.25rem;

					button {
						aspect-ratio: 1/1;
					}
				}

				@container (max-width: 520px) {
					.item-data {
						grid-template-columns: 1fr;
					}

					.item-btn {
						grid-template-columns: 1fr;
						button {
							aspect-ratio: auto;
						}
					}
				}
			}
		}
	}
</style>
