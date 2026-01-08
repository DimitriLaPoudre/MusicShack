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
					<a
						class="data"
						href="/artist/{item.provider}/{item.artistId}"
					>
						<img
							class="picture"
							src={item.artistPictureUrl}
							alt={item.artistName}
						/>
						<p class="artist">{item.artistName}</p>
					</a>
					<button
						class="hover-full"
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
	h1 {
		font-weight: bolder;
	}
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
		display: flex;
		flex-direction: column;
		gap: 0.25rem;

		.artist {
			display: grid;
			grid-template-columns: 1fr auto;
			gap: 0.75rem;

			.data {
				display: grid;
				grid-template-columns: auto 1fr;
				align-items: stretch;
				gap: 0.75rem;

				.picture {
					width: 58px;
					height: 58px;
					align-self: center;
				}

				.artist {
					font-style: italic;
					padding-left: 0.75rem;
					display: flex;
					align-items: center;
				}
			}

			@media not all and (pointer: coarse) and (hover: none) {
				.data:hover {
					.artist {
						color: inherit;
						background-color: inherit;
						outline: 1px solid var(--fg);
						outline-offset: -1px;
					}
				}
			}

			button {
				/* align-self: center; */
				/* aspect-ratio: 1/1; */
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
