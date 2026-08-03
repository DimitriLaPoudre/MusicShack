<script lang="ts">
	import { Accordion as AccordionPrimitive } from "bits-ui";

	let { ...restProps }: AccordionPrimitive.RootProps = $props();
</script>

<AccordionPrimitive.Root {...restProps} />
