<script lang="ts">
	import { loadFollows, removeFollow } from "$lib/functions/follow";
	import type { FollowsResponse } from "$lib/types/response";
	import { HeartIcon, HeartOff } from "lucide-svelte";
	import { onMount } from "svelte";

	let list = $state<null | FollowsResponse>(null);
	let error = $state<null | string>(null);

	onMount(() => {
		async function firstInterval() {
			({ list, error } = await loadFollows());
		}
		firstInterval();
		const interval = setInterval(async () => {
			({ list, error } = await loadFollows());
		}, 500);
		return () => clearInterval(interval);
	});
</script>

<div class="body">
	<h1>Followed Artist</h1>
	{#if error}
		<p class="error">{error}</p>
	{/if}
	{#if !list}
		<p class="loading">Loading...</p>
	{:else}
		<div class="list">
			{#each list as item}
				<div class="artist">
					<a class="data" href="/artist/{item.api}/{item.artist.id}">
						<img
							src={item.artist.pictureUrl}
							alt={item.artist.name}
						/>
						<p>{item.artist.name}</p>
					</a>
					<button
						onclick={async () => {
							await removeFollow(item.id);
						}}
					>
						<div class="nothover">
							<HeartIcon />
						</div>
						<div class="hover">
							<HeartOff />
						</div>
					</button>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.error {
		text-align: center;
		background-color: var(--err);
		padding: 0.5rem;
		margin: 0;
	}
	.loading {
		text-align: center;
	}
	.list {
		.artist {
			display: grid;
			grid-template-columns: 1fr auto;
			gap: 8px;

			.data {
				display: grid;
				grid-template-columns: auto 1fr;
				align-items: stretch;
				gap: 8px;

				img {
					width: 58px;
					height: auto;
					aspect-ratio: 1/1;
				}

				p {
					padding-left: 8px;
					display: flex;
					align-items: center;
				}
			}
			.data:hover {
				p {
					outline: 1px solid #ffffff;
					outline-offset: -1px;
				}
			}

			button {
				aspect-ratio: 1/1;
				.nothover {
					display: block;
				}
				.hover {
					display: none;
				}
			}
			button:hover {
				.nothover {
					display: none;
				}
				.hover {
					display: block;
				}
			}
		}
	}
</style>
