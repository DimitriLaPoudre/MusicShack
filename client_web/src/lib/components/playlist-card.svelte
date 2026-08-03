<script lang="ts">
	import { goto } from "$app/navigation";
	import { ApiError, downloadPlaylist, listDownloads } from "$lib/api";
	import { Button } from "$lib/components/ui/button/index.js";
	import { DiscAlbum, Download } from "@lucide/svelte";
	import { downloadList } from "$lib/stores";
	import type { PlaylistInfo } from "$lib/types";

	let {
		provider,
		playlist,
		onerror,
	}: {
		provider: string;
		playlist: PlaylistInfo;
		onerror?: (message: string | null) => void;
	} = $props();

	async function startDownload() {
		try {
			await downloadPlaylist({ provider, id: playlist.id });
			onerror?.(null);
		} catch (e) {
			onerror?.(
				e instanceof ApiError ? e.message : "Failed to download playlist",
			);
		}
		try {
			downloadList.set(await listDownloads());
		} catch {
			// keep cached list on refresh failure
		}
	}
</script>

<div class="w-50 h-auto">
	<Button
		variant="hover-full"
		class="flex w-50 h-auto flex-col items-center gap-2 py-5 overflow-hidden shadow-[inset_0_1px_0_var(--foreground),inset_1px_0_0_var(--foreground),inset_-1px_0_0_var(--foreground)]"
		onclick={() => goto(`/playlist/${provider}/${playlist.id}`)}
	>
		<div class="size-40">
			{#if playlist.cover_url}
				<img
					class="size-full object-cover"
					src={playlist.cover_url}
					alt={playlist.title}
				/>
			{:else}
				<div class="grid size-full place-items-center bg-muted">
					<DiscAlbum class="size-24" />
				</div>
			{/if}
		</div>
		<p class="break-words font-extrabold">{playlist.title}</p>
	</Button>
	<Button
		variant="hover-full"
		class="w-full p-3 shadow-[inset_0_-1px_0_var(--foreground),inset_1px_0_0_var(--foreground),inset_-1px_0_0_var(--foreground)]"
		onclick={startDownload}
	>
		<Download />
	</Button>
</div>
