<script lang="ts">
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Search } from "@lucide/svelte";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Button } from "$lib/components/ui/button/index.js";

	let {
		anchor = undefined,
	}: {
		anchor?: HTMLElement | null;
	} = $props();
</script>

<Popover.Root>
	<Popover.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="hover-full" size="icon-panel">
				<Search />
			</Button>
		{/snippet}
	</Popover.Trigger>
	<Popover.Content customAnchor={anchor} class="w-full">
		<form>
			<Input
				id="search"
				type="search"
				placeholder="Search some music..."
				required
			/>
		</form>
	</Popover.Content>
</Popover.Root>
