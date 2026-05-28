<script lang="ts">
	import { Dialog } from "bits-ui";
	import type { Snippet } from "svelte";

	type Props = {
		open?: boolean;
		title?: string;
		children?: Snippet;
		trigger?: Snippet;
		close?: Snippet;
		validate?: Snippet;
	};

	let {
		open = $bindable(false),
		title = "",
		children,
		trigger,
		close,
		validate,
	}: Props = $props();
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger>
		{@render trigger?.()}
	</Dialog.Trigger>

	<Dialog.Portal>
		<Dialog.Overlay class="overlay" />

		<Dialog.Content class="content">
			<Dialog.Title>{title}</Dialog.Title>

			{@render children?.()}

			<Dialog.Close>
				{@render close?.()}
			</Dialog.Close>

			<Dialog.Close>
				{@render validate?.()}
			</Dialog.Close>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
