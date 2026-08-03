<script lang="ts">
	import { Accordion as AccordionPrimitive } from "bits-ui";

	let { ...restProps }: AccordionPrimitive.ItemProps = $props();
</script>

<AccordionPrimitive.Item class="w-full" {...restProps} />
