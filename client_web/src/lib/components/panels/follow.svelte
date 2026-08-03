<script lang="ts">
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Error } from "$lib/components/ui/error/index.js";
	import { ApiError, deleteFollow, listFollows } from "$lib/api";
	import { Heart, HeartOff } from "@lucide/svelte";
	import { onMount } from "svelte";
	import { followList, openPanel } from "$lib/stores";

	let {
		anchor = undefined,
	}: {
		anchor?: HTMLElement | null;
	} = $props();

	let error = $state<string | null>(null);

	async function load() {
		try {
			followList.set(await listFollows());
			error = null;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to load follows";
		}
	}

	async function unfollow(id: string) {
		try {
			await deleteFollow(id);
			error = null;
		} catch (e) {
			error = e instanceof ApiError ? e.message : "Failed to remove follow";
		}
		await load();
	}

	onMount(() => {
		if ($followList === null) load();
	});
</script>

<Popover.Root
	open={$openPanel === "follow"}
	onOpenChange={(v) => openPanel.set(v ? "follow" : null)}
>
	<Popover.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="hover-full" size="icon-panel">
				<Heart class="fill-current" />
			</Button>
		{/snippet}
	</Popover.Trigger>
	<Popover.Content customAnchor={anchor} class="w-full">
		<h1 class="text-xs font-semibold tracking-widest uppercase">
			Followed Artists
		</h1>
		<Error text={error} />
		{#if $followList === null}
			<p class="py-2 text-center text-xs text-muted-foreground">
				Loading...
			</p>
		{:else if $followList.length === 0}
			<p class="py-2 text-center text-xs text-muted-foreground">
				No followed artists
			</p>
		{:else}
			<div class="flex flex-col gap-1">
				{#each $followList as item}
					<div class="grid grid-cols-[1fr_auto] gap-2">
						<a
							href="/artist/{item.provider}/{item.artist_id}"
							class="grid grid-cols-[auto_1fr] items-stretch gap-2 hover:shadow-[inset_0_0_0_1px_var(--foreground)]"
						>
							{#if item.artist_picture_url}
								<img
									src={item.artist_picture_url}
									alt={item.artist_name}
									class="size-14.5 self-center object-cover"
								/>
							{:else}
								<div class="size-14.5 self-center bg-muted"></div>
							{/if}
							<p class="italic flex items-center pl-3">
								{item.artist_name}
							</p>
						</a>
						<Button
							variant="hover-full"
							size="icon"
							class="group/remove relative self-center"
							onclick={() => unfollow(item.id)}
						>
							<span class="block group-hover/remove:hidden">
								<Heart class="fill-current" />
							</span>
							<span class="hidden group-hover/remove:block">
								<HeartOff />
							</span>
						</Button>
					</div>
				{/each}
			</div>
		{/if}
	</Popover.Content>
</Popover.Root>
