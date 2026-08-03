<script lang="ts">
	import { goto } from "$app/navigation";
	import { ApiError, addFollow, deleteFollow, listFollows } from "$lib/api";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Heart, User } from "@lucide/svelte";
	import { followList } from "$lib/stores";
	import type { ArtistInfo } from "$lib/types";

	let {
		provider,
		artist,
		onerror,
	}: {
		provider: string;
		artist: ArtistInfo;
		onerror?: (message: string | null) => void;
	} = $props();

	async function refreshFollows() {
		try {
			followList.set(await listFollows());
		} catch {
			// keep cached list on refresh failure
		}
	}

	async function toggleFollow() {
		try {
			if (artist.followed) {
				await deleteFollow(artist.followed);
				artist.followed = null;
			} else {
				const follow = await addFollow({
					provider,
					artist_id: artist.id,
					featuring: false,
				});
				artist.followed = follow.id;
			}
			onerror?.(null);
		} catch (e) {
			onerror?.(e instanceof ApiError ? e.message : "Failed to toggle follow");
		}
		await refreshFollows();
	}
</script>

<div class="w-50 h-auto">
	<Button
		variant="hover-full"
		class="flex w-50 h-auto flex-col items-center gap-2 py-5 overflow-hidden shadow-[inset_0_1px_0_var(--foreground),inset_1px_0_0_var(--foreground),inset_-1px_0_0_var(--foreground)]"
		onclick={() => goto(`/artist/${provider}/${artist.id}`)}
	>
		<div class="size-40">
			{#if artist.picture_url}
				<img
					class="size-full rounded-full object-cover"
					src={artist.picture_url}
					alt={artist.name}
				/>
			{:else}
				<div class="grid size-full place-items-center rounded-full bg-muted">
					<User class="size-24" />
				</div>
			{/if}
		</div>
		<p class="break-words font-extrabold">{artist.name}</p>
	</Button>
	<Button
		variant="hover-full"
		class="w-full p-3 shadow-[inset_0_-1px_0_var(--foreground),inset_1px_0_0_var(--foreground),inset_-1px_0_0_var(--foreground)]"
		onclick={toggleFollow}
	>
		<Heart class={artist.followed ? "fill-current" : ""} />
	</Button>
</div>
