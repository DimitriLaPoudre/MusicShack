<script lang="ts">
	import type { Snippet } from "svelte";
	import { Dialog, type WithoutChild } from "bits-ui";

	type Props = Dialog.RootProps & {
		buttonText: string;
		title: Snippet;
		description: Snippet;
		contentProps?: WithoutChild<Dialog.ContentProps>;
		// ...other component props if you wish to pass them
	};

	let {
		open = $bindable(false),
		children,
		buttonText,
		contentProps,
		title,
		description,
		...restProps
	}: Props = $props();
</script>

<Dialog.Root bind:open {...restProps}>
	<Dialog.Trigger class="hover-full">
		{buttonText}
	</Dialog.Trigger>
	<Dialog.Portal>
		<Dialog.Overlay />
		<Dialog.Content {...contentProps}>
			<Dialog.Title>
				{@render title()}
			</Dialog.Title>
			<Dialog.Description>
				{@render description()}
			</Dialog.Description>
			{@render children?.()}
			<Dialog.Close>Close</Dialog.Close>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
