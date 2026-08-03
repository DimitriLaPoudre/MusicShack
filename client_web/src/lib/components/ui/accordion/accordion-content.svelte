<script lang="ts">
	import { Accordion as AccordionPrimitive } from "bits-ui";

	let {
		children,
		...restProps
	}: AccordionPrimitive.ContentProps = $props();
</script>

<AccordionPrimitive.Content {...restProps} class="overflow-hidden">
	<div class="flex flex-col items-center gap-4 pt-4">
		{@render children?.()}
	</div>
</AccordionPrimitive.Content>
