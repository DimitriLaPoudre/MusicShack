<script lang="ts">
	import { Dialog } from "bits-ui";

	export let open = false;
	export let title = "";
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger>
		<slot name="trigger" />
	</Dialog.Trigger>

	<Dialog.Portal>
		<Dialog.Overlay class="overlay" />

		<Dialog.Content class="content">
			<Dialog.Title>{title}</Dialog.Title>

			<slot />

			<Dialog.Close>
				<slot name="close">Close</slot>
			</Dialog.Close>
			<Dialog.Close>
				<slot name="validate">Validate</slot>
			</Dialog.Close>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
