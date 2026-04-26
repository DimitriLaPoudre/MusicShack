<script lang="ts">
	import { Check, Minus } from "@lucide/svelte";
	import { Checkbox, Label } from "bits-ui";
	export let id: string;
	export let name: string = "";
	export let checked: boolean = false;
	export let disabled = false;
</script>

<div class="flex items-center gap-3">
	<Checkbox.Root
		{id}
		{name}
		bind:checked
		indeterminate={false}
		{disabled}
		aria-labelledby={`${id}-label`}
		class=""
	>
		{#snippet children({ checked, indeterminate })}
			<div class="">
				{#if indeterminate}
					<Minus class="size-4" />
				{:else if checked}
					<Check class="size-4" />
				{/if}
			</div>
		{/snippet}
	</Checkbox.Root>

	<Label.Root id={`${id}-label`} for={id} class="">
		<slot />
	</Label.Root>
</div>
