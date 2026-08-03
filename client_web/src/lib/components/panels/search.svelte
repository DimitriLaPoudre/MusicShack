<script lang="ts">
	import { goto } from "$app/navigation";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Search } from "@lucide/svelte";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { openPanel, searchInput } from "$lib/stores";

	let {
		anchor = undefined,
	}: {
		anchor?: HTMLElement | null;
	} = $props();

	async function searchFunction(e: SubmitEvent) {
		e.preventDefault();
		const encoded = encodeURIComponent($searchInput);
		$searchInput = "";
		openPanel.set(null);
		await goto(`/search?q=${encoded}`);
	}
</script>

<Popover.Root
	open={$openPanel === "search"}
	onOpenChange={(v) => openPanel.set(v ? "search" : null)}
>
	<Popover.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="hover-full" size="icon-panel">
				<Search />
			</Button>
		{/snippet}
	</Popover.Trigger>
	<Popover.Content customAnchor={anchor} class="w-full">
		<form onsubmit={searchFunction}>
			<Input
				id="search"
				type="search"
				placeholder="Search some music..."
				required
				autofocus
				bind:value={$searchInput}
			/>
		</form>
	</Popover.Content>
</Popover.Root>
