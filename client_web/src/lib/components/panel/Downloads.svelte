<script lang="ts">
	import { goto } from "$app/navigation";
	import { apiFetch } from "$lib/functions/apiFetch";
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

	let list = $state<null | any>(null);
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

	async function loadDownloads() {
		let list;
		let error;
		try {
			const res = await apiFetch(`/users/downloads`);
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to fetch downloads");
			}
			list = body.tasks
				.slice()
				.sort((a: any, b: any) => Number(b.Id) - Number(a.Id));
			error = null;
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Failed to reload download queue";
		}
		return { list, error };
	}

	async function retryDownload(id: string) {
		try {
			const res = await apiFetch(`/users/downloads/retry/${id}`, "POST");
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to retry download");
			}
		} catch (e) {
			error = e instanceof Error ? e.message : "Failed to retry download";
		}
	}

	async function cancelDownload(id: string) {
		try {
			const res = await apiFetch(`/users/downloads/cancel/${id}`, "POST");
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to cancel download");
			}
		} catch (e) {
			error =
				e instanceof Error ? e.message : "Failed to cancel download";
		}
	}

	async function deleteDownload(id: string) {
		try {
			const res = await apiFetch(`/users/downloads/${id}`, "DELETE");
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to delete download");
			}
		} catch (e) {
			error =
				e instanceof Error ? e.message : "Failed to delete download";
		}
	}
</script>

<div class="body">
	<h1>Downloads Queue</h1>
	{#if error}
		<p class="error">{error}</p>
	{/if}
	{#if !list}
		<p class="loading">Loading...</p>
	{:else}
		<div class="items">
			{#each list as download, index}
				<div class="item">
					{#if download.Data.Album.CoverUrl !== ""}
						<img
							src={download.Data.Album.CoverUrl}
							alt={download.Data.Album.CoverUrl}
						/>
					{:else}
						<Disc />
					{/if}
					<button
						class="item-data"
						onclick={(e) => {
							if (
								e.target instanceof Element &&
								e.target.closest("a")
							)
								return;
							goto(`/song/${download.Api}/${download.Data.Id}`);
						}}
					>
						<p>{download.Data.Title}</p>
						<a
							href="/artist/{download.Api}/{download.Data.Artist
								.Id}">{download.Data.Artist.Name}</a
						>
					</button>
					<div class="item-btn">
						{#if download.Status === "done"}
							<button>
								<CircleCheck />
							</button>
						{:else if download.Status === "pending" || download.Status === "running"}
							<button
								onmouseenter={() => (buttonHover = index)}
								onmouseleave={() => (buttonHover = null)}
								onclick={async () => {
									cancelDownload(download.Id);
									({ list, error } = await loadDownloads());
								}}
							>
								{#if buttonHover != null && buttonHover === index}
									<CircleX />
								{:else if download.Status === "pending"}
									<CircleDashed />
								{:else}
									<LoaderCircleIcon />
								{/if}
							</button>
						{:else if download.Status === "failed" || download.Status === "cancel"}
							<button
								onmouseenter={() => (buttonHover = index)}
								onmouseleave={() => (buttonHover = null)}
								onclick={async () => {
									retryDownload(download.Id);
									({ list, error } = await loadDownloads());
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
								deleteDownload(download.Id);
								({ list, error } = await loadDownloads());
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
	.loading {
		text-align: center;
	}
	.error {
		text-align: center;
		background-color: var(--err);
		padding: 0.5rem;
		margin: 0;
	}
	.items {
		display: flex;
		flex-direction: column;
		gap: 4px;
		.item {
			display: grid;
			grid-template-columns: auto 1fr auto;
			gap: 8px;
			align-items: stretch;
			container-type: inline-size;

			img {
				margin: auto;
				width: 58px;
				height: 58px;
				aspect-ratio: 1/1;
			}

			.item-data {
				display: grid;
				grid-template-columns: 1fr 1fr;
				align-items: center;
				justify-items: left;
				border: none;
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
				button {
					aspect-ratio: 1/1;
				}
			}

			@container (max-width: 420px) {
				.item-data {
					grid-template-columns: 1fr;
				}

				.item-btn {
					grid-template-columns: 1fr;
				}
			}
		}
	}
</style>
