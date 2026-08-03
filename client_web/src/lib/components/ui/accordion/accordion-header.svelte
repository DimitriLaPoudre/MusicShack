<script lang="ts">
	import { Accordion as AccordionPrimitive } from "bits-ui";

	let { ...restProps }: AccordionPrimitive.HeaderProps = $props();
</script>

<AccordionPrimitive.Header {...restProps} />
