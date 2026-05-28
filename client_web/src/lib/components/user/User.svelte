<script lang="ts">
	import type { UserResponse } from "$lib/types/response/user";
	import { Pencil, Trash } from "@lucide/svelte";
	import { DeleteUser, UpdateUser } from "$lib/usecases/user";
	import Button from "../bits-ui/Button.svelte";
	import { toast } from "svelte-sonner";
	import Dialog from "../bits-ui/Dialog.svelte";

	let updateUserOpen = false;
	let deleteUserOpen = false;

	export let user: UserResponse;
	export let reloadUsers: () => Promise<void>;
</script>

<div class="flex gap-1 border p-1">
	<div class="flex gap-1 border p-1">
		<p>Username</p>
		<p>{user.username}</p>
	</div>
	<div class="flex gap-1 border p-1">
		<p>HiRes</p>
		<p>{user.hi_res}</p>
	</div>
	<!-- <div class="border p-1"> -->
	<!-- 	<Dialog title={"Update User"} bind:open={updateUserOpen}> -->
	<!-- 		{#snippet trigger()} -->
	<!-- 			<Button> -->
	<!-- 				<Pencil /> -->
	<!-- 			</Button> -->
	<!-- 		{/snippet} -->
	<!-- 		<UpdateUser {afterUpdateUser} /> -->
	<!-- 	</Dialog> -->
	<!-- 	<Button -->
	<!-- 		onclick={async () => { -->
	<!-- 			try { -->
	<!-- 				// await UpdateUser(user.id); -->
	<!-- 				await reloadUsers(); -->
	<!-- 			} catch (e) { -->
	<!-- 				return toast.error( -->
	<!-- 					e instanceof Error ? e.message : "Unknown error", -->
	<!-- 				); -->
	<!-- 			} -->
	<!-- 		}} -->
	<!-- 	> -->
	<!-- 		<Pencil /> -->
	<!-- 	</Button> -->
	<!-- </div> -->
	<div class="border p-1">
		<Dialog
			title={`Delete user ${user.username} ?`}
			bind:open={deleteUserOpen}
		>
			{#snippet trigger()}
				<Button>
					<Trash />
				</Button>
			{/snippet}

			<form
				onsubmit={async (e: Event) => {
					e.preventDefault();

					try {
						await DeleteUser(user.id);
						await reloadUsers();
						deleteUserOpen = false;
					} catch (e) {
						return toast.error(
							e instanceof Error ? e.message : "Unknown error",
						);
					}
				}}
			>
				<Button>Confirm</Button>
			</form>
		</Dialog>
	</div>
</div>
