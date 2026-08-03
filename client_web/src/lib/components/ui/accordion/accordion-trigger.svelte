<script lang="ts">
	import { Accordion as AccordionPrimitive } from "bits-ui";

	let { ...restProps }: AccordionPrimitive.TriggerProps = $props();
</script>

<AccordionPrimitive.Trigger
	{...restProps}
	class="group flex w-full items-center justify-center gap-2 font-bold"
/>
