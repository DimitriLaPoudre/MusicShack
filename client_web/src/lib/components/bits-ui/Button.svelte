<script lang="ts">
	import { Button } from "bits-ui";

	export let children;
</script>

<Button.Root class="border" {...$$restProps}>
	{@render children?.()}
</Button.Root>
