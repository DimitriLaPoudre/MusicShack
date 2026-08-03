<script lang="ts">
	import { cn, type WithElementRef } from "$lib/utils.js";
	import type { HTMLAttributes } from "svelte/elements";

	let {
		ref = $bindable(null),
		class: className,
		text,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		text?: string | null;
	} = $props();
</script>

{#if text || children}
	<div
		bind:this={ref}
		data-slot="error"
		role="alert"
		class={cn(
			"rounded-none border border-destructive/30 bg-destructive/10 px-3 py-2 text-xs text-destructive",
			className
		)}
		{...restProps}
	>
		{#if text}
			{text}
		{/if}
		{@render children?.()}
	</div>
{/if}
