<script lang="ts">
	import { Button } from "$lib/components/ui/button/index.js";
	import { Loader2 } from "@lucide/svelte";

	let {
		hasMore,
		loading,
		onclick,
	}: {
		hasMore: boolean;
		loading: boolean;
		onclick: () => void;
	} = $props();
</script>

{#if hasMore}
	<Button
		variant="hover-full"
		size="default"
		class="w-full"
		disabled={loading}
		onclick={onclick}
	>
		{#if loading}
			<Loader2 class="animate-spin" />
			Loading...
		{:else}
			Load More
		{/if}
	</Button>
{/if}
