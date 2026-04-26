<script lang="ts">
	import type { UserResponse } from "$lib/types/response/user";
	import { onMount } from "svelte";
	import User from "./User.svelte";
	import { GetUsers } from "$lib/usecases/user";
	import { toast } from "svelte-sonner";
	import Dialog from "../bits-ui/Dialog.svelte";
	import Button from "../bits-ui/Button.svelte";
	import Checkbox from "../bits-ui/Checkbox.svelte";

	let users = $state<UserResponse[]>([]);

	onMount(async () => {
		try {
			users = await GetUsers();
		} catch (e) {
			return toast.error(
				e instanceof Error ? e.message : "Unknown error",
			);
		}
	});
</script>

<div class="flex flex-col border">
	<h2>Users</h2>

	<!-- <Dialog title={"test"}> -->
	<!-- 	<Button slot="trigger" disabled={true}>Add</Button> -->
	<!---->
	<!-- 	<input class="" name={"username"} placeholder="Username" /> -->
	<!-- 	<Checkbox checked={false} name={"hi_res"} id={"hi_res"}>HiRes</Checkbox> -->
	<!-- </Dialog> -->

	<div class="flex flex-col gap-2 border p-1">
		{#each users as user}
			<User
				{user}
				reloadUsers={async () => {
					console.log("ntm");
					users = await GetUsers();
				}}
			/>
		{/each}
	</div>
</div>
